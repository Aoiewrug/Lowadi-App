package orki

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Aoiewrug/Lowadi-App/lowadi-http-func-server/core/models"
	"github.com/go-rod/rod"
)

var Timeout = 400 * time.Millisecond // Wait element timeout
var HorseDetectionTimeout = Timeout + (400 * time.Millisecond)
var CancelationTimeout = HorseDetectionTimeout + (5000 * time.Millisecond)

var TestHorse = []string{
	"https://www.lowadi.com/elevage/chevaux/cheval?id=67016167",
}

// Async version runner
func RunOrkiHorse(chanStruct models.ChanStruct) (s string, err error) {
	// Async version:
	page := chanStruct.Page.MustWaitLoad()
	page.MustNavigate(chanStruct.SingleHorseLink).MustWaitLoad()

	//fmt.Println("Test horse link is: ", chanStruct.SingleHorseLink)

	//---- Horse Section --------------------------------

	// Check if it is a newborn horse ( we need to name it than)
	tryNewBornHorse0 := rod.Try(func() {
		page.Timeout(HorseDetectionTimeout).MustElement("#poulain-1").MustWaitLoad()

	})
	if errors.Is(tryNewBornHorse0, context.DeadlineExceeded) {
		//fmt.Println("no newborn horses here")
	} else {
		err := NewbornHorse(page, chanStruct.Account)
		if err != nil {
			return "can't born new horse (In KCK)", err
		} else {
			return "successfully get new baby horse!", nil
		}
	}

	// Check if there is a newborn horse ready --- (this case is is if the horse is in KCK I guess)
	tryNewBornHorse1 := rod.Try(func() {
		page.Timeout(HorseDetectionTimeout).MustElement("#boutonVeterinaire > span > span > span").MustClick().MustWaitLoad()

	})
	if errors.Is(tryNewBornHorse1, context.DeadlineExceeded) {
		//fmt.Println("no newborn horses here")
	} else {
		err := NewbornHorse(page, chanStruct.Account)
		if err != nil {
			return "can't born new horse (In KCK)", err
		} else {
			return "successfully get new baby horse!", nil
		}
	}

	// #alerteVeterinaireContent > table > tbody > tr > td.align-top (form it self)
	// #agi-80673576001669571882 (close button)
	//// Check if there is a newborn horse ready --- (this case is is if the horse is NOT in KCK I guess)
	tryNewBornHorse2 := rod.Try(func() {
		page.Timeout(HorseDetectionTimeout).MustElement("#alerteVeterinaireContent > table > tbody > tr > td.align-top").MustWaitLoad()

	})
	if errors.Is(tryNewBornHorse2, context.DeadlineExceeded) {
		//fmt.Println("no newborn horses here")
	} else {
		tryNewBornHorse3 := rod.Try(func() {
			// Click Equ value to drop the window
			//page.Timeout(HorseDetectionTimeout).MustElement("#header-hud > ul > li.level-1.float-right.hud-equus > a > span").MustClick().MustWaitLoad()

			// Try to close it with "X" button
			page.Timeout(HorseDetectionTimeout).MustElementR("#alerteVeterinaireBox > div",
				"close-popup right",
			)

		})
		if errors.Is(tryNewBornHorse3, context.DeadlineExceeded) {
			return "can't close born new horse button (Out of KCK)", err
		} else {
			fmt.Println("successfully closed born new baby horse button (Out of KCK). Please proceed manually")
		}

	}

	// Check current Equ value and converting it to int
	tryCheckEquVal := rod.Try(func() {
		page.Timeout(10000 * time.Millisecond).MustElement("#header-hud > ul > li.level-1.float-right.hud-equus > a > span > span:nth-child(2)").MustWaitLoad()

	})
	if errors.Is(tryCheckEquVal, context.DeadlineExceeded) {
		return "can't parse equ value", errors.New("can't parse equ value")
	}

	money := page.Timeout(10000 * time.Millisecond).MustElement("#header-hud > ul > li.level-1.float-right.hud-equus > a > span > span:nth-child(2)").MustWaitLoad().MustText()
	moneyInt := CurrentEquValue(money)

	// -------------------------------------------

	// Check if it is baby horse :)
	isBabyHorse := rod.Try(func() {
		page.Timeout(HorseDetectionTimeout).MustElement("#boutonAllaiter > span")

	})
	if errors.Is(isBabyHorse, context.DeadlineExceeded) {
		//fmt.Println("Not a baby horse")

		// This is an Old horse. Check KCK ->
		isOldHorseKCK := rod.Try(func() {
			// Checking if KCK sing in button/picture exists
			page.Timeout(HorseDetectionTimeout).MustElement("#cheval-inscription > a > span.img")
		})

		if errors.Is(isOldHorseKCK, context.DeadlineExceeded) {
			// ---- execute Old horse with KCK -------------------------------
			rtrn, err := RunOldHorseWithKCK(page, chanStruct, moneyInt)
			if err != nil {
				return "error finish old horse with KCK", err
			}

			return rtrn, nil

		} else {
			// ---- execute Old horse without KCK --------------------------------
			// Sing in-KCK
			if chanStruct.Account.MaxDailyPrice != 1 && moneyInt >= chanStruct.Account.MaxDailyPrice {
				pageAfterKCK, err := KCKsingin(page, chanStruct.Account)
				if err != nil {
					return "error sing-in KCK", err
				}

				rtrn, er := RunOldHorseWithKCK(pageAfterKCK, chanStruct, moneyInt)
				if er != nil {
					return "error finish old horse without KCK", er
				}

				return rtrn, nil
			} else {
				//fmt.Println("No money for KCK sing-in")

				rtrn, err := RunOldHorseWithKCK(page, chanStruct, moneyInt)
				if err != nil {
					return "error finish old horse without KCK", err
				}

				return rtrn, nil
			}

		}

	}

	// ---- execute Baby horse --------------------------------
	rtrn, err := RunBabyHorse(page)
	if err != nil {
		return "error finish the bb horse", err
	}

	return rtrn, nil

}
