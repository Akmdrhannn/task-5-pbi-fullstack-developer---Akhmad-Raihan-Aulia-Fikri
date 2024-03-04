package middlewares

import (
	"net/http"
	"time"

	"btpn-golang/database"
	"btpn-golang/helpers"
	"btpn-golang/models"

	"github.com/gin-gonic/gin"
)

func JwtCheck() gin.HandlerFunc {

	return func(context *gin.Context) {
		auth_token, err := context.Cookie("Authorization")

		if err != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  401,
				"message": "Tidak diizinkan. Token tidak ditemukan.",
			})
			return
		}

		claims, err := helpers.ParseToken(auth_token)

		if err != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  401,
				"message": "Tidak diizinkan",
			})
			return
		}

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  4000,
				"message": "Token hangus atau expired",
			})
			return
		}

		var user models.User
		conn, err := database.Connect()

		conn.Where("email = ?", claims["email"]).First(&user)

		if user.ID == 0 || err != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  4000,
				"message": "Token tidak valid",
			})
			return
		}

		context.Set("user", user)
		context.Next()
	}
}
