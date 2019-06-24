package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	// "strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func getUserPageHandler(c *gin.Context, db *sql.DB) {

	user, err := GetUser(db, "ac6f8b68-8f31-48ea-a436-05b9813b484b")

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.HTML(http.StatusOK, "user.html", user)

}

type IssueName struct {
	ID          int
	Title       string
	Username    string
	IssueNumber int
}

func getIssuesPageHandler(c *gin.Context, db *sql.DB) {

	rows, err := db.Query(
		`SELECT issues.id,issues.title, users.username ,issues.issue_number 
		FROM issues JOIN users ON issues.user_id = users.id WHERE repo_id = $1;`,
		"d360c6f3-60dc-4846-bb6a-0919a1817d5e")
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	issues := []IssueName{}

	var issueNumber, ID int
	var title, username string

	for rows.Next() {
		err = rows.Scan(&ID, &title, &username, &issueNumber)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		issue := IssueName{
			ID:          ID,
			Title:       title,
			Username:    username,
			IssueNumber: issueNumber,
		}

		issues = append(issues, issue)
	}

	c.HTML(http.StatusOK, "issues.html", issues)
}

func getIssuePageHandler(c *gin.Context, db *sql.DB) {

	Id := c.Param("id")

	var issueNumber, ID int
	var title, username string

	row := db.QueryRow(
		`SELECT issues.id,issues.title, users.username ,issues.issue_number 
		 FROM issues JOIN users ON issues.user_id = users.id WHERE issues.id = $1;`,
		Id)

	err := row.Scan(&ID, &title, &username, &issueNumber)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	issue := IssueName{
		ID:          ID,
		Title:       title,
		Username:    username,
		IssueNumber: issueNumber,
	}

	c.HTML(http.StatusOK, "issue.html", issue)
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

	router := gin.Default()

	router.LoadHTMLGlob("./templates/*")

	api := router.Group("/api")

	api.GET("/users/:id", func(c *gin.Context) { getUserHandler(c, db) })
	api.GET("/repos/:id", func(c *gin.Context) { getRepoHandler(c, db) })
	api.GET("/issues/:id", func(c *gin.Context) { getIssueHandler(c, db) })
	api.GET("/comments/:id", func(c *gin.Context) { getCommentHandler(c, db) })
	api.DELETE("/users/:id", func(c *gin.Context) { deleteUserHandler(c, db) })
	api.DELETE("/repos/:id", func(c *gin.Context) { deleteRepoHandler(c, db) })
	api.DELETE("/issues/:id", func(c *gin.Context) { deleteIssueHandler(c, db) })
	api.DELETE("/comments/:id", func(c *gin.Context) { deleteCommentHandler(c, db) })
	api.POST("/users", func(c *gin.Context) { postUserHandler(c, db) })
	api.POST("/repos", func(c *gin.Context) { postRepoHandler(c, db) })
	api.POST("/issues", func(c *gin.Context) { postIssueHandler(c, db) })
	api.POST("/comments", func(c *gin.Context) { postCommentHandler(c, db) })
	api.PUT("/users", func(c *gin.Context) { putUserHandler(c, db) })
	api.PUT("/repos", func(c *gin.Context) { putRepoHandler(c, db) })
	api.PUT("/issues", func(c *gin.Context) { putIssueHandler(c, db) })
	api.PUT("/comments", func(c *gin.Context) { putCommentHandler(c, db) })

	router.GET("/user", func(c *gin.Context) { getUserPageHandler(c, db) })
	router.GET("/issues", func(c *gin.Context) { getIssuesPageHandler(c, db) })
	router.GET("/issues/:id", func(c *gin.Context) { getIssuePageHandler(c, db) })

	err = router.Run(":8000")
	if err != nil {
		log.Fatal(err)
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
	// fmt.Println("updated issue")

	// err = DeleteIssue(db, 1)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("deleted issue")

	//CRUD functions for comments under this section

	// comment, err := GetComment(db, 1)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println(comment)

	// err = CreateComment(db, Comment{
	// 	UserId:    "ac6f8b68-8f31-48ea-a436-05b9813b484b",
	// 	Body:      "a new comment has been created in the database",
	// 	IssueId:   1,
	// 	RepoId:    "d360c6f3-60dc-4846-bb6a-0919a1817d5e",
	// 	CreatedAt: time.Now(),
	// 	UpdatedAt: time.Now()})

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println("created an Comment")

	// err = UpdateComment(db, Comment{
	// 	ID:        1,
	// 	Body:      "a new comment has been updated",
	// 	UpdatedAt: time.Now()})

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("updated comment")

	// err = DeleteComment(db, 1)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("deleted repo")

}
