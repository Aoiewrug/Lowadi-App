package updatekck

import (
	"fmt"
	"strings"

	"github.com/Aoiewrug/Lowadi-App/lowadi-http-api/core/initializers"
	"github.com/Aoiewrug/Lowadi-App/lowadi-http-api/core/models"
	"github.com/go-rod/rod"
)

// Grabs KCK names and links
func UpdateKCK(page *rod.Page, account *models.Account) {
	var array []models.KCK

	page.MustNavigate(account.GameWebside + "elevage/chevaux/?elevage=all-horses")
	section := page.MustElement("#tab-all-breeding")
	kckNames := section.MustElements("li")
	// In loop all KCK names and links extraction
	for _, stable := range kckNames {

		link := fmt.Sprintf("%s", stable.MustProperty("id"))
		link = strings.ReplaceAll(link, "tab-select-", "")

		local := models.KCK{
			UserID:     account.UserID,
			Login:      account.Login,
			StableLink: link,
			StableName: stable.MustText(),
		}

		array = append(array, local)

	}

	//fmt.Println(Array)

	// Batch add records
	result := initializers.DB.Create(&array)
	if result.Error != nil {
		fmt.Println("Can't add KCK list to db")
		return
	}

	fmt.Println("KCK have been updated")

	/*
		The output should be like this:
		2022/10/26 00:42:50 kck: Чебурашки (tab-select-1478452): 'Чебурашки'
		2022/10/26 00:42:50 kck: Остальные лошади (tab-select-all-horses): 'Остальные лошади'
		2022/10/26 00:42:50 kck: + (): '+'
	*/

	// Turn off KCK update
	initializers.DB.Model(&account).
		Where("user_id = ?", account.UserID). // checking for User ID
		Where("login = ?", account.Login).    // and checking for account login
		Updates(models.Account{UpdateKCK: 2})

	fmt.Println("UpdateKCK have been set to 2 (won't update in next run)")

}
