package ws_auth

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds the WS upload auth settings, all read from env.
type Config struct {
	AccessTTL         time.Duration // WS_ACCESS_TOKEN_TTL_MINUTES, default 5
	RefreshTTL        time.Duration // WS_REFRESH_TOKEN_TTL_MINUTES, default 30
	RenewGrace        time.Duration // WS_RENEW_GRACE_SECONDS, default 30
	MaxConcurrent     int64         // WS_UPLOAD_MAX_CONCURRENT, default 3 (per user)
	CounterTTL        time.Duration // WS_UPLOAD_COUNTER_TTL_MINUTES, default 120 (crash safety)
	TokenRatePerMin   int64         // WS_TOKEN_RATE_LIMIT_PER_MIN, default 10
	AccessSecret      string        // WS_JWT_ACCESS_SECRET
	RefreshSecret     string        // WS_JWT_REFRESH_SECRET
	SignMethod        string        // JWT_HMAC_HASH
	EncryptionKey     string        // WS_TOKEN_ENCRYPTION_KEY
	EncryptionKeyPrev string        // WS_TOKEN_ENCRYPTION_KEY_PREVIOUS (decrypt only, for key rotation)
}

func ConfigFromEnv() Config {
	return Config{
		AccessTTL:         time.Duration(envInt("WS_ACCESS_TOKEN_TTL_MINUTES", 5)) * time.Minute,
		RefreshTTL:        time.Duration(envInt("WS_REFRESH_TOKEN_TTL_MINUTES", 30)) * time.Minute,
		RenewGrace:        time.Duration(envInt("WS_RENEW_GRACE_SECONDS", 30)) * time.Second,
		MaxConcurrent:     int64(envInt("WS_UPLOAD_MAX_CONCURRENT", 3)),
		CounterTTL:        time.Duration(envInt("WS_UPLOAD_COUNTER_TTL_MINUTES", 120)) * time.Minute,
		TokenRatePerMin:   int64(envInt("WS_TOKEN_RATE_LIMIT_PER_MIN", 10)),
		AccessSecret:      os.Getenv("WS_JWT_ACCESS_SECRET"),
		RefreshSecret:     os.Getenv("WS_JWT_REFRESH_SECRET"),
		SignMethod:        os.Getenv("JWT_HMAC_HASH"),
		EncryptionKey:     os.Getenv("WS_TOKEN_ENCRYPTION_KEY"),
		EncryptionKeyPrev: os.Getenv("WS_TOKEN_ENCRYPTION_KEY_PREVIOUS"),
	}
}

func envInt(key string, def int) int {
	value, err := strconv.Atoi(strings.TrimSpace(os.Getenv(key)))
	if err != nil || value <= 0 {
		return def
	}
	return value
}
