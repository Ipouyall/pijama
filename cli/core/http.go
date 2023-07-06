package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"saaj/core/data"
)

const (
	Domain         = "127.0.0.1:8000"
	LoginPath      = "/api/v1/login"
	PackagesPath   = "/api/v1/packages"
	PackageReqPath = "/api/v1/package_requirements"
	UploadDocsPath = "/api/v1/upload_user_docs"
	HotelsPath     = "/api/v1/hotels"
)

func NewREST(domain string) *REST {
	return &REST{Domain: domain}
}

type REST struct {
	Domain             string
	Token              string
	TreatmentPackage   data.Package
	TreatmentPackageID int
}

func (R *REST) Authenticate(username, password string) (err error, prompt string) {
	if R.Token != "" {
		// Already authenticated
		prompt = "Already authenticated."
		return
	}
	// Prepare the request body
	requestBody := map[string]string{
		"username": username,
		"password": password,
	}
	requestBodyBytes, _ := json.Marshal(requestBody)
	prompt = "You logged in successfully"

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

		// Parse the response body
		var responseBody map[string]string
		_ = json.NewDecoder(resp.Body).Decode(&responseBody)

		// Extract the token
		R.Token = responseBody["token"]

		return
	}
	if resp.StatusCode == http.StatusUnauthorized {
		// Authentication failed
		err = fmt.Errorf("invalid credential")

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

func (R *REST) RequestPackage(packID int) (requirements []data.Requirement) {
	// Prepare the request body
	requestBody := map[string]interface{}{
		"id":    packID,
		"token": R.Token,
	}
	requestBodyBytes, _ := json.Marshal(requestBody)

	// Create the HTTP request
	url := R.Domain + PackageReqPath
	req, _ := http.NewRequest("GET", url, bytes.NewBuffer(requestBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Send the request and handle the response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return
	}
	// Parse the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &requirements)
	if err != nil {
		return
	}

	return
}

func (R *REST) SubmitDocuments(packID int, docs []data.Document) (err error) {
	// Prepare the request body
	requestBody := map[string]interface{}{
		"token":     R.Token,
		"pid":       packID,
		"documents": docs,
	}
	requestBodyBytes, _ := json.Marshal(requestBody)

	// Create the HTTP request
	url := R.Domain + UploadDocsPath
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Send the request and handle the response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		// Request failed
		return fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	// Request succeeded
	var response struct {
		TreatmentRequestID int `json:"tr_id"`
	}
	// Parse the response body
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return
	}
	R.TreatmentPackageID = response.TreatmentRequestID

	return
}

func (R *REST) GetHotels() (rooms []data.HotelRoom) {
	// Prepare the request body
	requestBody := map[string]interface{}{
		"package_id": R.TreatmentPackage.ID,
		"token":      R.Token,
		"city":       R.TreatmentPackage.City,
	}
	requestBodyBytes, _ := json.Marshal(requestBody)

	// Create the HTTP request
	url := R.Domain + HotelsPath
	req, _ := http.NewRequest("GET", url, bytes.NewBuffer(requestBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Send the request and handle the response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return
	}

	// Parse the response body
	_ = json.NewDecoder(resp.Body).Decode(&rooms)
	return
}

func (R *REST) ReserveHotel(hotelID int) error {
	requestBody := map[string]interface{}{
		"package_id": R.TreatmentPackage.ID,
		"tr_id":      R.TreatmentPackageID,
		"hotel_id":   hotelID,
		"token":      R.Token,
	}
	requestBodyBytes, _ := json.Marshal(requestBody)
	url := R.Domain + HotelsPath
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}
	return nil
}

func (R *REST) RequestVisa() []data.Requirement {
	//TODO implement me
	panic("implement me")
}

func (R *REST) SubmitVisa(visaID int) error {
	//TODO implement me
	panic("implement me")
}

func (R *REST) GetBill() data.Bill {
	//TODO implement me
	panic("implement me")
}

func (R *REST) PayBill(billID int, code string) error {
	//TODO implement me
	panic("implement me")
}
