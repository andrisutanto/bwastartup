package main

import (
	"bwastartup/user"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/bwastartup?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	} else {
		fmt.Println("Successfully connected!")
	}

	//menampilkan data dari DB
	//nama plural, karena jamak
	var users []user.User

	/*
		length := len(users)
		fmt.Println(length)
		length = len(users)
		fmt.Println(length)
	*/

	//tipe pointer, mencari data di table users
	db.Find(&users)

	for _, user := range users {
		fmt.Println(user.Name)
		fmt.Println(user.Email)
		fmt.Println("============")
	}

}
