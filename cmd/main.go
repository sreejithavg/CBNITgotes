package main

import (
	"OneDrive/Desktop/CBNIT/models"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	r.GET("/user/fetch/", models.FetchUsersHandler)
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print("Server Started")

	<-done
	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}
