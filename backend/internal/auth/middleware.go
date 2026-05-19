package auth

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
			})
			c.Abort()
			return
		}

		tokenString := strings.Replace(
			authHeader,
			"Bearer ",
			"",
			1,
		)

		token, err := jwt.ParseWithClaims(
			tokenString,
			&JWTClaims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(
					os.Getenv("JWT_ACCESS_SECRET"),
				), nil
			},
		)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*JWTClaims)

		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
			})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)

		c.Next()
	}
}
