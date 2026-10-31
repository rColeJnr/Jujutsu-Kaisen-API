package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all application configuration, loaded from env variables.
type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	JWT       JWTConfig
	RateLimit RateLimitConfig
	Cache     CacheConfig
}

type ServerConfig struct {
	Port            string
	ReadTimeoutSec  int
	WriteTimeoutSec int
}

type DatabaseConfig struct {
	DSN string // SQLite Data Source Name
}

type JWTConfig struct {
	Secret      string
	ExpireHours int
}

type RateLimitConfig struct {
	RPS         float64
	Burst       int
	AuthPerMin  float64
	AuthBurst   int
	TrustHeader string
}

type CacheConfig struct {
	MaxAge time.Duration
}

// Load reads environment variables and returns a populated Config struct.
// It fails fast if required variables are missing.
func Load() (*Config, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("config: JWT_SECRET environment variable is required.")
	}

	dbDSN := os.Getenv("DATABASE_URL")
	if dbDSN == "" {
		return nil, fmt.Errorf("config: DATABASE_URL environment variable is required.")
	}

	readTimeout, _ := strconv.Atoi(os.Getenv("SERVER_READ_TIMEOUT_SEC"))
	if readTimeout == 0 {
		readTimeout = 10
	}
	writeTimeout, _ := strconv.Atoi(os.Getenv("SERVER_WRITE_TIMEOUT_SEC"))
	if writeTimeout == 0 {
		writeTimeout = 30
	}

	jwtExpire, _ := strconv.Atoi(os.Getenv(("JWT_EXPIRE_HOURS")))
	if jwtExpire == 0 {
		jwtExpire = 24
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8989"
	}

	cacheSec, err := intEnv("CACHE_MAX_AGE_SEC", 60, 0, 24*3600)
	if err != nil {
		return nil, err
	}

	rps := envFloat("RATE_LIMIT_RPS", 10)
	burst := int(envFloat("RATE_LIMIT_BURST", 20))
	authPerMin := envFloat("AUTH_RATE_LIMIT_PER_MIN", 5)
	authBurst := int(envFloat("AUTH_RATE_LIMIT_BURST", 5))

	return &Config{
		Server: ServerConfig{
			Port:            port,
			ReadTimeoutSec:  readTimeout,
			WriteTimeoutSec: writeTimeout,
		},
		Database: DatabaseConfig{DSN: dbDSN},
		JWT:      JWTConfig{Secret: secret, ExpireHours: jwtExpire},
		RateLimit: RateLimitConfig{
			RPS:         rps,
			Burst:       burst,
			AuthPerMin:  authPerMin,
			AuthBurst:   authBurst,
			TrustHeader: os.Getenv("TRUSTED_PROXY_HEADER"),
		},
		Cache: CacheConfig{MaxAge: time.Duration(cacheSec) * time.Second},
	}, nil
}

// envFloat reads a positive number from the environment, or returns default.
func envFloat(key string, def float64) float64 {
	v, err := strconv.ParseFloat(os.Getenv(key), 64)
	if err != nil || v <= 0 {
		return def
	}
	return v
}

// intEnv reads an integer in [min, max]. An unset variable gives def; a set
// but invalid one is an error (it used to silently fall back to the default).
func intEnv(key string, def, min, max int) (int, error) {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return def, nil
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < min || n > max {
		return 0, fmt.Errorf("config: %s must be an integer between %d and %d, got %q", key, min, max, raw)
	}
	return n, nil
}
