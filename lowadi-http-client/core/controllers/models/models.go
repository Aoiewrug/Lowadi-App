package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Cookie struct {
	Token string `json:"token"`
}

// Convert json to struct here https://mholt.github.io/json-to-go/

type Subnet struct {
	gorm.Model
	UUID       string `json:"uuid"`
	ServerName string `json:"name"`
	Subnet     string `json:"subnet"`
}

type Trust struct {
	gorm.Model
	Link   string `json:"link"`
	Rating string `json:"rating"`
}

type PackageComp struct {
	Link []string `json:"link"`
}

type PackageCompResult struct {
	List []string `json:"list"`
}

type Purpose struct {
	gorm.Model
	UUID        string `json:"link"`
	PurposeName string `json:"rating"`
}

type CryptoPayment struct {
	ID             uint      `json:"id" gorm:"primarykey"`
	Active         bool      `json:"active" gorm:"column:active"`                  // Does he use the package? (1 if yes, 0 if no) Used for manual shutdown
	Email          string    `json:"email" gorm:"column:email"`                    // *ENTER* email
	Package        string    `json:"package" gorm:"column:package"`                // *ENTER* package link
	Coin           string    `json:"coin" gorm:"column:coin"`                      // *ENTER* BTC, ETH, USDT...
	TransactionID  string    `json:"transactionid" gorm:"column:transaction_id"`   // *ENTER* trx-id
	PayAmount      float32   `json:"payamount" gorm:"column:pay_amount"`           // *ENTER* how many $ did he pay
	Timer          int       `json:"timer" gorm:"column:timer"`                    // *ENTER* (days)
	Balance        float32   `json:"balance" gorm:"column:balance"`                // PayAmount - (AmountPerDay)*(howMany days gone)
	PackageStarts  time.Time `json:"packagestart" gorm:"column:package_starts"`    // when did he pay?
	PackageUpdate  time.Time `json:"packageupdate" gorm:"column:package_update"`   // when did we update his payment
	PackageEnds    time.Time `json:"packageend" gorm:"column:package_ends"`        // when package should expire?
	ExpiredIn      float64   `json:"expiredin" gorm:"column:expired_in"`           // Expired in how many days
	AlreadyExpired bool      `json:"alreadyexpired" gorm:"column:already_expired"` // Already expired?
	OverDueHours   float64   `json:"overdue" gorm:"column:over_due_hours"`         // For how many hours did he exceed payment threshhold
}

// =============== Test ==========================

type TestStruct struct {
	Message string `json:"message"`
	Text    string `json:"text"`
}
