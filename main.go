package main

import (
	"bwastartup/user"
	"log"

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

	userInput := user.RegisterUserInput{}
	userInput.Name = "Tes simpan dari service"
	userInput.Email = "contoh@gmail.com"
	userInput.Occupation = "Engineer"
	userInput.Password = "password"

	userService.RegisterUser(userInput)
}
