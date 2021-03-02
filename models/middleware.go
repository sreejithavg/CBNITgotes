package models

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

//QueryMethods method inteface with the posgres
type QueryMethods interface {
	createUser(data User) error
	getUser(username string) User
}

// PostgressMiddleware to setup the
func PostgressMiddleware(ps *Postgres) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("ps", ps)
	}
}

//CreateToken creates the jwt token
func CreateToken(userid uint64) (string, error) {
	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}
