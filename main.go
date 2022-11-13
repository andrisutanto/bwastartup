package main

import (
	"bwastartup/auth"
	"bwastartup/campaign"
	"bwastartup/handler"
	"bwastartup/helper"
	"bwastartup/user"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
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

	//campaign
	campaignRepository := campaign.NewRepository(db)

	/* ini untuk testing repository find campaign
	campaigns, err := campaignRepository.FindByUserID(1)
	fmt.Println("debug")
	fmt.Println("debug")
	fmt.Println("debug")
	fmt.Println(len(campaigns))

	for _, campaign := range campaigns {
		fmt.Println(campaign.Name)
		//ambil gambar

		if len(campaign.CampaignImages) > 0 {
			fmt.Println(campaign.CampaignImages[0].FileName)
		}
	}
	end testing repository */

	userService := user.NewService(userRepository)
	//campaign service
	campaignService := campaign.NewService(campaignRepository)
	//tambahkan auth service
	authService := auth.NewService()

	/*testing service campaign
	campaigns, _ := campaignService.FindCampaigns(1)
	fmt.Println(len(campaigns))
	end testing service campaign*/

	//validate token
	// token, err := authService.ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.FS47zAvPV0vYxBfVNZAhTO3qA5gHetGYc3_VjY19wLU")
	// if err != nil {
	// 	fmt.Println("ERROR")
	// 	fmt.Println("ERROR")
	// 	fmt.Println("ERROR")
	// }

	// if token.Valid {
	// 	fmt.Println("VALID")
	// 	fmt.Println("VALID")
	// 	fmt.Println("VALID")
	// } else {
	// 	fmt.Println("INVALID")
	// 	fmt.Println("INVALID")
	// 	fmt.Println("INVALID")
	// }

	//fmt.Println(authService.GenerateToken(1))

	//coba update avatar user 1
	//userService.SaveAvatar(1, "images/1-profile.png")

	/*START testing service create campaign
	input := campaign.CreateCampaignInput{}
	input.Name = "Penggalangan Dana Startup Mania"
	input.ShortDescription = "short Desc"
	input.Description = "Long Desc"
	input.GoalAmount = 100000000
	input.Perks = "hadiah satu, hadiah dua, hadiah tiga"
	//untuk input user seolah2 dari user ID 1
	inputUser, _ := userService.GetUserByID(2)
	input.User = inputUser
	//simpan campaign
	_, err = campaignService.CreateCampaign(input)
	if err != nil {
		log.Fatal(err.Error())
	}
	END testing service create campaign */

	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)

	router := gin.Default()

	//ini untuk gambar, jadi menggunakan static
	//parameter pertama adalah URL, parameter kedua adalah foldernya
	router.Static("/images", "./images")

	//artinya, kalau ada yang akses ke api/v1, maka akan dilarikan ke RegisterUser
	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	//tambahkan middleware di avatar
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)

	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)

	router.Run()
}

//buat middleware
//ambil nilai header authorization: bearer token
//dari header authoriztion, kita ambil nilai token saja
//validasi token tsb
//ambil user_id
//ambil user dari db berdasarkan user-id lewat service
//set context isinya user

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		//dapetin headernya
		authHeader := c.GetHeader("Authorization")

		//cari apakah ada bearer di header
		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			// kalau tidak ada, langsung abort dan beri status unauthorized
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		//kalau ada bearernya
		//bearer sprti ini: Bearer nilaiTokenDisini
		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		//validasi tokennya disini
		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		//ambil data di dalam token / payload / claim
		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		//by default kalau angka, di float64, harus diubah ke INT
		userID := int(claim["user_id"].(float64))

		//cari user nya disini
		user, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		//kalau tidak ada error, then set context user
		c.Set("currentUser", user)

	}
}
