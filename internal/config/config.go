package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all application configuration, loaded from env variables.
type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	JWT       JWTConfig
	RateLimit RateLimitConfig
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
