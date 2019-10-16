package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	// "github.com/dustin/go-humanize"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
	// "log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

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
	currentUser, err := authorize(c, db)
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
			"RepoName":    repoName,
			"Issues":      issues,
			"CurrentUser": currentUser,
			"Authorized":  authorized})
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

	writeAccess, err := hasWriteAccess(db, currentRepo, currentUser)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	changeStatusAccess := writeAccess || currentUser.Username == issue.Username

	c.HTML(http.StatusOK, "issue.html", gin.H{
		"CurrentUser":        currentUser,
		"Authorized":         authorized,
		"Issue":              issue,
		"Comments":           comments,
		"UserName":           c.Param("user_name"),
		"PinnedIssuesCount":  PinnedIssuesCount,
		"CommentedUsers":     CommentedUsersImages,
		"NumberOfCommented":  NumberOfCommented,
		"Locked":             locked,
		"WriteAccess":        writeAccess,
		"ChangeStatusAccess": changeStatusAccess,
		"RepoName":           repoName})
}

func createIssueComment(c *gin.Context, db *sql.DB) {
	currentUser, err := authorize(c, db)
	if err != nil {
		c.Redirect(http.StatusFound, os.Getenv("URL")+"/login")
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

	body := c.PostForm("body")
	_issueId := c.PostForm("issue_id")
	IssueNumber := c.PostForm("issue_number")

	issueId, err := strconv.Atoi(_issueId)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if c.PostForm("comment_and_open") == "1" {
		_, err := db.Exec(`UPDATE issues SET status = 'Open' WHERE id = $1`, _issueId)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		err = CommentNotifications(db, issueId, currentRepo, currentUser)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Redirect(http.StatusFound, os.Getenv("URL")+"/"+username+"/"+repoName+"/issues/"+IssueNumber)
		return
	}

	if c.PostForm("comment_and_close") == "1" {
		_, err := db.Exec(`UPDATE issues SET status = 'Closed' WHERE id = $1`, _issueId)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		err = CommentNotifications(db, issueId, currentRepo, currentUser)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Redirect(http.StatusFound, os.Getenv("URL")+"/"+username+"/"+repoName+"/issues/"+IssueNumber)
		return
	}

	comment := Comment{
		Body:    body,
		UserId:  currentRepo.UserId,
		RepoId:  currentRepo.RepoId,
		IssueId: issueId,
	}
	err = CreateComment(db, comment)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = CommentNotifications(db, issueId, currentRepo, currentUser)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusFound, os.Getenv("URL")+"/"+username+"/"+repoName+"/issues/"+IssueNumber)
}

func getNewIssuePageHandler(c *gin.Context, db *sql.DB) {

	currentUser, err := authorize(c, db)
	if err != nil {
		c.Redirect(http.StatusFound, os.Getenv("URL")+"/login")
		return
	}

	authorized := err == nil

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
			"Authorized":  authorized,
			"RepoName":    repoName})
}
func postNewIssuePageHandler(c *gin.Context, db *sql.DB) {
	currentUser, err := authorize(c, db)
	if err != nil {
		c.Redirect(http.StatusFound, os.Getenv("URL")+"/login")
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

	title := c.PostForm("title")
	body := c.PostForm("body")

	issue := Issue{
		Title:  title,
		RepoId: currentRepo.RepoId,
		Body:   body,
		UserId: currentUser.ID,
	}

	IssueId, issueNumber, err := CreateIssue(db, issue)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	notificationUserIds := map[string]int{}
	rows, err := db.Query(`SELECT User_id FROM collaborators WHERE repo_id = $1`, currentRepo.RepoId)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		var userId string
		err = rows.Scan(&userId)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		notificationUserIds[userId] = 0
	}
	// notificationUserIds[currentRepo.UserId] = 0

	for UserId, _ := range notificationUserIds {
		err = CreateNotification(db, IssueId, UserId, currentRepo.RepoId)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	c.Redirect(http.StatusFound, os.Getenv("URL")+"/"+username+"/"+repoName+"/issues/"+strconv.Itoa(issueNumber))
}

func getRepoNewPageHandler(c *gin.Context, db *sql.DB) {
	currentUser, err := authorize(c, db)
	if err != nil {
		c.Redirect(http.StatusFound, os.Getenv("URL")+"/login")
		return
	}
	c.HTML(http.StatusOK, "repo_new.html", gin.H{"CurrentUser": currentUser})
}

func PostRepoNewPageHandler(c *gin.Context, db *sql.DB) {
	_, err := authorize(c, db)
	if err != nil {
		c.Redirect(http.StatusFound, os.Getenv("URL")+"/login")
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
	c.Redirect(http.StatusFound, os.Getenv("URL")+"/"+userName+"/"+name+"/issues")
}

func getUserSigninPageHandler(c *gin.Context, db *sql.DB) {
	currentUser, err := authorize(c, db)
	if err == nil {
		c.Redirect(http.StatusFound, os.Getenv("URL")+"/"+currentUser.Username)
		return
	}
	c.HTML(http.StatusOK, "user_signin.html", gin.H{})
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

func IsCollaborator(db *sql.DB, userId string, repoId string) (bool, error) {

	var a, b string
	row := db.QueryRow(`SELECT * FROM collaborators WHERE repo_id = $1 AND user_id = $2`, repoId, userId)
	err := row.Scan(&a, &b)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func hasWriteAccess(db *sql.DB, currentRepo CurrentRepo, currentUser User) (bool, error) {
	//there is no cookie for user
	if currentUser.ID == "" {
		return false, nil
	}

	if currentRepo.UserId == currentUser.ID {
		return true, nil
	}

	return IsCollaborator(db, currentUser.ID, currentRepo.RepoId)
}

func PostUserSignOutHandler(c *gin.Context, db *sql.DB) {

	_, err := authorize(c, db)
	if err != nil {
		fmt.Println("error")
		c.Redirect(http.StatusFound, os.Getenv("URL")+"/login")
		return
	}

	session, err := c.Cookie("session")
	if err != nil {
		c.Redirect(http.StatusFound, os.Getenv("URL")+"/login")
		return
	}

	err = DeleteCookie(db, session)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusFound, os.Getenv("URL")+"/login")
}

func postPinPageHandler(c *gin.Context, db *sql.DB) {

	currentUser, err := authorize(c, db)
	if err != nil {
		c.Redirect(http.StatusFound, os.Getenv("URL")+"/login")
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

	writeAccess, err := hasWriteAccess(db, currentRepo, currentUser)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if !writeAccess {
		c.Redirect(http.StatusFound, os.Getenv("URL")+"/"+username+"/"+repoName+"/issues/"+IssueNumber)
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
		c.Redirect(http.StatusFound, os.Getenv("URL")+"/"+username+"/"+repoName+"/issues/"+IssueNumber)
		return
	}

	_, err = db.Exec(`UPDATE issues SET pinned = 'Pinned' WHERE id = $1`, _issueId)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	issueId, err := strconv.Atoi(_issueId)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = CollaboratorNotifications(db, currentRepo, issueId)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusFound, os.Getenv("URL")+"/"+username+"/"+repoName+"/issues/"+IssueNumber)

}

func postUnPinPageHandler(c *gin.Context, db *sql.DB) {

	currentUser, err := authorize(c, db)
	if err != nil {
		c.Redirect(http.StatusFound, os.Getenv("URL")+"/login")
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

	writeAccess, err := hasWriteAccess(db, currentRepo, currentUser)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if !writeAccess {
		c.Redirect(http.StatusFound, os.Getenv("URL")+"/"+username+"/"+repoName+"/issues/"+IssueNumber)
		return
	}

	_, err = db.Exec(`UPDATE issues SET pinned = 'Unpinned' WHERE id = $1`, _issueId)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	issueId, err := strconv.Atoi(_issueId)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = CollaboratorNotifications(db, currentRepo, issueId)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusFound, os.Getenv("URL")+"/"+username+"/"+repoName+"/issues/"+IssueNumber)

}

func postLockPageHandler(c *gin.Context, db *sql.DB) {

	currentUser, err := authorize(c, db)
	if err != nil {
		c.Redirect(http.StatusFound, os.Getenv("URL")+"/login")
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

	writeAccess, err := hasWriteAccess(db, currentRepo, currentUser)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if !writeAccess {
		c.Redirect(http.StatusFound, os.Getenv("URL")+"/"+username+"/"+repoName+"/issues/"+IssueNumber)
		return
	}

	_, err = db.Exec(`UPDATE issues SET pinned = 'Unpinned' WHERE id = $1`, _issueId)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	_, err = db.Exec(`UPDATE issues SET lock = 'Locked' WHERE id = $1`, _issueId)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	issueId, err := strconv.Atoi(_issueId)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = CollaboratorNotifications(db, currentRepo, issueId)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusFound, os.Getenv("URL")+"/"+username+"/"+repoName+"/issues/"+IssueNumber)

}

func postUnlockPageHandler(c *gin.Context, db *sql.DB) {

	currentUser, err := authorize(c, db)
	if err != nil {
		c.Redirect(http.StatusFound, os.Getenv("URL")+"/login")
		return
	}

	username := c.Param("user_name")

	repoName := c.Param("repo_name")

	_issueId := c.PostForm("issue_id")

	IssueNumber := c.Param("issue_number")

	currentRepo, err := getCurrentRepo(db, username, repoName)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	writeAccess, err := hasWriteAccess(db, currentRepo, currentUser)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if !writeAccess {
		c.Redirect(http.StatusFound, os.Getenv("URL")+"/"+username+"/"+repoName+"/issues/"+IssueNumber)
		return
	}

	_, err = db.Exec(`UPDATE issues SET pinned = 'Unpinned' WHERE id = $1`, _issueId)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	_, err = db.Exec(`UPDATE issues SET lock = 'Unlocked' WHERE id = $1`, _issueId)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	issueId, err := strconv.Atoi(_issueId)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = CollaboratorNotifications(db, currentRepo, issueId)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusFound, os.Getenv("URL")+"/"+username+"/"+repoName+"/issues/"+IssueNumber)
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
		c.Redirect(http.StatusFound, os.Getenv("URL")+currentUser.Username)
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
		c.Redirect(http.StatusFound, os.Getenv("URL")+"/login")
		return
	}

	username := c.Param("user_name")
	repoName := c.Param("repo_name")

	IsRepoOwner := currentUser.Username == username
	if !IsRepoOwner {
		c.Redirect(http.StatusFound, os.Getenv("URL")+currentUser.Username)
		return
	}

	currentRepo, err := getCurrentRepo(db, username, repoName)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	userName := c.PostForm("user_name")
	if userName == currentUser.Username {
		c.Redirect(http.StatusFound, os.Getenv("URL")+"/"+username+"/"+repoName+"/collaboration")
		return
	}

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

	c.Redirect(http.StatusFound, os.Getenv("URL")+"/"+username+"/"+repoName+"/collaboration")
}

func postRemoveCollaboratorHandler(c *gin.Context, db *sql.DB) {

	currentUser, err := authorize(c, db)
	if err != nil {
		c.Redirect(http.StatusFound, os.Getenv("URL")+"/"+currentUser.Username)
		return
	}

	username := c.Param("user_name")
	repoName := c.Param("repo_name")

	IsRepoOwner := currentUser.Username == username
	if !IsRepoOwner {
		c.Redirect(http.StatusFound, os.Getenv("URL")+"/"+currentUser.Username)
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

	c.Redirect(http.StatusFound, os.Getenv("URL")+"/"+username+"/"+repoName+"/collaboration")
}

type NotificationRequired struct {
	Title   string
	Read    string
	IssueId int
}

func getUnreadNotificationsHandler(c *gin.Context, db *sql.DB) {

	currentUser, err := authorize(c, db)
	authorized := err == nil

	rows, err := db.Query(
		`SELECT DISTINCT issues.title, notifications.read,issues.id, notifications.repo_id
		 FROM issues JOIN notifications ON issues.id = notifications.issue_id 
		 WHERE notifications.user_id = $1 AND notifications.read = $2`,
		currentUser.ID, "unread")

	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	UserNotifications := map[string][]NotificationRequired{}

	for rows.Next() {
		var issueId int
		var issueTitle, read, repoId, repoName, username string
		err = rows.Scan(&issueTitle, &read, &issueId, &repoId)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if read != read {

		}
		row := db.QueryRow(`SELECT repos.name,users.username FROM repos JOIN users 
			ON repos.user_id = users.id WHERE repos.id = $1`, repoId)
		err := row.Scan(&repoName, &username)
		if err != nil {
			fmt.Println(err)
			return
		}
		repoNotifications := UserNotifications[username+"/"+repoName]
		Notification := NotificationRequired{
			Title:   issueTitle,
			Read:    read,
			IssueId: issueId,
		}
		repoNotifications = append(repoNotifications, Notification)
		UserNotifications[username+"/"+repoName] = repoNotifications
	}

	c.HTML(http.StatusOK, "notifications.html",
		gin.H{
			"CurrentUser":          currentUser,
			"Authorized":           authorized,
			"NotificationRequired": UserNotifications,
		})

}

func getReadNotificationsHandler(c *gin.Context, db *sql.DB) {

	currentUser, err := authorize(c, db)
	authorized := err == nil

	rows, err := db.Query(
		`SELECT DISTINCT issues.title, notifications.read,issues.id, notifications.repo_id
		 FROM issues JOIN notifications ON issues.id = notifications.issue_id 
		 WHERE notifications.user_id = $1 AND notifications.read = $2`,
		currentUser.ID, "read")

	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	UserNotifications := map[string][]NotificationRequired{}

	for rows.Next() {
		var issueId int
		var issueTitle, read, repoId, repoName, username string
		err = rows.Scan(&issueTitle, &read, &issueId, &repoId)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if read != read {

		}
		row := db.QueryRow(`SELECT repos.name,users.username FROM repos JOIN users 
			ON repos.user_id = users.id WHERE repos.id = $1`, repoId)
		err := row.Scan(&repoName, &username)
		if err != nil {
			fmt.Println(err)
			return
		}
		repoNotifications := UserNotifications[username+"/"+repoName]
		Notification := NotificationRequired{
			Title:   issueTitle,
			Read:    read,
			IssueId: issueId,
		}
		repoNotifications = append(repoNotifications, Notification)
		UserNotifications[username+"/"+repoName] = repoNotifications
	}

	c.HTML(http.StatusOK, "notifications.html",
		gin.H{
			"CurrentUser":          currentUser,
			"Authorized":           authorized,
			"NotificationRequired": UserNotifications,
		})

}

func PostNotificationsHandler(c *gin.Context, db *sql.DB) {

	_, err := authorize(c, db)
	if err != nil {
		fmt.Println("error")
		c.Redirect(http.StatusFound, os.Getenv("URL")+"/login")
		return
	}

	IssueId := c.PostForm("issue_id")
	Read := c.PostForm("read")
	// fmt.Println(IssueId)

	if Read != "unread" {
		_, err = db.Exec(`UPDATE notifications SET read = 'unread' WHERE issue_id = $1`, IssueId)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

	} else {
		_, err = db.Exec(`UPDATE notifications SET read = 'read' WHERE issue_id = $1`, IssueId)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	c.Redirect(http.StatusFound, os.Getenv("URL")+"/notifications")
}

func postUserSignupHandler(c *gin.Context, db *sql.DB) {
	user := User{}
	err := c.BindJSON(&user)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	err = CreateUser(db, user)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusCreated)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	var i int
	for i = range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

type Login struct {
	UserId    string    `json:"user_id,omitempty"`
	Secret    string    `json:"secret,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func CreateLogin(db *sql.DB, userId string) (Login, error) {

	secret := RandStringBytes(32)
	_, err := db.Exec(`INSERT INTO logins(user_id,secret,created_at,updated_at)
									VALUES($1,$2,$3,$4)`,
		userId,
		secret,
		time.Now().Format(time.RFC3339),
		time.Now().Format(time.RFC3339))

	login := Login{
		UserId:    userId,
		Secret:    secret,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err != nil {
		return login, err
	}

	return login, nil
}

func PostUserSigninPageHandler(c *gin.Context, db *sql.DB) {

	user := User{}
	err := c.BindJSON(&user)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var Password, Id string
	row := db.QueryRow(`SELECT password,id FROM users WHERE username=$1`, user.Username)
	err = row.Scan(&Password, &Id)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if Password != user.Password {
		c.Redirect(http.StatusFound, os.Getenv("URL")+"/login")
		return
	}
	login, err := CreateLogin(db, Id)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(201, login)

}

func authorization(c *gin.Context, db *sql.DB) (string, error) {

	value := c.GetHeader("Authorization")
	secret := strings.TrimPrefix(value, "Bearer ")

	var userId string
	row := db.QueryRow(`SELECT user_id FROM logins WHERE secret=$1`, secret)
	err := row.Scan(&userId)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return "", err
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
			return "", err
		}
	}

	return userId, nil

}

type UserDetails struct {
	User      User     `json:"user,omitempty"`
	RepoNames []string `json:"reponames"`
}

func getUserPageHandler(c *gin.Context, db *sql.DB) {

	userId, err := authorization(c, db)
	if err != nil {
		return
	}

	var username string
	row := db.QueryRow(`select username from users where id=$1`, userId)
	err = row.Scan(&username)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
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

	userDetails := UserDetails{
		User:      user,
		RepoNames: repoNames,
	}

	c.JSON(201, userDetails)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env file not found")
	}
	connStr := fmt.Sprintf("user=%s dbname=%s host=%s password=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PASSWORD"),
	)
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
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "GET", "DELETE", "POST", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.POST("/signup", func(c *gin.Context) { postUserSignupHandler(c, db) })
	router.POST("/signin", func(c *gin.Context) { PostUserSigninPageHandler(c, db) })
	router.GET("/user", func(c *gin.Context) { getUserPageHandler(c, db) })

	// userGroup := router.Group("/:user_name")

	// userGroup.GET("", func(c *gin.Context) {
	// 	switch c.Param("user_name") {
	// 	case "new":
	// 		getRepoNewPageHandler(c, db)
	// 	case "login":
	// 		getUserSigninPageHandler(c, db)
	// 	case "favicon.ico":
	// 		c.Status(http.StatusOK)
	// 	case "notifications":
	// 		getUnreadNotificationsHandler(c, db)
	// 	case "notifications-read":
	// 		getReadNotificationsHandler(c, db)
	// 	default:
	// 		getUserPageHandler(c, db)
	// 	}
	// })
	// userGroup.POST("", func(c *gin.Context) {
	// 	switch c.Param("user_name") {
	// 	case "new":
	// 		PostRepoNewPageHandler(c, db)
	// 	case "join":
	// 		PostUserNewPageHandler(c, db)
	// 	case "login":
	// 		PostUserSigninPageHandler(c, db)
	// 	case "Signout":
	// 		PostUserSignOutHandler(c, db)
	// 	case "notifications":
	// 		PostNotificationsHandler(c, db)
	// 	default:
	// 		c.Status(500)
	// 	}
	// })

	// pages := userGroup.Group("/:repo_name")
	// pages.GET("/issues",
	// 	func(c *gin.Context) { getIssuesPageHandler(c, db) })
	// pages.POST("/issues/:issue_number", func(c *gin.Context) {
	// 	if c.Param("issue_number") == "new" {
	// 		postNewIssuePageHandler(c, db)
	// 	} else {
	// 		c.Status(http.StatusNotFound)
	// 	}
	// })

	// pages.GET("/issues/:issue_number",
	// 	func(c *gin.Context) {
	// 		if c.Param("issue_number") == "new" {
	// 			getNewIssuePageHandler(c, db)
	// 		} else {
	// 			getIssuePageHandler(c, db)
	// 		}
	// 	})

	// pages.POST("/issues/:issue_number/pin", func(c *gin.Context) { postPinPageHandler(c, db) })
	// pages.POST("/issues/:issue_number/unpin", func(c *gin.Context) { postUnPinPageHandler(c, db) })
	// pages.POST("/issues/:issue_number/lock", func(c *gin.Context) { postLockPageHandler(c, db) })
	// pages.POST("/issues/:issue_number/unlock", func(c *gin.Context) { postUnlockPageHandler(c, db) })

	// pages.GET("/collaboration", func(c *gin.Context) { getCollaboratorPageHandler(c, db) })
	// pages.POST("/collaboration", func(c *gin.Context) { postCollaboratorPageHandler(c, db) })
	// pages.POST("/removecollaborator", func(c *gin.Context) { postRemoveCollaboratorHandler(c, db) })

	// pages.POST("/comments", func(c *gin.Context) { createIssueComment(c, db) })

	stylesRouter := gin.Default()
	stylesRouter.Static("/styles", "./styles")

	http.Handle("/styles/", stylesRouter)

	http.Handle("/", router)

	port := os.Getenv("PORT")

	if port != "" {
		port = ":" + port
	} else {
		port = ":8000"
	}

	http.ListenAndServe(port, nil)
}
