package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type Comment struct {
	ID        int       `json:"id,omitempty"`
	UserId    string    `json:"user_id,omitempty"`
	Body      string    `json:"body,omitempty"`
	IssueId   int       `json:"issue_id,omitempty"`
	RepoId    string    `json:"repo_id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func GetComment(db *sql.DB, id int) (Comment, error) {

	var userId, body, repoId, createdAt, updatedAt string
	var issueId int

	row := db.QueryRow("SELECT  id,user_id,body,issue_id,repo_id,created_at, updated_at FROM comments WHERE id=$1", id)
	err := row.Scan(&id, &userId, &body, &issueId, &repoId, &createdAt, &updatedAt)
	if err != nil {
		return Comment{}, err
	}

	CreatedAt, err := time.Parse(time.RFC3339, createdAt)

	if err != nil {
		return Comment{}, err
	}

	UpdatedAt, err := time.Parse(time.RFC3339, updatedAt)

	if err != nil {
		return Comment{}, err
	}

	comment := Comment{
		ID:        id,
		UserId:    userId,
		Body:      body,
		IssueId:   issueId,
		RepoId:    repoId,
		CreatedAt: CreatedAt,
		UpdatedAt: UpdatedAt,
	}

	return comment, nil

}

func CreateComment(db *sql.DB, comment Comment) error {

	_, err := db.Exec(`INSERT INTO comments(user_id,body,issue_id,repo_id,created_at, updated_at)
						VALUES($1,$2,$3,$4,$5,$6)`,
		comment.UserId,
		comment.Body,
		comment.IssueId,
		comment.RepoId,
		time.Now().Format(time.RFC3339),
		time.Now().Format(time.RFC3339))
	if err != nil {
		return err
	}

	return nil
}

func UpdateComment(db *sql.DB, comment Comment) error {

	_, err := db.Exec(`UPDATE comments SET body =$1, updated_at =$2 WHERE id = $3`,
		comment.Body, time.Now().Format(time.RFC3339), comment.ID)

	if err != nil {
		return err
	}
	return nil
}

func DeleteComment(db *sql.DB, id int) error {

	_, err := db.Exec(`DELETE FROM comments WHERE id = $1`, id)

	if err != nil {
		return err
	}

	return nil
}

func getCommentHandler(c *gin.Context, db *sql.DB) {

	id := c.Param("id")

	Id, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	comment, err := GetComment(db, Id)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, comment)

	// c.Status(http.StatusNoContent)
}

func postCommentHandler(c *gin.Context, db *sql.DB) {

	comment := Comment{}
	err := c.BindJSON(&comment)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = CreateComment(db, comment)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)

}

func putCommentHandler(c *gin.Context, db *sql.DB) {

	comment := Comment{}
	err := c.BindJSON(&comment)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = UpdateComment(db, comment)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

func deleteCommentHandler(c *gin.Context, db *sql.DB) {

	id := c.Param("id")

	Id, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = DeleteComment(db, Id)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	// c.JSON(http.StatusOK, repo)

	c.Status(http.StatusNoContent)

}
