package goucherdb

import (
	_ "database/sql" // sql database
	"fmt"
	"os"

	// I don't actually need these stuffs, check Makefile for proxy
	// _ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres" // GCP cloud-proxy
	// _ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql" // mysql
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres" // gorm postgres
	// _ "github.com/jinzhu/gorm/dialects/mysql" // mysql
)

// Db is gorm.DB object
var Db *gorm.DB

// ConnectDB connects to DB
func ConnectDB() error {

	// POSGRES
	connectionString := fmt.Sprintf(
		"host=/cloudsql/%s user=%s dbname=%s password=%s",
		os.Getenv("INSTANCE_CONNECTION_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASS"),
	)

	// MySQL
	// connectionString := fmt.Sprintf(
	// 	"%s:%s@cloudsql(%s)/%s?charset=utf8&parseTime=True&loc=UTC",
	// 	os.Getenv("DB_USER"),
	// 	os.Getenv("DB_PASS"),
	// 	os.Getenv("INSTANCE_CONNECTION_NAME"),
	// 	os.Getenv("DB_NAME"),
	// )

	var err error
	Db, err = gorm.Open("postgres", connectionString)

	return err
}

// MigrateDB makes sure the tables are well created :)
func MigrateDB() {
	MigrateRecipient()
	MigrateOffer()
	MigrateVoucher()
}
