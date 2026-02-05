package middleware

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/bilalabdelkadir/chis/pkg/response"
	"github.com/redis/go-redis/v9"
)

func RateLimit(rdb *redis.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ip := extractIP(r)
			key := fmt.Sprintf("rl:%s:%d", ip, time.Now().Unix()/60)

			count, err := rdb.Incr(ctx, key).Result()
			if err != nil {
				slog.Warn("rate limiter redis error", "error", err)
				next.ServeHTTP(w, r)
				return
			}

			if count == 1 {
				rdb.Expire(ctx, key, 2*time.Minute)
			}

			if count > 100 {
				response.WriteJSON(w, http.StatusTooManyRequests, map[string]string{"message": "rate limit exceeded"})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func extractIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if ip := strings.TrimSpace(strings.SplitN(xff, ",", 2)[0]); ip != "" {
			return ip
		}
	}

	if xri := r.Header.Get("X-Real-Ip"); xri != "" {
		return strings.TrimSpace(xri)
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
