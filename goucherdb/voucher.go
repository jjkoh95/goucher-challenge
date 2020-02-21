package goucherdb

import (
	"time"
)

// Voucher ORM
type Voucher struct {
	RecipientID uint `gorm:"primary_key" json:"recipientID"`
	OfferID     uint `gorm:"primary_key" json:"offerID"`
	// https://github.com/jinzhu/gorm/issues/320
	// uuid_generate_v4() needs to be turned on manually
	// CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	Code         string     `gorm:"default:LEFT(uuid_generate_v4()::TEXT, 8)" json:"code"`
	SpecialOffer Offer      `gorm:"foreignkey:OfferID;association_foreignkey:OfferID" json:"specialOffer"`
	CreatedAt    time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"createdAt"`
	ExpiredAt    time.Time  `gorm:"notnull" json:"expiredAt"`
	RedeemedAt   *time.Time `json:"redeeemedAt"`
}

// MigrateVoucher is the migration script to generate Voucher table
func MigrateVoucher() {
	Db.AutoMigrate(&Voucher{})
	Db.Model(&Voucher{}).AddForeignKey("recipient_id", "recipients(recipient_id)", "RESTRICT", "RESTRICT")
	Db.Model(&Voucher{}).AddForeignKey("offer_id", "offers(offer_id)", "RESTRICT", "RESTRICT")
}

// CreateVoucher creates a new entry for voucher
func CreateVoucher(voucher *Voucher) error {
	return Db.Create(voucher).Error
}

// MassCreateVouchers credits all (existing) recipients with the special offer
func MassCreateVouchers(offer *Offer, expiredAt time.Time) error {
	// https://stackoverflow.com/questions/4101739/sql-server-select-into-existing-table
	sqlQuery := `
	INSERT INTO vouchers (recipient_id, offer_id, expired_at)
	SELECT recipient_id, ? as offer_id, ? as expired_at FROM recipients;
	`
	return Db.Exec(sqlQuery, offer.OfferID, expiredAt).Error
}

// RedeemVoucher should update the redeemed_at column
func RedeemVoucher(voucher *Voucher) error {
	return Db.Model(voucher).Update("redeemed_at", time.Now()).Error
}

// GetActiveVoucherByCode finds the active voucher by code
func GetActiveVoucherByCode(recipientID uint, code string) (Voucher, error) {
	var voucher Voucher
	err := Db.Where("code = ? AND expired_at > ? AND redeemed_at IS NULL", code, time.Now()).Find(&voucher).Error
	return voucher, err
}

// IsVoucherValid validates the voucher
func IsVoucherValid(email, code string) bool {
	recipient, err := GetRecipientByEmail(email)
	if err != nil {
		return false
	}

	_, err = GetActiveVoucherByCode(recipient.RecipientID, code)
	if err != nil {
		return false
	}

	return true
}
