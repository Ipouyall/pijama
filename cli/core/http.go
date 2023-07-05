package core

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"saaj/core/data"
)

const (
	Domain       = "127.0.0.1:8000"
	LoginPath    = "/api/v1/login"
	PackagesPath = "/api/v1/packages"
)

func NewREST(domain string) *REST {
	return &REST{Domain: domain}
}

type REST struct {
	Domain           string
	Token            string
	TreatmentPackage data.Package
}

func (R *REST) Authenticate(username, password string) (success bool, prompt string) {
	if R.Token != "" {
		// Already authenticated
		success = true
		prompt = "Already authenticated."
		return
	}
	// Prepare the request body
	requestBody := map[string]string{
		"username": username,
		"password": password,
	}
	requestBodyBytes, _ := json.Marshal(requestBody)

	// Create the HTTP request
	url := R.Domain + LoginPath
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Send the request and handle the response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		prompt = "An error occurred while making the request."
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode == http.StatusOK {
		// Authentication succeeded
		success = true

		// Parse the response body
		var responseBody map[string]string
		_ = json.NewDecoder(resp.Body).Decode(&responseBody)

		// Extract the token
		R.Token = responseBody["token"]

		return
	}
	if resp.StatusCode == http.StatusUnauthorized {
		// Authentication failed
		success = false

		// Parse the response body
		var responseBody map[string]string
		_ = json.NewDecoder(resp.Body).Decode(&responseBody)

		// Extract the error message
		prompt = responseBody["error"]

		return
	}
	// Unexpected response status code
	prompt = "An unexpected response was received."
	return
}

func (R *REST) GetPackage() []data.Package {
	// Send the GET request
	url := R.Domain + PackagesPath
	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return nil
	}

	// Parse the response body into a slice of Package
	var packages []data.Package
	err = json.Unmarshal(body, &packages)
	if err != nil {
		return nil
	}

	return packages
}

func (R *REST) RequestPackage(packID int) data.Requirement {
	//TODO implement me
	panic("implement me")
}

func (R *REST) SubmitDocument(docID int, name, content string) bool {
	//TODO implement me
	panic("implement me")
}

func (R *REST) GetHotels() []data.HotelRoom {
	//TODO implement me
	panic("implement me")
}

func (R *REST) ReserveHotel(hotelID int) bool {
	//TODO implement me
	panic("implement me")
}

func (R *REST) RequestVisa() bool {
	//TODO implement me
	panic("implement me")
}

func (R *REST) SubmitVisa(visaID int) bool {
	//TODO implement me
	panic("implement me")
}

func (R *REST) GetBill() data.Bill {
	//TODO implement me
	panic("implement me")
}

func (R *REST) PayBill(billID int, code string) bool {
	//TODO implement me
	panic("implement me")
}
