package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"saaj/api/data"
)

const (
	Domain          = "http://127.0.0.1:8000"
	LoginPath       = "/api/v1/login"
	PackagesPath    = "/api/v1/packages"
	PackageReqPath  = "/api/v1/package_requirements"
	UploadDocsPath  = "/api/v1/upload_user_docs"
	HotelsPath      = "/api/v1/hotels"
	VisaRequestPath = "/api/v1/handle_visa_request"
	VisaStatusPath  = "/api/v1/visa_status"
	BillingPath     = "/api/v1/handle_payment_bill_request"
	LogoutPath      = "/api/v1/logout"
)

func NewREST(domain string) *REST {
	return &REST{Domain: domain}
}

type REST struct {
	Domain             string
	Token              string
	TreatmentPackage   data.Package
	TreatmentPackageID int
	VisaSerial         string
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

func (R *REST) RequestPackage(pack data.Package) (requirements []data.Requirement) {
	// Prepare the request body
	requestBody := map[string]interface{}{
		"id":    pack.ID,
		"token": R.Token,
	}
	requestBodyBytes, _ := json.Marshal(requestBody)

	// Create the HTTP request
	url := R.Domain + PackageReqPath
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		log.Fatal(err)
	}
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
	R.TreatmentPackage = pack
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

func (R *REST) SubmitDocuments(packID int, docs []data.Document, dKind string) (err error, bill data.Bill) {
	// Prepare the request body
	requestBody := map[string]interface{}{
		"token":     R.Token,
		"pid":       packID,
		"documents": docs,
		"tr_id":     R.TreatmentPackageID,
	}
	requestBodyBytes, _ := json.Marshal(requestBody)

	fmt.Println(requestBody)

	// Create the HTTP request
	url := R.Domain + UploadDocsPath
	if dKind == "Visa" {
		url = R.Domain + VisaRequestPath
	}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	fmt.Println(requestBody, url)

	// Send the request and handle the response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var resppp map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&resppp)
	fmt.Println(resppp)

	// Check the response status code
// 	if resp.StatusCode != http.StatusOK {
// 		// Request failed
// 		return fmt.Errorf("request failed with status code: %d", resp.StatusCode), bill
// 	}

	// Request succeeded
	var response struct {
		TreatmentRequestID int    `json:"tr_id"`
		VisaSerial         string `json:"serial_no"`
		PaymentID          int    `json:"payment_request_id"`
		Cost               int    `json:"total_cost"`
	}
	// Parse the response body
	err = json.NewDecoder(resp.Body).Decode(&response)
	fmt.Println(response)
// 	if err != nil {
// 		return
// 	}
	fmt.Println(dKind)
	if dKind == "Treat" {
		R.TreatmentPackageID = response.TreatmentRequestID
	}
	if dKind == "Visa" {
		R.VisaSerial = response.VisaSerial
		bill = data.Bill{
			Title:     "Visa",
			Cost:      response.Cost,
			PaymentID: response.PaymentID,
		}
	}
	return
}

func (R *REST) GetHotels() (rooms []data.HotelRoom) {
	// Prepare the request body
	requestBody := map[string]interface{}{
		"package_id": R.TreatmentPackage.ID,
		"token":      R.Token,
		"city":       R.TreatmentPackage.City,
		"tr_id":      R.TreatmentPackageID,
	}
	requestBodyBytes, _ := json.Marshal(requestBody)

	fmt.Println(requestBody)

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

	var resppp map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&resppp)
	fmt.Println(resppp)

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
		"tr_id":      9,
		"hotel_id":   hotelID,
		"token":      R.Token,
	}
	requestBodyBytes, _ := json.Marshal(requestBody)
	url := R.Domain + HotelsPath
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	fmt.Println(requestBody)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var resppp map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&resppp)
	fmt.Println(resppp)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}
	return nil
}

func (R *REST) RequestVisa() (requirements []data.Requirement) {
	// Prepare the request body
	requestBody := map[string]interface{}{
		"token": R.Token,
		"tr_id": R.TreatmentPackageID,
	}
	requestBodyBytes, _ := json.Marshal(requestBody)

	fmt.Println(requestBody)

	// Create the HTTP request
	url := R.Domain + VisaRequestPath
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request and handle the response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var resppp map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&resppp)
	fmt.Println(resppp)

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return
	}
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

func (R *REST) VisaStatus() (vss []data.VisaStatus) {
	requestBody := map[string]interface{}{
		"token":     R.Token,
		"serial_no": R.VisaSerial,
		"tr_id":     R.TreatmentPackageID,
	}
	requestBodyBytes, _ := json.Marshal(requestBody)

	// Create the HTTP request
	url := R.Domain + VisaStatusPath
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request and handle the response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	_ = json.NewDecoder(resp.Body).Decode(&vss)

	return
}

func (R *REST) GetBill() (bill data.Bill) {
	requestBody := map[string]interface{}{
		"token":     R.Token,
		"serial_no": R.VisaSerial,
		"tr_id":     R.TreatmentPackageID,
	}
	requestBodyBytes, _ := json.Marshal(requestBody)

	url := R.Domain + BillingPath
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	_ = json.NewDecoder(resp.Body).Decode(&bill)
	bill.Title = "Treat-" + R.TreatmentPackage.Name
	return
}

func (R *REST) PayBill(billID int, code string) error {
	return nil
}

func (R *REST) Logout() (err error) {
	requestBody := map[string]interface{}{
		"token": R.Token,
	}
	requestBodyBytes, _ := json.Marshal(requestBody)

	url := R.Domain + LogoutPath
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("failed to logging out")
	}

	return
}
