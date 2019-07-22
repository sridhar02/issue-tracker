package main

import (
	"database/sql"
	"fmt"
	"time"
)

type Notification struct {
	ID        int       `json:"id,omitempty"`
	Read      string    `json:"read,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	IssueId   int       `json:"issue_id,omitempty"`
	UserId    string    `json:"user_id,omitempty"`
	RepoId    string    `json:"repo_id,omitempty"`
}

func CreateNotification(db *sql.DB, notification Notification) error {

	_, err := db.Exec(`INSERT INTO notifications(read,created_at,issue_id,user_id,repo_id)
						VALUES($1,$2,$3,$4,$5)`,
		"unread",
		time.Now().Format(time.RFC3339),
		notification.IssueId,
		notification.UserId,
		notification.RepoId)

	if err != nil {
		return err
	}

	return nil
}

func ReadNotification(db *sql.DB, id int) (Notification, error) {

	var read, createdAt, userId, repoId string

	var issueId int

	row := db.QueryRow("SELECT  id,read,created_at,issue_id,,user_id,repo_id FROM notifications WHERE id=$1", id)
	err := row.Scan(&id, &read, &createdAt, &issueId, &userId, &repoId)
	if err != nil {
		return Notification{}, err
	}

	fmt.Println(id)

	CreatedAt, err := time.Parse(time.RFC3339, createdAt)
	if err != nil {
		return Notification{}, err
	}

	notification := Notification{
		ID:        id,
		Read:      read,
		CreatedAt: CreatedAt,
		IssueId:   issueId,
		UserId:    userId,
		RepoId:    repoId,
	}
	return notification, nil
}
