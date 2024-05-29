package db

import (
	"github.com/manimovassagh/bank-go/types"
	"gorm.io/gorm"
)

func MigrateTables(db *gorm.DB) {
	// Auto-migrate the tables managed by GORM
	err := db.AutoMigrate(&types.Customer{}, &types.Account{}, &types.Transaction{})
	if err != nil {
		return
	}
}
