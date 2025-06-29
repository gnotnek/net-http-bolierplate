package auth

import (
	"context"
	"net-http-boilerplate/internal/api/resp"
	"net-http-boilerplate/internal/pkg/jwt"
	"net/http"
	"strings"
)

type MiddlewareService struct {
	jwtService *jwt.JWT
}

func NewMiddleware(jwtService *jwt.JWT) *MiddlewareService {
	return &MiddlewareService{jwtService: jwtService}
}

func (m *MiddlewareService) AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			resp.WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "Token is missing"})
			return
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(tokenString, bearerPrefix) {
			resp.WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "Token must start with Bearer"})
			return
		}

		token := strings.TrimPrefix(tokenString, bearerPrefix)
		if token == "" {
			resp.WriteJSON(w, http.StatusUnauthorized, map[string]interface{}{
				"message": "Token is missing",
			})
			return
		}

		claims, err := m.jwtService.ParseToken(token)
		if err != nil {
			resp.WriteJSON(w, http.StatusUnauthorized, map[string]interface{}{
				"message": "Invalid token",
			})
			return
		}

		ctx = context.WithValue(ctx, "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}
