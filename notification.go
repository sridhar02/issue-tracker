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

func CreateNotification(db *sql.DB, issueId int, userId string, repoId string) error {

	_, err := db.Exec(`INSERT INTO notifications(read,created_at,issue_id,user_id,repo_id)
						VALUES($1,$2,$3,$4,$5)`,
		"unread",
		time.Now().Format(time.RFC3339),
		issueId,
		userId,
		repoId)

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

func CommentNotifications(db *sql.DB, issueId int, currentRepo CurrentRepo, currentUser User) error {

	notificationUserIds := map[string]int{}
	rows, err := db.Query(`SELECT DISTINCT users.id FROM USERS JOIN comments ON 
						comments.user_id = users.id WHERE issue_id = $1`, issueId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	for rows.Next() {
		var UsersId string
		err = rows.Scan(&UsersId)
		if err != nil {
			fmt.Println(err)
			return err
		}
		notificationUserIds[UsersId] = 0
	}

	rows, err = db.Query(`SELECT user_id FROM collaborators WHERE repo_id = $1`, currentRepo.RepoId)
	if err != nil {
		fmt.Println(err)
		return err
	}
	var userId string
	for rows.Next() {
		err = rows.Scan(&userId)
		if err != nil {
			fmt.Println(err)
			return err
		}
		notificationUserIds[userId] = 0
	}
	notificationUserIds[userId] = 0
	notificationUserIds[currentRepo.UserId] = 0
	delete(notificationUserIds, currentUser.ID)

	for UserId, _ := range notificationUserIds {
		err = CreateNotification(db, issueId, UserId, currentRepo.RepoId)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

func CollaboratorNotifications(db *sql.DB, currentRepo CurrentRepo, IssueId int) error {

	notificationUserIds := map[string]int{}
	rows, err := db.Query(`SELECT user_id FROM collaborators WHERE repo_id = $1`, currentRepo.RepoId)
	if err != nil {
		fmt.Println(err)
		return err
	}
	for rows.Next() {
		var userId string
		err = rows.Scan(&userId)
		if err != nil {
			fmt.Println(err)
			return err
		}
		notificationUserIds[userId] = 0
	}
	notificationUserIds[currentRepo.UserId] = 0

	for UserId, _ := range notificationUserIds {
		err = CreateNotification(db, IssueId, UserId, currentRepo.RepoId)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}
