package controllers

import (
	"net/http"
	"strconv"

	"btpn-golang/app"
	"btpn-golang/database"
	"btpn-golang/helpers"
	"btpn-golang/models"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

func Register(context *gin.Context) {
	conn, error := database.Connect()

	if error != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "Database tidak terhubung",
		})
		return
	}

	var penggunaBaru app.UserData
	if error := context.BindJSON(&penggunaBaru); error != nil {
		context.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"status":  404,
			"message": "Invalid credentials",
		})
		return
	}

	insertUser := models.User{
		Username: penggunaBaru.Username,
		Email:    penggunaBaru.Email,
		Password: helpers.EncryptPassword(penggunaBaru.Password),
	}

	_, error = govalidator.ValidateStruct(insertUser)

	if error != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": error.Error(),
		})
		return
	}

	var emailCek models.User
	conn.Where("email = ?", penggunaBaru.Email).First(&emailCek)

	var unameCek models.User
	conn.Where("username = ?", penggunaBaru.Username).First(&unameCek)

	if emailCek.Email != "" || unameCek.Username != "" {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "Username atau email sudah terpakai",
		})
		return
	}

	result := conn.Create(&insertUser)

	if result.Error != nil {
		context.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"status":  404,
			"message": "Username atau email sudah terpakai",
		})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Berhasil mendaftar",
	})
}

func Login(context *gin.Context) {
	conn, error := database.Connect()

	if error != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "Database tidak terhubung",
		})

		return
	}

	var user models.User

	email := context.Query("email")
	password := context.Query("password")

	err := conn.Where("email = ?", email).First(&user).Error

	if err != nil || !helpers.CheckPassword(password, user.Password) {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  404,
			"message": "Invalid credentials",
		})
		return
	}

	token, err := helpers.GenerateToken(user)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "Error while generating token",
		})
		return
	}

	context.SetCookie("Authorization", token, 3600, "", "", true, true)

	context.IndentedJSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Berhasil login",
	})
}

func Logout(context *gin.Context) {

	_, err := context.Cookie("Authorization")

	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  401,
			"message": "Unauthorized",
		})
		return
	}

	context.SetCookie("Authorization", "", -1, "", "", true, true)

	context.IndentedJSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Berhasil logout, silahkan login kembali",
	})
}

func UpdateUser(context *gin.Context) {

	conn, error := database.Connect()

	if error != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "Database tidak terhubung",
		})
		return
	}

	updateID := context.Param("id")

	var penggunaBaru app.UserData
	if error := context.ShouldBindJSON(&penggunaBaru); error != nil {
		context.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"status":  404,
			"message": "Invalid credentials",
		})
		return
	}

	userData := context.MustGet("user").(models.User)

	if strconv.Itoa(userData.ID) != updateID {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  401,
			"message": "Tidak boleh mengakses aksi",
		})
		return
	}

	var user models.User
	conn.Where("id = ?", updateID).First(&user)

	var emailCek models.User
	conn.Where("email = ?", penggunaBaru.Email).First(&emailCek)

	var unameCek models.User
	conn.Where("username = ?", penggunaBaru.Username).First(&unameCek)

	if emailCek.Email != "" || unameCek.Username != "" {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "Username atau email sudah ada",
		})
		return
	}

	user.Username = penggunaBaru.Username
	user.Email = penggunaBaru.Email
	user.Password = helpers.EncryptPassword(penggunaBaru.Password)

	_, error = govalidator.ValidateStruct(user)

	if error != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": error.Error(),
		})
		return
	}

	conn.Save(&user)

	context.SetCookie("Authorization", "", -1, "", "", true, true)

	context.IndentedJSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Akun berhasil di update",
	})

}

func DeleteUser(context *gin.Context) {

	conn, error := database.Connect()

	if error != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "Database tidak terhubung",
		})
		return
	}

	deleteID := context.Param("id")

	var user models.User
	conn.Where("id = ?", deleteID).First(&user)

	userData := context.MustGet("user").(models.User)

	if strconv.Itoa(userData.ID) != deleteID {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  401,
			"message": "Tidak boleh mengakses aksi",
		})
		return
	}

	conn.Delete(&user)

	context.SetCookie("Authorization", "", -1, "", "", true, true)

	context.IndentedJSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Berhasil menghapus akun, silahkan daftar kembali atau login",
	})

}
