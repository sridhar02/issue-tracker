package main

import (
	"database/sql"
	"time"
)

type Issue struct {
	ID          int       `json:"id,omitempty"`
	Title       string    `json:"title,omitempty"`
	UserId      string    `json:"user_id,omitempty"`
	Body        string    `json:"body,omitempty"`
	RepoId      string    `json:"repo_id,omitempty"`
	IssueNumber int       `json:"issue_number,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

func GetIssue(db *sql.DB, id int) (Issue, error) {

	var title, userId, body, repoId, createdAt, updatedAt string
	var issueNumber int

	row := db.QueryRow("SELECT title, user_id,body,repo_id,issue_number,created_at, updated_at FROM issues WHERE id=$1", id)
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
		ID:          id,
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

func CreateIssue(db *sql.DB, issue Issue) error {

	var issueCount int

	row := db.QueryRow("SELECT issue_count FROM repos WHERE id=$1", issue.RepoId)
	err := row.Scan(&issueCount)
	if err != nil {
		return err
	}

	// fmt.Println(issue)

	_, err = db.Exec(`INSERT INTO issues(title, user_id,body,repo_id,issue_number, created_at, updated_at)
						VALUES($1,$2,$3,$4,$5,$6,$7)`,
		issue.Title,
		issue.UserId,
		issue.Body,
		issue.RepoId,
		issueCount+1,
		time.Now().Format(time.RFC3339),
		time.Now().Format(time.RFC3339))
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
