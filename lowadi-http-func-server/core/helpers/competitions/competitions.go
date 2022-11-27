package competitions

import (
	"fmt"

	"github.com/Aoiewrug/Lowadi-App/lowadi-http-func-server/core/models"
	"github.com/go-rod/rod"
)

var TestHorse = "https://www.lowadi.com/elevage/chevaux/cheval?id=62907056"

// Competition selection grid
var Ryc = "#competition-body-content > table > tbody > tr:nth-child(1) > td.first.top"    // 1st row 1st element
var Galop = "#competition-body-content > table > tbody > tr:nth-child(1) > td.middle.top" // 1st row 2nd element
var Vyezdka = "#competition-body-content > table > tbody > tr:nth-child(1) > td.last.top" // 1st row 3rd element

var Cross = "#competition-body-content > table > tbody > tr:nth-child(2) > td.first.top.bottom"  // 2nd row 1st element
var Dunno = "#competition-body-content > table > tbody > tr:nth-child(2) > td.middle.top.bottom" // 2nd row 2nd element
var Konkur = "#competition-body-content > table > tbody > tr:nth-child(2) > td.last.top.bottom"  // 2nd row 3rd element
// ==================

var Energy = "#energie"                // energy level
var OR = "#boutonVieillir > span > em" // grow points value

var StartButton = "#race > tbody > tr:nth-child(1) > td:nth-child(9)" // Press start 1st competition button

func StartCompetitionsTest(page *rod.Page, account *models.Account) {
	fmt.Println("I'm in competition section!")
	/*
		page = page.MustWaitLoad()
		page.MustNavigate(TestHorse)

		energy := page.MustElement(Energy).MustText()
		growPoints := page.MustElement(OR).MustText()

		fmt.Println(energy)
		fmt.Println(growPoints)
	*/

}

func StartCompetitions(page *rod.Page, account *models.Account) {

}
