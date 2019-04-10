package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/reznov53/law-cots2/mq"
	"github.com/reznov53/law-cots2/download"
)

type appError struct {
	Code	int    `json:"status"`
	Message	string `json:"message"`
}

// DlFile gin routing method to download file
func DlFile(c *gin.Context) {
	rKey := c.GetHeader("X-Routing-Key")

	url, found := c.GetPostForm("url")
	if !found {
		c.JSON(http.StatusBadRequest, appError{
			Code:		http.StatusBadRequest,
			Message:	"URL not found, did you guys specify the URL?",
		})
		return
	}

	filepath := "/dl/" + rKey
	go func(path string, ch *mq.Channel, rKey string, url string){
		err := download.File(path, url, ch, rKey)
		if err != nil {
			panic(err)
		}
	}(filepath, ch, rKey, url)

	c.JSON(http.StatusOK, appError{
		Code:		http.StatusOK,
		Message:	"Download starts",
	})
	return
}