package main

import (
	"github.com/Tak1za/tech-news/pkg/hn"
	"github.com/Tak1za/tech-news/pkg/reddit"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
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
