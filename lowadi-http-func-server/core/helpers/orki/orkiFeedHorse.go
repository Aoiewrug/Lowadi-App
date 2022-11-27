package orki

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/Aoiewrug/Lowadi-App/lowadi-http-func-server/core/models"
	"github.com/go-rod/rod"
)

// new version
func FeedTwoBarHorse(page *rod.Page, account *models.Account) (errors error) {
	page = page.MustWaitLoad()

	// fmt.Println("There is second feed bar")

	// 0 resp code meand didn't feed it. 1 = successfully feed
	response1, _ := FeedKormBar(page, account)
	//fmt.Println(er)

	// 0 resp code meand didn't feed it. 1 = successfully feed
	response2, _ := FeedOvecBar(page, account)
	//fmt.Println(err)

	if response1 == 0 && response2 == 0 {
		err2 := PressCloseFeedButton(page)
		if err2 != nil {
			return err2
		}
	} else {
		// Press "give food button"
		_, err3 := PressFinishFeedButton(page)
		if err3 != nil {
			return err3
		}
		//fmt.Println("feed status: ", x)
	}

	return nil
}

func FeedOneBarHorse(page *rod.Page, account *models.Account) (errors error) {
	page = page.MustWaitLoad()

	//fmt.Println("No second feed bar here")

	// If no error - feed, else close
	response, _ := FeedKormBar(page, account)
	//fmt.Println(err)
	if response == 1 {
		// Press "give food button"
		_, err := PressFinishFeedButton(page)
		if err != nil {
			return err
		}
		//fmt.Println("feed status: ", x)
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
			//fmt.Println("Fall into 5 child ovec section: ", currentKormLink)
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
	//fmt.Println("Starting with Ovec bar")
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
			//fmt.Println("Fall into 5 child ovec section: ", currentOvecLink)
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
