package models

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	bindingFailedMsg = "failed to bind the request"
	aKey             = "jdnfksdmfksd"
)

type status struct {
	StatusCode  int    `json:"status_code"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

var errorResponse struct {
	Status status
}

// SignupHandler to create a new user
func SignupHandler(c *gin.Context) {
	log.Println("signupHandler invoked")
	ps := c.MustGet("ps").(QueryMethods)
	type response struct {
		Status status `json:"status"`
	}
	var request User
	if err := c.ShouldBindJSON(&request); err != nil {
		errorResponder(err, "SignupHandler", bindingFailedMsg, http.StatusBadRequest, c)
		return
	}
	log.Println(request)
	err := ps.createUser(request)
	if err != nil {
		errorResponder(err, "SignupHandler", err.Error(), http.StatusInternalServerError, c)
		return
	}
	var res response
	res.Status.StatusCode = http.StatusOK
	res.Status.Description = "created the user successfully"
	res.Status.Status = http.StatusText(http.StatusOK)
	log.Println("response", res)
	c.JSON(http.StatusOK, res)
}

// LoginHandler to create a new user
func LoginHandler(c *gin.Context) {
	log.Println("LoginHandler invoked")
	ps := c.MustGet("ps").(QueryMethods)
	type response struct {
		Status status `json:"status"`
		Token  string `json:"token"`
	}

	var request User
	if err := c.ShouldBindJSON(&request); err != nil {
		errorResponder(err, "LoginHandler", bindingFailedMsg, http.StatusBadRequest, c)
		return
	}
	user, err := ps.getUser(request.Username, request.Password)
	if err != nil {
		errorResponder(err, "LoginHandler", err.Error(), http.StatusInternalServerError, c)
		return
	}
	maker, err := NewJWTMaker(aKey)
	if err != nil {
		log.Println("err at the new jwt token maket")
		errorResponder(err, "LoginHandler", err.Error(), http.StatusInternalServerError, c)
		return
	}
	duration := time.Minute
	token, err := maker.CreateToken(user.Username, duration)
	if err != nil {
		log.Println("error while creating token", err)
		errorResponder(err, "LoginHandler", err.Error(), http.StatusInternalServerError, c)
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	var res response
	res.Status.StatusCode = http.StatusOK
	res.Status.Description = "created the user successfully"
	res.Status.Status = http.StatusText(http.StatusOK)
	res.Token = token
	log.Println("response", res)
	c.JSON(http.StatusOK, res)
}

// FetchUsersHandler to get all user
func FetchUsersHandler(c *gin.Context) {
	log.Println("FetchUsersHandler invoked")
	token := c.GetHeader("access_token")
	ps := c.MustGet("ps").(QueryMethods)

	type response struct {
		Status status `json:"status"`
		Data   []User `json:"data"`
	}

	maker, err := NewJWTMaker(aKey)
	payload, err := maker.ValidateToken(token)
	if err != nil {
		log.Println("error at validating the token ", err)
		errorResponder(err, "FetchUsersHandler", err.Error(), http.StatusUnauthorized, c)
		return
	}
	log.Println(payload)
	data, err := ps.getUsers()
	if err != nil {
		errorResponder(err, "FetchUsersHandler", err.Error(), http.StatusInternalServerError, c)
		return
	}

	var res response
	res.Status.StatusCode = http.StatusOK
	res.Status.Description = "created the user successfully"
	res.Status.Status = http.StatusText(http.StatusOK)
	res.Data = data
	log.Println("response", res)
	c.JSON(http.StatusOK, res)
}

func errorResponder(err error, method string, description string, httpCode int, ctx *gin.Context) {
	log.Println("error ::", err, "occured at ", method)
	errorResponse.Status.Description = description
	errorResponse.Status.Status = "failed"
	errorResponse.Status.StatusCode = httpCode
	log.Println(errorResponse)
	ctx.JSON(httpCode, errorResponse)
}
