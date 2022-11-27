package main

import (
	"log"
	"os"

	"github.com/Aoiewrug/Lowadi-App/lowadi-http-api/core/controllers/active"
	"github.com/Aoiewrug/Lowadi-App/lowadi-http-api/core/controllers/middleware"
	"github.com/Aoiewrug/Lowadi-App/lowadi-http-api/core/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVar()
	initializers.ConnectToDB()
	initializers.SyncDB()
}

func main() {

	r := gin.New()

	// This webside user functionality
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/singup", active.SingUp)                                 // http://127.0.0.1:8001/auth/singup
		authRoutes.POST("/login", active.Login)                                   // http://127.0.0.1:8001/auth/login
		authRoutes.GET("/validate", middleware.ReqireAuth, active.Validate)       // http://127.0.0.1:8001/auth/validate
		authRoutes.PATCH("/balance", middleware.ReqireAuth, active.UpdateBalance) // http://127.0.0.1:8001/account/update/balance
	}

	// Game accounts functionality
	accountRoutes := r.Group("/account", middleware.ReqireAuth)
	{
		accountRoutes.POST("/add", active.AddAcc)                                // http://127.0.0.1:8001/account/add
		accountRoutes.GET("/get", active.GetAcc)                                 // http://127.0.0.1:8001/account/get
		accountRoutes.DELETE("/delete", middleware.ReqireAuth, active.DeleteAcc) // http://127.0.0.1:8001/account/delete

		accountRoutes.PATCH("/update/password", active.UpdatePass)         // http://127.0.0.1:8001/account/update/password
		accountRoutes.PATCH("/update/status", active.UpdateStatus)         // http://127.0.0.1:8001/account/update/status
		accountRoutes.PATCH("/update/version", active.UpdateVersion)       // http://127.0.0.1:8001/account/update/version
		accountRoutes.PATCH("/update/expiration", active.UpdateExpiration) // http://127.0.0.1:8001/account/update/expiration
		accountRoutes.PATCH("/update/settings", active.UpdateSettings)     // http://127.0.0.1:8001/account/update/settings
		// We can add HARD delete here db.Unscoped().Where("age = 20").Find(&users)
		// We can add here AUTOMATIC GAME ACCOUNT billing function
	}

	// Game functionality
	kckRoutes := r.Group("/kck", middleware.ReqireAuth)
	{
		// Post is done automatically UpdateKCK function
		kckRoutes.GET("/gets", active.GetKCKs)                 // http://127.0.0.1:8001/kck/gets
		kckRoutes.GET("/get", active.GetKCK)                   // http://127.0.0.1:8001/kck/get
		kckRoutes.PATCH("/update", active.SetKCKtoGameAccount) // http://127.0.0.1:8001/kck/update
		kckRoutes.DELETE("/delete", active.DeleteKCK)          // http://127.0.0.1:8001/kck/delete

	}

	// Game functionality
	competitionRoutes := r.Group("/copmetition", middleware.ReqireAuth)
	{
		// Post is done automatically UpdateKCK function
		competitionRoutes.GET("/get", active.GetCompetition)                   // http://127.0.0.1:8001/copmetition/get
		competitionRoutes.PATCH("/update", active.SetCompetitionToGameAccount) // http://127.0.0.1:8001/copmetition/update

	}

	// Game functionality
	gameRoutes := r.Group("/game", middleware.ReqireAuth)
	{
		gameRoutes.GET("/start", active.GameEnterPoint) // http://127.0.0.1:8001/game/start

	}

	// run on 8001 port
	if err := r.Run(os.Getenv("HTTP_SERVER_ADDR")); err != nil {
		log.Fatalf("error occured while running the HTTP server: %s", err.Error())
	}

}
