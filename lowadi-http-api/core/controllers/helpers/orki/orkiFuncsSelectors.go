package orki

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Aoiewrug/Lowadi-App/lowadi-http-api/core/models"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
)

// Get int val of pages with horses links ADD ERROR SECTION
func OrkiPasrePages(page *rod.Page, account *models.Account) (horseLinks models.OrkiHorses, err error) {
	var links models.OrkiHorses

	// Open choosen stable
	BTN := rod.Try(func() {
		page.Timeout(10000 * time.Millisecond).MustNavigate(account.GameWebside + "elevage/chevaux/?elevage=" + account.StableLink).MustWaitLoad()
	})
	if errors.Is(BTN, context.DeadlineExceeded) {
		return links, errors.New("can't open correct KCK page to parse page length")
	}

	// Check how may pages does KCK have
	pageSection := page.MustElement("#horseList > div > div.module-before > div:nth-child(2) > div > ul").MustWaitLoad()
	pages := pageSection.MustElements("li")
	fmt.Println("How may mages in this KCK?: ", len(pages))

	links, err = OrkiPasreHorsesLinks(page, len(pages))
	if err != nil {
		return links, err
	}

	return links, nil

	// #horseList > div > div.module-before > div:nth-child(2) > div > ul > li.page.selected - 1st page (li:nth-child(2))
	// #horseList > div > div.module-before > div:nth-child(2) > div > ul > li:nth-child(3) - 2nd page
	// #horseList > div > div.module-before > div:nth-child(2) > div > ul > li:nth-child(4) - 3nd page
	// #horseList > div > div.module-before > div:nth-child(2) > div > ul > li:nth-child(5) - 4nd page
	// ... 5 actual page = li:nth-child(6)

}

// Get all horses links ADD ERROR SECTION
func OrkiPasreHorsesLinks(page *rod.Page, lenPages int) (horseLinks models.OrkiHorses, err error) {
	var array models.OrkiHorses
	// Iterating over all pages. Starting from the end
	for i := lenPages; i > 1; i-- {
		time.Sleep(1000 * time.Millisecond)

		fmt.Println("Current page: ", i)

		openThisPage := fmt.Sprintf("#horseList > div > div.module-before > div:nth-child(2) > div > ul > li:nth-child(%v)", i)
		BTN := rod.Try(func() {
			page.Timeout(10000 * time.Millisecond).MustElement(openThisPage).MustClick().MustWaitLoad()
		})
		if errors.Is(BTN, context.DeadlineExceeded) {
			return array, errors.New("can't open correct page to parse horses")
		}

		// time.Sleep(1000 * time.Millisecond)

		// Parse horse IDs from current page:
		horseGrid := page.MustElement("#horseList > div > div.damier-table.grid-table.width-100")
		horses := horseGrid.MustElements("li.damier-horsename.nowrap")
		fmt.Println("How many horses on this page?: ", len(horses))

		// Loop over horses on this page
		for _, row := range horses {
			horses := row.MustElement("a")

			// Extract horse's direct link
			link := fmt.Sprintf("%s", horses.MustProperty("href"))
			array.IDs = append(array.IDs, link)

		}
	}

	//fmt.Println(array)

	return array, nil

}

// Unused for now
func NextPage(page *rod.Page) {
	//defer handleOutOfBounds()
	// Next page
	page = page.MustWaitLoad()
	page.Timeout(CancelationTimeout).MustElement("#nav-next").MustClick().MustWaitLoad()
	page.MustWaitLoad()
	time.Sleep(Timeout)
}

func RunBabyHorse(page *rod.Page) (s string, err error) {
	//defer handleOutOfBounds()
	// maybe we need to use page.Eval with JS
	page = page.MustWaitLoad()

	// Feed
	BTN := rod.Try(func() {
		page.Timeout(CancelationTimeout).MustElement("#care-tab-main > div > div.grid-row.row-0.even > div.grid-cell.even.first.top.width-33").MustClick().MustWaitLoad()

	})
	if errors.Is(BTN, context.DeadlineExceeded) {
		return "", errors.New("can't find feed bb horse element (soska)")
	}

	time.Sleep(Timeout)

	// Clear
	PressClearButton(page)

	// Sleep
	PressSleepButton(page)

	return "done with this bb horse", nil
}

func RunOldHorseWithKCK(page *rod.Page, chanStruct chanStruct, currentMoney int) (s string, err error) {
	page = page.MustWaitLoad()
	time.Sleep(Timeout)

	// Mission availability check FIRST!
	isMissionActive := rod.Try(func() {
		page.Timeout(HorseDetectionTimeout).MustElement("#mission-tab-0 > div > div > div:nth-child(2)")

	})
	if errors.Is(isMissionActive, context.DeadlineExceeded) {
		fmt.Println("No mission for this horse")
		// Parse optimum KORM value

	} else {
		// Start Mission
		BTN := rod.Try(func() {
			page.Timeout(CancelationTimeout).MustElement("#mission-tab-0 > div > div > div:nth-child(2)").MustClick().MustWaitLoad()
		})
		if errors.Is(BTN, context.DeadlineExceeded) {
			return "", errors.New("can't press mission button")
		}
		time.Sleep(Timeout)
	}

	// Propagate horse
	if chanStruct.account.BirthHorses == 1 {
		err := Propagate(page, chanStruct, currentMoney)
		if err != nil {
			return "", err
		}
	}

	time.Sleep(HorseDetectionTimeout)

	// Open feed menu
	BTNopenfeed := rod.Try(func() {
		page.Timeout(10000 * time.Millisecond).MustElement("#boutonNourrir > span").MustClick().MustWaitLoad()
	})
	if errors.Is(BTNopenfeed, context.DeadlineExceeded) {
		return "", errors.New("can't open feed menu")
	}

	time.Sleep(Timeout)

	// Check if second feed bar with OVEC exists
	isSecondFeedBar := rod.Try(func() {
		page.Timeout(Timeout).MustElement("#oatsSlider")

	})
	if errors.Is(isSecondFeedBar, context.DeadlineExceeded) {
		// one feed bar horses
		err := FeedOneBarHorse(page, chanStruct.account)

		time.Sleep(Timeout)

		// Clear
		PressClearButton(page)

		// Sleep
		PressSleepButton(page)
		if err != nil {
			return "", err
		}
	} else {
		// two feed bars horses
		er := FeedTwoBarHorse(page, chanStruct.account)

		time.Sleep(Timeout)

		// Clear
		PressClearButton(page)

		// Sleep
		PressSleepButton(page)
		if er != nil {
			return "", er
		}
	}

	/*
		ADD THIS ON TOP OF EACH HORSE AFTER OPEN HORSE PAGE!
			#boutonVeterinaire > span > span > span
		This is the Veterinaire button to born new babyhorse!

		page.Timeout(Timeout).MustElement("#poulain-1").MustClick() - Select horse name field
		page.MustElement("#admin_user_password").MustInput(body.BirthHorsesName) - Enter new horse name
		page.Timeout(Timeout).MustElement("#boutonChoisirNom > span > span > span").MustClick() - Confirm horse name, open new page

		=========================================================

		Selector for right bottom button:
			page.MustElementR("#reproduction-tab-0 > table > tbody > tr > td.last",
				"xxxx", // - RUSSIAN WORDS!
			).MustClick()



		male:
			page.MustElementR("#reproduction-tab-0 > table > tbody > tr > td.last",
				"Покрыть кобылу", // - RUSSIAN WORDS!
			).MustClick()

			page.Timeout(Timeout).MustElement("#formMalePublicTypePublic").MustClick() - click "Open world" filter
			page.Timeout(Timeout).MustElement("#formMalePublicPrice > option:nth-child(2)").MustClick() - click "500 equ" prive
			page.Timeout(Timeout).MustElement("#boutonMaleReproduction > span > span > span").MustClick() - click "Confirm" button
			time.Sleep(Timeout)

		=========================================================



		female: (start bith)
			try:
			page.MustElementR("#reproduction-tab-0 > table > tbody > tr > td.last",
					"Случить кобылу", // - RUSSIAN WORDS!
				).MustClick()

			!!! CHECK IF current moneyInt > (500 + 400) * thread number + 1 !!!
			!!! ONly 300 male horses per day. Set Up counter!!!

			--- Wait new page opens --- (select boyfriend)

			page.Timeout(Timeout).MustElement("#table-0 > tbody > tr:nth-child(1) > td.align-center.action > a").MustClick() - click "first male horse" button

			--- Wait new page opens --- (pay boyfriend)

			page.Timeout(Timeout).MustElement("#boutonDoReproduction > span > span > span").MustClick() - click "Confirm meeting" button

			time.Sleep(Timeout)

			GG

		female: (echo) (we can pass if we see this element)
			page.Timeout(Timeout).MustElement("#boutonEchographie").MustClick()

			try{
				page.Timeout(Timeout).MustElement("#boutonEchographie").MustClick()
				!!! Here you will see new horse specs. Pass this step now
			}
	*/

	return "done with this old horse horse", nil

}

func Propagate(page *rod.Page, chanStruct chanStruct, currentMoney int) (err error) {
	page = page.MustWaitLoad()

	// Try male selector 1st
	BTN1 := rod.Try(func() {
		page.Timeout(HorseDetectionTimeout).MustElementR("#reproduction-tab-0 > table > tbody > tr > td.last",
			"Покрыть кобылу", // - RUSSIAN WORDS!
		).MustClick()
	})
	if errors.Is(BTN1, context.DeadlineExceeded) {
		// Try female selector than
		fmt.Println("no male propagate selector")
		BTN2 := rod.Try(func() {
			page.Timeout(HorseDetectionTimeout).MustElementR("#reproduction-tab-0 > table > tbody > tr > td.last",
				"Случить кобылу", // - RUSSIAN WORDS!
			)
		})
		if errors.Is(BTN2, context.DeadlineExceeded) {
			fmt.Println("no both propagate selectors")
			return nil
		} else {
			// Proceed with female babymaking :D
			BTN1 := rod.Try(func() {
				page.Timeout(CancelationTimeout).MustElementR("#reproduction-tab-0 > table > tbody > tr > td.last",
					"Случить кобылу", // - RUSSIAN WORDS!
				).MustClick()
			})
			if errors.Is(BTN1, context.DeadlineExceeded) {
				return errors.New("error can't open propagation female section")
			}

			//!!! CHECK IF current moneyInt > (500 + 400) * thread number + 1 !!!
			if currentMoney > 40000 {

				// Filter with the lowest price
				BTN1 := rod.Try(func() {
					page.Timeout(CancelationTimeout).MustElement("#table-0 > thead > tr > td:nth-child(7) > a").MustClick() // click "filter with the lowest price" button
				})
				if errors.Is(BTN1, context.DeadlineExceeded) {
					return errors.New("error can't click lowest price for propagation")
				}

				page.MustWaitLoad()

				BTN2 := rod.Try(func() {
					page.Timeout(CancelationTimeout).MustElement("#table-0 > tbody > tr:nth-child(1) > td.align-center.action > a").MustClick() // click "first male horse" button
				})
				if errors.Is(BTN2, context.DeadlineExceeded) {
					return errors.New("error can't select first male horse for propagation")
				}

				page.MustWaitLoad()

				//--- Wait pay boyfriend ---

				BTN3 := rod.Try(func() {
					page.Timeout(CancelationTimeout).MustElement("#boutonDoReproduction > span > span > span").MustClick() // click "Confirm meeting" button
				})
				if errors.Is(BTN3, context.DeadlineExceeded) {
					return errors.New("error can't confirm meeting horses for propagation")
				}

				time.Sleep(Timeout)

				page.MustWaitLoad()

				fmt.Println("successfull propagation female")

				return nil

			} else {
				fmt.Println("no money for female propagation")
				return nil

			}
		}

	} else {
		// Proceed with male babymaking :D
		//!!! ONLY 300 male horses per day. Set Up counter!!!
		//if chanStruct.counter < 301 {
		BTN1 := rod.Try(func() {
			page.Timeout(CancelationTimeout).MustElement("#formMalePublicTypePublic").MustClick() // click "Open world" filter
		})
		if errors.Is(BTN1, context.DeadlineExceeded) {
			fmt.Println("can't open propagation male section") // This could be okay cuz there are only 300 operations per day
			return nil
		}

		BTN0 := rod.Try(func() {
			page.Timeout(CancelationTimeout).MustElement("#formMalePublicPrice").MustClick() // click "select price" button
		})
		if errors.Is(BTN0, context.DeadlineExceeded) {
			return errors.New("error can't click price for propagation male")
		}

		// Select 500 equ price
		page.KeyActions().Press(input.ArrowDown).Type(input.Enter).MustDo()

		/*
			BTN2 := rod.Try(func() {
				page.Timeout(CancelationTimeout).MustElement("#formMalePublicPrice > option:nth-child(2)").MustClick() // click "500 equ" price
			})
			if errors.Is(BTN2, context.DeadlineExceeded) {
				return errors.New("error can't select price for propagation male")
			}
		*/
		BTN3 := rod.Try(func() {
			page.Timeout(CancelationTimeout).MustElement("#boutonMaleReproduction > span > span > span").MustClick().MustWaitLoad() // click "Confirm" button
		})
		if errors.Is(BTN3, context.DeadlineExceeded) {
			return errors.New("error can't select price for propagation male")
		}

		time.Sleep(Timeout)

		fmt.Println("successfull propagation male")

		return nil
		//}

		// fmt.Println("done with male propagations today")

		// return nil

	}

}

func NewbornHorse(page *rod.Page, account *models.Account) (err error) {
	BTN1 := rod.Try(func() {
		page.Timeout(CancelationTimeout).MustElement("#poulain-1").MustClick() // Select horse name field
	})
	if errors.Is(BTN1, context.DeadlineExceeded) {
		return errors.New("can't select horse name field")
	}

	BTN2 := rod.Try(func() {
		page.Timeout(CancelationTimeout).MustElement("#admin_user_password").MustInput(account.BirthHorsesName) // Enter new horse name
	})
	if errors.Is(BTN2, context.DeadlineExceeded) {
		return errors.New("can't enter newborn horse name")
	}

	BTN3 := rod.Try(func() {
		page.Timeout(CancelationTimeout).MustElement("#boutonChoisirNom > span > span > span").MustClick() // Confirm horse name, open new page
	})
	if errors.Is(BTN3, context.DeadlineExceeded) {
		return errors.New("can't confirm horse name, open new page")
	}
	//time.Sleep(Timeout * 2)
	page.MustWaitLoad()
	time.Sleep(Timeout)
	return nil
}

// If we succeded here we have (5) instead of (6) child selector for korm and oves selectors
// this means we need to return it!
func KCKsingin(page *rod.Page, account *models.Account) (pg *rod.Page, err error) {
	page = page.MustWaitLoad()
	//time.Sleep(CancelationTimeout)

	// Press "follow kck sing-in button"
	BTN := rod.Try(func() {
		page.Timeout(10000 * time.Millisecond).MustElement("#cheval-inscription > a > span.img").MustClick().MustWaitLoad()
		// or #cheval-inscription
		// or #cheval-inscription > a > span.img

	})
	if errors.Is(BTN, context.DeadlineExceeded) {
		return page, errors.New("can't find KCK sing-in button")
	}

	time.Sleep(Timeout * 2)

	/*
		Advanteges:
		#fourrageCheckbox - furaj
		#avoineCheckbox - ovec
		#carotteCheckbox - carrot
	*/

	// Check Furage box
	if account.AdvantagesFuraj == 1 {
		BTN := rod.Try(func() {
			page.MustWaitLoad().Timeout(CancelationTimeout).MustElement("#fourrageCheckbox").MustClick().MustWaitLoad()
		})
		if errors.Is(BTN, context.DeadlineExceeded) {
			return page, errors.New("can't find AdvantagesFuraj button")
		}
		time.Sleep(Timeout)
	}

	// Check Ovec box
	if account.AdvantagesOvec == 1 {
		BTN := rod.Try(func() {
			page.MustWaitLoad().Timeout(CancelationTimeout).MustElement("#avoineCheckbox").MustClick().MustWaitLoad()
		})
		if errors.Is(BTN, context.DeadlineExceeded) {
			return page, errors.New("can't find AdvantagesOvec button")
		}
		time.Sleep(Timeout)
	}

	// Check Carrot box
	if account.AdvantagesCarrot == 1 {
		BTN := rod.Try(func() {
			page.MustWaitLoad().Timeout(CancelationTimeout).MustElement("#carotteCheckbox").MustClick().MustWaitLoad()
		})
		if errors.Is(BTN, context.DeadlineExceeded) {
			return page, errors.New("can't find AdvantagesCarrot button")
		}
		time.Sleep(Timeout)
	}

	//time.Sleep(Timeout)

	// Press search button
	BTNtimeout := rod.Try(func() {
		page.MustWaitLoad().Timeout(CancelationTimeout).MustElement("#cheval-centre-inscription > button:nth-child(6) > span").MustClick().MustWaitLoad()
	})
	if errors.Is(BTNtimeout, context.DeadlineExceeded) {
		return page, errors.New("can't find search button")
	}

	time.Sleep(Timeout * 2)
	//page.MustWaitIdle()

	// Press filter 3 days
	BTNfilter := rod.Try(func() {
		page.MustWaitLoad().Timeout(CancelationTimeout).MustElement("#table-0 > thead > tr > td:nth-child(7) > span.grid-table.align-middle.font-small > span > span:nth-child(3)").MustClick().MustWaitLoad()
	})
	if errors.Is(BTNfilter, context.DeadlineExceeded) {
		return page, errors.New("can't find filter 3 days button")
	}

	time.Sleep(Timeout)

	// 3 days minimum price cell
	BTNmediumPrice := rod.Try(func() {
		page.MustWaitLoad().Timeout(CancelationTimeout).MustElement("#table-0 > tbody > tr:nth-child(1) > td:nth-child(8)").MustText()
	})
	if errors.Is(BTNmediumPrice, context.DeadlineExceeded) {
		return page, errors.New("can't find lowcoster KCK sing-in price button")
	}
	mediumPrice := page.MustElement("#table-0 > tbody > tr:nth-child(1) > td:nth-child(8)").MustText()
	fmt.Println("raw price per 3 days is ", mediumPrice)
	mediumPrice = strings.ReplaceAll(mediumPrice, " ", "") //trim spaces

	intConv, err := strconv.Atoi(mediumPrice) // conver to int

	if err != nil {
		return page, errors.New("can't convert equ value: " + mediumPrice)
	}

	mediumPriceInt := intConv / 3 // get price per 1 day
	fmt.Println("current price per 1 day is", mediumPriceInt)
	// raw price per 3 days is  180
	//current price per 1 day is 0

	if mediumPriceInt < account.MaxDailyPrice && mediumPriceInt != 1 {
		// Press the lowest 3 day price point
		BTNcofirm := rod.Try(func() {
			page.MustWaitLoad().Timeout(CancelationTimeout).MustElement("#table-0 > tbody > tr:nth-child(1) > td:nth-child(8)").MustClick().MustWaitLoad()
		})
		if errors.Is(BTNcofirm, context.DeadlineExceeded) {
			return page, errors.New("can't press lowest 3day price button")
		}
		fmt.Println("current price per 1 day is Okay", mediumPriceInt)
		time.Sleep(Timeout)

		return page, nil
	} else {
		// reject this kck
		BTNreject := rod.Try(func() {
			page.MustWaitLoad().Timeout(CancelationTimeout).MustElement("#page-contents > header > div > div > div:nth-child(1) > a > span.img").MustClick().MustWaitLoad()
		})
		if errors.Is(BTNreject, context.DeadlineExceeded) {
			return page, errors.New("can't press reject kck sing-in button")
		}

		fmt.Println("current price per 1 day is too High", mediumPriceInt)
		time.Sleep(Timeout)

		return page, nil
	}

}

// Press "close feed menu button"
func PressCloseFeedButton(page *rod.Page) (err error) {
	page = page.MustWaitLoad()
	time.Sleep(HorseDetectionTimeout)

	BTN := rod.Try(func() {
		page.MustWaitLoad().Timeout(CancelationTimeout).MustElement("#care-tab-feed > table > tbody > tr:nth-child(1) > td.float-right").MustClick().MustWaitLoad()

	})
	if errors.Is(BTN, context.DeadlineExceeded) {
		return errors.New("can't find close finish feed button")
	}

	return nil
}

// Press "feed horse button"
func PressFinishFeedButton(page *rod.Page) (s string, err error) {
	page = page.MustWaitLoad()
	time.Sleep(Timeout)

	BTN := rod.Try(func() {
		page.MustWaitLoad().Timeout(CancelationTimeout).MustElement("#feed-button > span").MustClick().MustWaitLoad()

	})
	if errors.Is(BTN, context.DeadlineExceeded) {
		return "", errors.New("can't press finish feed button-1")
	}

	//page.Timeout(Timeout).MustElement("#feed-button").MustClick()
	time.Sleep(Timeout)
	//page.Timeout(Timeout).ustElement("#feed-button > span").MustClick()

	// KORM error selector: #fieldError-eat-fourrage-inventaire
	// OVEC error selector:

	// If there is an error after pressing feed button // #care-tab-feed > table > tbody > tr:nth-child(2) > td > #messageBoxInline
	isFeedError := rod.Try(func() {
		page.MustWaitLoad().Timeout(HorseDetectionTimeout).MustElement("#fieldError-eat-fourrage-inventaire")

	})
	if errors.Is(isFeedError, context.DeadlineExceeded) {
		fmt.Println("success after feeding")
	} else {

		// Checking if the horse if too fat rn
		isFatHorseError := rod.Try(func() {
			page.MustWaitLoad().Timeout(HorseDetectionTimeout).MustElement("#care-tab-feed > table > tbody > tr:nth-child(2) > td #messageBoxInline")

		})
		if errors.Is(isFatHorseError, context.DeadlineExceeded) {
			// Close feed menu, not enough Korm
			err := PressCloseFeedButton(page)
			return "not enough korm", err

		} else {
			err := PressCloseFeedButton(page)
			return "she is too fat", err
		}
	}

	return "done with feed button", nil

}

// Press "clear horse button"
func PressClearButton(page *rod.Page) (s string, err error) {
	page = page.MustWaitLoad()
	BTN := rod.Try(func() {
		page.MustWaitLoad().Timeout(CancelationTimeout).MustElement("#boutonPanser > span").MustClick().MustWaitLoad()

	})
	if errors.Is(BTN, context.DeadlineExceeded) {
		return "", errors.New("can't press clear button")
	}

	time.Sleep(Timeout)

	return "done with clear button", nil
}

// Press "send horse sleeping button"
func PressSleepButton(page *rod.Page) (s string, err error) {
	page = page.MustWaitLoad()
	BTN := rod.Try(func() {
		page.MustWaitLoad().Timeout(CancelationTimeout).MustElement("#boutonCoucher > span").MustClick().MustWaitLoad()

	})
	if errors.Is(BTN, context.DeadlineExceeded) {
		return "", errors.New("can't press sleep button")
	}

	time.Sleep(Timeout)

	return "done with sleep button", nil
}

func FeedStatusArraySplit(s string) int {
	//defer handleOutOfBounds()
	// Get this: 0 / 4

	x := strings.Split(s, " ") // trying to split it to different arrays with " " key
	y := x[0]
	currentFeedValue, err := strconv.Atoi(y)

	if err != nil {
		fmt.Println("Wrong current feed status selector")
	}

	return currentFeedValue
}

func CurrentEquValue(s string) int {
	//defer handleOutOfBounds()
	/*
		Example:
			Parsed Equ value is:  37 297
			Экю
		integer Equ value is:  37
	*/
	re := regexp.MustCompile(`\r?\n`) // \r - IS only wor Windows?! (remove new lines)
	x := re.ReplaceAllString(s, "_")
	x = strings.ReplaceAll(x, " ", "")

	z := strings.Split(x, "_") // trying to split it to different arrays with "_" key
	y := z[0]
	currentEqu, err := strconv.Atoi(y)

	if err != nil {
		fmt.Println("Wrong Equ selector")
		// if credit
		return 0
	}

	return currentEqu
}
