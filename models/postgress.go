package models

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	//to support required driver for database/sql
	_ "github.com/lib/pq"
)

var (
	instance *Postgres
	once     sync.Once
)

//Postgres db pointer
type Postgres struct {
	db *sql.DB
}

//Connection function establish the connection to the postgressDB
func connection() *Postgres {
	log.Println("establishing connection")
	instance := new(Postgres)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalln("Connection failed :: ", err)
	}
	_, err = db.Exec("CREATE TABLE $1", "users_info")
	if err != nil {
		log.Println("error at creationg user table", err)
	}
	log.Println("created the table user_info")
	instance.db = db
	log.Println("established the connection", db, instance.db)
	return instance
}

// GetPostgressInstance get the signleton instance of the postgress
func GetPostgressInstance() *Postgres {
	// Singleton instance
	once.Do(func() {
		instance = connection()
	})
	return instance
}

//Close function closing the postgress DB
func (pSQL *Postgres) Close() {
	pSQL.db.Close()
	log.Println("Stopping postgress DB")
}
