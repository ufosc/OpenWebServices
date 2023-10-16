package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ufosc/OpenWebServices/websmtp/websmtp"
	"net/http"
)

// TODO: Currently, the purpose of this program is to provide a mechanism for
// sending email verification emails. Authentication should eventually be
// implemented, but otherwise this is sufficient as a demo.
func main() {
	send := websmtp.NewSender("notifications@ufosc.org")
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

	go send.Run()
	r.Run()
}
