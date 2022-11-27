package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Aoiewrug/Lowadi-App/lowadi-http-func-server/core/handlers"
	"github.com/Aoiewrug/Lowadi-App/lowadi-http-func-server/core/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVar()
	initializers.ConnectToDB()

}

func main() {

	fmt.Println("Starting the FUNC server")

	r := gin.New()

	funcRoutes := r.Group("/func")
	{
		funcRoutes.POST("/run", handlers.StartRemoteOrki) // http://127.0.0.1:8001/auth/singup

	}

	// run on 8002 port
	if err := r.Run(os.Getenv("FUNC_SERVER_ADDR")); err != nil {
		log.Fatalf("error occured while running the FUNC server: %s", err.Error())
	}

}
