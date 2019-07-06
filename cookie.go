package main

import (
	// "fmt"
	"database/sql"
	// "github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

type Cookie struct {
	Id        string    `json:"id,omitempty"`
	UserId    string    `json:"user_id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

func CreateCookie(db *sql.DB, userId string) (string, error) {

	ID := uuid.New().String()

	_, err := db.Exec(`INSERT INTO cookies(id,user_id,created_at,expires_at)VALUES($1,$2,$3,$4)`,

		ID,
		userId,
		time.Now().UTC().Format(time.RFC3339),
		time.Now().UTC().Add(time.Minute*5).Format(time.RFC3339))

	if err != nil {
		return "", err
	}

	return ID, nil

}

func ReadCookie(db *sql.DB, id string) (Cookie, error) {

	var userId, createdAt, expiresAt string

	row := db.QueryRow("SELECT user_id,created_at, expires_at FROM cookies WHERE id=$1", id)
	err := row.Scan(&userId, &createdAt, &expiresAt)
	if err != nil {
		return Cookie{}, err
	}

	CreatedAt, err := time.Parse(time.RFC3339, createdAt)

	if err != nil {
		return Cookie{}, err
	}

	ExpiresAt, err := time.Parse(time.RFC3339, expiresAt)

	if err != nil {
		return Cookie{}, err
	}

	cookie := Cookie{
		Id:        id,
		UserId:    userId,
		CreatedAt: CreatedAt,
		ExpiresAt: ExpiresAt,
	}

	return cookie, nil
}

func DeleteCookie(db *sql.DB, id string) error {

	_, err := db.Exec(`DELETE FROM cookies WHERE id = $1`, id)

	if err != nil {
		return err
	}

	return nil
}
