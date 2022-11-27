package active

import (
	"fmt"
	"net/http"

	"github.com/Aoiewrug/Lowadi-App/lowadi-http-api/core/initializers"
	"github.com/Aoiewrug/Lowadi-App/lowadi-http-api/core/models"
	"github.com/gin-gonic/gin"
)

// Auto collected after the test run
// Get all accs with specific user
func GetCompetition(c *gin.Context) {
	// Show lowadi account
	body := &models.Account{}

	// SELECT * FROM `accounts` WHERE 'user_id' = 1
	initializers.DB.Where("user_id = ?", fmt.Sprint(GetUserID(c))).Find(&body)
	//initializers.DB.Find(&body)
	c.JSON(200, &body)
}

func SetCompetitionToGameAccount(c *gin.Context) {
	account := &models.Account{}

	if c.BindJSON(&account) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to read account update body",
			"payload": "Please check syntax",
		})

		return
	}

	initializers.DB.Model(account).
		Where("user_id = ?", fmt.Sprint(GetUserID(c))). // checking for User ID
		Where("login = ?", account.Login).              // and checking for account login
		Updates(models.Account{
			CompetitionsLink:  account.CompetitionsLink,
			CompetitionsEvent: account.CompetitionsEvent,
		}) //

	// Swap to somehting else
	c.JSON(200, gin.H{
		"message": "Updated account settings are: ",
		"payload": account,
	})

}
