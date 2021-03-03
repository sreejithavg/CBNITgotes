package models

import (
	"log"
	"strings"

	"github.com/google/uuid"
)

func (ps *Postgres) createUser(data User) error {

	id := uuid.New()
	insertStatement := "INSERT INTO users ( id,first_name,last_name,username,password) VALUES ($1, $2, $3, $4,$5)"
	data.ID = strings.Replace(id.String(), "-", "", -1)
	log.Println(data)
	res, err := ps.db.Exec(insertStatement, data.ID, data.FirstName, data.LastName, data.Username, data.Password)
	log.Print(res)

	if err != nil {
		log.Println("failed to create user", err)
		return err
	}
	return nil
}

func (ps *Postgres) getUser(username string, password string) (User, error) {
	log.Println("invoked handlers::getUser")
	log.Println(username)
	selectQuery := `SELECT * FROM users WHERE username = $1 AND password =$2`
	var user User
	err := ps.db.QueryRow(selectQuery, username, password).Scan(&user.LastName, &user.Password, &user.Username, &user.ID, &user.FirstName)
	if err != nil {
		log.Println("Failed to execute query: ", err)
		return user, err
	}
	return user, nil
}

func (ps *Postgres) getUsers() ([]User, error) {
	log.Println("invoked models :: getUser")
	var users []User
	rows, err := ps.db.Query("SELECT * FROM users WHERE id IS NOT NULL")
	if err != nil {
		log.Println("error at getUSers :: ", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		err = rows.Scan(&user.LastName, &user.Password, &user.Username, &user.ID, &user.FirstName)
		if err != nil {
			// handle this error
			log.Println("error at scanning the user", err)
			return nil, err
		}
		users = append(users, user)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		log.Println("error at the time of iteration", err)
		return nil, err
	}
	return users, nil
}
