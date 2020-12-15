package reddit

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	redditItemsURL string = "https://www.reddit.com/r/technology.json"
	redditAppender string = "https://www.reddit.com"
	userAgent      string = "Tech-News-API/1.0"
)

var (
	errGetItems      = errors.New("could not fetch items")
	errBadQueryParam = errors.New("bad query paramter")
)

// Story struct
type Story struct {
	Name  string `json:"name"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

// Reddit struct
type Reddit struct {
	Data Data `json:"data"`
}

// Data struct
type Data struct {
	Children []Children `json:"children"`
}

// Children struct
type Children struct {
	Data Story `json:"data"`
}

// GetAll stories handler
func GetAll(c *gin.Context) {
	var stories []Story
	redditStories, err := getRedditItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errGetItems.Error()})
		return
	}

	for _, redditItem := range redditStories {
		if redditItem.Data.URL[0] == '/' {
			redditItem.Data.URL = redditAppender + redditItem.Data.URL
		}
		stories = append(stories, redditItem.Data)
	}

	c.JSON(http.StatusOK, gin.H{"data": stories})
}

func getRedditItems() ([]Children, error) {
	var redditResponse Reddit

	client := &http.Client{}

	req, err := http.NewRequest("GET", redditItemsURL, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = json.Unmarshal(data, &redditResponse)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return redditResponse.Data.Children, nil
}
