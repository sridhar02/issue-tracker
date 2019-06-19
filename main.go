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
	ID          string    `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	UserId      string    `json:"user_id,omitempty"`
	IssuesCount int       `json:"issue_count,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

func GetRepo(db *sql.DB, id string) (Repo, error) {

	var name, userId, createdAt, updatedAt string
	var issueCount int

	row := db.QueryRow("SELECT name, user_id,issue_count,created_at, updated_at FROM repos WHERE id=$1", id)
	err := row.Scan(&name, &userId, &issueCount, &createdAt, &updatedAt)
	if err != nil {
		return Repo{}, err
	}

	CreatedAt, err := time.Parse(time.RFC3339, createdAt)

	if err != nil {
		return Repo{}, err
	}

	UpdatedAt, err := time.Parse(time.RFC3339, updatedAt)

	if err != nil {
		return Repo{}, err
	}

	repo := Repo{
		Name:        name,
		UserId:      userId,
		IssuesCount: issueCount,
		CreatedAt:   CreatedAt,
		UpdatedAt:   UpdatedAt,
	}

	return repo, nil
}

func CreateRepo(db *sql.DB, repo Repo) error {

	_, err := db.Exec(`INSERT INTO repos(id,name, user_id, issue_count, created_at, updated_at)
						VALUES($1,$2,$3,$4,$5,$6)`,
		repo.ID,
		repo.Name,
		repo.UserId,
		repo.IssuesCount,
		repo.CreatedAt.Format(time.RFC3339),
		repo.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return err
	}
	return nil
}

func UpdateRepo(db *sql.DB, repo Repo) error {

	_, err := db.Exec(`UPDATE repos SET name = $1,user_id = $2,issue_count =$3, updated_at = $4 WHERE id = $5`,
		repo.Name, repo.UserId, repo.IssuesCount, time.Now().Format(time.RFC3339), repo.ID)

	if err != nil {
		return err
	}
	return nil
}

func DeleteRepo(db *sql.DB, id string) error {

	_, err := db.Exec(`DELETE FROM repos WHERE id = $1`, id)

	if err != nil {
		return err
	}

	return nil
}

type Issue struct {
	ID          int       `json:"id,omitempty"`
	Title       string    `json:"title,omitempty"`
	UserId      string    `json:"user_id"`
	Body        string    `json:body,omitempty`
	RepoId      string    `repos_id,omitempty`
	IssueNumber int       `issue_number,omitempty`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

func GetIssue(db *sql.DB, id int) (Issue, error) {

	var title, userId, body, repoId, createdAt, updatedAt string
	var issueNumber int

	row := db.QueryRow("SELECT title, user_id,body,repos_id,issue_number,created_at, updated_at FROM issues WHERE id=$1", id)
	err := row.Scan(&title, &userId, &body, &repoId, &issueNumber, &createdAt, &updatedAt)
	if err != nil {
		return Issue{}, err
	}

	CreatedAt, err := time.Parse(time.RFC3339, createdAt)

	if err != nil {
		return Issue{}, err
	}

	UpdatedAt, err := time.Parse(time.RFC3339, updatedAt)

	if err != nil {
		return Issue{}, err
	}

	issue := Issue{
		Title:       title,
		UserId:      userId,
		Body:        body,
		RepoId:      repoId,
		IssueNumber: issueNumber,
		CreatedAt:   CreatedAt,
		UpdatedAt:   UpdatedAt,
	}

	return issue, nil
}

func CreateIssue(db *sql.DB, issue Issue, repoId string) error {

	var issueCount int

	row := db.QueryRow("SELECT issue_count FROM repos WHERE id=$1", issue.RepoId)
	err := row.Scan(&issueCount)
	if err != nil {
		return err
	}

	_, err = db.Exec(`INSERT INTO issues(title, user_id,body,repos_id,issue_number, created_at, updated_at)
						VALUES($1,$2,$3,$4,$5,$6,$7)`,
		issue.Title,
		issue.UserId,
		issue.Body,
		issue.RepoId,
		issueCount+1,
		issue.CreatedAt.Format(time.RFC3339),
		issue.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return err
	}

	_, err = db.Exec(`UPDATE repos SET issue_count =$1, updated_at = $2 WHERE id = $3`,
		issueCount+1, time.Now().Format(time.RFC3339), issue.RepoId)

	if err != nil {
		return err
	}

	return nil
}

func UpdateIssue(db *sql.DB, issue Issue) error {

	_, err := db.Exec(`UPDATE issues SET title = $1,body =$2, updated_at =$3 WHERE id = $4`,
		issue.Title, issue.Body, time.Now().Format(time.RFC3339), issue.ID)

	if err != nil {
		return err
	}
	return nil
}

func DeleteIssue(db *sql.DB, id int) error {

	_, err := db.Exec(`DELETE FROM issues WHERE id = $1`, id)

	if err != nil {
		return err
	}

	return nil
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

	//CRUD functions for user under this statement

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

	// CRUD functions for repos under this statement

	// repo, err := GetRepo(db, "d360c6f3-60dc-4846-bb6a-0919a1817d5e")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println(repo)

	// err = CreateRepo(db, Repo{
	// 	ID:          "d360c6f3-60dc-4846-bb6a-0919a1817d5e",
	// 	Name:        "sridhar",
	// 	UserId:      "ac6f8b68-8f31-48ea-a436-05b9813b484b",
	// 	IssuesCount: 0,
	// 	CreatedAt:   time.Now(),
	// 	UpdatedAt:   time.Now()})

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("created an Repo")

	// err = UpdateRepo(db, Repo{
	// 	ID:          "d360c6f3-60dc-4846-bb6a-0919a1817d5e",
	// 	Name:        "vramana",
	// 	UserId:      "ac6f8b68-8f31-48ea-a436-05b9813b484b",
	// 	IssuesCount: 10,
	// 	// CreatedAt: time.Now()
	// 	UpdatedAt: time.Now()})

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("updated repo")

	// err = DeleteRepo(db, "d360c6f3-60dc-4846-bb6a-0919a1817d5e")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("deleted repo")

	//CRUD functions for issues under this section

	// issue, err := GetIssue(db, 1)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println(issue)

	// err = CreateIssue(db, Issue{
	// 	Title:       "new issue",
	// 	UserId:      "ac6f8b68-8f31-48ea-a436-05b9813b484b",
	// 	Body:        "a new has been created in the data base of issues",
	// 	RepoId:      "d360c6f3-60dc-4846-bb6a-0919a1817d5e",
	// 	IssueNumber: 0,
	// 	CreatedAt:   time.Now(),
	// 	UpdatedAt:   time.Now()},
	// 	"d360c6f3-60dc-4846-bb6a-0919a1817d5e")

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("created an Issue")

	// err = UpdateIssue(db, Issue{
	// 	ID:    1,
	// 	Title: "vramana",
	// 	Body:  "a new issue has been updated",
	// 	// IssuesCount: 10,
	// 	// CreatedAt: time.Now()
	// 	UpdatedAt: time.Now()})

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("updated repo")

	// err = DeleteIssue(db, 1)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("deleted repo")

}
