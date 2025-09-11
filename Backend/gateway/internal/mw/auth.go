package mw

import (
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type JWTConfig struct {
	Secret   []byte
	Issuer   string
	Audience string
}

// AuthOptional lets requests pass through but, if a token is present, validates it
// and injects X-User-Id / X-Roles.
func AuthOptional(cfg JWTConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" || !strings.HasPrefix(strings.ToLower(auth), "bearer ") {
				next.ServeHTTP(w, r)
				return
			}

			tokenStr := strings.TrimSpace(auth[len("Bearer "):])
			tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				return cfg.Secret, nil
			}, jwt.WithAudience(cfg.Audience), jwt.WithIssuer(cfg.Issuer))
			if err == nil && tok != nil && tok.Valid {
				if claims, ok := tok.Claims.(jwt.MapClaims); ok {
					// Map your claim keys as used by Auth service
					if sub, ok := claims["sub"].(string); ok {
						r.Header.Set("X-User-Id", sub)
					} else if uid, ok := claims["uid"].(string); ok {
						r.Header.Set("X-User-Id", uid)
					}
					// roles could be array or string
					if roles, ok := claims["role"]; ok {
						switch v := roles.(type) {
						case []any:
							var s []string
							for _, it := range v {
								if str, ok := it.(string); ok {
									s = append(s, str)
								}
							}
							r.Header.Set("X-Roles", strings.Join(s, ","))
						case string:
							r.Header.Set("X-Roles", v)
						}
					}
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

// AuthRequired rejects if there is no valid JWT.
func AuthRequired(cfg JWTConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" || !strings.HasPrefix(strings.ToLower(auth), "bearer ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			tokenStr := strings.TrimSpace(auth[len("Bearer "):])
			_, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				return cfg.Secret, nil
			}, jwt.WithAudience(cfg.Audience), jwt.WithIssuer(cfg.Issuer))
			if err != nil {
				log.Printf("JWT parse error: %v", err) // Add this line
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
