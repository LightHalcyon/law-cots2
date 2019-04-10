package main

import (
	"github.com/gin-gonic/gin"
	"github.com/reznov53/law-cots2/mq"
	"github.com/reznov53/law-cots2/download"
)

// DlFile gin routing method to download file
func DlFile(c *gin.Context) {
	rKey := mq.TokenGenerator()

	url, found := c.GetPostForm("url")
	if !found {
		return
	}

	filepath := "/dl/" + rKey
	go func(path string, ch *mq.Channel, rKey string, url string){
		err := download.File(path, url, ch, rKey)
		if err != nil {
			panic(err)
		}
	}(filepath, ch, rKey, url)

	
}