package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var myClient = &http.Client{Timeout: 30 * time.Second}

// Calls GET method for my api server :8000
// NO cookie
func GetJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		fmt.Println("Error with sending GET request")
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

// Calls POST method for my api server :8000
// NO cookie (used for sending 1st post to authenticate yourself)
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

// With cookie, no timiout. Works with non array structs
func GetRequest(url string, target interface{}, r *http.Request) {
	//fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Can't get request 1 ", err)
	}

	cookie, err := r.Cookie("Authorization")
	if err != nil {
		fmt.Println("Bad cookie, Sir ", err)
	}

	token := &http.Cookie{
		Name:  cookie.Name,
		Value: cookie.Value,
	}
	req.AddCookie(token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Can't get request 2 ", err)
	}
	defer resp.Body.Close()

	//fmt.Printf("Status code is: %d\n", resp.StatusCode)

	b, err := io.ReadAll(resp.Body)
	// b, err := ioutil.ReadAll(resp.Body)  Go.1.15 and earlier
	if err != nil {
		log.Fatalln(err)
	}

	//fmt.Println(string(b))
	er := json.Unmarshal([]byte(b), target)
	if er != nil {
		fmt.Println("Can't Unmarshal GET response ", err)
		//fmt.Println(resp.Body)
		//return json.NewDecoder(resp.Body).Decode(target)
	}

}

// Post with cookie and timeout
func PostRequest(url string, requestBody []byte, r *http.Request) (res string, err error) {

	cookie, err := r.Cookie("Authorization")
	if err != nil {
		fmt.Println("Bad cookie, Sir ", err)
	}

	cook := &http.Cookie{
		Name:  cookie.Name,
		Value: cookie.Value,
	}

	//fmt.Println(fmt.Sprintf("COOCKKIIEE IS %s", cook))
	fmt.Println(bytes.NewBuffer(requestBody))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Can't do POST request")
		return "", err
	}
	//req.Header.Add("Content-Type", "application/json")
	//req.Header.Add("Cookie", fmt.Sprintf("name=%s; count=%s", name, value))
	req.AddCookie(cook)

	resp, err := myClient.Do(req)
	if err != nil {
		fmt.Println("Can't get response")
		return "", err
	}

	//cock.Header.Set("Cookie", fmt.Sprintf("name=%s; count=%s",))
	//r.Header.Set("Cookie", fmt.Sprintf("name=%s; count=%s", name, value))

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Can't read response")
		return "", err
	}

	//fmt.Println(string(content))

	return string(content), nil
}
