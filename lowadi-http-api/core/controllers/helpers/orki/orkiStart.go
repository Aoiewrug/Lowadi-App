package orki

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Aoiewrug/Lowadi-App/lowadi-http-api/core/models"
	"github.com/go-rod/rod"
)

var Timeout = 400 * time.Millisecond // Wait element timeout
var HorseDetectionTimeout = Timeout + (400 * time.Millisecond)
var CancelationTimeout = HorseDetectionTimeout + (5000 * time.Millisecond)

var TestHorse = []string{
	"https://www.lowadi.com/elevage/chevaux/cheval?id=67016167",
}

const goroutines = 20 // Total number of threads to use, excluding the main() thread

type chanStruct struct {
	browser *rod.Browser
	account *models.Account
	link    string
	counter int
}

// Async test
func TestStartOrkiAsync(browser *rod.Browser, account *models.Account) {
	var ch = make(chan chanStruct) // buffered channel with 0 beffer :D
	var wg sync.WaitGroup

	// This starts goroutines number of goroutines that wait for something to do
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			for {
				chanStruct, ok := <-ch
				if !ok { // if there is nothing to do and the channel has been closed then end the goroutine
					wg.Done()
					return
				}
				// =====

				// create 1 page per thread on this task
				page := chanStruct.browser.MustPage(chanStruct.link)

				start := time.Now()
				fmt.Println()
				fmt.Println("Starting horse: ", chanStruct.counter)
				result, err := RunOrkiHorseAsync(page, chanStruct) // do the thing
				if err != nil {
					fmt.Println("Error occured on this iteration: ", chanStruct.counter)
					fmt.Println("Error occured on this horse link: ", chanStruct.link)
					fmt.Println("Error is: ", err)
				} else {
					fmt.Println("Succes horse iteration: ", chanStruct.counter)
					fmt.Println("Succes horse link: ", chanStruct.link)
					fmt.Println(result)
				}
				elapsed := time.Since(start)
				fmt.Println("It took ", elapsed, " time")
				page.MustClose()
				// =====
			}
		}()
	}

	// Now the jobs can be added to the channel, which is used as a queue
	for counter, lnk := range TestHorse {
		ch <- chanStruct{browser, account, lnk, counter} // add job to the queue
		time.Sleep(1000 * time.Millisecond)
		fmt.Println("I've sleeped 1 second while starting this goroutine")
	}

	close(ch) // This tells the goroutines there's nothing else to do
	wg.Wait() // Wait for the threads to finish

}

// Test call :8002 lowadi-http-func-server
func TestCallFuncServer(browser *rod.Browser, page *rod.Page, account *models.Account) {
	/*horseLinks, err := OrkiPasrePages(page, account)
	if err != nil {
		fmt.Println("can't get page length : ", err)
	}
	*/
}

// Async version starter
func StartOrkiAsync(browser *rod.Browser, page *rod.Page, account *models.Account) {
	horseLinks, err := OrkiPasrePages(page, account)

	if err != nil {
		fmt.Println("can't get page length : ", err)
	}

	var ch = make(chan chanStruct) // buffered channel with 0 beffer :D
	var wg sync.WaitGroup

	// This starts goroutines number of goroutines that wait for something to do
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			for {
				chanStruct, ok := <-ch
				if !ok { // if there is nothing to do and the channel has been closed then end the goroutine
					wg.Done()
					return
				}
				// =====
				// create 1 page per thread on this task
				pageTwo := chanStruct.browser.MustPage(chanStruct.link)

				start := time.Now()
				fmt.Println()
				fmt.Println("Starting horse: ", chanStruct.counter)
				result, err := RunOrkiHorseAsync(pageTwo, chanStruct) // do the thing
				if err != nil {
					fmt.Println("Error occured on this iteration: ", chanStruct.counter)
					fmt.Println("Error occured on this horse link: ", chanStruct.link)
					fmt.Println("Error is: ", err)
				} else {
					fmt.Println("Succes horse iteration: ", chanStruct.counter)
					fmt.Println("Succes horse link: ", chanStruct.link)
					fmt.Println(result)
				}
				elapsed := time.Since(start)
				fmt.Println("It took ", elapsed, " time")
				pageTwo.MustClose()
				// =====
			}
		}()
	}

	// Now the jobs can be added to the channel, which is used as a queue
	for counter, lnk := range horseLinks.IDs {
		// we need to run 1000 horses max
		if counter >= 1001 {
			break
		}

		ch <- chanStruct{browser, account, lnk, counter} // add job to the queue
		time.Sleep(1000 * time.Millisecond)
		//fmt.Println("I've sleeped 1 second while starting this goroutine")
	}

	close(ch) // This tells the goroutines there's nothing else to do
	wg.Wait() // Wait for the threads to finish

}

// Async version runner
func RunOrkiHorseAsync(page *rod.Page, chanStruct chanStruct) (s string, err error) {
	// Async version:
	page = page.MustWaitLoad()
	page.MustNavigate(chanStruct.link).MustWaitLoad()

	fmt.Println("Test horse link is: ", chanStruct.link)

	//---- Horse Section --------------------------------
	// Check current Equ value and converting it to int
	tryCheckEquVal := rod.Try(func() {
		page.Timeout(10000 * time.Millisecond).MustElement("#header-hud > ul > li.level-1.float-right.hud-equus > a > span > span:nth-child(2)").MustWaitLoad()

	})
	if errors.Is(tryCheckEquVal, context.DeadlineExceeded) {
		return "can't parse equ value", errors.New("can't parse equ value")
	}

	money := page.Timeout(10000 * time.Millisecond).MustElement("#header-hud > ul > li.level-1.float-right.hud-equus > a > span > span:nth-child(2)").MustWaitLoad().MustText()
	moneyInt := CurrentEquValue(money)

	// Check if there is a newborn horse ready ---
	tryNewBornHorse := rod.Try(func() {
		page.Timeout(HorseDetectionTimeout).MustElement("#boutonVeterinaire > span > span > span").MustClick().MustWaitLoad()

	})
	if errors.Is(tryNewBornHorse, context.DeadlineExceeded) {
		fmt.Println("no newborn horses here")
	} else {
		err := NewbornHorse(page, chanStruct.account)
		if err != nil {
			return "can't born new horse", err
		} else {
			return "successfully get new baby horse!", nil
		}
	}
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
			if chanStruct.account.MaxDailyPrice != 1 && moneyInt >= chanStruct.account.MaxDailyPrice {
				pageAfterKCK, err := KCKsingin(page, chanStruct.account)
				if err != nil {
					return "error sing-in KCK", err
				}

				rtrn, er := RunOldHorseWithKCK(pageAfterKCK, chanStruct, moneyInt)
				if er != nil {
					return "error finish old horse without KCK", er
				}

				return rtrn, nil
			} else {
				fmt.Println("No money for KCK sing-in")

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
