package models

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "sree@11081995"
	dbname   = "postgres"
)

// User type is default request info
type User struct {
	ID        string `json:"ID"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}
