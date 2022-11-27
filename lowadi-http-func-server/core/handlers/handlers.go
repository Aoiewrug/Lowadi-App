package handlers

import (
	"fmt"
	"net/http"

	"github.com/Aoiewrug/Lowadi-App/lowadi-http-func-server/core/dispatch"
	"github.com/Aoiewrug/Lowadi-App/lowadi-http-func-server/core/models"
	"github.com/gin-gonic/gin"
)

func StartRemoteOrki(c *gin.Context) {
	fmt.Println("StartRemoteOrki: ")
	account := &models.Account{}

	if c.BindJSON(&account) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to read account update body",
			"payload": "Please check syntax",
		})

		return
	}

	dispatch.AddAccount(account)

	c.JSON(200, gin.H{
		"message": "Successfully sended this acc to the workers queue:",
		"payload": account.Login,
	})

}

func TestStartRemoteOrki(c *gin.Context) {
	fmt.Println("StartRemoteOrki: ")
	account := &models.Account{}

	if c.BindJSON(&account) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to read account update body",
			"payload": "Please check syntax",
		})

		return
	}

	dispatch.AddAccount(account)

	c.JSON(200, gin.H{
		"message": "Successfully sended this acc to the workers queue:",
		"payload": account.Login,
	})

}
