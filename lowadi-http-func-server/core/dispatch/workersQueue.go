package dispatch

import (
	"fmt"
	"time"

	"github.com/Aoiewrug/Lowadi-App/lowadi-http-func-server/core/helpers"
	"github.com/Aoiewrug/Lowadi-App/lowadi-http-func-server/core/helpers/dbupdate"
	"github.com/Aoiewrug/Lowadi-App/lowadi-http-func-server/core/helpers/orki"
	"github.com/Aoiewrug/Lowadi-App/lowadi-http-func-server/core/models"
)

const maxAccountPoolBuffer = 100
const maxThreadsCount = 12
const HowManyFunctionsPerAccount = 4

const AccountChannelUpdateInterval = 1000 * time.Millisecond
const JobChannelUpdateInterval = 1000 * time.Millisecond

// const WorkChannelUpdateInterval = 500 * time.Millisecond

var AccountQueue = make(chan models.ChanStruct, maxAccountPoolBuffer)
var JobQueue = make(chan models.ChanStruct, (maxAccountPoolBuffer * HowManyFunctionsPerAccount)) // maxAccountPoolBuffer * 4 tasks for each horse
var WorkQueue = make(chan models.ChanStruct)

// Add account to the AccountQueue (called via external request)
func AddAccount(account *models.Account) {
	//fmt.Println("How big is the buffer right now?", len(AccountQueue))
	AccountQueue <- models.ChanStruct{Account: account}
}

// Counstantly listening for new Accounts
// New accs could come frome outside, from another server
func ReadAccountChannel() {
	for structure := range AccountQueue {
		// Pass if there are no accounts
		if structure.Account == nil {
			fmt.Println("There is nothing to process in Account channel")
		} else {
			// 	fmt.Println("How big is the buffer right now?", len(AccountQueue))
			// 	fmt.Println("Counter from channel: ", structure.Counter)
			// 	fmt.Println("here is the account from channel: ", structure.Account)

			// Send this account to the working queue
			fmt.Println("AccountQueue: Sending this account ot the Job queue")
			JobQueue <- structure
		}

		time.Sleep(AccountChannelUpdateInterval)
	}

}

// Counstantly listening for new Jobs
// New jobs could be added from:
//(new short job) Account queue
//(old long job) From Workers
func ReadJobChannel() {
	for structure := range JobQueue {
		// If the structure is send from the worker
		// len(structure.HorseArrayLinks) will be more than (0 and 1)
		// And we need to loop over all orki links and send it back to the worker pool
		if len(structure.HorseArrayLinks) != 0 {

			// Now the jobs will be added to the worker channel
			for counter, lnk := range structure.HorseArrayLinks {
				// we need to run 1000 horses max and 50 for any ocasions
				if counter >= 1050 {
					break
				}

				page := structure.Browser.MustPage(structure.SingleHorseLink)

				WorkQueue <- models.ChanStruct{
					Account:         structure.Account,
					Page:            page,
					HorseArrayLinks: nil, // skip this to reduce some mem usage :D
					SingleHorseLink: lnk,
					Counter:         counter,
				} // add job to the worker queue

				time.Sleep(JobChannelUpdateInterval)
				//time.Sleep(300 * time.Millisecond)
				//fmt.Println("I've sleeped 1 second while starting this goroutine")
			}

			// Timeout to avoid crushing already openned pages
			time.Sleep(30 * time.Second)
			structure.Browser.MustClose()

			// add DB status "User can log-in the account"
			dbupdate.UpdateRunningAccountStatus(structure.Account, 1)

		} else {
			// In this case the request came from the Account chennel and it is 1st run of this account
			fmt.Println("JobQueue: Sending this account to the Worker queue")
			//defer structure.Browser.MustClose()
			WorkQueue <- structure
			time.Sleep(JobChannelUpdateInterval)
		}

	}

}

// Create fixed amount of constantly working workers
// Which are listening for WorkQueue channel
func Worker() {
	for i := 0; i < maxThreadsCount; i++ {
		n := i
		go func() {
			for {
				chanStruct, ok := <-WorkQueue
				if !ok { // pass if there is nothing to do and the channel has been closed then end the goroutine
					return
				}
				// ===== Execute =====
				defer chanStruct.Browser.MustClose()
				// If we have link we need to run orki async
				if chanStruct.SingleHorseLink != "" {
					//fmt.Println("Starting async worker = ", n)
					// fmt.Println(chanStruct.Account)
					// fmt.Println(chanStruct.Browser)
					//fmt.Println(chanStruct.Counter)
					//fmt.Println(chanStruct.SingleHorseLink)

					// Do work
					result, err := orki.RunOrkiHorse(chanStruct) // do the thing
					if err != nil {
						fmt.Println(
							" User: ", chanStruct.Account.Login,
							" Iteration: ", chanStruct.Counter,
							" Error link: ", chanStruct.SingleHorseLink,
							" Error is: ", err)

					} else {
						fmt.Println(
							" User: ", chanStruct.Account.Login,
							" Iteration: ", chanStruct.Counter,
							" Succes link: ", chanStruct.SingleHorseLink,
							" Result is: ", result)

					}
					chanStruct.Page.MustClose()

				} else {
					// If Links are empty we can run lineary without Orki spread out

					fmt.Println("Starting linear worker = ", n)
					// Do work
					statusCode, chanStruct := helpers.GameLogIn(chanStruct.Account)
					fmt.Println(statusCode, chanStruct.Error)
					//fmt.Println(chanStruct)

					// If we need to run Orki - Parse orki horses links and  Send it back on the upper lvl
					// With added browser to the structure
					if chanStruct.Account.RunOrki == 1 && statusCode == 1 {
						JobQueue <- chanStruct
					} else {
						fmt.Println("I've passed async section")
						// add DB status "User can log-in the account"
						dbupdate.UpdateRunningAccountStatus(chanStruct.Account, 1)

						// !!!!!!!!!!!!!!!!!!!!
						//  Should I exit here?
						// !!!!!!!!!!!!!!!!!!!!

					}

					//time.Sleep(4000 * time.Millisecond)
				}
			}
		}()
	}

}
