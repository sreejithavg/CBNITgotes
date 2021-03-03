package models

import (
	"errors"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

//different types of error return during the token validation
var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

//QueryMethods method inteface with the posgres
type QueryMethods interface {
	createUser(data User) error
	getUser(username string, password string) (User, error)
	getUsers() ([]User, error)
}

//Payload contains the payload data of token
type Payload struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// JWTTokenMaker struct for JWT token maker
type JWTTokenMaker struct {
	SecretKey string
}

//Maker interface for creating and validating the token
type Maker interface {
	// CreateToken create JWT token
	CreateToken(username string, duration time.Duration) (string, error)
	//ValidateToken validate the provided token
	ValidateToken(token string) (*Payload, error)
}

// PostgressMiddleware to setup the
func PostgressMiddleware(ps *Postgres) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("ps", ps)
	}
}

// NewPayload creates new tokens payload with a specific username and duration
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		log.Println("error at creating the tokenID : ", err)
		return nil, err
	}
	Payload := &Payload{
		ID:        tokenID.String(),
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return Payload, nil
}

// NewJWTMaker creates new JwtTokenMaker
func NewJWTMaker(secretKey string) (Maker, error) {
	return &JWTTokenMaker{secretKey}, nil
}

//Valid validate the token payload
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

// CreateToken create JWT token
func (maker *JWTTokenMaker) CreateToken(username string, duration time.Duration) (string, error) {
	Payload, err := NewPayload(username, duration)
	if err != nil {
		log.Println("error at CreateToken", err)
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Payload)
	log.Println("created the jwtToken")
	return jwtToken.SignedString([]byte(maker.SecretKey))
}

//ValidateToken validate the provided token
func (maker *JWTTokenMaker) ValidateToken(token string) (*Payload, error) {
	keyFunction := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.SecretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunction)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}
