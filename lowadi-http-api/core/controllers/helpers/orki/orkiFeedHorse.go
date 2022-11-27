package orki

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/Aoiewrug/Lowadi-App/lowadi-http-api/core/models"
	"github.com/go-rod/rod"
)

/*
// press KCK button
func FeedTwoBarHorseOld(page *rod.Page, account *models.Account) (errors error) {
	var intShouldKorm int
	var intShouldOvec int

	var currentKormValue string
	var currentOvecValue string

	time.Sleep(Timeout)

	fmt.Println("We have second feed bar here")
	// #feeding > table:nth-child(5) > tbody > tr.dashed.section-fourrage.section-fourrage-content > td:nth-child(1) > span.float-right.section-fourrage.section-fourrage-quantity > strong
	currentKormLink6 := "#feeding > table:nth-child(6) > tbody > tr.dashed.section-fourrage.section-fourrage-content > td:nth-child(1) > span.float-right.section-fourrage.section-fourrage-quantity"
	currentOvecLink6 := "#feeding > table:nth-child(6) > tbody > tr.dashed.section-avoine.section-avoine-content > td:nth-child(1) > span.float-right.section-avoine.section-avoine-quantity"
	shouldKormLink6 := currentKormLink6 + " > strong"
	shouldOvecLink6 := currentOvecLink6 + " > strong"

	// #feeding > table:nth-child(5) > tbody > tr.dashed.section-fourrage.section-fourrage-content > td:nth-child(1) > span.float-right.section-fourrage.section-fourrage-quantity > strong
	currentKormLink5 := "#feeding > table:nth-child(5) > tbody > tr.dashed.section-fourrage.section-fourrage-content > td:nth-child(1) > span.float-right.section-fourrage.section-fourrage-quantity"
	currentOvecLink5 := "#feeding > table:nth-child(5) > tbody > tr.dashed.section-avoine.section-avoine-content > td:nth-child(1) > span.float-right.section-avoine.section-avoine-quantity"
	shouldKormLink5 := currentKormLink5 + " > strong"
	shouldOvecLink5 := currentOvecLink5 + " > strong"

	time.Sleep(Timeout)

	// Give KORM -0/4- selecting "4"
	// Try table:nth-child(6) selector first
	BTNkormvalue6 := rod.Try(func() {
		page.Timeout(CancelationTimeout).MustElement(shouldKormLink6)
	})
	if errors.Is(BTNkormvalue6, context.DeadlineExceeded) {
		//Try table:nth-child(5) selector than
		BTNkormvalue5 := rod.Try(func() {
			page.Timeout(CancelationTimeout).MustElement(shouldKormLink6)
		})
		if errors.Is(BTNkormvalue5, context.DeadlineExceeded) {
			err := PressCloseFeedButton(page)
			if err != nil {
				return err
			}
			return errors.New("can't find max korm current slider value")
		} else {
			shouldKorm := page.Timeout(CancelationTimeout).MustElement(shouldKormLink5).MustText()
			intShouldKorm, _ = strconv.Atoi(shouldKorm)
		}

	} else {
		shouldKorm := page.Timeout(CancelationTimeout).MustElement(shouldKormLink6).MustText()
		intShouldKorm, _ = strconv.Atoi(shouldKorm)
	}

	// Give OVEC -0/4- selecting "4"
	// Try table:nth-child(6) selector first, than table:nth-child(5)
	BTNovecvalue6 := rod.Try(func() {
		page.Timeout(CancelationTimeout).MustElement(shouldOvecLink6)
	})
	if errors.Is(BTNovecvalue6, context.DeadlineExceeded) {
		//Try table:nth-child(5) selector than
		BTNovecvalue5 := rod.Try(func() {
			page.Timeout(CancelationTimeout).MustElement(shouldOvecLink5)
		})
		if errors.Is(BTNovecvalue5, context.DeadlineExceeded) {
			err := PressCloseFeedButton(page)
			if err != nil {
				return err
			}
			return errors.New("can't find max ovec current slider value")
		} else {
			shouldOvec := page.Timeout(CancelationTimeout).MustElement(shouldOvecLink5).MustText()
			intShouldOvec, _ = strconv.Atoi(shouldOvec)
		}

	} else {
		shouldOvec := page.Timeout(CancelationTimeout).MustElement(shouldOvecLink6).MustText()
		intShouldOvec, _ = strconv.Atoi(shouldOvec)
	}

	// Check korm slider value (-0/4- selecting "0") ==================================================
	// Try table:nth-child(6) selector first, than table:nth-child(5)
	BTNcurrentkormvalue6 := rod.Try(func() {
		page.Timeout(CancelationTimeout).MustElement(currentKormLink6)
	})
	if errors.Is(BTNcurrentkormvalue6, context.DeadlineExceeded) {
		//Try table:nth-child(5) selector than
		BTNcurrentkormvalue5 := rod.Try(func() {
			page.Timeout(CancelationTimeout).MustElement(currentKormLink5)
		})
		if errors.Is(BTNcurrentkormvalue5, context.DeadlineExceeded) {
			err := PressCloseFeedButton(page)
			if err != nil {
				return err
			}
			return errors.New("can't find current korm slider value")
		} else {
			currentKormValue = page.Timeout(CancelationTimeout).MustElement(currentKormLink5).MustText()
		}

	} else {
		currentKormValue = page.Timeout(CancelationTimeout).MustElement(currentKormLink6).MustText()
	}

	// check ovec slider value
	// Try table:nth-child(6) selector first, than table:nth-child(5)
	BTNcurrentovecvalue6 := rod.Try(func() {
		page.Timeout(CancelationTimeout).MustElement(currentOvecLink6)
	})
	if errors.Is(BTNcurrentovecvalue6, context.DeadlineExceeded) {
		//Try table:nth-child(5) selector than
		BTNcurrentovecvalue5 := rod.Try(func() {
			page.Timeout(CancelationTimeout).MustElement(currentOvecLink5)
		})
		if errors.Is(BTNcurrentovecvalue5, context.DeadlineExceeded) {
			err := PressCloseFeedButton(page)
			if err != nil {
				return err
			}
			return errors.New("can't find current ovec slider value")
		} else {
			currentOvecValue = page.Timeout(CancelationTimeout).MustElement(currentOvecLink5).MustText()
		}

	} else {
		currentOvecValue = page.Timeout(CancelationTimeout).MustElement(currentOvecLink6).MustText()
	}

	currentKormValueInt := FeedStatusArraySplit(currentKormValue)
	currentOvecValueInt := FeedStatusArraySplit(currentOvecValue)

	// If we already feeded the horse we won't feed it again
	if currentKormValueInt == 0 {
		if intShouldKorm == 0 {
			time.Sleep(Timeout)
			// Do nothing
			fmt.Println("Do nothing")
			// If we've already feeded the horse we won't feed it again
			if currentOvecValueInt == 0 {
				if intShouldOvec == 0 {
					// Press "close button"
					err := PressCloseFeedButton(page)
					if err != nil {
						return err
					}
				} else {
					correctOvecValue := fmt.Sprintf("#oatsSlider > ol > li:nth-child(%v)", (intShouldOvec + 1))
					BTN := rod.Try(func() {
						page.Timeout(CancelationTimeout).MustElement(correctOvecValue).MustWaitLoad().MustClick().MustWaitLoad()
					})
					if errors.Is(BTN, context.DeadlineExceeded) {
						err := PressCloseFeedButton(page)
						if err != nil {
							return err
						}
						return errors.New("can't select correct slider-1")
					} else {
						// Press "give food button"
						x, err := PressFinishFeedButton(page)
						if err != nil {
							return err
						}
						fmt.Println("feed status: ", x)
					}

				}
			} else {
				// Press "close button"
				err := PressCloseFeedButton(page)
				if err != nil {
					return err
				}
			}
		} else {
			time.Sleep(Timeout)
			// Here we should feed the horse in any case cuz KORM will be setted up as not ZERO value!
			// Set up KORM value
			correctKormValue := fmt.Sprintf("#haySlider > ol > li:nth-child(%v)", (intShouldKorm + 1))
			BTN := rod.Try(func() {
				page.Timeout(CancelationTimeout).MustElement(correctKormValue).MustWaitLoad().MustClick().MustWaitLoad()
			})
			if errors.Is(BTN, context.DeadlineExceeded) {
				err := PressCloseFeedButton(page)
				if err != nil {
					return err
				}
				return errors.New("can't select correct slider-2")
			}
			if currentOvecValueInt == 0 {
				if intShouldOvec == 0 {
					// Press "give food button"
					x, err := PressFinishFeedButton(page)
					if err != nil {
						return err
					}
					fmt.Println("feed status: ", x)
					return nil
				} else {
					correctOvecValue := fmt.Sprintf("#oatsSlider > ol > li:nth-child(%v)", (intShouldOvec + 1))
					BTN := rod.Try(func() {
						page.Timeout(CancelationTimeout).MustElement(correctOvecValue).MustWaitLoad().MustClick().MustWaitLoad()
					})
					if errors.Is(BTN, context.DeadlineExceeded) {
						// Press "give food button"
						x, err := PressFinishFeedButton(page)
						if err != nil {
							return err
						}
						fmt.Println("feed status: ", x)
						return nil
					}
					// Press "give food button"
					x, err := PressFinishFeedButton(page)
					if err != nil {
						return err
					}
					fmt.Println("feed status: ", x)
					return nil
				}
			} else {
				// Press "close button"
				// Press "give food button"
				x, err := PressFinishFeedButton(page)
				if err != nil {
					return err
				}
				fmt.Println("feed status: ", x)
				return nil
			}
		}
	}

	/*else {
		// Press "close button"
		err := PressCloseFeedButton(page)
		if err != nil {
			return err
		}
	}


	return nil

}
*/
// new version
func FeedTwoBarHorse(page *rod.Page, account *models.Account) (errors error) {
	page = page.MustWaitLoad()

	fmt.Println("There is second feed bar")

	// 0 resp code meand didn't feed it. 1 = successfully feed
	response1, er := FeedKormBar(page, account)
	fmt.Println(er)

	// 0 resp code meand didn't feed it. 1 = successfully feed
	response2, err := FeedOvecBar(page, account)
	fmt.Println(err)

	if response1 == 0 && response2 == 0 {
		err2 := PressCloseFeedButton(page)
		if err2 != nil {
			return err2
		}
	} else {
		// Press "give food button"
		x, err3 := PressFinishFeedButton(page)
		if err3 != nil {
			return err3
		}
		fmt.Println("feed status: ", x)
	}

	return nil
}

func FeedOneBarHorse(page *rod.Page, account *models.Account) (errors error) {
	page = page.MustWaitLoad()

	fmt.Println("No second feed bar here")

	// If no error - feed, else close
	response, err := FeedKormBar(page, account)
	fmt.Println(err)
	if response == 1 {
		// Press "give food button"
		x, err := PressFinishFeedButton(page)
		if err != nil {
			return err
		}
		fmt.Println("feed status: ", x)
	} else {
		err2 := PressCloseFeedButton(page)
		if err2 != nil {
			return err2
		}
	}

	return nil
}

// 0 resp code meand didn't feed it. 1 = successfully feed
func FeedKormBar(page *rod.Page, account *models.Account) (resp int, err error) {
	page = page.MustWaitLoad()

	/* remove
	selectorValueString := strconv.Itoa(selectorValue)

	fmt.Println("Starting with Korm bar")
	maxKormLink := "#feeding > table:nth-child(" + selectorValueString + ") > tbody > tr.dashed.section-fourrage.section-fourrage-content > td:nth-child(1) > span.float-right.section-fourrage.section-fourrage-quantity > strong"
	currentKormLink := "#feeding > table:nth-child(" + selectorValueString + ") > tbody > tr.dashed.section-fourrage.section-fourrage-content > td:nth-child(1) > span.float-right.section-fourrage.section-fourrage-quantity"
	*/

	//var childSelector string
	var currentKormLink string // "#feeding > table:nth-child(" + childSelector + ") > tbody > tr.dashed.section-fourrage.section-fourrage-content > td:nth-child(1) > span.float-right.section-fourrage.section-fourrage-quantity"
	var maxKormLink string

	// Give KORM -0/4- selecting 4 value
	BTNkormvalue6 := rod.Try(func() {
		// Try table:nth-child(6)  first
		currentKormLink = "#feeding > table:nth-child(6) > tbody > tr.dashed.section-fourrage.section-fourrage-content > td:nth-child(1) > span.float-right.section-fourrage.section-fourrage-quantity"
		maxKormLink = currentKormLink + " > strong"
		page.Timeout(CancelationTimeout).MustElement(currentKormLink)
	})
	if errors.Is(BTNkormvalue6, context.DeadlineExceeded) {

		BTNkormvalue5 := rod.Try(func() {
			// Try table:nth-child(5)  second
			currentKormLink = "#feeding > table:nth-child(5) > tbody > tr.dashed.section-fourrage.section-fourrage-content > td:nth-child(1) > span.float-right.section-fourrage.section-fourrage-quantity"
			maxKormLink = currentKormLink + " > strong"
			fmt.Println("Fall into 5 child ovec section: ", currentKormLink)
			page.Timeout(CancelationTimeout).MustElement(currentKormLink)
		})
		if errors.Is(BTNkormvalue5, context.DeadlineExceeded) {
			//
			return 0, errors.New("can't find max korm value")
		}

	}

	shouldKorm := page.Timeout(CancelationTimeout).MustElement(maxKormLink).MustText()
	intShouldKorm, _ := strconv.Atoi(shouldKorm)

	// How feed is horse now? -0/4- selecting 0 value
	BTNcurrentkormvalue := rod.Try(func() {
		page.Timeout(CancelationTimeout).MustElement(currentKormLink)
	})
	if errors.Is(BTNcurrentkormvalue, context.DeadlineExceeded) {
		return 0, errors.New("can't find current korm value")
	}
	currentKormValue := page.Timeout(CancelationTimeout).MustElement(currentKormLink).MustText()
	currentKormValueInt := FeedStatusArraySplit(currentKormValue)

	// If we already feeded the horse we won't feed it again
	if currentKormValueInt == intShouldKorm {

		// close option
		//fmt.Println("there is no need to feed it (0/0)")
		return 0, errors.New("no need to feed it")
	} else {
		correctKormValue := fmt.Sprintf("#haySlider > ol > li:nth-child(%v)", (intShouldKorm + 1))
		BTNsliderkormvalue := rod.Try(func() {
			page.Timeout(CancelationTimeout).MustElement(correctKormValue).MustClick().MustWaitLoad()
		})
		if errors.Is(BTNsliderkormvalue, context.DeadlineExceeded) {
			return 0, errors.New("err can't find correct korm slider value value")
		} else {
			// feed option
			return 1, nil
		}

	}
}

// 0 resp code meand didn't feed it. 1 = successfully feed
func FeedOvecBar(page *rod.Page, account *models.Account) (resp int, err error) {
	fmt.Println("Starting with Ovec bar")
	/*
			selectorValueString := strconv.Itoa(selectorValue)

			maxOvecLink := "#feeding > table:nth-child(" + selectorValueString + ") > tbody > tr.dashed.section-avoine.section-avoine-content > td:nth-child(1) > span.float-right.section-avoine.section-avoine-quantity > strong"

			currentOvecLink := "#feeding > table:nth-child(" + selectorValueString + ") > tbody > tr.dashed.section-avoine.section-avoine-content > td:nth-child(1) > span.float-right.section-avoine.section-avoine-quantity"

		childSelector := "6"
		maxOvecLink := "#feeding > table:nth-child(" + childSelector + ") > tbody > tr.dashed.section-avoine.section-avoine-content > td:nth-child(1) > span.float-right.section-avoine.section-avoine-quantity > strong"
		currentOvecLink := "#feeding > table:nth-child(" + childSelector + ") > tbody > tr.dashed.section-avoine.section-avoine-content > td:nth-child(1) > span.float-right.section-avoine.section-avoine-quantity"
	*/
	var currentOvecLink string // "#feeding > table:nth-child(" + childSelector + ") > tbody > tr.dashed.section-fourrage.section-fourrage-content > td:nth-child(1) > span.float-right.section-fourrage.section-fourrage-quantity"
	var maxOvecLink string

	// Give Ovec -0/4- selecting 4 value
	BTNovecvalue6 := rod.Try(func() {
		// Try table:nth-child(6)  first
		currentOvecLink = "#feeding > table:nth-child(6) > tbody > tr.dashed.section-avoine.section-avoine-content > td:nth-child(1) > span.float-right.section-avoine.section-avoine-quantity"
		maxOvecLink = currentOvecLink + " > strong"
		page.Timeout(CancelationTimeout).MustElement(currentOvecLink)
	})
	if errors.Is(BTNovecvalue6, context.DeadlineExceeded) {
		BTNovecvalue5 := rod.Try(func() {
			// Try table:nth-child(5)  second
			currentOvecLink = "#feeding > table:nth-child(5) > tbody > tr.dashed.section-avoine.section-avoine-content > td:nth-child(1) > span.float-right.section-avoine.section-avoine-quantity"
			maxOvecLink = currentOvecLink + " > strong"
			fmt.Println("Fall into 5 child ovec section: ", currentOvecLink)
			page.Timeout(CancelationTimeout).MustElement(currentOvecLink)
		})
		if errors.Is(BTNovecvalue5, context.DeadlineExceeded) {
			return 0, errors.New("err can't find max ovec value")
		}

	}

	shouldOvec := page.Timeout(CancelationTimeout).MustElement(maxOvecLink).MustText()
	intShouldOvec, _ := strconv.Atoi(shouldOvec)

	// How feed is horse now? -0/4- selecting 0 value
	BTNcurrentovecvalue := rod.Try(func() {
		page.Timeout(CancelationTimeout).MustElement(currentOvecLink)
	})
	if errors.Is(BTNcurrentovecvalue, context.DeadlineExceeded) {
		return 0, errors.New("can't find current ovec value")
	}
	currentOvecValue := page.Timeout(CancelationTimeout).MustElement(currentOvecLink).MustText()
	currentOvecValueInt := FeedStatusArraySplit(currentOvecValue)

	// If we already feeded the horse we won't feed it again
	if intShouldOvec == currentOvecValueInt {
		// close option
		//fmt.Println("there is no need to feed it (0/0)")
		return 0, errors.New("no need to feed it with ovec")
	} else {
		correctOvecValue := fmt.Sprintf("#oatsSlider > ol > li:nth-child(%v)", (intShouldOvec + 1))
		BTNsliderovecvalue := rod.Try(func() {
			page.Timeout(CancelationTimeout).MustElement(correctOvecValue).MustClick().MustWaitLoad()
		})
		if errors.Is(BTNsliderovecvalue, context.DeadlineExceeded) {
			return 0, errors.New("can't find correct ovec slider value value")
		} else {
			// feed option
			return 1, nil
		}

	}
}
