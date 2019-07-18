package main

import (
	"database/sql"
	"fmt"
	// "github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
	"time"
)

func getUserPageHandler(c *gin.Context, db *sql.DB) {

	_, err := authorize(c, db)
	authorized := true
	if err != nil {
		authorized = false
	}

	username := c.Param("user_name")
	user, err := GetUserByUserName(db, username)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	rows, err := db.Query(
		`SELECT repos.name FROM repos JOIN users ON repos.user_id = users.id WHERE users.username = $1;`,
		username)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	repoNames := []string{}

	var repoName string

	for rows.Next() {
		err = rows.Scan(&repoName)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		repoNames = append(repoNames, repoName)
	}

	c.HTML(http.StatusOK, "user.html", gin.H{"User": user, "RepoNames": repoNames, "Authorized": authorized})
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
	Pinned      string
	Lock        string
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
		`SELECT issues.id,issues.title, users.username ,issues.issue_number,issues.pinned,issues.status
		 FROM issues JOIN users ON issues.user_id = users.id WHERE repo_id = $1;`,
		currentRepo.RepoId)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	issues := []IssueName{}
	var issueNumber, ID int
	var title, username, pinned, status string
	for rows.Next() {
		err = rows.Scan(&ID, &title, &username, &issueNumber, &pinned, &status)
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
			Pinned:      pinned,
			Status:      status,
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
	var title, body, repoId, userId, status, image, pinned, lock string
	row := db.QueryRow(
		`SELECT issues.id,issues.title, issues.body, users.username ,issues.issue_number,
		issues.repo_id, issues.user_id ,issues.status,users.image,issues.pinned,issues.lock
		FROM issues JOIN users ON issues.user_id = users.id WHERE issues.issue_number = $1 AND issues.repo_id=$2;`,
		issueNumberStr, currentRepo.RepoId)
	err = row.Scan(&ID, &title, &body, &username, &issueNumber, &repoId, &userId, &status, &image, &pinned, &lock)
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
		Pinned:      pinned,
		Lock:        lock,
	}

	var PinnedIssuesCount int

	row = db.QueryRow(`SELECT count(*) FROM issues WHERE pinned = 'Pinned' AND repo_id = $1`, repoId)
	err = row.Scan(&PinnedIssuesCount)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
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

	IsRepoOwner := c.Param("user_name") == currentUser.Username

	CommentedUsersImages := []string{}
	var UsersImages string
	rows, err = db.Query(`SELECT DISTINCT users.image FROM USERS JOIN comments ON 
						comments.user_id = users.id WHERE issue_id = $1`, IssueId)

	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		err = rows.Scan(&UsersImages)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		CommentedUsersImages = append(CommentedUsersImages, UsersImages)
	}

	contains := false
	for _, item := range CommentedUsersImages {
		if item == issue.Image {
			contains = true
		}
	}

	if !contains {
		CommentedUsersImages = append(CommentedUsersImages, issue.Image)
	}

	NumberOfCommented := len(CommentedUsersImages)

	locked := lock == "Locked" && IsRepoOwner == false

	IsRepoOwnerOrIssueUser := currentUser.Username == c.Param("user_name") ||
		currentUser.Username == issue.Username

	c.HTML(http.StatusOK, "issue.html", gin.H{
		"CurrentUser":            currentUser,
		"Authorized":             authorized,
		"Issue":                  issue,
		"Comments":               comments,
		"UserName":               c.Param("user_name"),
		"PinnedIssuesCount":      PinnedIssuesCount,
		"CommentedUsers":         CommentedUsersImages,
		"RepoOwner":              IsRepoOwner,
		"NumberOfCommented":      NumberOfCommented,
		"Locked":                 locked,
		"IsRepoOwnerOrIssueUser": IsRepoOwnerOrIssueUser,
		"RepoName":               repoName})
}

func createIssueComment(c *gin.Context, db *sql.DB) {
	_, err := authorize(c, db)
	if err != nil {
		c.Redirect(http.StatusFound, "http://localhost:8000/login")
		return
	}
	username := c.Param("user_name")
	repoName := c.Param("repo_name")

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
	fmt.Println(username)
	title := c.PostForm("title")
	// repoId := c.PostForm("repo_id")
	body := c.PostForm("body")
	// userId := c.PostForm("user_id")
	issue := Issue{
		Title:  title,
		RepoId: currentRepo.RepoId,
		Body:   body,
		UserId: currentUser.ID,
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
	c.Redirect(http.StatusFound, "http://localhost:8000/login")
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
	fmt.Println(currentUser)
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
		c.Redirect(http.StatusFound, "http://localhost:8000/login")
		return
	}
	cookie, err := CreateCookie(db, Id)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.SetCookie("session", cookie, 300, "/", "localhost:8000", false, true)
	c.Redirect(http.StatusFound, "http://localhost:8000/"+username)
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

func PostUserSignOutHandler(c *gin.Context, db *sql.DB) {

	_, err := authorize(c, db)
	if err != nil {
		c.Redirect(http.StatusFound, "http://localhost:8000/login")
		return
	}

	session, err := c.Cookie("session")
	if err != nil {
		c.Redirect(http.StatusFound, "http://localhost:8000/login")
		return
	}

	err = DeleteCookie(db, session)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusFound, "http://localhost:8000/login")
}

func postPinPageHandler(c *gin.Context, db *sql.DB) {

	currentUser, err := authorize(c, db)
	if err != nil {
		c.Redirect(http.StatusFound, "http://localhost:8000/login")
		return
	}

	username := c.Param("user_name")

	repoName := c.Param("repo_name")

	IssueNumber := c.Param("issue_number")

	_issueId := c.PostForm("issue_id")

	currentRepo, err := getCurrentRepo(db, username, repoName)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	IsRepoOwner := currentUser.Username == c.Param("user_name")
	if !IsRepoOwner {
		c.Redirect(http.StatusFound, "http://localhost:8000/"+username+"/"+repoName+"/issues/"+IssueNumber)
		return
	}

	var count int

	row := db.QueryRow(`SELECT count(*) FROM issues WHERE pinned = 'Pinned' AND repo_id = $1`, currentRepo.RepoId)
	err = row.Scan(&count)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if count >= 3 {
		c.Redirect(http.StatusFound, "http://localhost:8000/"+username+"/"+repoName+"/issues/"+IssueNumber)
		return
	}

	_, err = db.Exec(`UPDATE issues SET pinned = 'Pinned' WHERE id = $1`, _issueId)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Redirect(http.StatusFound, "http://localhost:8000/"+username+"/"+repoName+"/issues/"+IssueNumber)
}

func postUnPinPageHandler(c *gin.Context, db *sql.DB) {

	currentUser, err := authorize(c, db)
	if err != nil {
		c.Redirect(http.StatusFound, "http://localhost:8000/login")
		return
	}

	username := c.Param("user_name")

	repoName := c.Param("repo_name")

	IssueNumber := c.Param("issue_number")

	_issueId := c.PostForm("issue_id")

	IsRepoOwner := currentUser.Username == c.Param("user_name")
	if !IsRepoOwner {
		c.Redirect(http.StatusFound, "http://localhost:8000/"+username+"/"+repoName+"/issues/"+IssueNumber)
		return
	}

	_, err = db.Exec(`UPDATE issues SET pinned = 'Unpinned' WHERE id = $1`, _issueId)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusFound, "http://localhost:8000/"+username+"/"+repoName+"/issues/"+IssueNumber)

}

func postLockPageHandler(c *gin.Context, db *sql.DB) {

	currentUser, err := authorize(c, db)
	if err != nil {
		c.Redirect(http.StatusFound, "http://localhost:8000/login")
		return
	}

	username := c.Param("user_name")

	repoName := c.Param("repo_name")

	IssueNumber := c.Param("issue_number")

	_issueId := c.PostForm("issue_id")

	IsRepoOwner := currentUser.Username == c.Param("user_name")
	if !IsRepoOwner {
		c.Redirect(http.StatusFound, "http://localhost:8000/"+username+"/"+repoName+"/issues/"+IssueNumber)
		return
	}

	_, err = db.Exec(`UPDATE issues SET lock = 'Locked' WHERE id = $1`, _issueId)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusFound, "http://localhost:8000/"+username+"/"+repoName+"/issues/"+IssueNumber)

}

func postUnlockPageHandler(c *gin.Context, db *sql.DB) {

	currentUser, err := authorize(c, db)
	if err != nil {
		c.Redirect(http.StatusFound, "http://localhost:8000/login")
		return
	}

	username := c.Param("user_name")

	repoName := c.Param("repo_name")

	_issueId := c.PostForm("issue_id")

	IssueNumber := c.Param("issue_number")

	IsRepoOwner := currentUser.Username == c.Param("user_name")
	if !IsRepoOwner {
		c.Redirect(http.StatusFound, "http://localhost:8000/"+username+"/"+repoName+"/issues/"+IssueNumber)
		return
	}

	_, err = db.Exec(`UPDATE issues SET lock = 'Unlocked' WHERE id = $1`, _issueId)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusFound, "http://localhost:8000/"+username+"/"+repoName+"/issues/"+IssueNumber)
}

type Collaborator struct {
	Username  string
	UserImage string
	Name      string
}

func getCollaboratorPageHandler(c *gin.Context, db *sql.DB) {

	currentUser, err := authorize(c, db)
	authorized := err == nil

	username := c.Param("user_name")
	repoName := c.Param("repo_name")

	IsRepoOwner := currentUser.Username == username
	if !IsRepoOwner {
		c.Redirect(http.StatusFound, "http://localhost:8000/"+currentUser.Username)
		return
	}

	currentRepo, err := getCurrentRepo(db, username, repoName)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	rows, err := db.Query(
		`SELECT users.username,users.image,users.name 
		 FROM users JOIN collaborators ON users.id = collaborators.user_id 
		 WHERE collaborators.repo_id = $1;`,
		currentRepo.RepoId)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	collaborators := []Collaborator{}

	for rows.Next() {
		var UserName, image, name string
		err = rows.Scan(&UserName, &image, &name)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		collaborator := Collaborator{
			Username:  UserName,
			UserImage: image,
			Name:      name,
		}
		collaborators = append(collaborators, collaborator)
	}

	IsCollaboratorsAvailable := len(collaborators) >= 1

	c.HTML(http.StatusOK, "collaboration.html",
		gin.H{"Username": username,
			"RepoName":                 repoName,
			"CurrentRepo":              currentRepo,
			"Collaborators":            collaborators,
			"Authorized":               authorized,
			"IsRepoOwner":              IsRepoOwner,
			"IsCollaboratorsAvailable": IsCollaboratorsAvailable,
		})

}

func postCollaboratorPageHandler(c *gin.Context, db *sql.DB) {

	currentUser, err := authorize(c, db)
	if err != nil {
		c.Redirect(http.StatusFound, "http://localhost:8000/login")
		return
	}

	username := c.Param("user_name")
	repoName := c.Param("repo_name")

	IsRepoOwner := currentUser.Username == username
	if !IsRepoOwner {
		c.Redirect(http.StatusFound, "http://localhost:8000/"+currentUser.Username)
		return
	}

	currentRepo, err := getCurrentRepo(db, username, repoName)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	userName := c.PostForm("user_name")

	user, err := GetUserByUserName(db, userName)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	_, err = db.Exec(`INSERT INTO collaborators (repo_id,user_id)VALUES($1,$2)`,
		currentRepo.RepoId, user.ID)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusFound, "http://localhost:8000/"+username+"/"+repoName+"/collaboration")
}

func postRemoveCollaboratorHandler(c *gin.Context, db *sql.DB) {

	currentUser, err := authorize(c, db)
	if err != nil {
		c.Redirect(http.StatusFound, "http://localhost:8000/"+currentUser.Username)
		return
	}

	username := c.Param("user_name")
	repoName := c.Param("repo_name")

	IsRepoOwner := currentUser.Username == username
	if !IsRepoOwner {
		c.Redirect(http.StatusFound, "http://localhost:8000/"+currentUser.Username)
		return
	}

	currentRepo, err := getCurrentRepo(db, username, repoName)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	userName := c.PostForm("user_name")

	user, err := GetUserByUserName(db, userName)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("DELETE FROM collaborators WHERE user_id = $1 AND repo_id = $2",
		user.ID, currentRepo.RepoId)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusFound, "http://localhost:8000/"+username+"/"+repoName+"/collaboration")
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

	router.GET("", func(c *gin.Context) { getUserNewPageHandler(c, db) })
	router.POST("", func(c *gin.Context) { PostUserNewPageHandler(c, db) })

	userGroup := router.Group("/:user_name")

	userGroup.GET("", func(c *gin.Context) {
		switch c.Param("user_name") {
		case "new":
			getRepoNewPageHandler(c, db)
		case "login":
			getUserSigninPageHandler(c, db)
		case "favicon.ico":
			c.Status(http.StatusOK)
		default:
			getUserPageHandler(c, db)
		}
	})
	userGroup.POST("", func(c *gin.Context) {
		switch c.Param("user_name") {
		case "new":
			PostRepoNewPageHandler(c, db)
		case "login":
			PostUserSigninPageHandler(c, db)
		case "Signout":
			PostUserSignOutHandler(c, db)
		default:
			c.Status(500)
		}
	})

	pages := userGroup.Group("/:repo_name")
	pages.GET("/issues",
		func(c *gin.Context) { getIssuesPageHandler(c, db) })
	pages.POST("/issues/:issue_number", func(c *gin.Context) {
		if c.Param("issue_number") == "new" {
			postNewIssuePageHandler(c, db)
		} else {
			c.Status(http.StatusNotFound)
		}
	})

	pages.GET("/issues/:issue_number",
		func(c *gin.Context) {
			if c.Param("issue_number") == "new" {
				getNewIssuePageHandler(c, db)
			} else {
				getIssuePageHandler(c, db)
			}
		})

	pages.POST("/issues/:issue_number/pin", func(c *gin.Context) { postPinPageHandler(c, db) })
	pages.POST("/issues/:issue_number/unpin", func(c *gin.Context) { postUnPinPageHandler(c, db) })
	pages.POST("/issues/:issue_number/lock", func(c *gin.Context) { postLockPageHandler(c, db) })
	pages.POST("/issues/:issue_number/unlock", func(c *gin.Context) { postUnlockPageHandler(c, db) })

	pages.GET("/collaboration", func(c *gin.Context) { getCollaboratorPageHandler(c, db) })
	pages.POST("/collaboration", func(c *gin.Context) { postCollaboratorPageHandler(c, db) })
	pages.POST("/removecollaborator", func(c *gin.Context) { postRemoveCollaboratorHandler(c, db) })

	pages.POST("/comments", func(c *gin.Context) { createIssueComment(c, db) })

	stylesRouter := gin.Default()
	stylesRouter.Static("/styles", "./styles")

	http.Handle("/styles/", stylesRouter)

	http.Handle("/", router)

	http.ListenAndServe(":8000", nil)
}
