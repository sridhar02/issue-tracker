package main

import (
	"database/sql"
	"time"
)

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

	row := db.QueryRow("SELECT id,name, user_id,issue_count,created_at, updated_at FROM repos WHERE id=$1", id)
	err := row.Scan(&id, &name, &userId, &issueCount, &createdAt, &updatedAt)
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
