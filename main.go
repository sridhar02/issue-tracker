package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func getUserHandler(c *gin.Context, db *sql.DB) {

	id := c.Param("id")

	user, err := GetUser(db, id)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, user)

	// c.Status(http.StatusNoContent)

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

func getIssueHandler(c *gin.Context, db *sql.DB) {

	id := c.Param("id")

	Id, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	issue, err := GetIssue(db, Id)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, issue)

	// c.Status(http.StatusNoContent)

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

func deleteUserHandler(c *gin.Context, db *sql.DB) {

	id := c.Param("id")

	err := DeleteUser(db, id)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	// c.JSON(http.StatusOK, repo)

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

func postUserHandler(c *gin.Context, db *sql.DB) {

	user := User{}
	err := c.BindJSON(&user)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = CreateUser(db, user)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)

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

func postIssueHandler(c *gin.Context, db *sql.DB) {

	issue := Issue{}
	err := c.BindJSON(&issue)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = CreateIssue(db, issue)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)

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

func putUserHandler(c *gin.Context, db *sql.DB) {

	user := User{}
	err := c.BindJSON(&user)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = UpdateUser(db, user)
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

	c.Status(http.StatusCreated)

}

func putIssueHandler(c *gin.Context, db *sql.DB) {

	issue := Issue{}
	err := c.BindJSON(&issue)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = UpdateIssue(db, issue)
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

	c.Status(http.StatusCreated)
}

func main() {

	connStr := "user=postgres dbname=issue_tracker host=localhost password=test1234 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
		return
	}

	// defer fmt.Println("succesfully closed end")
	defer db.Close()
	// defer fmt.Println("succesfully closed")

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return
	}

	router := gin.Default()

	router.GET("/users/:id", func(c *gin.Context) { getUserHandler(c, db) })
	router.GET("/repos/:id", func(c *gin.Context) { getRepoHandler(c, db) })
	router.GET("/issues/:id", func(c *gin.Context) { getIssueHandler(c, db) })
	router.GET("/comments/:id", func(c *gin.Context) { getCommentHandler(c, db) })
	router.DELETE("/users/:id", func(c *gin.Context) { deleteUserHandler(c, db) })
	router.DELETE("/repos/:id", func(c *gin.Context) { deleteRepoHandler(c, db) })
	router.DELETE("/issues/:id", func(c *gin.Context) { deleteIssueHandler(c, db) })
	router.DELETE("/comments/:id", func(c *gin.Context) { deleteCommentHandler(c, db) })
	router.POST("/users", func(c *gin.Context) { postUserHandler(c, db) })
	router.POST("/repos", func(c *gin.Context) { postRepoHandler(c, db) })
	router.POST("/issues", func(c *gin.Context) { postIssueHandler(c, db) })
	router.POST("/comments", func(c *gin.Context) { postCommentHandler(c, db) })
	router.PUT("/users", func(c *gin.Context) { putUserHandler(c, db) })
	router.PUT("/repos", func(c *gin.Context) { putRepoHandler(c, db) })
	router.PUT("/issues", func(c *gin.Context) { putIssueHandler(c, db) })
	router.PUT("/comments", func(c *gin.Context) { putCommentHandler(c, db) })

	err = router.Run(":8000")
	if err != nil {
		log.Fatal(err)
	}

	//CRUD functions for user under this statement

	// user, err := GetUser(db, "ac6f8b68-8f31-48ea-a436-05b9813b484b")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println(user)

	// err = CreateUser(db, User{
	// 	ID:        "ac6f8b68-8f31-48ea-a436-05b9813b484b",
	// 	Name:      "sridhar",
	// 	Username:  "sridhar02",
	// 	Email:     "kattasridhar02@gmail.com",
	// 	CreatedAt: time.Now(),
	// 	UpdatedAt: time.Now()})

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("created an user")

	// err = UpdateUser(db, User{
	// 	ID:       "ac6f8b68-8f31-48ea-a436-05b9813b484b",
	// 	Name:     "vramana",
	// 	Username: "vramana08",
	// 	Email:    "vramana@gmail.com",
	// 	// CreatedAt: time.Now()
	// 	UpdatedAt: time.Now()})

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("updated issue")

	// err = DeleteUser(db, "ac6f8b68-8f31-48ea-a436-05b9813b484b")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("deleted issue")

	// CRUD functions for repos under this statement

	// repo, err := GetRepo(db, "d360c6f3-60dc-4846-bb6a-0919a1817d5e")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println(repo)

	// err = CreateRepo(db, Repo{
	// 	ID:          "d360c6f3-60dc-4846-bb6a-0919a1817d5e",
	// 	Name:        "sridhar",
	// 	UserId:      "ac6f8b68-8f31-48ea-a436-05b9813b484b",
	// 	IssuesCount: 0,
	// 	CreatedAt:   time.Now(),
	// 	UpdatedAt:   time.Now()})

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("created an Repo")

	// err = UpdateRepo(db, Repo{
	// 	ID:          "d360c6f3-60dc-4846-bb6a-0919a1817d5e",
	// 	Name:        "vramana",
	// 	UserId:      "ac6f8b68-8f31-48ea-a436-05b9813b484b",
	// 	IssuesCount: 10,
	// 	// CreatedAt: time.Now()
	// 	UpdatedAt: time.Now()})

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("updated repo")

	// err = DeleteRepo(db, "d360c6f3-60dc-4846-bb6a-0919a1817d5e")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("deleted repo")

	//CRUD functions for issues under this section

	// issue, err := GetIssue(db, 1)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println(issue)

	// err = CreateIssue(db, Issue{
	// 	Title:       "new issue",
	// 	UserId:      "ac6f8b68-8f31-48ea-a436-05b9813b484b",
	// 	Body:        "a new has been created in the data base of issues",
	// 	RepoId:      "d360c6f3-60dc-4846-bb6a-0919a1817d5e",
	// 	IssueNumber: 0,
	// 	CreatedAt:   time.Now(),
	// 	UpdatedAt:   time.Now()},
	// 	"d360c6f3-60dc-4846-bb6a-0919a1817d5e")

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("created an Issue")

	// err = UpdateIssue(db, Issue{
	// 	ID:    1,
	// 	Title: "vramana",
	// 	Body:  "a new issue has been updated",
	// 	// IssuesCount: 10,
	// 	// CreatedAt: time.Now()
	// 	UpdatedAt: time.Now()})

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("updated issue")

	// err = DeleteIssue(db, 1)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("deleted issue")

	//CRUD functions for comments under this section

	// comment, err := GetComment(db, 1)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println(comment)

	// err = CreateComment(db, Comment{
	// 	UserId:    "ac6f8b68-8f31-48ea-a436-05b9813b484b",
	// 	Body:      "a new comment has been created in the database",
	// 	IssueId:   1,
	// 	RepoId:    "d360c6f3-60dc-4846-bb6a-0919a1817d5e",
	// 	CreatedAt: time.Now(),
	// 	UpdatedAt: time.Now()})

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println("created an Comment")

	// err = UpdateComment(db, Comment{
	// 	ID:        1,
	// 	Body:      "a new comment has been updated",
	// 	UpdatedAt: time.Now()})

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("updated comment")

	// err = DeleteComment(db, 1)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("deleted repo")

}
