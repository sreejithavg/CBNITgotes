package main

import (
	"fmt"

	"OneDrive/Desktop/CBNIT/models"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("main functions")
	postgressDB := models.GetPostgressInstance()
	defer postgressDB.Close()

	r := gin.Default()
	r.Use(models.PostgressMiddleware(postgressDB))
	r.POST("/user/signup/", models.SignupHandler)
	r.POST("/user/login/", models.LoginHandler)
	// r.POST("/user/fetch/",models)
	r.Run()
}
