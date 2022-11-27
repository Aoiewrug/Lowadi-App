package active

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Aoiewrug/Lowadi-App/lowadi-http-api/core/initializers"
	"github.com/Aoiewrug/Lowadi-App/lowadi-http-api/core/models"
	"github.com/gin-gonic/gin"
)

// ADD PROXY ASSIGNMENT !!!
// Create an account with 1 day trial period
func AddAcc(c *gin.Context) {
	body := &models.Account{}
	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to read account creation body",
			"payload": "Please check syntax",
		})

		return
	}

	// Account settings
	body.UserID = GetUserID(c)                    // Bind to a specific user
	body.AccCreated = time.Now()                  // Set the creation time
	body.AccEnds = time.Now().Add(time.Hour * 24) // Add trial period? (24 hours default)
	body.AccUpdated = time.Now()                  // Last user $ activity
	body.AlreadyExpired = 1                       // 2=expired, 1=active                 (Manual turn off/on)
	body.Active = 1                               // 2=expired, 1=active add 1 day trial (Auto turn off (balance exceeded))
	body.OverDueHours = 0                         //
	body.CostPerDay = 5                           // How many rub per day we want per account
	body.GameWebside = "https://www.lowadi.com/"  // Default website
	body.LoggedIn = 1                             // Status: 1 - User can enter, 2 - The bot is running, 3 - The bot is running...

	// game settings
	body.UpdateKCK = 1            // 2=no, 1=yes Update stables list?
	body.RunOrki = 2              // 2=no, 1=yes Run Orki?
	body.Competitions = 2         // 2=no, 1=yes Run competitions?
	body.BirthHorses = 2          // 2=no, 1=yes Should we born new horses?
	body.BirthHorsesName = "male" // default horse name

	// KCK sing-in settings
	body.MaxAge = "30"        // Difines the peak year of the running horses (selecting it via filter)
	body.MaxDailyPrice = 30   // Max daily price for KCK sing-in
	body.AdvantagesFuraj = 2  // 1=use this check-mark. 2=don't
	body.AdvantagesOvec = 2   // 1=use this check-mark. 2=don't
	body.AdvantagesCarrot = 2 // 1=use this check-mark. 2=don't

	result := initializers.DB.Create(body)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Can't add new account",
			"payload": body.Login,
		})
		return
	}

	// Swap to somehting else
	c.JSON(200, body)
}

// Hard delete of an account
func DeleteAcc(c *gin.Context) {
	body := &models.Account{}

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
		"message": "Successfully removed this account: ",
		"payload": body.Login,
	})
}

// Get all accs with specific user
func GetAcc(c *gin.Context) {
	// Show lowadi account
	body := &[]models.Account{}

	// SELECT * FROM `accounts` WHERE 'user_id' = 1
	initializers.DB.Where("user_id = ?", fmt.Sprint(GetUserID(c))).Find(&body)
	//initializers.DB.Find(&body)
	c.JSON(200, &body)
}

// Updating account's password
func UpdatePass(c *gin.Context) {
	body := &models.Account{}

	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to read account update body",
			"payload": "Please check syntax",
		})

		return
	}

	initializers.DB.Model(&body).
		Where("user_id = ?", fmt.Sprint(GetUserID(c))). // checking for User ID
		Where("login = ?", body.Login).                 // and checking for account login
		Updates(models.Account{Pass: body.Pass})        // updating with new pass

	// Swap to somehting else
	c.JSON(200, gin.H{
		"message": "New password is: ",
		"payload": body.Pass,
	})
}

// Turn off/on the account manually
func UpdateStatus(c *gin.Context) {
	body := &models.Account{}
	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to read update status body",
			"payload": "Please check syntax",
		})

		return
	}

	initializers.DB.Model(&body).
		Where("user_id = ?", fmt.Sprint(GetUserID(c))). // checking for User ID
		Where("login = ?", body.Login).                 // and checking for account login
		Updates(models.Account{Active: body.Active})

	c.JSON(200, gin.H{
		"message": "Your account status is: ",
		"payload": body.Active,
	})
}

// Updates website account belongs to - https://www.lowadi.com/
func UpdateVersion(c *gin.Context) {
	body := &models.Account{}
	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to read update status body",
			"payload": "Please check syntax",
		})

		return
	}

	initializers.DB.Model(&body).
		Where("user_id = ?", fmt.Sprint(GetUserID(c))). // checking for User ID
		Where("login = ?", body.Login).                 // and checking for account login
		Updates(models.Account{GameWebside: body.GameWebside})

	c.JSON(200, gin.H{
		"message": "New website is : ",
		"payload": body.GameWebside,
	})
}

// Updates expiration date
// Need to agjust it with timer, balance. If balbance is to low - error
func UpdateExpiration(c *gin.Context) {
	// Temporary save the request's body

	// To get user balance from DB
	user := &models.User{}
	// To get cost per day from DB
	account := &models.Account{}
	// To get request body
	body := &models.Account{}
	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to read update expiration date body",
			"payload": "Please check syntax",
		})

		return
	}

	// Go to user accounts via ID and get balance
	getUserBalance := initializers.DB.Model(&user).
		Where("id = ?", fmt.Sprint(GetUserID(c))).
		Find(&user)

	if getUserBalance.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Can't check the user's balance",
			"payload": "",
		})
		return

	}

	// Go to game accounts via ID and get balance
	initializers.DB.Model(&account).
		Where("user_id = ?", fmt.Sprint(GetUserID(c))). // checking for User ID
		Where("login = ?", body.Login).                 // and checking for account login
		Find(&account)

	// Calculations. If balance < (body.Timer * cost_per_day) - reject
	// 		We use field models.User.Timer as a buffer for setting up for how
	//		many days should we increase the expitaration date
	sum := user.Balance - body.Timer*account.CostPerDay

	if sum < 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "You wand to deduct too much. Please add this amount to your balance: ",
			"payload": -1 * sum,
		})
		// If success assign new expiration date and deduct the balance
	} else {
		// Update game account
		initializers.DB.Model(&body).
			Where("user_id = ?", fmt.Sprint(GetUserID(c))). // checking for User ID
			Where("login = ?", body.Login).                 // and checking for account login
			Updates(models.Account{
				AccUpdated:     time.Now(),
				AccEnds:        time.Now().Add(time.Hour * 24 * time.Duration(body.Timer)),
				AlreadyExpired: 2,
			})

		// Update user account
		initializers.DB.Model(&user).
			Where("id = ?", fmt.Sprint(GetUserID(c))). // checking for User ID
			Updates(models.User{
				Balance: sum,
			})

		c.JSON(200, gin.H{
			"message": "Succesfully activated account for this amount of days: ",
			"payload": body.Timer,
		})
	}
}

func UpdateSettings(c *gin.Context) {
	body := &models.Account{}
	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to read update status body",
			"payload": "Please check syntax",
		})

		return
	}

	initializers.DB.Model(&body).
		Where("user_id = ?", fmt.Sprint(GetUserID(c))). // checking for User ID
		Where("login = ?", body.Login).                 // and checking for account login
		Updates(models.Account{
			UpdateKCK:    body.UpdateKCK,
			RunOrki:      body.RunOrki,
			Competitions: body.Competitions,
		})

	c.JSON(200, gin.H{
		"message": "New settings are : ",
		"payload": body,
	})

}
