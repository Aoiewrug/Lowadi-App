package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Aoiewrug/Lowadi-App/lowadi-http-func-server/core/dispatch"
	"github.com/Aoiewrug/Lowadi-App/lowadi-http-func-server/core/handlers"
	"github.com/Aoiewrug/Lowadi-App/lowadi-http-func-server/core/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	fmt.Println(dispatch.MaxParallelism())
	initializers.LoadEnvVar()
	initializers.ConnectToDB()

}

func main() {
	fmt.Println("Starting ReadAccountChannel")
	go dispatch.ReadAccountChannel()
	fmt.Println("Starting ReadJobsChannel")
	go dispatch.ReadJobChannel()

	fmt.Println("Starting Worker pool")
	go dispatch.Worker()

	/*
		go func() {

			for i := 1; i < 1000; i++ {
				dispatch.AccountQueue <- models.ChanStruct{Counter: i}
				//fmt.Println("Here is the counter sended", i)
				time.Sleep(1000 * time.Millisecond)

			}

		}()
	*/

	fmt.Println("Starting the FUNC server")
	r := gin.New()

	funcRoutes := r.Group("/func")
	{
		funcRoutes.POST("/run", handlers.StartRemoteOrki)  // http://127.0.0.1:8002/func/run
		funcRoutes.POST("/test", handlers.StartRemoteOrki) // http://127.0.0.1:8002/func/test

	}

	// run on 8002 port
	if err := r.Run(os.Getenv("FUNC_SERVER_ADDR")); err != nil {
		log.Fatalf("error occured while running the FUNC server: %s", err.Error())
	}

}
