package models

import (
	"fmt"
	"log"
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	bindingFailedMsg = "failed to bing the request"
)

type status struct {
	statusCode  int
	description string
	status      string
}

var errorResponse struct {
	Status status
}

// SignupHandler to create a new user
func SignupHandler(c *gin.Context) {
	log.Println("signupHandler invoked")
	ps := c.MustGet("ps").(QueryMethods)
	type response struct {
		status status
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
	res.status.statusCode = http.StatusOK
	res.status.description = "created the user successfully"
	res.status.status = http.StatusText(http.StatusOK)
	log.Println("response", res)
	c.JSON(http.StatusOK, res)
}

// LoginHandler to create a new user
func LoginHandler(c *gin.Context) {
	log.Println("LoginHandler invoked")
	ps := c.MustGet("ps").(QueryMethods)
	fmt.Println(ps)
	type response struct {
		status status
		token  string
	}

	var request User
	if err := c.ShouldBindJSON(&request); err != nil {
		errorResponder(err, "LoginHandler", bindingFailedMsg, http.StatusBadRequest, c)
		return
	}

	err := ps.createUser(request)
	if err != nil {
		errorResponder(err, "LoginHandler", err.Error(), http.StatusInternalServerError, c)
		return
	}
	id, err := strconv.ParseUint(request.id, 10, 64)
	if err != nil {
		errorResponder(err, "LoginHandler", err.Error(), http.StatusInternalServerError, c)
		return
	}
	token, err := CreateToken(id)
	if err != nil {
		errorResponder(err, "LoginHandler", err.Error(), http.StatusInternalServerError, c)
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	var res response
	res.status.statusCode = http.StatusOK
	res.status.description = "created the user successfully"
	res.status.status = http.StatusText(http.StatusOK)
	res.token = token
	log.Println("response", res)
	c.JSON(http.StatusOK, res)
}

// // LoginHandler to create a new user
// func FetchUsersHandler(c *gin.Context) {
// 	log.Println("signupHandler invoked")
// 	ps := c.MustGet("ps").(QueryMethods)
// 	fmt.Println(ps)
// 	type response struct {
// 		status status
// 	}

// 	var request User
// 	if err := c.ShouldBindJSON(&request); err != nil {
// 		errorResponder(err, "SignupHandler", bindingFailedMsg, http.StatusBadRequest, c)
// 		return
// 	}

// 	err := ps.createUser(request)
// 	if err != nil {
// 		errorResponder(err, "SignupHandler", err.Error(), http.StatusInternalServerError, c)
// 		return
// 	}

// 	var res response
// 	res.status.statusCode = http.StatusOK
// 	res.status.description = "created the user successfully"
// 	res.status.status = http.StatusText(http.StatusOK)
// 	log.Println("response", res)
// 	c.JSON(http.StatusOK, res)
// }

func errorResponder(err error, method string, description string, httpCode int, ctx *gin.Context) {
	log.Fatalln("error ::", err, "occured at ", method)
	errorResponse.Status.description = description
	errorResponse.Status.description = "failed"
	errorResponse.Status.statusCode = httpCode
	ctx.JSON(httpCode, errorResponse)
}
