package helpers

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Aoiewrug/Lowadi-App/lowadi-http-func-server/core/helpers/competitions"
	"github.com/Aoiewrug/Lowadi-App/lowadi-http-func-server/core/helpers/dbupdate"
	"github.com/Aoiewrug/Lowadi-App/lowadi-http-func-server/core/helpers/orki"
	"github.com/Aoiewrug/Lowadi-App/lowadi-http-func-server/core/helpers/updatekck"
	"github.com/Aoiewrug/Lowadi-App/lowadi-http-func-server/core/models"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
)

// Log IN
func GameLogIn(account *models.Account) (statusCode int, chanStruct models.ChanStruct) {
	//handling
	// forming proxy like: 198.165.0.23:5432
	proxyString := account.ProxyIP + ":" + strconv.Itoa(account.ProxyPort)
	fmt.Println(proxyString)

	url := launcher.New().
		Proxy(proxyString).          // set flag "--proxy-server=127.0.0.1:8080"
		Delete("use-mock-keychain"). // delete flag "--use-mock-keychain"
		MustLaunch()

	// Using incognito sessions to isolate it better
	browser := rod.New().ControlURL(url).MustConnect().MustIncognito()

	// defer browser.MustClose()

	// Proxy auth
	//Async version without err handling
	go browser.MustHandleAuth(account.ProxyLogin, account.ProxyPass)()

	// Timeout if there is no such element
	// Did they change elements?
	Page := browser.MustPage(account.GameWebside)

	rejectCookieTimeout := rod.Try(func() {
		Page.Timeout(30*time.Second).MustElementR("button",
			"Я отказываюсь от использования файлов cookie", // - RUSSIAN WORDS!
		)
	})
	if errors.Is(rejectCookieTimeout, context.DeadlineExceeded) {
		fmt.Println("|Reject cookies| button timeout error. Proxy or website layout error")
		return 2, models.ChanStruct{
			Account: account,
			Error:   fmt.Sprintf("User: %v - Can't load the website. Bad proxy or Game server is down", account.UserID),
		}

	} else {
		// Trying to Log-in
		Page.MustElementR("button",
			"Я отказываюсь от использования файлов cookie", // - RUSSIAN WORDS!
		).MustClick()

		time.Sleep(4 * time.Second)

		Page.MustElementR("#header-login-label",
			"Войти", // - RUSSIAN WORDS!
		).MustClick()

		Page.MustElement("#login").MustInput(account.Login).MustType(input.Tab)
		Page.MustElement("#password").MustInput(account.Pass).MustType(input.Enter)
		time.Sleep(2 * time.Second)

		// Checking for successfull Log-in (equ value)
		loginTimeout := rod.Try(func() {
			Page.Timeout(30 * time.Second).MustElement("#header-hud > ul > li.level-1.float-right.hud-equus > a > span > span:nth-child(2) > span").MustWaitLoad()

		})
		if errors.Is(loginTimeout, context.DeadlineExceeded) {
			fmt.Sprintf("User: %v, Account: %v - Failed to Log-in", account.UserID, account.Login)
			//time.Sleep(10 * time.Second)
			return 2, models.ChanStruct{
				Account: account,
				Error:   fmt.Sprintf("User: %v, Account: %v- Failed to Log-in", account.UserID, account.Login),
			}
		} else {
			//
			// Successfully log-in
			fmt.Println("I'm Inside!!")

			// Extract KCK names and links to DB
			// Should be enabled by defauldt after creating an account
			if account.UpdateKCK == 1 {
				// add DB status "Updating KCKs" - 2
				dbupdate.UpdateRunningAccountStatus(account, 2)

				updatekck.UpdateKCK(Page, account)

				// add return values to DB?
			}

			if account.Competitions == 1 {
				// add DB status "Running competitions" - 3
				dbupdate.UpdateRunningAccountStatus(account, 3)

				competitions.StartCompetitionsTest(Page, account)

				// add return values to DB?
			}

			//
			// Add functionality here
			//

			// Run orki
			if account.RunOrki == 1 {
				// add DB status "Running orki" - 4
				dbupdate.UpdateRunningAccountStatus(account, 4)

				horseLinks, err := orki.OrkiPasrePages(Page, account)

				if err != nil {
					fmt.Println("can't get page length : ", err)
					return 2, models.ChanStruct{
						Account: account,
						Error:   fmt.Sprintf("Login: %v - can't transfer orki array", account.Login),
					}
				}

				x := fmt.Sprintf("Login: %v - Transferring Orki array", account.Login)
				fmt.Println(x)

				return 1,
					models.ChanStruct{
						Account:         account,
						Browser:         browser,
						HorseArrayLinks: horseLinks.IDs,
						Error:           "",
					}

				//StartOrkiLinear(Page, account) // Linear (halted)
				//orki.TestStartOrkiAsync(browser, account) // Async test
				//orki.StartOrkiAsync(browser, Page, account) // Async production
				//orki.TestCallFuncServer(browser, Page, account)

				// add return values to DB?
			}

			// add DB status "User can log-in the account"
			dbupdate.UpdateRunningAccountStatus(account, 1)
			//time.Sleep(10000 * time.Second)

			y := fmt.Sprintf("User: %v - I've Logged-in successfully", account.Login)
			fmt.Println(y)

			Page.MustClose()

			return 3, models.ChanStruct{
				Account: account,
				Error:   "",
			}
		}
	}
}
