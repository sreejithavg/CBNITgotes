package models

import (
	"log"

	guuid "github.com/google/uuid"
)

func (ps *Postgres) createUser(data User) error {

	id := guuid.New()
	insertStatement := `INSERT INTO users ( id ,firstName ,LastName,username,password) VALUES ($1, $2, $3, $4,$5)`
	data.id = id.String()
	log.Println(data)
	_, err := ps.db.Exec(insertStatement, data.id, data.firstName, data.LastName, data.username, data.password)
	if err != nil {
		log.Println("failed to create user", err)
		return err
	}
	return nil
}

func (ps *Postgres) getUser(id string) User {
	selectQuery := `SELECT * FROM users WHERE id = $1`
	var user User
	err := ps.db.QueryRow(selectQuery, id).Scan(&user.id, &user.username, &user.password, &user.firstName, &user.LastName)
	if err != nil {
		log.Fatal("Failed to execute query: ", err)
	}
	return user
}

// func (ps *Postgres) getUsers(data User) error {
// 	ps.db.AutoMigrate(&data)
// 	ps.db.Create(&data)
// 	return nil
// }
