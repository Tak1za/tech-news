package story

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Story struct
type Story struct {
	By          string  `json:"by"`
	Descendants int     `json:"descendants"`
	ID          int     `json:"id"`
	Kids        []int   `json:"kids"`
	Time        float64 `json:"time"`
	Title       string  `json:"title"`
	Type        string  `json:"type"`
	URL         string  `json:"url"`
}

const (
	defaultSize   int    = 10
	defaultPage   int    = 1
	maxPageSize   int    = 500
	getStoriesURL string = "https://hacker-news.firebaseio.com/v0/newstories.json"
	getItemURL    string = "https://hacker-news.firebaseio.com/v0/item"
)

var (
	errBadQueryParam = errors.New("bad query paramter")
	errGetItems      = errors.New("could not fetch items")
)

type res struct {
	Data Story
	Err  error
}

// GetAll stories handler
func GetAll(c *gin.Context) {
	pageSize, pageNumber, err := getAndValidateQueryParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errBadQueryParam.Error()})
		return
	}

	items, err := getItemIDs(pageSize, pageNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errGetItems.Error()})
		return
	}

	dataChan := make(chan res)

	stories := make([]Story, 0, len(items))
	for _, itemID := range items {
		go getItem(&dataChan, itemID)
	}

	for range items {
		item := <-dataChan
		if item.Err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errGetItems.Error()})
			return
		}

		stories = append(stories, item.Data)
	}

	c.JSON(http.StatusOK, gin.H{"data": stories})
}

func getAndValidateQueryParams(c *gin.Context) (int, int, error) {
	size, ok := c.GetQuery("size")
	if !ok {
		size = strconv.Itoa(defaultSize)
	}
	pageSize, err := strconv.Atoi(size)
	if err != nil {
		log.Println(err)
		return 0, 0, err
	}

	number, ok := c.GetQuery("page")
	if !ok {
		number = strconv.Itoa(defaultPage)
	}
	pageNumber, err := strconv.Atoi(number)
	if err != nil {
		log.Println(err)
		return 0, 0, err
	}

	if pageSize <= 0 || pageNumber <= 0 {
		log.Println("query parameters can't be less than or equal to zero")
		return 0, 0, errBadQueryParam
	}

	if pageSize >= maxPageSize {
		pageSize = maxPageSize
		pageNumber = defaultPage
	}

	return pageSize, pageNumber, nil
}

func getItemIDs(pageSize, pageNumber int) ([]int, error) {
	var items []int
	response, err := http.Get(getStoriesURL)
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

func getItem(resChan *chan res, itemID int) {
	var story Story
	response, err := http.Get(getItemURL + fmt.Sprintf("/%d.json", itemID))
	if err != nil {
		log.Println(err)
		*resChan <- res{
			Data: story,
			Err:  err,
		}
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		*resChan <- res{
			Data: story,
			Err:  err,
		}
	}

	err = json.Unmarshal(data, &story)
	if err != nil {
		log.Println(err)
		*resChan <- res{
			Data: story,
			Err:  err,
		}
	}

	*resChan <- res{
		Data: story,
		Err:  nil,
	}
}
