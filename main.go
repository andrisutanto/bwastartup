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

	userByEmail, err := userRepository.FindByEmail("me@asndrisutanto.com")

	if err != nil {
		fmt.Println(err.Error())
	}

	if userByEmail.ID == 0 {
		fmt.Println("User tidak ditemukan")
	} else {
		fmt.Println(userByEmail.Name)
	}

	userHandler := handler.NewUserHandler(userService)
	router := gin.Default()
	//artinya, kalau ada yang akses ke api/v1, maka akan dilarikan ke RegisterUser
	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	router.Run()
}
