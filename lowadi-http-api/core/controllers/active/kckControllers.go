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
func GetKCKs(c *gin.Context) {
	// Show lowadi account
	body := &[]models.KCK{}

	// SELECT * FROM `accounts` WHERE 'user_id' = 1
	initializers.DB.Where("user_id = ?", fmt.Sprint(GetUserID(c))).Find(&body)

	c.JSON(200, &body)
}

// Get all KCKs of the exact use(r/rs)
func GetKCK(c *gin.Context) {
	// get login and user id
	body := &models.StartupAccount{}

	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to read account update body",
			"payload": "Please check syntax",
		})

		return
	}

	res := &[]models.KCK{}
	// Find KCK link/id in db
	initializers.DB.
		Where("user_id = ?", fmt.Sprint(GetUserID(c))).
		Where("login = ?", body.Logins).
		Find(&res)

	fmt.Println(res)

	c.JSON(200, &res)
}

// Hard delete of an account
func DeleteKCK(c *gin.Context) {
	body := &models.KCK{}

	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to read account update body",
			"payload": "Please check syntax",
		})

		return
	}

	initializers.DB.Model(&body).
		Where("login = ?", body.Login).
		Where("user_id = ?", fmt.Sprint(GetUserID(c))).
		Delete(&body) // checking for User ID

	// Swap to somehting else
	c.JSON(http.StatusOK, gin.H{
		// Rename to "Invalid user or password"
		"message": "Successfully removed this KCK: ",
		"payload": body.StableName,
	})
}

func SetKCKtoGameAccount(c *gin.Context) {
	kck := &models.KCK{}
	account := &models.Account{}

	if c.BindJSON(&account) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to read account update body",
			"payload": "Please check syntax",
		})

		return
	}

	// Find KCK link/id in db
	initializers.DB.
		Where("user_id = ?", fmt.Sprint(GetUserID(c))).
		Where("login = ?", account.Login).
		Where("stable_name = ?", &account.StableName).
		Find(&kck)

	fmt.Println(kck)

	initializers.DB.Model(account).
		Where("user_id = ?", fmt.Sprint(GetUserID(c))). // checking for User ID
		Where("login = ?", account.Login).              // and checking for account login
		Updates(models.Account{
			StableName:       account.StableName,
			StableLink:       kck.StableLink,
			MaxAge:           account.MaxAge,
			AdvantagesFuraj:  account.AdvantagesFuraj,
			AdvantagesOvec:   account.AdvantagesOvec,
			AdvantagesCarrot: account.AdvantagesCarrot,
			MaxDailyPrice:    account.MaxDailyPrice,
			BirthHorses:      account.BirthHorses,
			BirthHorsesName:  account.BirthHorsesName,
		}) //

	// Swap to somehting else
	c.JSON(200, gin.H{
		"message": "Updated account settings are: ",
		"payload": account,
	})

}
