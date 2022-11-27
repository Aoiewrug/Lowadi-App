package active

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/Aoiewrug/Lowadi-app/lowadi-http-client/core/controllers/helpers"
	"github.com/Aoiewrug/Lowadi-app/lowadi-http-client/core/controllers/models"
	"github.com/Aoiewrug/Lowadi-app/lowadi-http-client/core/controllers/render"
)

// Test func
func Handler(w http.ResponseWriter, r *http.Request) {
	// Leggo to API exterenal server
	test := &models.TestStruct{}
	helpers.GetJson(os.Getenv("SERVER_ADDR")+"auth/test2", test)

	println(test.Message)
	println(test.Text)

	render.RenderTemplate(w, "1-test.html", test)

}

// Test
func Test(w http.ResponseWriter, r *http.Request) {
	render.RenderMultipleTemplates(w, "test", nil)
}

// Test
func Test2(w http.ResponseWriter, r *http.Request) {
	var localBuffer []string
	url := os.Getenv("SERVER_ADDR") + "info/compare"

	// 1st form value
	str1 := []string{r.FormValue("package1")}
	localBuffer = append(localBuffer, str1...)

	// 2nd form value
	str2 := r.FormValue("package2")
	str2 = strings.ReplaceAll(str2, " ", "") //trim spaces
	str2 = strings.ReplaceAll(str2, ",", "")

	re := regexp.MustCompile(`\r\n`)
	array := re.ReplaceAllString(str2, " ")

	array2 := strings.Split(array, " ")
	localBuffer = append(localBuffer, array2...)

	requestBody, err := json.Marshal(&models.PackageComp{
		Link: localBuffer,
	})

	if err != nil {
		fmt.Println("Can't Marshal -3", err)
	}
	/*	// Input example
		body := map[string]interface{}{
			"link": []string{
				"http://api.privateproxy.me:3000/admin/user_packages/43043",
				"http://api.privateproxy.me:3000/admin/user_packages/43042",
			},
		}
	*/

	response, err := helpers.PostRequest(url, requestBody, r)
	if err != nil {
		fmt.Println("Can't read POST response", err)
	}

	//fmt.Println(response)

	result := strings.ReplaceAll(response, `"`, "")
	result = strings.ReplaceAll(result, "[", "")
	result = strings.ReplaceAll(result, "]", "")
	result2 := strings.Split(result, ",")

	render.RenderMultipleTemplates(w, "test", result2)
}

// GET Log in HTML page
func GetLogin(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "2-login.html", nil)

}

// POST Send request to server api to get cookie
func PostLogin(w http.ResponseWriter, r *http.Request) {
	// Past form values
	requestBody, err := json.Marshal(&models.User{Email: r.FormValue("email"), Password: r.FormValue("password")})
	if err != nil {
		fmt.Println("Can't marshal JSON")
	}

	tokenString, _ := helpers.PostJson(os.Getenv("SERVER_ADDR")+"auth/login", requestBody)
	fmt.Println(tokenString)
	var cookie models.Cookie
	json.Unmarshal([]byte(tokenString), &cookie)

	// Drop if there is no token from the API
	if cookie.Token == "" {
		w.WriteHeader(404)
		w.Write([]byte("Wrong cookie"))
	} else {
		cookieToSet := &http.Cookie{
			Name:     "Authorization",
			Value:    cookie.Token,
			Path:     "/",
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
		}
		http.SetCookie(w, cookieToSet)

		//w.WriteHeader(200)
		//w.Write([]byte("Cookie has been set successfully"))

		// We can add redirect here or a forward button
		render.RenderMultipleTemplates(w, "mainPage", nil)

	}
}

// Broken a lil. Can't properly unmarshal models.ValidateTest
func Validate(w http.ResponseWriter, r *http.Request) {
	message := &models.User{}
	url := os.Getenv("SERVER_ADDR") + "auth/validate"
	helpers.GetRequest(url, message, r)
	//fmt.Println("response is: ", message)
	render.RenderMultipleTemplates(w, "validate", message)
}

// Show subnets from MySQL DB
func GETSubnetList(w http.ResponseWriter, r *http.Request) {
	subnet := &[]models.Subnet{}
	url := os.Getenv("SERVER_ADDR") + "info/subnets"
	helpers.GetRequest(url, subnet, r)
	//fmt.Println(subnet)
	render.RenderMultipleTemplates(w, "showSubnet", subnet)

}

// Show purposes from MySQL DB
func GETPurposeList(w http.ResponseWriter, r *http.Request) {
	purposes := &[]models.Purpose{}
	url := os.Getenv("SERVER_ADDR") + "info/purposes"
	helpers.GetRequest(url, purposes, r)
	//fmt.Println(subnet)
	render.RenderMultipleTemplates(w, "showPurpose", purposes)

}

// GET Log in HTML page
func GETComparePackagesHandler(w http.ResponseWriter, r *http.Request) {
	render.RenderMultipleTemplates(w, "packageCompare", nil)

}

func POSTComparePackagesHandler(w http.ResponseWriter, r *http.Request) {
	var localBuffer []string
	url := os.Getenv("SERVER_ADDR") + "info/compare"

	// 1st form value
	str1 := []string{r.FormValue("package1")}
	localBuffer = append(localBuffer, str1...)

	// 2nd form value
	str2 := r.FormValue("package2")
	str2 = strings.ReplaceAll(str2, " ", `\r\n`) //trim spaces
	str2 = strings.ReplaceAll(str2, ",", "")

	re := regexp.MustCompile(`\r\n`) // rmv new lines (win version)
	array := re.ReplaceAllString(str2, " ")

	array2 := strings.Split(array, " ")
	localBuffer = append(localBuffer, array2...)

	requestBody, err := json.Marshal(&models.PackageComp{
		Link: localBuffer,
	})

	if err != nil {
		fmt.Println("Can't Marshal -3", err)
	}
	/*	// Input example
		body := map[string]interface{}{
			"link": []string{
				"http://api.privateproxy.me:3000/admin/user_packages/43043",
				"http://api.privateproxy.me:3000/admin/user_packages/43042",
			},
		}
	*/

	response, err := helpers.PostRequest(url, requestBody, r)
	if err != nil {
		fmt.Println("Can't read POST response", err)
	}

	//fmt.Println(response)

	result := strings.ReplaceAll(response, `"`, "")
	result = strings.ReplaceAll(result, "[", "")
	result = strings.ReplaceAll(result, "]", "")
	result2 := strings.Split(result, ",")

	render.RenderMultipleTemplates(w, "packageCompare", result2)

}

// GET Log in HTML page
func GETFindUniqueIPsHandler(w http.ResponseWriter, r *http.Request) {
	render.RenderMultipleTemplates(w, "packageUnique", nil)

}

func POSTFindUniqueIPsHandler(w http.ResponseWriter, r *http.Request) {
	var localBuffer []string
	url := os.Getenv("SERVER_ADDR") + "info/unique"

	// Form value
	str2 := r.FormValue("packages")
	str2 = strings.ReplaceAll(str2, " ", `\r\n`) //trim spaces
	str2 = strings.ReplaceAll(str2, ",", "")

	re := regexp.MustCompile(`\r\n`) // rmv new lines (win version)
	array := re.ReplaceAllString(str2, " ")

	array2 := strings.Split(array, " ")
	localBuffer = append(localBuffer, array2...)

	requestBody, err := json.Marshal(&models.PackageComp{
		Link: localBuffer,
	})

	if err != nil {
		fmt.Println("Can't Marshal -3", err)
	}
	/*	// Input example
		body := map[string]interface{}{
			"link": []string{
				"http://api.privateproxy.me:3000/admin/user_packages/43043",
				"http://api.privateproxy.me:3000/admin/user_packages/43042",
			},
		}
	*/

	response, err := helpers.PostRequest(url, requestBody, r)
	if err != nil {
		fmt.Println("Can't read POST response", err)
	}

	//fmt.Println(response)

	result := strings.ReplaceAll(response, `"`, "")
	result = strings.ReplaceAll(result, "[", "")
	result = strings.ReplaceAll(result, "]", "")
	result2 := strings.Split(result, ",")

	render.RenderMultipleTemplates(w, "packageUnique", result2)

}

// GET Log in HTML page
func GetCryptoPaymentHandler(w http.ResponseWriter, r *http.Request) {
	payments := &[]models.CryptoPayment{}
	url := os.Getenv("SERVER_ADDR") + "actions/crypto"
	helpers.GetRequest(url, payments, r)

	render.RenderMultipleTemplates(w, "cryptoPayments", payments)

}

func AddNewCryptoPaymentHandler(w http.ResponseWriter, r *http.Request) {
	// Add logic
}
