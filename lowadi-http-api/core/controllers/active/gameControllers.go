package active

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Aoiewrug/Lowadi-App/lowadi-http-api/core/controllers/external"
	"github.com/Aoiewrug/Lowadi-App/lowadi-http-api/core/initializers"
	"github.com/Aoiewrug/Lowadi-App/lowadi-http-api/core/models"
	"github.com/gin-gonic/gin"
)

// Main functions to interact with lowady website
// Accepts multiple user logins from the front-end
func GameEnterPoint(c *gin.Context) {
	body := models.StartupAccount{}

	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to read update status body",
			"payload": "Please check syntax",
		})

		return
	}
	/*
		fmt.Println()
		fmt.Println(body.Logins)
		fmt.Println(len(body.Logins))
		fmt.Println(body.Logins[0])
		fmt.Println()
	*/
	userID := GetUserID(c)

	for _, login := range body.Logins {
		account := &models.Account{}
		initializers.DB.Model(&models.Account{}).
			Where("user_id = ?", fmt.Sprint(userID)). // checking for User ID
			Where("login = ?", login).                // and checking for account login
			Find(account)                             // Get this acc info

		fmt.Println("User ID: ", userID, " Login is: ", account.Login, " Password is: ", account.Pass)
		//fmt.Println(account)
		fmt.Println("Sending this account to the workers queue of the remote server ")

		requestBody, err := json.Marshal(&account) //(&models.User{Email: r.FormValue("email"), Password: r.FormValue("password")})
		if err != nil {
			fmt.Println("Can't marshal JSON")
		}

		statusCode, err := external.PostJson(os.Getenv("FUNC_SERVER_ADDR")+"/func/run", requestBody) // http://127.0.0.1:8002/func/run
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Error occured on remote call with this account",
				"payload": account.Login,
			})

			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Successfully send this account",
			"payload": statusCode,
		})

		// fmt.Println(statusCode)
		// fmt.Println(err)
		/*
			statusCode, message := helpers.GameLogIn(account)
			fmt.Println("Status code is: ", statusCode)
			fmt.Println("message is: ", message)

		*/
	}

	// Check expiration date
	// Check "Active" status
	// Check "AlreadyExpired" status
	// If all 3 parameters are okay -> pass forward
	// else -> reject

	// ADD PROXY CONFIG FROM DB

	// If it is:
	// 1 = Entered
	// 2 = Error

}
