package main

import (
	"time"

	"github.com/Tak1za/tech-news/pkg/hn"
	"github.com/Tak1za/tech-news/pkg/reddit"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	hackerNewsGroup := r.Group("/hn")
	{
		hackerNewsGroup.GET("/stories", hn.GetAll)
	}
	redditGroup := r.Group("/r")
	{
		redditGroup.GET("/stories", reddit.GetAll)
	}
	r.Run()
}
