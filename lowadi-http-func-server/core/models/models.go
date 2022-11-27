package models

import (
	"time"

	"gorm.io/gorm"
)

type OrkiParameters struct {
	Account *Account
	Links   []string
}

// Game account info
type Account struct {
	ID               uint           `json:"id" gorm:"primarykey"`                          // Id
	UserID           uint           `json:"user_id"  gorm:"column:user_id"`                // userID
	Nickname         string         `json:"nick_name" gorm:"column:nick_name"`             // manually entered nickname of the acc
	Login            string         `json:"login" gorm:"unique"`                           // Login from the acc
	Pass             string         `json:"pass" gorm:"column:pass"`                       // Password
	VerifiedNickname string         `json:"verified_nick" gorm:"column:verified_nick"`     // Parsed nickname from the game (auto)
	Active           int            `json:"active" gorm:"column:active"`                   // 1=expired, 2=active Manual turn off/on (1 day trial by default)
	AccCreated       time.Time      `json:"acc_created" gorm:"column:acc_created"`         //
	AccDeleted       gorm.DeletedAt `json:"acc_deleted" gorm:"column:acc_deleted"`         // SOFT DELETION
	AccEnds          time.Time      `json:"acc_ends" gorm:"column:acc_ends"`               //
	AccUpdated       time.Time      `json:"acc_udpated" gorm:"column:acc_udpated"`         //
	AlreadyExpired   int            `json:"already_expired" gorm:"column:already_expired"` // 1=expired, 2=active  Auto turn off (balance exceeded)
	OverDueHours     float32        `json:"overdue_hours" gorm:"column:overdue_hours"`     // How many hours from the last active time
	ProxyIP          string         `json:"proxy_ip" gorm:"column:proxy_ip"`               // Add proxy?
	ProxyPort        int            `json:"proxy_port" gorm:"column:proxy_port"`           // Add proxy?
	ProxyLogin       string         `json:"proxy_login" gorm:"column:proxy_login"`         // Add proxy?
	ProxyPass        string         `json:"proxy_pass" gorm:"column:proxy_pass"`           // Add proxy?
	Timer            int32          `json:"timer" gorm:"column:timer"`                     // Used to increase AccEnds time
	CostPerDay       int32          `json:"cost_per_day" gorm:"column:cost_per_day"`       // How much $ per day does 1 acc costs
	GameWebside      string         `json:"game_website" gorm:"column:game_website"`       // International users use not lowadi.com, nl.horse.com etc.
	LoggedIn         int            `json:"logged_in" gorm:"column:logged_in"`             // Status: 1 - The bot is running, 2 - User can enter

	// Game accounts start up settings. Which modules will run?
	// 1 = true, 2 = false
	UpdateKCK    int `json:"update_kck" gorm:"column:update_kck"`
	RunOrki      int `json:"run_orki" gorm:"column:run_orki"`
	Competitions int `json:"competitions" gorm:"column:competitions"`

	// KCK - Default config to RunOrki
	// 1 = true, 2 = false
	StableName       string `json:"stable_name" gorm:"column:stable_name"`             // Using this to run orki in this KCK
	StableLink       string `json:"stable_link" gorm:"column:stable_link"`             //
	AdvantagesFuraj  int    `json:"advantages_furaj" gorm:"column:advantages_furaj"`   // 1=use this check-mark. 2=don't
	AdvantagesOvec   int    `json:"advantages_ovec" gorm:"column:advantages_ovec"`     // 1=use this check-mark. 2=don't
	AdvantagesCarrot int    `json:"advantages_carrot" gorm:"column:advantages_carrot"` // 1=use this check-mark. 2=don't
	MaxDailyPrice    int    `json:"max_daily_price" gorm:"column:max_daily_price"`     // price per 1 day in KCK
	BirthHorses      int    `json:"birth_horses" gorm:"column:birth_horses"`           // 1=allow, 2=pass this step
	BirthHorsesName  string `json:"birth_horses_name" gorm:"column:birth_horses_name"` // How we should name the horse?

	// Competitions section
	CompetitionsLink  string `json:"competitions_link" gorm:"column:competitions_link"`   // The link of a horse we should upgrade
	CompetitionsEvent int    `json:"competitions_event" gorm:"column:competitions_event"` // 1-6 Which competition we should start?
}

// Stores KCK list
type KCK struct {
	gorm.Model
	UserID     uint   `json:"user_id"  gorm:"column:user_id"`        // userID
	Login      string `json:"login" gorm:"login"`                    // Login from the acc
	StableLink string `json:"stable_link" gorm:"column:stable_link"` // Stores account's stable link
	StableName string `json:"stable_name" gorm:"column:stable_name"` // Stores account's stable name
}
