package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"sugurta/internal/pkg/config"
	"sugurta/internal/pkg/helper"

	"github.com/casbin/casbin"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Authorizer(e *casbin.CachedEnforcer, logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data, ok := r.Context().Value(RequestAuthCtx).(map[string]string)
			if !ok {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			ok = e.Enforce(data["sub"], r.URL.Path, r.Method)
			// if err != nil {
			// 	logger.Error("middleware authorizer",
			// 		zap.Error(err),
			// 		zap.String("sub", data["sub"]),
			// 		zap.String("path", r.URL.Path),
			// 		zap.String("method", r.Method),
			// 	)
			// }
			if ok {
				next.ServeHTTP(w, r)
				return
			}
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		})
	}
}

func AuthMiddleware(e *casbin.Enforcer, config config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			userRole string
			act      = c.Request.Method
			obj      = c.FullPath()
		)

		token := c.GetHeader("Authorization")
		if token == "" {
			userRole = "unauthorized"
		}
		
		if userRole == "" {
			token = strings.TrimPrefix(token, "Bearer ")

			claims, err := helper.ParseJWT(token, config.JWT.Secret)
			if err != nil {
				userRole = "unauthorized"
			}

			v, ok := claims["role"].(string)
			if !ok {
				userRole = "unauthorized"
			} else {
				userRole = v
			}

			for key, value := range claims {
				c.Request.Header.Set(key, fmt.Sprintf("%v", value))
			}
		}

		// TO DO: Check if session is valid

		// if userRole != "unauthorized" {
		// 	session, err := h.UseCase.SessionRepo.GetSingle(c, entity.Id{ID: c.GetHeader("session_id")})
		// 	if err != nil {
		// 		fmt.Println("error while gettign single session", err)
		// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Session is invalid"})
		// 		return
		// 	}

		// 	if !session.IsActive {
		// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Session is not active"})
		// 		return
		// 	}
		// }

		ok := e.Enforce(userRole, obj, act)
		fmt.Println("role: ", userRole)
		fmt.Println("path: ", obj)
		fmt.Println("method: ", act)

		if !ok {

			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})

			return
		}

		c.Next()
	}
}
