package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

type User struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:email,omitempty`
	// Password   string `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func GetUser(db *sql.DB, id string) (User, error) {

	var name, username, email, createdAt, updatedAt string

	row := db.QueryRow("SELECT name, username, email,created_at, updated_at FROM users WHERE id=$1", id)
	err := row.Scan(&name, &username, &email, &createdAt, &updatedAt)
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
		Name:      name,
		Username:  username,
		Email:     email,
		CreatedAt: CreatedAt,
		UpdatedAt: UpdatedAt,
	}

	return user, nil
}

func CreateUser(db *sql.DB, user User) error {

	_, err := db.Exec(`INSERT INTO users(id,name, username, email, created_at, updated_at)
						VALUES($1,$2,$3,$4,$5,$6)`,
		user.ID,
		user.Name,
		user.Username,
		user.Email,
		user.CreatedAt.Format(time.RFC3339),
		user.UpdatedAt.Format(time.RFC3339))
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

type Repo struct {
	Id         string    `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	UserId     string    `json:"user_id,omitempty"`
	IssueCount int       `json:"issue_count,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

func GetRepo(db *sql.DB, id string) (repo, error) {

	var name, userId, createdAt, updatedAt string
	var issueCount int

	row := db.QueryRow("SELECT name, user_id,issueCount,created_at, updated_at FROM repos WHERE id=$1", id)
	err := row.Scan(&name, &userId, &issueCount, &createdAt, &updatedAt)
	if err != nil {
		return User{}, err
	}

	repo := Repos{
		Name:        name,
		UserId:      user_id,
		IssuesCount: issue_count,
		CreatedAt:   created_at,
		UpdatedAt:   updated_at,
	}

}

func main() {

	connStr := "user=postgres dbname=issue_tracker host=localhost password=test1234 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
		return
	}

	// defer fmt.Println("succesfully closed end")
	defer db.Close()
	// defer fmt.Println("succesfully closed")

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return
	}

	// user, err := GetUser(db, "ac6f8b68-8f31-48ea-a436-05b9813b484b")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println(user)

	// err = CreateUser(db, User{
	// 	ID:        "ac6f8b68-8f31-48ea-a436-05b9813b484b",
	// 	Name:      "sridhar",
	// 	Username:  "sridhar02",
	// 	Email:     "kattasridhar02@gmail.com",
	// 	CreatedAt: time.Now(),
	// 	UpdatedAt: time.Now()})

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("created an user")

	// err = UpdateUser(db, User{
	// 	ID:       "ac6f8b68-8f31-48ea-a436-05b9813b484b",
	// 	Name:     "vramana",
	// 	Username: "vramana08",
	// 	Email:    "vramana@gmail.com",
	// 	// CreatedAt: time.Now()
	// 	UpdatedAt: time.Now()})

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("updated issue")

	// err = DeleteUser(db, "ac6f8b68-8f31-48ea-a436-05b9813b484b")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("deleted issue")

}
