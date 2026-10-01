package ws_auth

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/fari-99/go-helper/token_generator"
	"github.com/redis/go-redis/v9"
)

func newTestService(t *testing.T) (Service, *miniredis.Miniredis) {
	t.Helper()
	t.Setenv("USER_DETAILS_PASSPHRASE", "test-passphrase")

	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	cfg := Config{
		AccessTTL: 5 * 60e9, RefreshTTL: 30 * 60e9, RenewGrace: 30e9,
		MaxConcurrent: 3, CounterTTL: 3600e9, TokenRatePerMin: 2,
		AccessSecret: "access-secret", RefreshSecret: "refresh-secret", SignMethod: "HS256",
		EncryptionKey: "encryption-key",
	}

	return NewService(client, cfg), mr
}

var testUser = token_generator.UserDetails{ID: "42", Username: "fari", Email: "f@x.io"}

func TestIssueAuthenticate(t *testing.T) {
	svc, _ := newTestService(t)
	ctx := context.Background()

	pair, err := svc.Issue(ctx, nil, testUser)
	if err != nil {
		t.Fatalf("issue: %v", err)
	}

	if strings.Count(pair.AccessToken, ".") != 0 {
		t.Error("access token looks like a plain JWT, it must be an opaque encrypted blob")
	}

	session, exp, err := svc.Authenticate(ctx, pair.AccessToken)
	if err != nil {
		t.Fatalf("authenticate: %v", err)
	}
	if session.UserID != "42" || session.FamilyID != session.Uuid || exp.IsZero() {
		t.Errorf("unexpected session %+v exp %v", session, exp)
	}

	if _, _, err = svc.Authenticate(ctx, pair.RefreshToken); err == nil {
		t.Error("refresh token must not authenticate as an access token")
	}
	if _, _, err = svc.Authenticate(ctx, pair.AccessToken+"x"); err == nil {
		t.Error("tampered token must be rejected")
	}
	if _, _, err = svc.Authenticate(ctx, ""); err == nil {
		t.Error("empty token must be rejected")
	}
}

func TestAccessTokenRevokedWhenRedisRecordGone(t *testing.T) {
	svc, mr := newTestService(t)
	ctx := context.Background()

	pair, _ := svc.Issue(ctx, nil, testUser)
	session, _, _ := svc.Authenticate(ctx, pair.AccessToken)

	mr.Del(keyAccess(session.Uuid))
	if _, _, err := svc.Authenticate(ctx, pair.AccessToken); !errors.Is(err, ErrInvalidToken) {
		t.Errorf("err = %v, want ErrInvalidToken after redis record removed", err)
	}
}

func TestRenewRotatesAndKeepsFamily(t *testing.T) {
	svc, _ := newTestService(t)
	ctx := context.Background()

	first, _ := svc.Issue(ctx, nil, testUser)
	firstSession, _, _ := svc.Authenticate(ctx, first.AccessToken)

	second, session, err := svc.Renew(ctx, nil, first.RefreshToken)
	if err != nil {
		t.Fatalf("renew: %v", err)
	}

	if session.FamilyID != firstSession.FamilyID || session.Uuid == firstSession.Uuid {
		t.Errorf("family must be kept and uuid rotated: %+v vs %+v", session, firstSession)
	}
	if session.UserID != "42" || session.Username != "fari" {
		t.Errorf("user details lost on renew: %+v", session)
	}

	if _, _, err = svc.Authenticate(ctx, first.AccessToken); err == nil {
		t.Error("old access token must stop working after rotation")
	}
	if _, _, err = svc.Authenticate(ctx, second.AccessToken); err != nil {
		t.Errorf("new access token must work: %v", err)
	}
}

func TestRefreshReuseRevokesWholeFamily(t *testing.T) {
	svc, _ := newTestService(t)
	ctx := context.Background()

	first, _ := svc.Issue(ctx, nil, testUser)
	second, session, err := svc.Renew(ctx, nil, first.RefreshToken)
	if err != nil {
		t.Fatalf("renew: %v", err)
	}

	// attacker (or buggy client) replays the already used refresh token
	if _, _, err = svc.Renew(ctx, nil, first.RefreshToken); !errors.Is(err, ErrTokenReused) {
		t.Fatalf("err = %v, want ErrTokenReused", err)
	}

	if svc.SessionAlive(ctx, session.FamilyID) {
		t.Error("family must be revoked after reuse")
	}
	if _, _, err = svc.Authenticate(ctx, second.AccessToken); err == nil {
		t.Error("the legitimately rotated access token must be dead after reuse detection")
	}
	if _, _, err = svc.Renew(ctx, nil, second.RefreshToken); err == nil {
		t.Error("the rotated refresh token must be dead after reuse detection")
	}
}

func TestEncryptionKeyRotation(t *testing.T) {
	svc, _ := newTestService(t)
	ctx := context.Background()
	pair, _ := svc.Issue(ctx, nil, testUser)

	rotated := svc.(*service)
	rotated.cfg.EncryptionKeyPrev = rotated.cfg.EncryptionKey
	rotated.cfg.EncryptionKey = "brand-new-key"

	if _, _, err := svc.Authenticate(ctx, pair.AccessToken); err != nil {
		t.Errorf("token encrypted with previous key must still authenticate: %v", err)
	}

	rotated.cfg.EncryptionKeyPrev = ""
	if _, _, err := svc.Authenticate(ctx, pair.AccessToken); err == nil {
		t.Error("token must fail once previous key is dropped")
	}
}

func TestUploadConcurrencyCap(t *testing.T) {
	svc, mr := newTestService(t)
	ctx := context.Background()

	var releases []func()
	for i := 0; i < 3; i++ {
		release, ok, err := svc.AcquireUpload(ctx, "42")
		if err != nil || !ok {
			t.Fatalf("slot %d: ok=%v err=%v", i, ok, err)
		}
		releases = append(releases, release)
	}

	if _, ok, _ := svc.AcquireUpload(ctx, "42"); ok {
		t.Fatal("4th concurrent upload must be rejected")
	}
	if _, ok, _ := svc.AcquireUpload(ctx, "43"); !ok {
		t.Fatal("another user must have its own cap")
	}

	releases[0]()
	releases[0]() // double release must not free a second slot
	if _, ok, _ := svc.AcquireUpload(ctx, "42"); !ok {
		t.Fatal("slot must be free after release")
	}
	if _, ok, _ := svc.AcquireUpload(ctx, "42"); ok {
		t.Fatal("double release freed an extra slot")
	}

	for _, r := range releases {
		r()
	}
	if v, _ := mr.Get("ws_upload:active:42"); v != "" && v != "1" {
		t.Logf("counter after releases: %q", v)
	}
}

func TestTokenRateLimit(t *testing.T) {
	svc, mr := newTestService(t)
	ctx := context.Background()

	for i := 0; i < 2; i++ {
		if ok, _ := svc.AllowTokenRequest(ctx, "u1"); !ok {
			t.Fatalf("request %d should pass", i)
		}
	}
	if ok, _ := svc.AllowTokenRequest(ctx, "u1"); ok {
		t.Fatal("3rd request in the window must be limited")
	}

	mr.FastForward(61e9)
	if ok, _ := svc.AllowTokenRequest(ctx, "u1"); !ok {
		t.Fatal("limit must reset after the window")
	}
}

func TestGarbageTokensNeverPanic(t *testing.T) {
	svc, _ := newTestService(t)
	ctx := context.Background()

	for _, token := range []string{"", "a", "AAAA", strings.Repeat("A", 5), strings.Repeat("_-", 40), "not base64 !!"} {
		if _, _, err := svc.Authenticate(ctx, token); err == nil {
			t.Errorf("token %q must be rejected", token)
		}
		if _, _, err := svc.Renew(ctx, nil, token); err == nil {
			t.Errorf("refresh token %q must be rejected", token)
		}
	}
}
