package initializers

import "github.com/Aoiewrug/Lowadi-App/lowadi-http-api/core/models"

// Sync SQL tables
func SyncDB() {
	DB.AutoMigrate(
		&models.User{},
		&models.Account{},
		&models.KCK{},
	)
}
