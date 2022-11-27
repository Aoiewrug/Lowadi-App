package dbupdate

import (
	"fmt"

	"github.com/Aoiewrug/Lowadi-App/lowadi-http-func-server/core/initializers"
	"github.com/Aoiewrug/Lowadi-App/lowadi-http-func-server/core/models"
)

func UpdateRunningAccountStatus(account *models.Account, statusCode int) {
	fmt.Println("I'm in db update section")

	// Status: 1 - User can enter ,
	// 2 - The bot is in Update KCK section,
	// 3 - The bot is in Update Competition section,
	// 4 - The bot is in Orki Competition section,
	// X - Other status codes...
	initializers.DB.Model(account).
		Where("user_id = ?", account.UserID). // checking for User ID
		Where("login = ?", account.Login).    // and checking for account login
		Updates(models.Account{LoggedIn: statusCode})

	//time.Sleep(10 * time.Second)

}
