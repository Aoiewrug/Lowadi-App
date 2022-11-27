package handlers

import (
	"net/http"

	"github.com/Aoiewrug/Lowadi-App/lowadi-http-func-server/core/models"
	"github.com/gin-gonic/gin"
)

func StartRemoteOrki(c *gin.Context) {
	account := &models.Account{}

	if c.BindJSON(&account) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to read account update body",
			"payload": "Please check syntax",
		})

		return
	}

	// Swap to somehting else
	c.JSON(200, gin.H{
		"message": "Updated account settings are: ",
		"payload": account,
	})

}
