package controllers

import (
	"btpn-golang/app"
	"btpn-golang/database"
	"btpn-golang/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPhoto(context *gin.Context) {

	conn, error := database.Connect()

	if error != nil {
		context.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"status":  404,
			"message": "Invalid credentials",
		})

		return
	}

	var photos []models.Photo
	conn.Find(&photos)

	context.IndentedJSON(http.StatusOK, gin.H{
		"result":  photos,
		"status":  200,
		"message": "Success",
	})
}
func CreatePhoto(context *gin.Context) {

	conn, error := database.Connect()

	if error != nil {
		context.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"status":  404,
			"message": "Invalid credentials",
		})
		return
	}

	var fotoBaru app.PhotoData
	if error := context.BindJSON(&fotoBaru); error != nil {
		context.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"status":  404,
			"message": "Invalid credentials",
		})
		return
	}

	userData := context.MustGet("user").(models.User)

	insertPhoto := models.Photo{
		Title:    fotoBaru.Title,
		Caption:  fotoBaru.Caption,
		PhotoUrl: fotoBaru.PhotoUrl,
		UserID:   userData.ID,
	}

	conn.Create(&insertPhoto)

	context.IndentedJSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Berhasil menambahkan foto baru",
	})
}

func UpdatePhoto(context *gin.Context) {

	updateID := context.Param("id")

	conn, error := database.Connect()

	if error != nil {
		context.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"status":  404,
			"message": "Invalid credentials",
		})
		return
	}

	var fotoBaru app.PhotoData
	if error := context.BindJSON(&fotoBaru); error != nil {
		context.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"status":  404,
			"message": "Invalid credentials",
		})
		return
	}

	var photo models.Photo
	conn.Where("id = ?", updateID).First(&photo)

	photo.Title = fotoBaru.Title
	photo.Caption = fotoBaru.Caption
	photo.PhotoUrl = fotoBaru.PhotoUrl

	conn.Save(&photo)

	context.IndentedJSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Deskripsi Foto berhasil di update",
	})
}

func DeletePhoto(context *gin.Context) {

	deleteID := context.Param("id")

	conn, error := database.Connect()

	if error != nil {
		context.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"status":  404,
			"message": "Invalid credentials",
		})
		return
	}

	var photo models.Photo
	conn.Where("id = ?", deleteID).First(&photo)

	conn.Delete(&photo)

	context.IndentedJSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Foto berhasil dihapus",
	})
}
