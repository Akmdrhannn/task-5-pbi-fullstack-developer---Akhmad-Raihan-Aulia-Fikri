package middlewares

import (
	"btpn-golang/database"
	"btpn-golang/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PhotoAuthorization() gin.HandlerFunc {

	return func(context *gin.Context) {

		photoId := context.Param("id")

		conn, error := database.Connect()

		if error != nil {
			context.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
				"status":  500,
				"message": "Tidak bisa mengakses database",
			})
			return
		}

		var photo models.Photo
		conn.Where("id = ?", photoId).First(&photo)

		if photo.ID == 0 {
			context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"status":  404,
				"message": "Foto tidak ditemukan",
			})
			return
		}

		userData := context.MustGet("user").(models.User)

		if photo.UserID != userData.ID {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  401,
				"message": "Tidak boleh melakukan aksi",
			})
			return
		}

		context.Next()
	}
}
