package active

import (
	"fmt"
	"net/http"

	"github.com/Aoiewrug/Lowadi-App/lowadi-http-api/core/controllers/helpers"
	"github.com/Aoiewrug/Lowadi-App/lowadi-http-api/core/initializers"
	"github.com/Aoiewrug/Lowadi-App/lowadi-http-api/core/models"
	"github.com/gin-gonic/gin"
)

// Main functions to interact with lowady website

func GameEnterPoint(c *gin.Context) {
	body := &models.Account{}
	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to read update status body",
			"payload": "Please check syntax",
		})

		return
	}

	userID := GetUserID(c)

	initializers.DB.Model(&body).
		Where("user_id = ?", fmt.Sprint(userID)). // checking for User ID
		Where("login = ?", body.Login).           // and checking for account login
		Find(&body)                               // Get this acc info

	fmt.Println("User ID: ", userID, " Password is: ", body.Pass, " Login is: ", body.Login)

	// Check expiration date
	// Check "Active" status
	// Check "AlreadyExpired" status
	// If all 3 parameters are okay -> pass forward
	// else -> reject

	// ADD PROXY CONFIG FROM DB

	// If it is:
	// 1 = Entered
	// 2 = Error
	statusCode, message := helpers.GameLogIn(body)
	fmt.Println("Status code is: ", statusCode)
	fmt.Println("message is: ", message)

}
