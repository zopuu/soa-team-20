package config

import (
	"log"
	"os"
	"strings"
	"time"
)

type Config struct {
	Port string

	AuthBase  string
	StakeBase string
	BlogBase  string
	TourBase  string

	JwtSecret   string
	JwtIssuer   string
	JwtAudience string

	CorsOrigins []string

	// Timeouts
	DialTimeout       time.Duration
	ProxyTimeout      time.Duration
	FollowersGRPCAddr string
	TourGRPCAddr      string
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func New() *Config {
	c := &Config{
		Port:        getenv("GATEWAY_PORT", "7000"),
		AuthBase:    getenv("AUTH_BASE", "http://localhost:7001"),
		StakeBase:   getenv("STAKE_BASE", "http://stakeholders:8080"),
		BlogBase:    getenv("BLOG_BASE", "http://blogservice:8080"),
		TourBase:    getenv("TOUR_BASE", "http://localhost:5200"),
		JwtSecret:   getenv("JWT_SECRET", "CHANGE_ME"),
		JwtIssuer:   getenv("JWT_ISSUER", "AuthService"),
		JwtAudience: getenv("JWT_AUDIENCE", "AuthServiceClient"),
		CorsOrigins: strings.Split(getenv("CORS_ORIGINS",
			"http://localhost:4200,http://localhost:5173,http://localhost:8080"), ","),
		DialTimeout:       5 * time.Second,
		ProxyTimeout:      30 * time.Second,
		FollowersGRPCAddr: getenv("FOLLOWERS_GRPC_ADDR", "followers-service:50051"),
		TourGRPCAddr:      getenv("TOUR_GRPC_ADDR", "tourservice:50052"),
	}

	if c.JwtSecret == "CHANGE_ME" {
		log.Println("[WARN] JWT_SECRET is default; set a strong secret in env.")
	}
	return c
}
