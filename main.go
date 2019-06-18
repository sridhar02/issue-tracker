package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Users struct {
	Id       int    `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:email,omitempty`
	// Password   string `json:"password,omitempty"`
	Created_at string `json:"created_at,omitempty"`
	Updated_at string `json:"updated_at,omitempty"`
}

func GetUser(db *sql.DB, id int) (User, error) {

	var name, username, email, created_at, updated_at string

	row := db.QueryRow("SELECT name, username, email,created_at, updated_at FROM users WHERE id=$1", id)
	err := row.Scan(&name, &username, &email, &created_at, &updated_at)
	if err != nil {
		return User{}, err
	}

	user := User{
		Name:       name,
		Username:   username,
		Email:      email,
		Created_at: created_at,
		Updated_at: updated_at,
	}

	return user, nil
}

func CreateIssue(db *sql.DB, user User) error {

	_, err := db.Exec(`INSERT INTO users(username, email, created_at, updated_at)
													VALUES($1,$2,$3,$4)`, User.Username, User.Email, User.Created_at, User.Updated_at)
	if err != nil {
		return err
	}
	return nil
}

func UpdateIssue(db *sql.DB, user User, id int) error {

	_, err := db.Exec(`UPDATE issues SET username = $1,email =$2,created_at = $3,updated_at = $4 WHERE number = $5`, User.Username, User.Email, User.Created_at, User.Updated_at)

	if err != nil {
		return err
	}
	return nil
}

func DeletedIssue(db *sql.DB, id int) error {

	_, err := db.Exec(`DELETE FROM users WHERE id = $1`, id)

	if err != nil {
		return err
	}

	return nil
}
