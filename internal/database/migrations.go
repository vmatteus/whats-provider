package database

import (
	"github.com/your-org/boilerplate-go/internal/user/domain"
	"github.com/your-org/boilerplate-go/internal/whatsapp/infrastructure"
	"gorm.io/gorm"
)

// MigrateWithUsers runs migrations including user table
func MigrateWithUsers(db *gorm.DB) error {
	return db.AutoMigrate(&domain.User{})
}

// MigrateWhatsApp runs migrations for WhatsApp tables
func MigrateWhatsApp(db *gorm.DB) error {
	return db.AutoMigrate(
		&infrastructure.GormInstance{},
		&infrastructure.GormMessage{},
	)
}

// MigrateAll runs all migrations
func MigrateAll(db *gorm.DB) error {
	if err := MigrateWithUsers(db); err != nil {
		return err
	}

	if err := MigrateWhatsApp(db); err != nil {
		return err
	}

	return nil
}

// Seed seeds the database with initial data
func Seed(db *gorm.DB) error {
	// Add your seed logic here
	// Example: Create default admin user, default settings, etc.

	// Check if users table exists and is empty
	var count int64
	if err := db.Model(&domain.User{}).Count(&count).Error; err != nil {
		return err
	}

	// Seed admin user if no users exist
	if count == 0 {
		adminUser := &domain.User{
			Name:  "Admin User",
			Email: "admin@example.com",
		}

		if err := db.Create(adminUser).Error; err != nil {
			return err
		}
	}

	return nil
}
