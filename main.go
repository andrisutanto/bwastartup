package main

import (
	"bwastartup/handler"
	"bwastartup/user"
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

	//coba update avatar user 1
	userService.SaveAvatar(1, "images/1-profile.png")

	userHandler := handler.NewUserHandler(userService)
	router := gin.Default()
	//artinya, kalau ada yang akses ke api/v1, maka akan dilarikan ke RegisterUser
	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	router.Run()
}
