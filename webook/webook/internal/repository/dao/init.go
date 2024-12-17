package dao

import "gorm.io/gorm"

func InitTableUser(db *gorm.DB) {
	err := db.AutoMigrate(User{})
	if err != nil {
		panic(err)
	}
}
