package orki

//====== Stable linear functions ===========================================================================================================================================
/*
// Linear version starter
func StartOrkiLinear(page *rod.Page, account *models.Account) {
	horseLinks, err := OrkiPasrePages(page, account)

	if err != nil {
		fmt.Println("can't get page length : ", err)
	}

	for counter, lnk := range horseLinks.IDs {

		start := time.Now()
		fmt.Println()
		fmt.Println("Starting horse: ", counter)
		// Linear version:
		//result, err := RunOrkiHorse(page, account, lnk)
		//Async version:
		result, err := RunOrkiHorse(page, account, lnk)
		if err != nil {

			fmt.Println("Error occured on this iteration: ", counter)
			fmt.Println("Error occured on this horse link: ", lnk)
			fmt.Println("Error is: ", err)

		} else {

			fmt.Println("Succes horse iteration: ", counter)
			fmt.Println("Succes horse link: ", lnk)
			fmt.Println(result)
		}
		elapsed := time.Since(start)
		fmt.Println("It took ", elapsed, " time")

		// add it to any variable to return to a user?

	}

}

// Linear version runner
func RunOrkiHorse(page *rod.Page, account *models.Account, horseLink string) (s string, err error) {
	// Open choosen horse (linear version)
	page.MustNavigate(horseLink).MustWaitLoad()

	fmt.Println("Test horse link is: ", horseLink)

	//---- Horse Section --------------------------------
	// Check current Equ value and converting it to int
	money := fmt.Sprintf("%s", page.Timeout(CancelationTimeout).MustElement("#header-hud > ul > li.level-1.float-right.hud-equus > a > span > span:nth-child(2)").MustText())
	moneyInt := CurrentEquValue(money)
	//fmt.Println("Parsed Equ value is: ", money)
	//fmt.Println("integer Equ value is: ", moneyInt)

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
			rtrn, err := RunOldHorseWithKCK(page, account, moneyInt)
			if err != nil {
				return "error finish old horse with KCK", err
			}

			return rtrn, nil

		} else {
			// ---- execute Old horse without KCK --------------------------------
			// Sing in-KCK
			if account.MaxDailyPrice != 1 && moneyInt >= account.MaxDailyPrice {
				pageAfterKCK, err := KCKsingin(page, account)
				if err != nil {
					return "error sing-in KCK", err
				}

				rtrn, er := RunOldHorseWithKCK(pageAfterKCK, account, moneyInt)
				if er != nil {
					return "error finish old horse without KCK", er
				}

				return rtrn, nil
			}

			fmt.Println("No money for KCK sing-in")

			rtrn, err := RunOldHorseWithKCK(page, account, moneyInt)
			if err != nil {
				return "error finish old horse without KCK", err
			}

			return rtrn, nil
		}

	}

	// ---- execute Baby horse --------------------------------
	rtrn, err := RunBabyHorse(page)
	if err != nil {
		return "error finish the bb horse", err
	}

	return rtrn, nil

}

*/
