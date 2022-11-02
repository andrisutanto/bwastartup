package main

import (
	"bwastartup/handler"
	"bwastartup/user"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/bwastartup?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	//contoh buat user
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	input := user.LoginInput{
		Email:    "me@andrisutanto.com",
		Password: "password",
	}

	//ini untuk cek login berdasarkan inputan user, passing ke service utk cek
	user, err := userService.Login(input)
	if err != nil {
		fmt.Println("Terjadi kesalahan")
		fmt.Println(err.Error())
	}

	fmt.Println(user.Email)
	fmt.Println(user.Name)

	userHandler := handler.NewUserHandler(userService)
	router := gin.Default()
	//artinya, kalau ada yang akses ke api/v1, maka akan dilarikan ke RegisterUser
	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	router.Run()
}
