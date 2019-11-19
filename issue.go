package main

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type Issue struct {
	ID          int        `json:"id,omitempty"`
	Title       string     `json:"title,omitempty"`
	UserId      string     `json:"user_id,omitempty"`
	Body        string     `json:"body,omitempty"`
	RepoId      string     `json:"repo_id,omitempty"`
	IssueNumber int        `json:"issue_number,omitempty"`
	Status      string     `json:"status,omitempty"`
	Pinned      string     `json:"pinned,omitempty"`
	Lock        string     `json:"lock,omitempty"`
	User        User       `json:"user,omitempty"`
	CreatedAt   time.Time  `json:"created_at,omitempty"`
	UpdatedAt   time.Time  `json:"updated_at,omitempty"`
	Assignees   []Assignee `json:"assignees"`
}

func GetIssue(db *sql.DB, Id int) (Issue, error) {

	var title, body, repoId, status, createdAt, userId, updatedAt, pinned, lock string
	var issueNumber int

	row := db.QueryRow(`SELECT
		                     title,body,repo_id,issue_number,status,
		                     pinned,lock,user_id,created_at,
						     updated_at
						     FROM
						     issues
						     WHERE id=$1`, Id)
	err := row.Scan(&title, &body, &repoId, &issueNumber, &status, &pinned, &lock, &userId, &createdAt, &updatedAt)
	if err != nil {
		fmt.Println(err)
		return Issue{}, err
	}
	CreatedAt, err := time.Parse(time.RFC3339, createdAt)
	if err != nil {
		fmt.Println(createdAt)
		fmt.Println(err)
		return Issue{}, err
	}

	UpdatedAt, err := time.Parse(time.RFC3339, updatedAt)
	if err != nil {
		fmt.Println(err)
		return Issue{}, err
	}

	user, err := GetUser(db, userId)
	if err != nil {
		fmt.Println(err)
		return Issue{}, err
	}

	rows, err := db.Query("SELECT user_id FROM assignees WHERE issue_id=$1", Id)
	if err != nil {
		fmt.Println(err)
		return Issue{}, err
	}

	assignees := []Assignee{}

	for rows.Next() {
		var userID string
		err = rows.Scan(&userID)
		if err != nil {
			fmt.Println(err)
			return Issue{}, err
		}
		user, err := GetUser(db, userID)
		if err != nil {
			fmt.Println(err)
			return Issue{}, err

		}
		assignee := Assignee{
			User: user,
		}
		assignees = append(assignees, assignee)
	}
	issue := Issue{
		ID:          Id,
		Title:       title,
		Body:        body,
		RepoId:      repoId,
		IssueNumber: issueNumber,
		Status:      status,
		CreatedAt:   CreatedAt,
		UpdatedAt:   UpdatedAt,
		Pinned:      pinned,
		Lock:        lock,
		User:        user,
		Assignees:   assignees,
	}

	return issue, nil
}

func CreateIssue(db *sql.DB, issue Issue) (int, int, error) {

	var issueId int
	var issueCount int

	row := db.QueryRow("SELECT issue_count FROM repos WHERE id=$1", issue.RepoId)
	err := row.Scan(&issueCount)
	if err != nil {
		return 0, 0, err
	}

	// fmt.Println(issue)

	err = db.QueryRow(`INSERT INTO
		                        issues (
		                        title, user_id,body,repo_id,issue_number,status,
		                        pinned,lock,created_at, updated_at
		                        )
						        VALUES(
						        $1,$2,$3,$4,$5,$6,$7,$8,$9,$10
						        ) RETURNING id`,
		issue.Title,
		issue.UserId,
		issue.Body,
		issue.RepoId,
		issueCount+1,
		"Open",
		"Unpinned",
		"Unlocked",
		time.Now().Format(time.RFC3339),
		time.Now().Format(time.RFC3339)).Scan(&issueId)
	if err != nil {
		return 0, 0, err
	}

	_, err = db.Exec(`UPDATE repos SET issue_count =$1, updated_at = $2 WHERE id = $3`,
		issueCount+1, time.Now().Format(time.RFC3339), issue.RepoId)

	if err != nil {
		return 0, 0, err
	}

	return issueId, issueCount + 1, nil
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

func postIssueHandler(c *gin.Context, db *sql.DB) {

	userId, err := authorization(c, db)
	if err != nil {
		return
	}
	username := c.Param("owner")
	repoName := c.Param("repo")
	var repoId string
	row := db.QueryRow(`SELECT repos.id FROM repos JOIN users ON repos.user_id= users.id
						       WHERE repos.name=$1 AND users.username=$2`, repoName, username)
	err = row.Scan(&repoId)
	if err != nil {
		fmt.Println(err)
		return
	}

	issue := Issue{}
	err = c.BindJSON(&issue)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	issue.UserId = userId
	issue.RepoId = repoId

	_, _, err = CreateIssue(db, issue)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)

}

func putIssueHandler(c *gin.Context, db *sql.DB) {

	_, err := authorization(c, db)
	if err != nil {
		return
	}
	username := c.Param("owner")
	repoName := c.Param("repo")
	issueNumber := c.Param("issue_number")

	var Id int
	row := db.QueryRow(`WITH repo_cte AS (
		                    SELECT repos.id FROM repos JOIN users ON repos.user_id= users.id
                		    WHERE repos.name= $1 AND users.username= $2
                		)
                		SELECT id FROM issues
                		WHERE  repo_id in (select id from repo_cte) AND issue_number=$3`,
		repoName, username, issueNumber)
	err = row.Scan(&Id)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	issue := Issue{}
	err = c.BindJSON(&issue)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	query := psql.Update("issues")
	if issue.Status != "" {
		query = query.Set("status", issue.Status)
	}
	if issue.Title != "" {
		query = query.Set("title", issue.Title)
	}
	if issue.Body != "" {
		query = query.Set("body", issue.Body)
	}
	query = query.Where(sq.Eq{"id": Id})

	q, args, err := query.ToSql()

	_, err = db.Exec(q, args...)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}


	c.Status(http.StatusNoContent)
}

func deleteIssueHandler(c *gin.Context, db *sql.DB) {

	id := c.Param("id")

	Id, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = DeleteIssue(db, Id)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	// c.JSON(http.StatusOK, repo)

	c.Status(http.StatusNoContent)

}

func putLockHandler(c *gin.Context, db *sql.DB) {

	_, err := authorization(c, db)
	if err != nil {
		return
	}
	username := c.Param("owner")
	repoName := c.Param("repo")
	issueNumber := c.Param("issue_number")

	issue := Issue{}
	err = c.BindJSON(&issue)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	_, err = db.Exec(`WITH repo_cte AS (
		                     SELECT repos.id  FROM repos JOIN users ON repos.user_id= users.id
                		     WHERE repos.name= $1 AND users.username= $2
                		     ) ,
                		issue_CTE AS (
                			 SELECT id FROM issues
                		     WHERE repo_id in (SELECT id from repo_cte) and issue_number=$3
                		     )
                		UPDATE issues SET lock= $4 WHERE id in ( select id from issue_cte )`,
		repoName, username, issueNumber, issue.Lock)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)

}

func putPinHandler(c *gin.Context, db *sql.DB) {

	_, err := authorization(c, db)
	if err != nil {
		return
	}
	username := c.Param("owner")
	repoName := c.Param("repo")
	issueNumber := c.Param("issue_number")

	var Id int
	var repoId string
	row := db.QueryRow(`WITH repo_cte AS (
		                    SELECT repos.id FROM repos JOIN users ON repos.user_id= users.id
                		    WHERE repos.name= $1 AND users.username= $2
                	        )
                	    SELECT id,repo_id FROM issues
                	    WHERE repo_id in (SELECT id FROM repo_cte) AND issue_number=$3`,
		repoName, username, issueNumber)
	err = row.Scan(&Id, &repoId)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	issue := Issue{}
	err = c.BindJSON(&issue)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var count int

	row = db.QueryRow(` SELECT count(*) FROM issues WHERE pinned = 'Pinned' AND repo_id = $1`, repoId)
	err = row.Scan(&count)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if count >= 3 {
		c.Status(422)
		return
	}

	_, err = db.Exec(`UPDATE issues SET pinned = $1 WHERE id = $2`, issue.Pinned, Id)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
