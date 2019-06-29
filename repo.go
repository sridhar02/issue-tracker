package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type Repo struct {
	ID          string    `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	UserId      string    `json:"user_id,omitempty"`
	IssuesCount int       `json:"issue_count,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	Description string    `json:"description,omitempty"`
	Type        string    `json:"type,omitempty"`
}

func GetRepo(db *sql.DB, id string) (Repo, error) {

	var name, userId, createdAt, updatedAt, description, TYPE string
	var issueCount int

	row := db.QueryRow("SELECT id,name, user_id,issue_count,created_at, updated_at,description,type FROM repos WHERE id=$1", id)
	err := row.Scan(&id, &name, &userId, &issueCount, &createdAt, &updatedAt, &description, &TYPE)
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
		ID:          id,
		Name:        name,
		UserId:      userId,
		IssuesCount: issueCount,
		CreatedAt:   CreatedAt,
		UpdatedAt:   UpdatedAt,
		Description: description,
		Type:        TYPE,
	}

	return repo, nil
}

func CreateRepo(db *sql.DB, repo Repo) error {

	ID := uuid.New().String()

	_, err := db.Exec(`INSERT INTO repos(id,name, user_id, issue_count, created_at, updated_at,description,type)
						VALUES($1,$2,$3,$4,$5,$6,$7,$8)`,
		ID,
		repo.Name,
		repo.UserId,
		repo.IssuesCount,
		time.Now().Format(time.RFC3339),
		time.Now().Format(time.RFC3339),
		repo.Description,
		repo.Type)
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

func getRepoHandler(c *gin.Context, db *sql.DB) {

	id := c.Param("id")

	repo, err := GetRepo(db, id)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, repo)

	// c.Status(http.StatusNoContent)

}

func postRepoHandler(c *gin.Context, db *sql.DB) {

	repo := Repo{}
	err := c.BindJSON(&repo)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = CreateRepo(db, repo)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)

}

func putRepoHandler(c *gin.Context, db *sql.DB) {

	repo := Repo{}
	err := c.BindJSON(&repo)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = UpdateRepo(db, repo)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)

}

func deleteRepoHandler(c *gin.Context, db *sql.DB) {

	id := c.Param("id")

	err := DeleteRepo(db, id)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	// c.JSON(http.StatusOK, repo)

	c.Status(http.StatusNoContent)

}
