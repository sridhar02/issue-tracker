package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func getUserPageHandler(c *gin.Context, db *sql.DB) {

	fmt.Println("user")

	username := c.Param("user_name")

	user, err := GetUserByUserName(db, username)

	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.HTML(http.StatusOK, "user.html", gin.H{"User": user})

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

type CurrentRepo struct {
	UserId string
	RepoId string
}

func getCurrentRepo(db *sql.DB, username string, repoName string) (CurrentRepo, error) {

	var userId, repoId string

	row := db.QueryRow(`SELECT users.id,repos.id FROM repos JOIN users ON 
		     repos.user_id = users.id WHERE users.username= $1 AND repos.name = $2`, username, repoName)
	err := row.Scan(&userId, &repoId)
	if err != nil {
		return CurrentRepo{}, err
	}

	currentRepo := CurrentRepo{
		UserId: userId,
		RepoId: repoId,
	}

	return currentRepo, nil
}

func getIssuesPageHandler(c *gin.Context, db *sql.DB) {

	_, err := authorize(c, db)

	authorized := true
	if err != nil {
		authorized = false
	}

	userName := c.Param("user_name")
	repoName := c.Param("repo_name")

	currentRepo, err := getCurrentRepo(db, userName, repoName)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	rows, err := db.Query(
		`SELECT issues.id,issues.title, users.username ,issues.issue_number 
		FROM issues JOIN users ON issues.user_id = users.id WHERE repo_id = $1;`,
		currentRepo.RepoId)
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

	c.HTML(http.StatusOK, "issues.html",
		gin.H{"UserName": userName,
			"RepoName":   repoName,
			"Issues":     issues,
			"Authorized": authorized})
}

type CommentsIssue struct {
	ID       int
	Username string
	Body     string
	issueId  int
	Image    string
}

func getIssuePageHandler(c *gin.Context, db *sql.DB) {
	currentUser, err := authorize(c, db)

	authorized := true
	if err != nil {
		authorized = false
	}

	username := c.Param("user_name")
	repoName := c.Param("repo_name")

	currentRepo, err := getCurrentRepo(db, username, repoName)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	issueNumberStr := c.Param("issue_number")

	var issueNumber, ID int
	var title, body, repoId, userId, status, image string

	row := db.QueryRow(
		`SELECT issues.id,issues.title, issues.body, users.username ,issues.issue_number,
		issues.repo_id, issues.user_id ,issues.status,users.image
		 FROM issues JOIN users ON issues.user_id = users.id WHERE issues.issue_number = $1 AND issues.repo_id=$2;`,
		issueNumberStr, currentRepo.RepoId)

	err = row.Scan(&ID, &title, &body, &username, &issueNumber, &repoId, &userId, &status, &image)
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
							comments.user_id = users.id WHERE comments.issue_id=$1`, ID)

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

	c.HTML(http.StatusOK, "issue.html", gin.H{
		"CurrentUser": currentUser,
		"Authorized":  authorized,
		"Issue":       issue,
		"Comments":    comments,
		"UserName":    c.Param("user_name"),
		"RepoName":    repoName})
}

func createIssueComment(c *gin.Context, db *sql.DB) {
	_, err := authorize(c, db)
	if err != nil {
		c.Redirect(http.StatusFound, "http://localhost:8000/login")
		return
	}

	username := c.Param("user_name")
	repoName := c.Param("repo_name")

	// currentRepo, err := getCurrentRepo(db, username, repoName)
	// if err != nil {
	// 	fmt.Println(err)
	// 	c.AbortWithStatus(http.StatusInternalServerError)
	// 	return
	// }

	body := c.PostForm("body")
	repoId := c.PostForm("repo_id")
	_issueId := c.PostForm("issue_id")
	userId := c.PostForm("user_id")
	IssueNumber := c.PostForm("issue_number")

	fmt.Println(c.PostForm("comment_and_close"))

	if c.PostForm("comment_and_close") == "1" {
		_, err := db.Exec(`UPDATE issues SET status = 'Closed' WHERE id = $1`, _issueId)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Redirect(http.StatusFound, "http://localhost:8000/"+username+"/"+repoName+"/issues/"+IssueNumber)
		return
	}

	if c.PostForm("comment_and_open") == "1" {
		_, err := db.Exec(`UPDATE issues SET status = 'Open' WHERE id = $1`, _issueId)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Redirect(http.StatusFound, "http://localhost:8000/"+username+"/"+repoName+"/issues/"+IssueNumber)
		return
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

	c.Redirect(http.StatusFound, "http://localhost:8000/"+username+"/"+repoName+"/issues/"+IssueNumber)
}

func getNewIssuePageHandler(c *gin.Context, db *sql.DB) {

	currentUser, err := authorize(c, db)
	if err != nil {
		c.Redirect(http.StatusFound, "http://localhost:8000/login")
		return
	}

	username := c.Param("user_name")
	repoName := c.Param("repo_name")

	currentRepo, err := getCurrentRepo(db, username, repoName)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.HTML(http.StatusOK, "issue_new.html",
		gin.H{"UserId": currentRepo.UserId,
			"RepoId":      currentRepo.RepoId,
			"UserName":    username,
			"CurrentUser": currentUser,
			"RepoName":    repoName})

}

func postNewIssuePageHandler(c *gin.Context, db *sql.DB) {

	_, err := authorize(c, db)

	if err != nil {
		c.Redirect(http.StatusFound, "http://localhost:8000/login")
		return
	}

	username := c.Param("user_name")
	repoName := c.Param("repo_name")

	currentRepo, err := getCurrentRepo(db, username, repoName)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	fmt.Println(username)

	title := c.PostForm("title")
	// repoId := c.PostForm("repo_id")
	body := c.PostForm("body")
	// userId := c.PostForm("user_id")

	issue := Issue{
		Title:  title,
		RepoId: currentRepo.RepoId,
		Body:   body,
		UserId: currentRepo.UserId,
	}

	issueNumber, err := CreateIssue(db, issue)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusFound, "http://localhost:8000/"+username+"/"+repoName+"/issues/"+strconv.Itoa(issueNumber))

}

func getRepoNewPageHandler(c *gin.Context, db *sql.DB) {
	currentUser, err := authorize(c, db)
	if err != nil {
		c.Redirect(http.StatusFound, "http://localhost:8000/login")
		return
	}

	c.HTML(http.StatusOK, "repo_new.html", gin.H{"CurrentUser": currentUser})
}

func PostRepoNewPageHandler(c *gin.Context, db *sql.DB) {
	_, err := authorize(c, db)

	if err != nil {
		c.Redirect(http.StatusFound, "http://localhost:8000/login")
		return
	}

	name := c.PostForm("name")
	userId := c.PostForm("user_id")
	userName := c.PostForm("user_name")
	description := c.PostForm("description")
	Type := c.PostForm("type")

	repo := Repo{
		Name:        name,
		UserId:      userId,
		Description: description,
		Type:        Type,
	}

	err = CreateRepo(db, repo)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusFound, "http://localhost:8000/"+userName+"/"+name+"/issues")

}

func getUserNewPageHandler(c *gin.Context, db *sql.DB) {

	currentUser, err := authorize(c, db)
	if err == nil {
		c.Redirect(http.StatusFound, "http://localhost:8000/"+currentUser.Username)
		return
	}

	c.HTML(http.StatusOK, "user_signup.html", gin.H{})
}

func PostUserNewPageHandler(c *gin.Context, db *sql.DB) {

	currentUser, err := authorize(c, db)
	if err == nil {
		c.Redirect(http.StatusFound, "http://localhost:8000/"+currentUser.Username)
		return
	}

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

	err = CreateUser(db, user)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusFound, "http://localhost:8000/")

}

func getUserSigninPageHandler(c *gin.Context, db *sql.DB) {

	currentUser, err := authorize(c, db)
	if err == nil {
		c.Redirect(http.StatusFound, "http://localhost:8000/"+currentUser.Username)
		return
	}

	c.HTML(http.StatusOK, "user_signin.html", gin.H{})

}

func PostUserSigninPageHandler(c *gin.Context, db *sql.DB) {

	currentUser, err := authorize(c, db)
	if err == nil {
		c.Redirect(http.StatusFound, "http://localhost:8000/"+currentUser.Username)
		return
	}

	username := c.PostForm("username")
	password := c.PostForm("password")

	// user := User{
	// 	Username: username,
	// 	Password: password,
	// }

	var Password, Id string
	row := db.QueryRow(`SELECT password,id FROM users WHERE username=$1`, username)

	err = row.Scan(&Password, &Id)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if Password != password {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	cookie, err := CreateCookie(db, Id)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.SetCookie("session", cookie, 300, "/", "localhost:8000", false, true)

	c.Redirect(http.StatusFound, "http://localhost:8000/new")

}

func authorize(c *gin.Context, db *sql.DB) (User, error) {
	session, err := c.Cookie("session")
	if err != nil {
		return User{}, err
	}

	cookie, err := ReadCookie(db, session)
	if err != nil {
		return User{}, err
	}
	if cookie.ExpiresAt.Before(time.Now()) {
		return User{}, fmt.Errorf("cookie expired")
	}
	currentUser, err := GetUser(db, cookie.UserId)
	if err != nil {
		return User{}, err
	}

	return currentUser, nil

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

	userGroup := router.Group("/:user_name")
	pages := userGroup.Group("/:repo_name")

	// router.GET("/user", func(c *gin.Context) { getUserPageHandler(c, db) })
	pages.GET("/issues",
		func(c *gin.Context) { getIssuesPageHandler(c, db) })
	pages.POST("/issues/new", func(c *gin.Context) { postNewIssuePageHandler(c, db) })
	pages.GET("/issues/:issue_number",
		func(c *gin.Context) {
			if c.Param("issue_number") == "new" {
				getNewIssuePageHandler(c, db)
			} else {
				getIssuePageHandler(c, db)
			}
		})

	pages.POST("/comments", func(c *gin.Context) { createIssueComment(c, db) })
	userGroup.GET("", func(c *gin.Context) {
		if c.Param("user_name") == "new" {
			getRepoNewPageHandler(c, db)
		} else if c.Param("user_name") == "login" {
			getUserSigninPageHandler(c, db)
		} else {
			getUserPageHandler(c, db)
		}
	})

	userGroup.POST("", func(c *gin.Context) {
		if c.Param("user_name") == "new" {
			PostRepoNewPageHandler(c, db)
		} else if c.Param("user_name") == "login" {
			PostUserSigninPageHandler(c, db)
		} else {
			c.Status(500)
		}
	})

	router.GET("", func(c *gin.Context) { getUserNewPageHandler(c, db) })
	router.POST("", func(c *gin.Context) { PostUserNewPageHandler(c, db) })

	// api := router.Group("/api")
	// api.GET("/users/:id", func(c *gin.Context) { getUserHandler(c, db) })
	// api.GET("/repos/:id", func(c *gin.Context) { getRepoHandler(c, db) })
	// api.GET("/issues/:id", func(c *gin.Context) { getIssueHandler(c, db) })
	// api.GET("/comments/:id", func(c *gin.Context) { getCommentHandler(c, db) })
	// api.DELETE("/users/:id", func(c *gin.Context) { deleteUserHandler(c, db) })
	// api.DELETE("/repos/:id", func(c *gin.Context) { deleteRepoHandler(c, db) })
	// api.DELETE("/issues/:id", func(c *gin.Context) { deleteIssueHandler(c, db) })
	// api.DELETE("/comments/:id", func(c *gin.Context) { deleteCommentHandler(c, db) })
	// api.POST("/users", func(c *gin.Context) { postUserHandler(c, db) })
	// api.POST("/repos", func(c *gin.Context) { postRepoHandler(c, db) })
	// api.POST("/issues", func(c *gin.Context) { postIssueHandler(c, db) })
	// api.POST("/comments", func(c *gin.Context) { postCommentHandler(c, db) })
	// api.PUT("/users", func(c *gin.Context) { putUserHandler(c, db) })
	// api.PUT("/repos", func(c *gin.Context) { putRepoHandler(c, db) })
	// api.PUT("/issues", func(c *gin.Context) { putIssueHandler(c, db) })
	// api.PUT("/comments", func(c *gin.Context) { putCommentHandler(c, db) })

	err = router.Run(":8000")
	if err != nil {
		log.Fatal(err)
	}

}
