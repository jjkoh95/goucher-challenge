package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/jjkoh95/goucher-challenge/goucherdb"
	"github.com/jjkoh95/goucher-challenge/goucherutils"
)

// return all recipients without preload (associations)
func getAllRecipient(w http.ResponseWriter, r *http.Request) {
	rs, err := goucherdb.GetAllRecipients()
	if err != nil {
		http.Error(w, "Unable to retrieve recipients", http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(rs)
	if err != nil {
		http.Error(w, "Unable to get recipients", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

type createRecipientPayload struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// create recipient to database
func createRecipient(w http.ResponseWriter, r *http.Request) {
	// decode payload
	var payload createRecipientPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	if !goucherutils.IsEmail(payload.Email) {
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}

	recipient := goucherdb.Recipient{
		Email: payload.Email,
		Name:  payload.Name,
	}
	err = goucherdb.CreateRecipient(&recipient)
	if err != nil {
		http.Error(w, "Unable to create recipient", http.StatusBadRequest)
		return
	}

	w.Write([]byte("Recipient created"))
}

// get recipient with all preload/association data
func getRecipient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if !goucherutils.IsEmail(vars["email"]) {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	recipient, err := goucherdb.GetRecipientWithValidVouchers(vars["email"])
	if err != nil || recipient.Email == "" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	responseBytes, _ := json.Marshal(recipient)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)
}

type createSpecialOfferPayload struct {
	OfferName          string    `json:"offerName"`
	DiscountPercentage uint      `json:"discountPercentage"`
	ExpiredAt          time.Time `json:"expiredAt"`
}

// create special offer
func createSpecialOffer(w http.ResponseWriter, r *http.Request) {
	// decode payload
	var payload createSpecialOfferPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	specialOffer := goucherdb.Offer{
		Name:               payload.OfferName,
		DiscountPercentage: payload.DiscountPercentage,
	}

	err = goucherdb.CreateOffer(&specialOffer)
	if err != nil {
		http.Error(w, "Unable to create offer", http.StatusBadRequest)
		return
	}

	// distribute to every recipient
	err = goucherdb.MassCreateVouchers(&specialOffer, payload.ExpiredAt)
	if err != nil {
		http.Error(w, "Error distributing vouchers", http.StatusBadRequest)
		return
	}

	w.Write([]byte("Special offer created"))
}

type redeemVoucherRequestPayload struct {
	Email       string `json:"email"`
	VoucherCode string `json:"voucherCode"`
}

type redeemVoucherResponsePayload struct {
	Message            string `json:"message"`
	OfferName          string `json:"offerName"`
	PercentageDiscount uint   `json:"percentageDiscount"`
}

// redeem voucher
func redeemVoucher(w http.ResponseWriter, r *http.Request) {
	// decode payload
	var payload redeemVoucherRequestPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	// find recipient_id
	// since our code only uses 8 chars
	// this helps to avoid collision of voucher code
	recipient, err := goucherdb.GetRecipientByEmail(payload.Email)
	if err != nil {
		http.Error(w, "Unable to get recipient", http.StatusNotFound)
		return
	}

	// get voucher
	voucher, err := goucherdb.GetActiveVoucherByCode(recipient.RecipientID, payload.VoucherCode)
	if err != nil {
		http.Error(w, "Invalid voucher code", http.StatusBadRequest)
		return
	}

	// get offer
	// this is to retrieve the percentage
	offer, err := goucherdb.GetOffer(voucher.OfferID)
	if err != nil {
		http.Error(w, "Offer not found", http.StatusInternalServerError)
	}

	// update voucher redeemed time
	err = goucherdb.RedeemVoucher(&voucher)
	if err != nil {
		http.Error(w, "Unable to redeem voucher", http.StatusInternalServerError)
		return
	}

	responsePayload := redeemVoucherResponsePayload{
		Message:            "Successfully redeemed",
		OfferName:          offer.Name,
		PercentageDiscount: offer.DiscountPercentage,
	}
	responseBytes, _ := json.Marshal(responsePayload)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)
}

// healthcheck
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Ok!"))
}

func handleRequest() {
	// create router
	router := mux.NewRouter().StrictSlash(true)

	// health check
	router.HandleFunc("/", healthCheck)
	router.HandleFunc("/health", healthCheck)

	router.HandleFunc("/recipient", getAllRecipient).Methods(http.MethodGet)
	router.HandleFunc("/recipient", createRecipient).Methods(http.MethodPost)
	router.HandleFunc("/recipient/{email}", getRecipient).Methods(http.MethodGet)

	router.HandleFunc("/offer", createSpecialOffer).Methods(http.MethodPost)

	router.HandleFunc("/voucher/redeem", redeemVoucher).Methods(http.MethodPost)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Println("Serving at port", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}

func main() {
	log.Println("Connecting to DB")
	goucherdb.ConnectDB()
	defer goucherdb.Db.Close()
	goucherdb.MigrateDB()

	log.Println("Creating server")
	handleRequest()
}
