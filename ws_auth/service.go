package ws_auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/fari-99/go-helper/crypts"
	"github.com/fari-99/go-helper/token_generator"
	"github.com/redis/go-redis/v9"
)

const (
	accessTokenType  = "access_token"
	refreshTokenType = "refresh_token"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenReused  = errors.New("refresh token already used")
)

// Session is the value stored in Redis for a WS token pair.
// FamilyID is the uuid of the first pair of a rotation chain and stays the same across renewals.
type Session struct {
	Uuid     string    `json:"uuid"`
	FamilyID string    `json:"family_id"`
	UserID   string    `json:"user_id"`
	Username string    `json:"username"`
	IssuedAt time.Time `json:"issued_at"`
}

type TokenPair struct {
	AccessToken      string    `json:"access_token"`
	RefreshToken     string    `json:"refresh_token"`
	AccessExpiredAt  time.Time `json:"access_expired_at"`
	RefreshExpiredAt time.Time `json:"refresh_expired_at"`
}

type Service interface {
	Config() Config

	// Issue creates a new token pair (new family) for the user.
	Issue(ctx context.Context, req *http.Request, user token_generator.UserDetails) (*TokenPair, error)
	// Authenticate decrypts and validates an access token against JWT and the live Redis record.
	Authenticate(ctx context.Context, encryptedAccess string) (*Session, time.Time, error)
	// Renew rotates the pair using an encrypted refresh token, with reuse detection.
	Renew(ctx context.Context, req *http.Request, encryptedRefresh string) (*TokenPair, *Session, error)
	// SessionAlive reports whether the family has not been revoked or expired.
	SessionAlive(ctx context.Context, familyID string) bool

	// AllowTokenRequest is a fixed window (1 minute) limiter for token issue/refresh.
	AllowTokenRequest(ctx context.Context, key string) (bool, error)
	// AcquireUpload takes one concurrent-upload slot for the user. ok is false when the cap is reached.
	// release is safe to call more than once.
	AcquireUpload(ctx context.Context, userID string) (release func(), ok bool, err error)
}

type service struct {
	redis redis.UniversalClient
	cfg   Config
}

// NewService panics on missing secrets/keys so a misconfigured deploy fails at startup.
func NewService(redisClient redis.UniversalClient, cfg Config) Service {
	if cfg.AccessSecret == "" || cfg.RefreshSecret == "" || cfg.EncryptionKey == "" {
		panic("ws_auth: WS_JWT_ACCESS_SECRET, WS_JWT_REFRESH_SECRET and WS_TOKEN_ENCRYPTION_KEY must be set")
	}

	if cfg.AccessTTL > cfg.RefreshTTL {
		panic("ws_auth: WS_ACCESS_TOKEN_TTL_MINUTES must not exceed WS_REFRESH_TOKEN_TTL_MINUTES")
	}

	return &service{redis: redisClient, cfg: cfg}
}

func (s *service) Config() Config { return s.cfg }

func keyAccess(uuid string) string  { return uuid + ":ws_access_token" }
func keyRefresh(uuid string) string { return uuid + ":ws_refresh_token" }
func keyUsed(uuid string) string    { return uuid + ":ws_used" }
func keyFamily(id string) string    { return id + ":ws_family" }

func (s *service) Issue(ctx context.Context, req *http.Request, user token_generator.UserDetails) (*TokenPair, error) {
	signed, err := s.sign(req, user)
	if err != nil {
		return nil, err
	}

	// the first uuid is the family id for the whole rotation chain
	return s.persist(ctx, signed, signed.Uuid, user)
}

func (s *service) Authenticate(ctx context.Context, encryptedAccess string) (*Session, time.Time, error) {
	claims, err := s.parse(accessTokenType, encryptedAccess)
	if err != nil {
		return nil, time.Time{}, ErrInvalidToken
	}

	session, err := s.loadSession(ctx, keyAccess(claims.Uuid))
	if err != nil {
		return nil, time.Time{}, err
	}

	if !s.SessionAlive(ctx, session.FamilyID) {
		return nil, time.Time{}, ErrInvalidToken
	}

	var expiresAt time.Time
	if claims.ExpiresAt != nil {
		expiresAt = claims.ExpiresAt.Time
	}

	return session, expiresAt, nil
}

func (s *service) Renew(ctx context.Context, req *http.Request, encryptedRefresh string) (*TokenPair, *Session, error) {
	claims, err := s.parse(refreshTokenType, encryptedRefresh)
	if err != nil || claims.UserDetails == nil {
		return nil, nil, ErrInvalidToken
	}

	session, err := s.loadSession(ctx, keyRefresh(claims.Uuid))
	if errors.Is(err, ErrInvalidToken) {
		// record is gone: either expired, or this refresh token was already rotated (reuse)
		if familyID, getErr := s.redis.Get(ctx, keyUsed(claims.Uuid)).Result(); getErr == nil {
			s.revokeFamily(ctx, familyID)
			return nil, nil, ErrTokenReused
		}
		return nil, nil, ErrInvalidToken
	} else if err != nil {
		return nil, nil, err
	}

	// atomically claim the refresh token; losing the race means it is being reused
	claimed, err := s.redis.SetNX(ctx, keyUsed(claims.Uuid), session.FamilyID, s.cfg.RefreshTTL).Result()
	if err != nil {
		return nil, nil, err
	}
	if !claimed {
		s.revokeFamily(ctx, session.FamilyID)
		return nil, nil, ErrTokenReused
	}

	if !s.SessionAlive(ctx, session.FamilyID) {
		return nil, nil, ErrInvalidToken
	}

	signed, err := s.sign(req, *claims.UserDetails)
	if err != nil {
		return nil, nil, err
	}

	pair, err := s.persist(ctx, signed, session.FamilyID, *claims.UserDetails)
	if err != nil {
		return nil, nil, err
	}

	s.redis.Del(ctx, keyAccess(claims.Uuid), keyRefresh(claims.Uuid))

	newSession, err := s.loadSession(ctx, keyAccess(signed.Uuid))
	if err != nil {
		return nil, nil, err
	}

	return pair, newSession, nil
}

func (s *service) SessionAlive(ctx context.Context, familyID string) bool {
	exists, err := s.redis.Exists(ctx, keyFamily(familyID)).Result()
	return err == nil && exists == 1
}

func (s *service) AllowTokenRequest(ctx context.Context, key string) (bool, error) {
	redisKey := "ws_token_rate:" + key
	count, err := s.redis.Incr(ctx, redisKey).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		s.redis.Expire(ctx, redisKey, time.Minute)
	}

	return count <= s.cfg.TokenRatePerMin, nil
}

func (s *service) AcquireUpload(ctx context.Context, userID string) (func(), bool, error) {
	key := "ws_upload:active:" + userID

	count, err := s.redis.Incr(ctx, key).Result()
	if err != nil {
		return nil, false, err
	}

	// TTL is a crash safety net so a killed process can't lock a user out forever
	s.redis.Expire(ctx, key, s.cfg.CounterTTL)

	if count > s.cfg.MaxConcurrent {
		s.decrement(ctx, key)
		return nil, false, nil
	}

	var once sync.Once
	release := func() {
		once.Do(func() {
			releaseCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			s.decrement(releaseCtx, key)
		})
	}

	return release, true, nil
}

// decrement never leaves the counter below zero.
func (s *service) decrement(ctx context.Context, key string) {
	if count, err := s.redis.Decr(ctx, key).Result(); err == nil && count <= 0 {
		s.redis.Del(ctx, key)
	}
}

func (s *service) sign(req *http.Request, user token_generator.UserDetails) (*token_generator.SignedToken, error) {
	helper := token_generator.NewJwt(s.cfg.AccessSecret, s.cfg.RefreshSecret, s.cfg.SignMethod)
	if req != nil {
		helper = helper.SetCtx(req)
	}

	helper, err := helper.SetExpired(int(s.cfg.AccessTTL/time.Minute), time.Minute, int(s.cfg.RefreshTTL/time.Minute), time.Minute)
	if err != nil {
		return nil, err
	}

	helper, err = helper.SetClaim(user)
	if err != nil {
		return nil, err
	}

	return helper.SignClaims()
}

// persist stores the pair in Redis and returns it encrypted for the client.
func (s *service) persist(ctx context.Context, signed *token_generator.SignedToken, familyID string, user token_generator.UserDetails) (*TokenPair, error) {
	session := Session{
		Uuid:     signed.Uuid,
		FamilyID: familyID,
		UserID:   user.ID,
		Username: user.Username,
		IssuedAt: time.Now(),
	}

	value, err := json.Marshal(session)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.encrypt(signed.AccessToken)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.encrypt(signed.RefreshToken)
	if err != nil {
		return nil, err
	}

	pipe := s.redis.TxPipeline()
	pipe.Set(ctx, keyAccess(signed.Uuid), value, time.Until(signed.AccessExpiredAt))
	pipe.Set(ctx, keyRefresh(signed.Uuid), value, time.Until(signed.RefreshExpiredAt))
	pipe.Set(ctx, keyFamily(familyID), signed.Uuid, time.Until(signed.RefreshExpiredAt))
	if _, err = pipe.Exec(ctx); err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		AccessExpiredAt:  signed.AccessExpiredAt,
		RefreshExpiredAt: signed.RefreshExpiredAt,
	}, nil
}

// revokeFamily kills the current pair of the family; open sockets notice via SessionAlive.
func (s *service) revokeFamily(ctx context.Context, familyID string) {
	if currentUuid, err := s.redis.Get(ctx, keyFamily(familyID)).Result(); err == nil {
		s.redis.Del(ctx, keyAccess(currentUuid), keyRefresh(currentUuid))
	}

	s.redis.Del(ctx, keyFamily(familyID))
}

func (s *service) loadSession(ctx context.Context, key string) (*Session, error) {
	value, err := s.redis.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, ErrInvalidToken
	} else if err != nil {
		return nil, err
	}

	var session Session
	if err = json.Unmarshal([]byte(value), &session); err != nil {
		return nil, ErrInvalidToken
	}

	return &session, nil
}

func (s *service) encrypt(jwtToken string) (string, error) {
	encrypted, err := crypts.NewEncryptionBase().SetPassphrase(s.cfg.EncryptionKey).Encrypt([]byte(jwtToken))
	if err != nil {
		return "", fmt.Errorf("failed to encrypt token, err := %w", err)
	}

	return string(encrypted), nil
}

// decrypt tries the current key, then the previous one so key rotation doesn't drop live sessions.
// The token is attacker controlled, so recover from any panic in the crypto path.
func (s *service) decrypt(token string) (result string, err error) {
	defer func() {
		if recover() != nil {
			result, err = "", ErrInvalidToken
		}
	}()

	if token == "" {
		return "", ErrInvalidToken
	}

	keys := []string{s.cfg.EncryptionKey}
	if s.cfg.EncryptionKeyPrev != "" {
		keys = append(keys, s.cfg.EncryptionKeyPrev)
	}

	for _, key := range keys {
		decrypted, decryptErr := crypts.NewEncryptionBase().SetPassphrase(key).Decrypt([]byte(token))
		if decryptErr == nil {
			return string(decrypted), nil
		}
	}

	return "", ErrInvalidToken
}

func (s *service) parse(typeClaims, encryptedToken string) (*token_generator.JwtMapClaims, error) {
	rawToken, err := s.decrypt(encryptedToken)
	if err != nil {
		return nil, err
	}

	claims, err := token_generator.NewJwt(s.cfg.AccessSecret, s.cfg.RefreshSecret, s.cfg.SignMethod).ParseToken(typeClaims, rawToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
