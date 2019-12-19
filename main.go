package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	// "github.com/dustin/go-humanize"
	sq "github.com/Masterminds/squirrel"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
	// "log"
	"net/http"
	"os"
	// "strconv"
	"strings"
	"time"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

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
	_, err := db.Exec(`INSERT INTO logins (user_id, secret, created_at, updated_at ) VALUES ( $1, $2, $3, $4)`,
		userId, secret, time.Now().Format(time.RFC3339), time.Now().Format(time.RFC3339))

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
	row := db.QueryRow(`SELECT password, id FROM users WHERE username = $1`, user.Username)
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

	c.JSON(http.StatusCreated, login)

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

func getAuthenticatedUserHandler(c *gin.Context, db *sql.DB) {

	userId, err := authorization(c, db)
	if err != nil {
		return
	}

	user, err := GetUser(db, userId)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, user)
}

func getAuthenticatedUserReposHandler(c *gin.Context, db *sql.DB) {

	userId, err := authorization(c, db)
	if err != nil {
		return
	}

	rows, err := db.Query(`SELECT id, name, issue_count, created_at, updated_at, description, type FROM repos
						                    WHERE user_id = $1 ORDER BY id ASC `, userId)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	Repos := []Repo{}
	for rows.Next() {
		var id, name, createdAt, updatedAt, description, Type string
		var issueCount int
		err = rows.Scan(&id, &name, &issueCount, &updatedAt, &createdAt, &description, &Type)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		CreatedAt, err := time.Parse(time.RFC3339, createdAt)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		UpdatedAt, err := time.Parse(time.RFC3339, updatedAt)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		user, err := GetUser(db, userId)
		if err != nil {
			fmt.Println(err)
			return
		}

		repo := Repo{
			ID:          id,
			Name:        name,
			IssuesCount: issueCount,
			CreatedAt:   CreatedAt,
			UpdatedAt:   UpdatedAt,
			Description: description,
			Type:        Type,
			User:        user,
		}
		Repos = append(Repos, repo)
	}
	c.JSON(http.StatusOK, Repos)
}

func getIssuesHandler(c *gin.Context, db *sql.DB) {

	username := c.Param("owner")
	repoName := c.Param("repo")

	rows, err := db.Query(`WITH repoCTE AS (
							                SELECT repos.id FROM repos JOIN users ON repos.user_id= users.id
						                  WHERE repos.name=$1 AND users.username=$2
						                  )
						             SELECT id,title,user_id ,body, created_at, updated_at,issue_number,pinned,
						             status, lock FROM issues WHERE repo_id IN
						             (SELECT id FROM repoCTE) ORDER BY id DESC`, repoName, username)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	issues := []Issue{}
	for rows.Next() {
		var title, userId, body, createdAt, updatedAt, pinned, status, lock string
		var id, issueNumber int
		err = rows.Scan(&id, &title, &userId, &body, &createdAt, &updatedAt, &issueNumber, &pinned, &status, &lock)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		CreatedAt, err := time.Parse(time.RFC3339, createdAt)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		UpdatedAt, err := time.Parse(time.RFC3339, updatedAt)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		user, err := GetUser(db, userId)
		if err != nil {
			fmt.Println(err)
			return
		}

		issue := Issue{
			ID:          id,
			Title:       title,
			UserId:      userId,
			Body:        body,
			IssueNumber: issueNumber,
			Status:      status,
			Pinned:      pinned,
			Lock:        lock,
			CreatedAt:   CreatedAt,
			UpdatedAt:   UpdatedAt,
			User:        user,
		}
		issues = append(issues, issue)
	}
	c.JSON(http.StatusOK, issues)
}
func getIssueHandler(c *gin.Context, db *sql.DB) {

	username := c.Param("owner")
	repoName := c.Param("repo")
	issueNumber := c.Param("issue_number")

	var Id int
	row := db.QueryRow(`WITH repo_cte AS (
		                       SELECT repos.id FROM repos JOIN users ON repos.user_id= users.id
                		       WHERE repos.name = $1 AND users.username = $2
                		       )
                		  SELECT id from issues WHERE repo_id
                		  IN (SELECT id FROM repo_cte) and issue_number=$3 `, repoName, username, issueNumber)
	err := row.Scan(&Id)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	issue, err := GetIssue(db, Id)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, issue)
}

func getCommentsHandler(c *gin.Context, db *sql.DB) {

	username := c.Param("owner")
	repoName := c.Param("repo")
	issueNumber := c.Param("issue_number")

	rows, err := db.Query(`WITH repo_cte AS (
		                          SELECT repos.id FROM repos JOIN users ON repos.user_id= users.id
         				              WHERE repos.name= $1 AND users.username= $2
         				              ),
         				         issue_cte AS (
         				    	         SELECT id FROM issues WHERE repo_id IN
         				               (select id from repo_cte) AND issue_number=$3
         				                )
         				        SELECT id ,user_id,body,created_at,updated_at FROM comments
         				        WHERE issue_id in (select id from issue_cte)`, repoName, username, issueNumber)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	comments := []Comment{}
	for rows.Next() {

		var userId, body, createdAt, updatedAt string
		var id int
		err = rows.Scan(&id, &userId, &body, &createdAt, &updatedAt)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		CreatedAt, err := time.Parse(time.RFC3339, createdAt)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		UpdatedAt, err := time.Parse(time.RFC3339, updatedAt)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		user, err := GetUser(db, userId)
		if err != nil {
			fmt.Println(err)
			return
		}

		comment := Comment{
			ID:        id,
			UserId:    userId,
			Body:      body,
			User:      user,
			CreatedAt: CreatedAt,
			UpdatedAt: UpdatedAt,
		}
		comments = append(comments, comment)
	}

	c.JSON(http.StatusOK, comments)
}

type CurrentRepo struct {
	UserId string
	RepoId string
}

func getCurrentRepo(db *sql.DB, username string, repoName string) (CurrentRepo, error) {

	var userId, repoId string
	row := db.QueryRow(`SELECT users.id, repos.id FROM repos JOIN users ON
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

type Collaborator struct {
	User User `json:"user"`
}

func getCollaborators(c *gin.Context, db *sql.DB) {

	_, err := authorization(c, db)
	if err != nil {
		return
	}

	username := c.Param("owner")
	repoName := c.Param("repo")

	currentRepo, err := getCurrentRepo(db, username, repoName)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	rows, err := db.Query(`SELECT users.id FROM users JOIN collaborators
                                ON users.id = collaborators.user_id  WHERE collaborators.repo_id = $1`, currentRepo.RepoId)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	collaborators := []Collaborator{}

	for rows.Next() {
		var userId string
		err = rows.Scan(&userId)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		user, err := GetUser(db, userId)
		if err != nil {
			fmt.Println(err)
			return
		}
		collaborator := Collaborator{
			User: user,
		}
		collaborators = append(collaborators, collaborator)
	}

	c.JSON(http.StatusOK, collaborators)

}
func postCollaborator(c *gin.Context, db *sql.DB) {
	userId, err := authorization(c, db)
	if err != nil {
		return
	}

	username := c.Param("owner")
	collaboratorUsername := c.Param("username")
	repoName := c.Param("repo")

	collaboratorUser, err := GetUserByUserName(db, collaboratorUsername)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var repoId, repoOwnerId string
	row := db.QueryRow(`SELECT repos.id, repos.user_id FROM repos JOIN users ON repos.user_id= users.id
                		         WHERE repos.name= $1 AND users.username= $2`, repoName, username)
	err = row.Scan(&repoId, &repoOwnerId)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if userId != repoOwnerId {
		c.Status(http.StatusForbidden)
		return
	}

	_, err = db.Exec(`INSERT INTO collaborators ( repo_id , user_id ) VALUES ( $1 , $2) `, repoId, collaboratorUser.ID)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusCreated)

}

type Assignee struct {
	Usernames []string `json:"usernames"`
	Image     string   `json:"image"`
	User      User     `json:"user"`
}

func postAssignee(c *gin.Context, db *sql.DB) {

	authorizedUserId, err := authorization(c, db)
	if err != nil {
		return
	}

	username := c.Param("owner")
	repoName := c.Param("repo")
	issueNumber := c.Param("issue_number")

	currentRepo, err := getCurrentRepo(db, username, repoName)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var issueId int
	var repoId string
	row := db.QueryRow(`WITH repo_cte AS (
		                     SELECT repos.id FROM repos JOIN users ON repos.user_id= users.id
                		     WHERE repos.name= $1 AND users.username= $2
                		     )
                		SELECT id, repo_id from issues WHERE repo_id
                		IN (SELECT id FROM repo_cte) and issue_number=$3 `, repoName, username, issueNumber)
	err = row.Scan(&issueId, &repoId)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	assignee := Assignee{}
	err = c.BindJSON(&assignee)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if authorizedUserId != currentRepo.UserId {
		c.Status(http.StatusForbidden)
		return
	}

	for _, assigneeUsername := range assignee.Usernames {
		assigneeUser, err := GetUserByUserName(db, assigneeUsername)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		_, err = db.Exec(`INSERT INTO assignees ( issue_id , user_id ) VALUES ( $1 , $2) `, issueId, assigneeUser.ID)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		err = CreateNotification(db, issueId, assigneeUser.ID, repoId)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	c.Status(http.StatusCreated)

}

func deleteCollaborator(c *gin.Context, db *sql.DB) {
	userId, err := authorization(c, db)
	if err != nil {
		return
	}
	username := c.Param("owner")
	collaboratorUsername := c.Param("username")
	repoName := c.Param("repo")

	collaboratorUser, err := GetUserByUserName(db, collaboratorUsername)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	var repoId, repoOwnerId string
	row := db.QueryRow(`SELECT repos.id, repos.user_id FROM repos JOIN users ON repos.user_id= users.id
                		         WHERE repos.name= $1 AND users.username= $2`, repoName, username)
	err = row.Scan(&repoId, &repoOwnerId)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if userId != repoOwnerId {
		c.Status(http.StatusForbidden)
		return
	}
	_, err = db.Exec(`DELETE FROM collaborators WHERE repo_id = $1 AND user_id= $2 `, repoId, collaboratorUser.ID)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusNoContent)
}

func deleteAssignee(c *gin.Context, db *sql.DB) {

	authorizedUserId, err := authorization(c, db)
	if err != nil {
		fmt.Println(err)
		return
	}

	username := c.Param("owner")
	repoName := c.Param("repo")
	issueNumber := c.Param("issue_number")
	currentRepo, err := getCurrentRepo(db, username, repoName)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	var issueId int
	var repoId string

	row := db.QueryRow(`WITH repo_cte AS (
		                     SELECT repos.id FROM repos JOIN users ON repos.user_id= users.id
                		     WHERE repos.name= $1 AND users.username= $2
                		     )
                		SELECT id, repo_id from issues
										WHERE repo_id IN (SELECT id FROM repo_cte) and issue_number = $3 `, repoName, username, issueNumber)
	err = row.Scan(&issueId, &repoId)
	if err != nil {
		fmt.Println(err)
		return
	}
	assignee := Assignee{}
	err = c.BindJSON(&assignee)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if authorizedUserId != currentRepo.UserId {
		c.Status(http.StatusForbidden)
		return
	}

	for _, assigneeUsername := range assignee.Usernames {
		assigneeUser, err := GetUserByUserName(db, assigneeUsername)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		_, err = db.Exec(`DELETE FROM assignees WHERE issue_id = $1 AND user_id= $2 `, issueId, assigneeUser.ID)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	c.Status(http.StatusNoContent)

}

type Subject struct {
	Title string `json:"title,omitempty"`
	Type  string `json:"type,omitempty"`
}

type UserNotification struct {
	ID        int       `json:"id,omitempty"`
	Repo      Repo      `json:"repo,omitempty"`
	Read      string    `json:"read,omitempty"`
	Subject   Subject   `json:"subject,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func getNotifications(c *gin.Context, db *sql.DB) {

	userId, err := authorization(c, db)
	if err != nil {
		fmt.Println(err)
		return
	}
	rows, err := db.Query(`SELECT DISTINCT
                                issues.id, issues.title, notifications.id, notifications.repo_id, notifications.read, notifications.created_at, notifications.updated_at
		                       FROM issues
													      JOIN notifications ON issues.id = notifications.issue_id WHERE notifications.user_id = $1 AND notifications.read = $2`, userId, "unread")

	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	notifications := []UserNotification{}
	for rows.Next() {
		var issueId, id int
		var issueTitle, read, repoId, createdAt, updatedAt string
		err = rows.Scan(&issueId, &issueTitle, &id, &repoId, &read, &createdAt, &updatedAt)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		CreatedAt, err := time.Parse(time.RFC3339, createdAt)
		if err != nil {
			fmt.Println(err)
			return
		}
		UpdatedAt, err := time.Parse(time.RFC3339, updatedAt)
		if err != nil {
			fmt.Println(err)
			return
		}
		repo, err := GetRepo(db, repoId)
		if err != nil {
			fmt.Println(err)
			return
		}
		subject := Subject{
			Title: issueTitle,
			Type:  "issue",
		}
		notification := UserNotification{
			ID:        id,
			Repo:      repo,
			Read:      read,
			Subject:   subject,
			CreatedAt: CreatedAt,
			UpdatedAt: UpdatedAt,
		}
		notifications = append(notifications, notification)
	}
	c.JSON(http.StatusOK, notifications)
}

func getRepoNotifications(c *gin.Context, db *sql.DB) {

}

func getLabelsHandler(c *gin.Context, db *sql.DB) {
	_, err := authorization(c, db)
	if err != nil {
		return
	}
}

func getLabelHandler(c *gin.Context, db *sql.DB) {
	_, err := authorization(c, db)
	if err != nil {
		return
	}
}

func postLabelHandler(c *gin.Context, db *sql.DB) {

	_, err := authorization(c, db)
	if err != nil {
		return
	}
	username := c.Param("owner")
	repoName := c.Param("repo")

	currentRepo, err := getCurrentRepo(db, username, repoName)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	label := Label{}
	err = c.BindJSON(&label)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	label.RepoId = currentRepo.RepoId

	err = CreateLabel(db, label)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)

}

func getIssueLabels(c *gin.Context, db *sql.DB) {

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
	router.GET("/user", func(c *gin.Context) { getAuthenticatedUserHandler(c, db) })
	router.GET("/users/:username", func(c *gin.Context) { getUserHandler(c, db) })
	router.GET("/users/:username/repos", func(c *gin.Context) { getReposHandler(c, db) })
	router.GET("/user/repos", func(c *gin.Context) { getAuthenticatedUserReposHandler(c, db) })
	router.POST("/user/repos", func(c *gin.Context) { postRepoHandler(c, db) })
	router.GET("/repos/:owner/:repo", func(c *gin.Context) { getRepoHandler(c, db) })
	router.GET("/repos/:owner/:repo/labels", func(c *gin.Context) { getLabelsHandler(c, db) })
	router.GET("/repos/:owner/:repo/labels/:name", func(c *gin.Context) { getLabelHandler(c, db) })
	router.POST("/repos/:owner/:repo/labels", func(c *gin.Context) { postLabelHandler(c, db) })
	router.GET("/repos/:owner/:repo/issues", func(c *gin.Context) { getIssuesHandler(c, db) })
	router.POST("/repos/:owner/:repo/issues", func(c *gin.Context) { postIssueHandler(c, db) })
	router.GET("/repos/:owner/:repo/collaborators", func(c *gin.Context) { getCollaborators(c, db) })
	router.POST("repos/:owner/:repo/collaborators/:username", func(c *gin.Context) { postCollaborator(c, db) })
	router.DELETE("repos/:owner/:repo/collaborators/:username", func(c *gin.Context) { deleteCollaborator(c, db) })
	router.POST("/repos/:owner/:repo/issues/:issue_number/assignees", func(c *gin.Context) { postAssignee(c, db) })
	router.DELETE("/repos/:owner/:repo/issues/:issue_number/assignees", func(c *gin.Context) { deleteAssignee(c, db) })
	router.GET("/repos/:owner/:repo/issues/:issue_number", func(c *gin.Context) { getIssueHandler(c, db) })
	router.PUT("/repos/:owner/:repo/issues/:issue_number", func(c *gin.Context) { putIssueHandler(c, db) })
	router.PUT("/repos/:owner/:repo/issues/:issue_number/pin", func(c *gin.Context) { putPinHandler(c, db) })
	router.PUT("/repos/:owner/:repo/issues/:issue_number/lock", func(c *gin.Context) { putLockHandler(c, db) })
	router.GET("/repos/:owner/:repo/issues/:issue_number/comments", func(c *gin.Context) { getCommentsHandler(c, db) })
	router.POST("/repos/:owner/:repo/issues/:issue_number/comments", func(c *gin.Context) { postCommentHandler(c, db) })
	router.GET("/repos/:owner/:repo/issues/:issue_number/labels", func(c *gin.Context) { getIssueLabels(c, db) })
	router.GET("/notifications", func(c *gin.Context) { getNotifications(c, db) })
	router.GET("/repos/:owner/:repo/notifications", func(c *gin.Context) { getRepoNotifications(c, db) })
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
