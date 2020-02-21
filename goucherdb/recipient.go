package goucherdb

import "time"

// Recipient ORM
type Recipient struct {
	RecipientID uint      `gorm:"primary_key;AUTO_INCREMENT" json:"recipientID"`
	Email       string    `gorm:"unique;notnull" json:"email"`
	Name        string    `gorm:"notnull" json:"name"`
	Vouchers    []Voucher `gorm:"foreignkey:RecipientID" json:"vouchers,omitempty"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"createdAt"`
}

// MigrateRecipient is migration script to generate Recipient table
func MigrateRecipient() {
	Db.AutoMigrate(&Recipient{})
}

// CreateRecipient creates a new entry for recipient
func CreateRecipient(recipient *Recipient) error {
	return Db.Create(recipient).Error
}

// GetAllRecipients returns all the active recipients
func GetAllRecipients() ([]Recipient, error) {
	var recipients []Recipient
	err := Db.Find(&recipients).Error
	return recipients, err
}

// GetRecipientWithValidVouchers returns the recipient with all associations
func GetRecipientWithValidVouchers(email string) (Recipient, error) {
	var recipient Recipient
	err := Db.Where("email = ?", email).Preload("Vouchers", "expired_at > ? AND redeemed_at IS NULL", time.Now()).Preload("Vouchers.SpecialOffer").Find(&recipient).Error
	return recipient, err
}

// GetRecipientByEmail returns the recipient without associations
func GetRecipientByEmail(email string) (Recipient, error) {
	var recipient Recipient
	err := Db.Where("email = ?", email).Find(&recipient).Error
	return recipient, err
}
