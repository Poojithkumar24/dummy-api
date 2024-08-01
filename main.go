package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type SoftOffer struct {
	IsOfferFound  bool   `json:"is_offer_found"`
	AaFlag        bool   `json:"aa_flag"`
	OfferAmount   string `json:"offer_amount"`
	Tenure        string `json:"tenure"`
	Roi           string `json:"roi"`
	ProcessingFee string `json:"processing_fee"`
}

type Lender struct {
	LenderCode string    `json:"lenderCode"`
	SoftOffer  SoftOffer `json:"softOffer"`
}

type Data struct {
	CustomerId string   `json:"customerId"`
	Lenders    []Lender `json:"lenders"`
}

type Response struct {
	Status string  `json:"status"`
	Data   Data    `json:"data"`
	Errors *string `json:"errors"`
}

func dummyAPIHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Status: "Success",
		Data: Data{
			CustomerId: "3eaa8b97bd8efa6e9fd74506d11e1237",
			Lenders: []Lender{
				{
					LenderCode: "abfl",
					SoftOffer: SoftOffer{
						IsOfferFound:  false,
						AaFlag:        false,
						OfferAmount:   "",
						Tenure:        "",
						Roi:           "",
						ProcessingFee: "",
					},
				},
			},
		},
		Errors: nil,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

func main() {
	http.HandleFunc("/dummy-api", dummyAPIHandler)

	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
