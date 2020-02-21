package goucherdb

import "time"

// Offer ORM
type Offer struct {
	OfferID            uint      `gorm:"primary_key;AUTO_INCREMENT" json:"offerID"`
	Name               string    `gorm:"notnull" json:"name"`
	DiscountPercentage uint      `gorm:"notnull" json:"discountPercentage"`
	CreatedAt          time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"createdAt"`
}

// MigrateOffer is the migration script to generate Offer table
func MigrateOffer() {
	Db.AutoMigrate(&Offer{})
}

// CreateOffer creates new entry Offer
func CreateOffer(offer *Offer) error {
	return Db.Create(offer).Error
}

// GetOffer finds the offer entry by offer_id
func GetOffer(offerID uint) (Offer, error) {
	var offer Offer
	err := Db.Where("offer_id = ?", offerID).Find(&offer).Error
	return offer, err
}
