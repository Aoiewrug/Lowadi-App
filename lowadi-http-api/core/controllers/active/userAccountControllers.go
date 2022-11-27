package active

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Aoiewrug/Lowadi-App/lowadi-http-api/core/initializers"
	"github.com/Aoiewrug/Lowadi-App/lowadi-http-api/core/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var Body struct {
	Email    string
	Password string
}

// Shhould be turned off to reg manually
func SingUp(c *gin.Context) {
	// Get the email/pass of request body
	body := Body
	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	// Hash the pass
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})

		return
	}

	// Create the user
	user := models.User{Email: body.Email, Password: string(hash)}

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})

		return
	}

	// Respond OK
	c.JSON(http.StatusOK, gin.H{})

}

func Login(c *gin.Context) {
	// Get email and pass of request body
	var body struct {
		Email    string
		Password string
	}

	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}
	// Lookup requested user
	var user models.User

	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			// Rename to "Invalid user or password"
			"error": "Invalid user or password -1",
		})

		return
	}

	// Compare sent in pass with saved user pass hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			// Rename to "Invalid user or password"
			"error": "Invalid user or password -2",
		})

		return

	}

	// Generate jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// sub = subject, exp = expiration
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SALT")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			// Rename to "Invalid user or password"
			"error": "Unable to generate token",
		})

		return

	}

	// FOR TEST PURPOSES WITH POSTMAN
	// Set coockie locally

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{})
	/*
		// FOR CLIENT SIDE
		// Send token to user
		c.JSON(http.StatusOK, gin.H{
			// Rename to "Invalid user or password"
			"token": tokenString,
		})
	*/
}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, user)
}

func GetUserID(c *gin.Context) (userId uint) {
	// Here we solved empty interface case, woohoo!
	user, _ := c.Get("user")
	x := user.(models.User)
	//fmt.Println(x.ID)
	return x.ID
}

// Add Balance
// Should be used by admin only !
func UpdateBalance(c *gin.Context) {

	if GetUserID(c) == 1 {

		user := &models.User{}

		if c.BindJSON(&user) != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Failed to read user1 body",
				"payload": "Please check syntax",
			})

			return
		}

		// Buffering the desired balance to be added
		balanceBuffer := user.Balance

		// Getting existing balance
		admin := initializers.DB.Where("id = ?", fmt.Sprint(GetUserID(c))).First(&user)

		if admin.Error != nil {
			fmt.Println("Can't find such E-mail")
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Can't find such User",
				"payload": "-//-",
			})
			return
		}

		// ========================================================

		initializers.DB.Model(&user).
			Where("email = ?", user.Email). // checking for User ID
			Updates(models.User{Balance: (balanceBuffer + user.Balance)})

		c.JSON(200, gin.H{
			"message": "Current balance is",
			"payload": user.Balance,
		})

	} else {
		fmt.Println("Sorry, you are not admin")
		c.JSON(200, gin.H{
			"message": "Sorry, you are not admin",
			"payload": "shiet :(",
		})
	}

}
