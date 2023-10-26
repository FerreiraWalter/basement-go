package main

import (
    "fmt"
    "net/http"
    "github.com/gocql/gocql"
    "github.com/gin-gonic/gin"
)

type Post struct {
    PostID  gocql.UUID `json:"post_id"`
    Title   string     `json:"title"`
    Summary string     `json:"summary"`
    Body    string     `json:"body"`
}

func main() {
    r := gin.Default()

    cluster := gocql.NewCluster("localhost")
    cluster.Keyspace = "my_keyspace"
    session, _ := cluster.CreateSession()
    defer session.Close()

    r.POST("/post", func(c *gin.Context) {
        var post Post
        if err := c.BindJSON(&post); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        if post.Title == "" || post.Summary == "" || post.Body == "" {
            c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Incomplete data"})
            return
        }

        post.PostID = gocql.TimeUUID()
        if err := session.Query(
            "INSERT INTO posts (post_id, title, summary, body) VALUES (?, ?, ?, ?)",
            post.PostID, post.Title, post.Summary, post.Body).Exec(); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusCreated, post)
    })

    r.GET("/healthz", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "Health check OK"})
    })

    r.Run(":8080")
}
