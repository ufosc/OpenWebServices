package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ufosc/OpenWebServices/pkg/websmtp"
	"net/http"
	"strconv"
)

func main() {
	var send *websmtp.Sender
	config := GetConfig()

	// Decode port number.
	port, err := strconv.ParseInt(config.SMTP_PORT, 10, 32)
	if err != nil {
		panic("Invalid outbound port number")
	}

	if config.SMTP_SERVER != "" {
		send = websmtp.NewSenderRelay(int(port),
			&websmtp.RelayConfig{
				Host:      config.SMTP_SERVER,
				Port:      int(port),
				Username:  config.SMTP_USER,
				Password:  config.SMTP_PWD,
				UseSecure: false,
			})
	} else {
		send = websmtp.NewSender(int(port))
	}

	// Decode thread number.
	threads, err := strconv.ParseInt(config.THREADS, 10, 32)
	if err != nil {
		panic("Invalid Thread Number")
	}

	// Set server mode ("debug" or "release").
	gin.SetMode(config.GIN_MODE)
	r := gin.Default()

	r.POST("/mail/send", func(c *gin.Context) {
		var req websmtp.SendRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"ref": send.Enqueue(req),
		})
	})

	r.GET("/mail/status/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, send.GetStatus(id))
	})

	// Status check.
	r.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})

	send.Start(int(threads))
	defer send.Stop()
	r.Run("0.0.0.0:" + config.PORT)
}
