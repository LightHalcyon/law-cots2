package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/reznov53/law-cots2/mq"
)

var ch *mq.Channel
var err error

func main() {
	// url := "amqp://" + os.Getenv("UNAME") + ":" + os.Getenv("PW") + "@" + os.Getenv("URL") + ":" + os.Getenv("PORT") + "/"
	url := "amqp://1406568753:167664@152.118.148.103:5672/"
	// vhost := os.Getenv("VHOST")
	vhost := "1406568753"
	// exchangeName := os.Getenv("EXCNAME")
	exchangeName := "1406568753"
	exchangeType := "direct"
	ch, err = mq.InitMQ(url, vhost, exchangeName, exchangeType)
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	r.Static("/asset", "./asset")
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Upload Page",
		})
	})
	r.POST("/dl", DlFile)
	
	r.Run("0.0.0.0:20609")
	ch.Conn.Close()
	ch.Ch.Close()
}