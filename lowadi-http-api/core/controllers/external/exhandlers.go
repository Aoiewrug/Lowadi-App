package external

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var myClient = &http.Client{Timeout: 5 * time.Second}

// Calls POST method for my workers queue server :8002
// NO cookie
func PostJson(url string, requestBody []byte) (res string, err error) {

	// Test incoming json
	//fmt.Println(bytes.NewBuffer(requestBody))

	r, err := myClient.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error with sending POST request")
		return "", err
	}
	defer r.Body.Close()
	//r.Header.Set("Cookie", "name=xxxx; count=x")

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Can't read response")
		return "", err
	}

	//fmt.Println(string(content))

	return string(content), nil

}
