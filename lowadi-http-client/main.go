package main

import (
	"net/http"
	"os"

	//C:\Users\3\Documents\GitHub\Lowadi-app\lowadi-http-client\main.go
	"github.com/Aoiewrug/Lowadi-app/lowadi-http-client/core/controllers/active"
	"github.com/Aoiewrug/Lowadi-app/lowadi-http-client/core/initializers"
)

func init() {
	initializers.LoadEnvVar()
}

func main() {
	handleFunc()
}

func handleFunc() {
	// Test func
	//http.HandleFunc("/", active.Handler)
	//http.HandleFunc("/auth/test", active.Test)     // GET TEST page
	//http.HandleFunc("/info/compare", active.Test2) // post TEST page

	// Turned off to registration for users. Only manually from the backend side
	//http.HandleFunc("/auth/singup", active.SingUp) 		// POST req to http://127.0.0.1:8000/auth/singup

	http.HandleFunc("/", active.GetLogin) // GET render page FOR NGROK
	//http.HandleFunc("/auth/enter", active.GetLogin)    // GET render page
	http.HandleFunc("/auth/login", active.PostLogin)   // POST http://127.0.0.1:8000/auth/login
	http.HandleFunc("/auth/validate", active.Validate) // GET http://127.0.0.1:8000/auth/validate

	http.HandleFunc("/info/subnets", active.GETSubnetList)              // GET  http://127.0.0.1:8000/info/subnets
	http.HandleFunc("/info/purposes", active.GETPurposeList)            // GET  http://127.0.0.1:8000/info/purposes
	http.HandleFunc("/info/dupes", active.GETComparePackagesHandler)    // GET render page
	http.HandleFunc("/info/compare", active.POSTComparePackagesHandler) // POST http://127.0.0.1:8000/info/compare
	http.HandleFunc("/info/uniq", active.GETFindUniqueIPsHandler)       // GET render page
	http.HandleFunc("/info/unique", active.POSTFindUniqueIPsHandler)    // POST http://127.0.0.1:8000/info/unique

	http.HandleFunc("/actions/crypto", active.GetCryptoPaymentHandler)       // http://127.0.0.1:8000/actions/crypto
	http.HandleFunc("/actions/cryptoadd", active.AddNewCryptoPaymentHandler) // http://127.0.0.1:8000/actions/cryptopay
	// add get purpose IPs and Users

	// add purposes

	// Add certificates for https later
	/*
		err := http.ListenAndServeTLS(os.Getenv("CLIENT_PORT_HTTPS"), "server.crt", "server.key", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	*/

	http.ListenAndServe(os.Getenv("CLIENT_PORT_TEST"), nil)
}
