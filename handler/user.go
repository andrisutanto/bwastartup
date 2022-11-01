package handler

import (
	"bwastartup/helper"
	"bwastartup/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
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

	// token, err := h.jwtService.GenerateToken()

	formatter := user.FormatUser(newUser, "tokentes")

	//tambahkan helper response
	response := helper.APIResponse("Account has been registered", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}
