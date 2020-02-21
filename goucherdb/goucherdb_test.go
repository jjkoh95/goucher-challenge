package goucherdb_test

import (
	"os"
	"testing"
	"time"

	"github.com/jjkoh95/goucher-challenge/goucherdb"
)

var seedRecipient *goucherdb.Recipient = &goucherdb.Recipient{Email: "jjkoh95@gmail.com", Name: "Koh Jia Jun"}
var seedRecipient2 *goucherdb.Recipient = &goucherdb.Recipient{Email: "jjkoh@gmail.com", Name: "Koh Jia Jun"}
var seedOffer *goucherdb.Offer = &goucherdb.Offer{Name: "50% off", DiscountPercentage: 50}
var seedOffer2 *goucherdb.Offer = &goucherdb.Offer{Name: "30% off", DiscountPercentage: 30}
var seedVoucher *goucherdb.Voucher

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setup() {
	connect()
	migrate()
}

func connect() {
	err := goucherdb.ConnectDB()
	if err != nil {
		panic("Unable to connect to database")
	}
}

func migrate() {
	goucherdb.MigrateRecipient()
	goucherdb.MigrateOffer()
	goucherdb.MigrateVoucher()
}

func TestCreateRecipient(t *testing.T) {
	err := goucherdb.CreateRecipient(seedRecipient)
	if err != nil {
		t.Error("Expected to create Recipient seeders without error")
	}
}

func TestGetAllActiveRecipients(t *testing.T) {
	rs, err := goucherdb.GetAllRecipients()
	if err != nil {
		t.Error("Expected to get recipients without error")
	}

	if len(rs) < 1 {
		t.Error("Expected to get at least one recipient")
	}
}

func TestCreateOffer(t *testing.T) {
	err := goucherdb.CreateOffer(seedOffer)
	if err != nil {
		t.Error("Expected to create Offer seeders without error")
	}
}

func TestCreateVoucher(t *testing.T) {
	seedVoucher = &goucherdb.Voucher{
		RecipientID: seedRecipient.RecipientID,
		OfferID:     seedOffer.OfferID,
		ExpiredAt:   time.Now().AddDate(0, 0, 7),
	}
	err := goucherdb.CreateVoucher(seedVoucher)
	if err != nil {
		t.Error("Expected to create Voucher seeders without error")
	}
}

func TestGetRecipientPreload(t *testing.T) {
	r, err := goucherdb.GetRecipientWithValidVouchers("jjkoh95@gmail.com")
	if err != nil {
		t.Error("Expected to get recipient preload without errors")
	}

	// expect to have nested voucher
	if len(r.Vouchers) < 1 {
		t.Error("Expected to have nested vouchers attached")
	}

	// b, _ := json.Marshal(r)
	// fmt.Println(string(b))
}

func TestMassCreateVouchers(t *testing.T) {
	// create a separate dummy recipient to test "mass"
	err := goucherdb.CreateRecipient(seedRecipient2)
	if err != nil {
		t.Error("Expected to create Recipient seeders without error")
	}

	// create a dedicated seed offer for this
	err = goucherdb.CreateOffer(seedOffer2)
	if err != nil {
		t.Error("Expected to create Offer seeders without error")
	}

	expiredAt := time.Now().AddDate(0, 0, 7)
	err = goucherdb.MassCreateVouchers(seedOffer2, expiredAt)
	if err != nil {
		t.Error("Expect to mass create vouchers without error")
	}
}

func TestRedeemVoucher(t *testing.T) {
	err := goucherdb.RedeemVoucher(seedVoucher)
	if err != nil {
		t.Error("Expected to redeem voucher without error")
	}

	if seedVoucher.RedeemedAt == nil {
		t.Error("Expected to update expiredAt")
	}
}

func TestGetActiveVoucherByCode(t *testing.T) {
	voucher, err := goucherdb.GetActiveVoucherByCode(seedVoucher.RecipientID, seedVoucher.Code)
	if err == nil {
		t.Error("Expected to throw error given invalid voucher")
	}
	if voucher.Code != "" {
		t.Error("Expected to not find any voucher")
	}
}

func TestGetRecipientByEmail(t *testing.T) {
	r, err := goucherdb.GetRecipientByEmail(seedRecipient.Email)
	if err != nil {
		t.Error("Expected to get recipient without error")
	}
	if r.RecipientID == 0 {
		t.Error("Expected to get recipient given valid email")
	}

	r, err = goucherdb.GetRecipientByEmail("iuaghsdiua")
	if err == nil {
		t.Error("Expected to get recipient with error given invalid email")
	}
	if r.RecipientID != 0 {
		t.Error("Expected to not get recipient given invalid email")
	}
}

func TestGetOffer(t *testing.T) {
	offer, err := goucherdb.GetOffer(seedOffer.OfferID)
	if err != nil {
		t.Error("Expected to get offer without error")
	}

	if offer.Name != seedOffer.Name {
		t.Error("Expected to get the same offer name")
	}
}

func tearDown() {
	// Unscoped to hard delete
	// delete all
	// goucherdb.Db.Unscoped().Delete(&goucherdb.Voucher{})
	// goucherdb.Db.Unscoped().Delete(&goucherdb.Offer{})
	// goucherdb.Db.Unscoped().Delete(&goucherdb.Recipient{})
	goucherdb.Db.Close()
}
