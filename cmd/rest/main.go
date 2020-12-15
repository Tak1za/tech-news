package main

import (
	"github.com/Tak1za/tech-news/pkg/story"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/stories", story.GetAll)
	r.Run()
}
