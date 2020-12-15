package story

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	defaultSize   int    = 10
	defaultPage   int    = 1
	maxPageSize   int    = 500
	newStoriesURL string = "https://hacker-news.firebaseio.com/v0/newstories.json"
)

var (
	errBadQueryParam = errors.New("bad query paramter")
	errGetItems      = errors.New("could not fetch items")
)

// GetAll stories handler
func GetAll(c *gin.Context) {
	size, ok := c.GetQuery("size")
	if !ok {
		size = strconv.Itoa(defaultSize)
	}
	pageSize, err := strconv.Atoi(size)
	if err != nil {
		log.Println(errBadQueryParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": errBadQueryParam})
		return
	}

	number, ok := c.GetQuery("page")
	if !ok {
		number = strconv.Itoa(defaultPage)
	}
	pageNumber, err := strconv.Atoi(number)
	if err != nil {
		log.Println(errBadQueryParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": errBadQueryParam})
		return
	}

	if pageSize <= 0 {
		pageSize = defaultPage
	}

	if pageNumber <= 0 {
		pageNumber = defaultSize
	}

	if pageSize >= maxPageSize {
		pageSize = maxPageSize
		pageNumber = defaultPage
	}

	items, err := getItems(pageSize, pageNumber)
	if err != nil {
		log.Println(errGetItems)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errGetItems})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": items})
}

func getItems(pageSize, pageNumber int) ([]int, error) {
	var items []int
	response, err := http.Get(newStoriesURL)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = json.Unmarshal(data, &items)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	startIndex := 0
	if pageNumber > 1 {
		startIndex = pageNumber * pageSize
		if startIndex > 500 {
			startIndex = 0
		}
	}

	endIndex := startIndex + pageSize
	if endIndex > maxPageSize {
		endIndex = maxPageSize
	}

	return items[startIndex:endIndex], nil
}
