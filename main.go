package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	Body        string
	Username    string
	UserId      string
	RepoId      string
	IssueNumber int
	Status      string
	Image       string
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

type CommentsIssue struct {
	ID       int
	Username string
	Body     string
	issueId  int
	Image    string
}

func getIssuePageHandler(c *gin.Context, db *sql.DB) {

	Id := c.Param("id")

	var issueNumber, ID int
	var title, username, body, repoId, userId, status, image string

	row := db.QueryRow(
		`SELECT issues.id,issues.title, issues.body, users.username ,issues.issue_number,
		issues.repo_id, issues.user_id ,issues.status,users.image
		 FROM issues JOIN users ON issues.user_id = users.id WHERE issues.id = $1;`,
		Id)

	err := row.Scan(&ID, &title, &body, &username, &issueNumber, &repoId, &userId, &status, &image)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	issue := IssueName{
		ID:          ID,
		Title:       title,
		Body:        body,
		Username:    username,
		IssueNumber: issueNumber,
		RepoId:      repoId,
		UserId:      userId,
		Status:      status,
		Image:       image,
	}

	rows, err := db.Query(`SELECT comments.id,users.username,comments.body,
							comments.issue_id,users.image FROM comments JOIN users ON 
							comments.user_id = users.id WHERE comments.issue_id=$1`, Id)

	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	//gin.H = map[string]interface{}
	comments := []CommentsIssue{}

	var IssueId int
	var IMAGE string

	for rows.Next() {
		err = rows.Scan(&ID, &username, &body, &IssueId, &IMAGE)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		comment := CommentsIssue{
			ID:       ID,
			Username: username,
			Body:     body,
			issueId:  IssueId,
			Image:    IMAGE,
		}

		comments = append(comments, comment)
	}

	c.HTML(http.StatusOK, "issue.html", gin.H{"Issue": issue, "Comments": comments})
}

func createIssueComment(c *gin.Context, db *sql.DB) {
	body := c.PostForm("body")
	repoId := c.PostForm("repo_id")
	_issueId := c.PostForm("issue_id")
	userId := c.PostForm("user_id")

	fmt.Println(c.PostForm("comment_and_close"))

	if c.PostForm("comment_and_close") == "1" {
		_, err := db.Exec(`UPDATE issues SET status = 'Closed' WHERE id = $1`, _issueId)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Redirect(http.StatusFound, "http://localhost:8000/issues/"+_issueId)
	}

	if c.PostForm("comment_and_open") == "1" {
		_, err := db.Exec(`UPDATE issues SET status = 'Open' WHERE id = $1`, _issueId)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Redirect(http.StatusFound, "http://localhost:8000/issues/"+_issueId)
	}

	issueId, err := strconv.Atoi(_issueId)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	comment := Comment{
		Body:    body,
		UserId:  userId,
		RepoId:  repoId,
		IssueId: issueId,
	}

	err = CreateComment(db, comment)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusFound, "http://localhost:8000/issues/"+_issueId)
	c.Header("Content-Type", "application/html")
}

func getNewIssuePageHandler(c *gin.Context, db *sql.DB) {

	c.HTML(http.StatusOK, "issue_new.html", gin.H{"UserId": "ac6f8b68-8f31-48ea-a436-05b9813b484b",
		"RepoId": "d360c6f3-60dc-4846-bb6a-0919a1817d5e"})

}

func postNewIssuePageHandler(c *gin.Context, db *sql.DB) {

	title := c.PostForm("title")
	repoId := c.PostForm("repo_id")
	body := c.PostForm("body")
	userId := c.PostForm("user_id")

	issue := Issue{
		Title:  title,
		RepoId: repoId,
		Body:   body,
		UserId: userId,
	}

	id, err := CreateIssue(db, issue)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusFound, "http://localhost:8000/issues/"+strconv.Itoa(id))

}

// type RepoNew struct {
// 	Id     string
// 	Name   string
// 	UserId string
// }

func getRepoNewPageHandler(c *gin.Context, db *sql.DB) {

	c.HTML(http.StatusOK, "repo_new.html", gin.H{"UserId": "ac6f8b68-8f31-48ea-a436-05b9813b484b"})

}

func PostRepoNewPageHandler(c *gin.Context, db *sql.DB) {

	name := c.PostForm("name")
	userId := c.PostForm("user_id")
	description := c.PostForm("description")
	TYpe := c.PostForm("type")

	repo := Repo{
		Name:        name,
		UserId:      userId,
		Description: description,
		Type:        TYpe,
	}

	err := CreateRepo(db, repo)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusFound, "http://localhost:8000/issues")

}

func getUserNewPageHandler(c *gin.Context, db *sql.DB) {

	c.HTML(http.StatusOK, "user_signup.html", gin.H{})
}

func PostUserNewPageHandler(c *gin.Context, db *sql.DB) {

	name := c.PostForm("name")
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")

	user := User{
		Name:     name,
		Username: username,
		Email:    email,
		Password: password,
	}

	err := CreateUser(db, user)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusFound, "http://localhost:8000/issues")

}

func getUserSigninPageHandler(c *gin.Context, db *sql.DB) {

	c.HTML(http.StatusOK, "user_signin.html", gin.H{})

}

func PostUserSigninPageHandler(c *gin.Context, db *sql.DB) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// user := User{
	// 	Username: username,
	// 	Password: password,
	// }
	var Password string
	row := db.QueryRow(`SELECT password FROM users WHERE username=$1`, username)

	err := row.Scan(&Password)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if Password == password {
		c.Redirect(http.StatusFound, "http://localhost:8000/issues")
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connStr := "user=postgres dbname=issue_tracker host=localhost password=test1234 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
		return
	}

	// defecr fmt.Println("succesfully closed end")
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
	router.POST("/issues/new", func(c *gin.Context) { postNewIssuePageHandler(c, db) })
	router.GET("/issues/:id", func(c *gin.Context) {
		if c.Param("id") == "new" {
			getNewIssuePageHandler(c, db)
		} else {
			getIssuePageHandler(c, db)
		}
	})

	router.POST("/comments", func(c *gin.Context) { createIssueComment(c, db) })
	router.GET("/repos/new", func(c *gin.Context) { getRepoNewPageHandler(c, db) })
	router.POST("/repos/new", func(c *gin.Context) { PostRepoNewPageHandler(c, db) })
	router.GET("/user/sign-up", func(c *gin.Context) { getUserNewPageHandler(c, db) })
	router.POST("/user/sign-up", func(c *gin.Context) { PostUserNewPageHandler(c, db) })
	router.GET("/user/sign-in", func(c *gin.Context) { getUserSigninPageHandler(c, db) })
	router.POST("/user/sign-in", func(c *gin.Context) { PostUserSigninPageHandler(c, db) })

	err = router.Run(":8000")
	if err != nil {
		log.Fatal(err)
	}

}
