package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type User struct {
	ID        string    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Username  string    `json:"username,omitempty"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func GetUser(db *sql.DB, id string) (User, error) {

	var name, username, email, createdAt, updatedAt, password string

	row := db.QueryRow("SELECT id,name, username, email,created_at, updated_at,password FROM users WHERE id=$1", id)
	err := row.Scan(&id, &name, &username, &email, &createdAt, &updatedAt, &password)
	if err != nil {
		return User{}, err
	}

	CreatedAt, err := time.Parse(time.RFC3339, createdAt)

	if err != nil {
		return User{}, err
	}

	UpdatedAt, err := time.Parse(time.RFC3339, updatedAt)

	if err != nil {
		return User{}, err
	}

	user := User{
		ID:        id,
		Name:      name,
		Username:  username,
		Email:     email,
		CreatedAt: CreatedAt,
		UpdatedAt: UpdatedAt,
		Password:  password,
	}

	return user, nil
}

func CreateUser(db *sql.DB, user User) error {

	ID := uuid.New().String()

	_, err := db.Exec(`INSERT INTO users(id,name, username, email, created_at, updated_at,password)
						VALUES($1,$2,$3,$4,$5,$6,$7)`,
		ID,
		user.Name,
		user.Username,
		user.Email,
		time.Now().Format(time.RFC3339),
		time.Now().Format(time.RFC3339),
		user.Password)
	if err != nil {
		return err
	}
	return nil
}

func UpdateUser(db *sql.DB, user User) error {

	_, err := db.Exec(`UPDATE users SET name = $1,username = $2,email =$3, updated_at = $4 WHERE id = $5`,
		user.Name, user.Username, user.Email, time.Now().Format(time.RFC3339), user.ID)

	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(db *sql.DB, id string) error {

	_, err := db.Exec(`DELETE FROM users WHERE id = $1`, id)

	if err != nil {
		return err
	}

	return nil
}

func getUserHandler(c *gin.Context, db *sql.DB) {

	id := c.Param("id")

	user, err := GetUser(db, id)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, user)

	// c.Status(http.StatusNoContent)

}

func postUserHandler(c *gin.Context, db *sql.DB) {

	user := User{}
	err := c.BindJSON(&user)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = CreateUser(db, user)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)

}

func putUserHandler(c *gin.Context, db *sql.DB) {

	user := User{}
	err := c.BindJSON(&user)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = UpdateUser(db, user)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

func deleteUserHandler(c *gin.Context, db *sql.DB) {

	id := c.Param("id")

	err := DeleteUser(db, id)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	// c.JSON(http.StatusOK, repo)

	c.Status(http.StatusNoContent)

}
