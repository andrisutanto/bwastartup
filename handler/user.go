package handler

import (
	"bwastartup/auth"
	"bwastartup/helper"
	"bwastartup/user"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	//tangkap input dari user
	//map input dari user ke struct RegisterUserInput
	//struct diatas kita passing sebagai parameter service

	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		//untuk menampilkan pesan error lebih rapi
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Register Account Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return //kalau ada return, berarti akan di stop saat ada error
	}

	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		response := helper.APIResponse("Register Account Failed", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return //kalau ada return, berarti akan di stop saat ada error
	}

	token, err := h.authService.GenerateToken(newUser.ID)
	if err != nil {
		response := helper.APIResponse("Register Account Failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return //kalau ada return, berarti akan di stop saat ada error
	}

	formatter := user.FormatUser(newUser, token)

	//tambahkan helper response
	response := helper.APIResponse("Account has been registered", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(c *gin.Context) {
	//user memasukkan input (email & password)
	//input ditangkap oleh handler
	//mapping dari input user ke input struct
	//input struct passing ke service
	//di service mencari dgn bantuan repository user dengan email X
	//mencocokkan dengan password

	var input user.LoginInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		//untuk menampilkan pesan error lebih rapi
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Login Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return //kalau ada return, berarti akan di stop saat ada error
	}

	//buat variable utk simpan login
	loggedinUser, err := h.userService.Login(input)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := helper.APIResponse("Login Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return //kalau ada return, berarti akan di stop saat ada error
	}

	token, err := h.authService.GenerateToken(loggedinUser.ID)
	if err != nil {
		response := helper.APIResponse("Login Failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return //kalau ada return, berarti akan di stop saat ada error
	}

	formatter := user.FormatUser(loggedinUser, token)

	//tambahkan helper response
	response := helper.APIResponse("Success login", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) CheckEmailAvailability(c *gin.Context) {
	// proses cek email
	// ada input email dari user
	// input email di mapping ke struct input - di handler
	// struct input di passing ke service
	// service akan manggil repository - email sudah ada atau belum
	// repository - db
	// end cek email

	// ambil inputan dari user
	var input user.CheckEmailInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		//untuk menampilkan pesan error lebih rapi
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Email Check Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return //kalau ada return, berarti akan di stop saat ada error
	}

	// check dari service
	isEmailAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		errorMessage := gin.H{"errors": "Server error"}
		response := helper.APIResponse("Email Check Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return //kalau ada return, berarti akan di stop saat ada error
	}

	// pakai gin.H utk beri response data
	data := gin.H{
		"is_available": isEmailAvailable,
	}

	// buat variable baru utk cek dan beri response json
	var metaMessage string

	if isEmailAvailable {
		metaMessage = "Email is available"
	} else {
		metaMessage = "Email has been registered"
	}

	// ini kalau balasan berhasil / tidak ada email kembar
	response := helper.APIResponse(metaMessage, http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) UploadAvatar(c *gin.Context) {
	//input dari user
	//simpan gambar di folder images/
	//di service kita panggil repo
	//JWT, sementara hardcode, ambil yang ID login=1
	//repo ambil data user ID=1
	//repo update data user simpan lokasi file

	//untuk menangkap file
	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	//dapatkan user ID dari Contetxt currentUser
	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID

	//tambahkan user ID pada filename agar tidak kembar
	//path := "images/" + file.Filename
	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.userService.SaveAvatar(userID, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	//kalau tidak ada error, maka
	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Avatar successfully uploaded", http.StatusOK, "success", data)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) FetchUser(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(user.User)
	formatter := user.FormatUser(currentUser, "")
	response := helper.APIResponse("Successfuly fetch user data", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}
