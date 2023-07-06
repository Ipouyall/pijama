package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"saaj/api/data"
	"testing"
)

func TestAuthenticate(t *testing.T) {
	// Create a mocked server for testing
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		success := true
		// Check the request URL
		if r.URL.Path != LoginPath {
			t.Errorf("Unexpected URL: %s", r.URL.Path)
			success = false
		}

		// Check the request method
		if r.Method != "POST" {
			t.Errorf("Unexpected request method: %s", r.Method)
			success = false
		}

		// Check the request body
		var requestBody map[string]string
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}
		if username, ok := requestBody["username"]; !ok || username != "<test-username>" {
			success = false
		}
		if password, ok := requestBody["password"]; !ok || password != "<test-password>" {
			success = false
		}

		// Send a mocked response based on the test case
		if success {
			w.WriteHeader(http.StatusOK)
			response := map[string]string{
				"token": "<auth-token>",
			}
			_ = json.NewEncoder(w).Encode(response)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			response := map[string]string{
				"error": "Invalid credentials",
			}
			_ = json.NewEncoder(w).Encode(response)
		}
	}))
	defer mockServer.Close()

	// Create an instance of the REST class and set the domain to the mocked server URL
	rest := NewREST(mockServer.URL)

	// Test case 1: Successful authentication
	err, prompt := rest.Authenticate("<test-username>", "<test-password>")
	if err != nil {
		t.Errorf("Expected successful authentication, got error: %v", err)
	}
	if prompt != "You logged in successfully" {
		t.Errorf("Expected prompt to be empty, got '%s'", prompt)
	}
	if rest.Token != "<auth-token>" {
		t.Errorf("Unexpected token value: %s", rest.Token)
	}

	rest = NewREST(mockServer.URL)
	// Test case 2: Authentication failure
	err, prompt = rest.Authenticate("<test-username-wrong>", "<test-password-wrong>")
	if err == nil {
		t.Errorf("Expected failed authentication, got no error")
	}
	expectedPrompt := "Invalid credentials"
	if prompt != expectedPrompt {
		t.Errorf("Expected prompt to be '%s', got '%s'", expectedPrompt, prompt)
	}
	if rest.Token != "" {
		t.Errorf("Expected token to be empty, got '%s'", rest.Token)
	}
}

func TestGetPackage(t *testing.T) {
	// Create a mocked server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the request details
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != PackagesPath {
			t.Errorf("Expected URL /api/v1/packages, got %s", r.URL.Path)
		}

		// Prepare a sample response
		response := []map[string]any{
			{
				"id":             1,
				"package_name":   "<disease-name>",
				"category":       "<category>",
				"description":    "<description>",
				"estimated_cost": 1200,
				"city":           "<city-name>",
				"doctor":         "<doctor-name>",
				"hospital":       "<hospital-name>",
				"package_class":  "<p-class>",
			},
		}

		// Send the response
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Set the REST domain to the mocked server's URL
	R := NewREST(server.URL)

	// Call the GetPackage method
	packages := R.GetPackage()

	// Check the returned packages
	expectedPackages := []data.Package{
		{
			ID:           1,
			Name:         "<disease-name>",
			Category:     "<category>",
			PDescription: "<description>",
			Cost:         1200,
			City:         "<city-name>",
			Doctor:       "<doctor-name>",
			Hospital:     "<hospital-name>",
			Class:        "<p-class>",
		},
		// Add more expected packages if needed
	}

	if !reflect.DeepEqual(packages, expectedPackages) {
		t.Errorf("Unexpected packages: got %+v, expected %+v", packages, expectedPackages)
	}
}

func TestRequestPackage(t *testing.T) {
	// Create a mocked server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the request details
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != PackageReqPath {
			t.Errorf("Expected URL %s, got %s", PackageReqPath, r.URL.Path)
		}

		// Check the request body
		expectedBody := map[string]interface{}{
			"id":    1,
			"token": "<auth-token>",
		}
		var requestBody map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}
		if id, ok := requestBody["id"].(float64); !ok || int(id) != expectedBody["id"] {
			t.Errorf("Unexpected request body field 'id': %v", id)
		}
		if token, ok := requestBody["token"].(string); !ok || token != expectedBody["token"] {
			t.Errorf("Unexpected request body field 'token': %v", token)
		}

		// Prepare a sample response
		response := []map[string]interface{}{
			{
				"id":          1,
				"name":        "<dummy1>",
				"description": "<dummy2>",
			},
			{
				"id":          2,
				"name":        "<dummy3>",
				"description": "<dummy4>",
			},
		}

		// Send the response
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	// Create an instance of the REST class and set the domain to the mocked server URL
	R := &REST{
		Domain: mockServer.URL,
		Token:  "<auth-token>",
	}

	// Call the RequestPackage method
	requirements := R.RequestPackage(data.Package{ID: 1})

	// Check the returned requirements
	expectedRequirements := []data.Requirement{
		{
			ID:          1,
			Name:        "<dummy1>",
			Description: "<dummy2>",
		},
		{
			ID:          2,
			Name:        "<dummy3>",
			Description: "<dummy4>",
		},
	}

	if !reflect.DeepEqual(requirements, expectedRequirements) {
		t.Errorf("Unexpected requirements: got %+v, expected %+v", requirements, expectedRequirements)
	}
}

func TestSubmitDocuments(t *testing.T) {
	// Create a mocked server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the request details
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != UploadDocsPath {
			t.Errorf("Expected URL %s, got %s", UploadDocsPath, r.URL.Path)
		}

		// Read the request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Failed to read request body: %v", err)
		}

		type docs struct {
			Id      int    `json:"related_req_id"`
			Name    string `json:"title"`
			Content string `json:"content"`
		}

		// Unmarshal the request body
		var requestBody struct {
			Token     string `json:"token"`
			PackID    int    `json:"pid"`
			Documents []docs `json:"documents"`
		}
		err = json.Unmarshal(body, &requestBody)
		if err != nil {
			t.Errorf("Failed to unmarshal request body: %v", err)
		}

		// Check the request body fields
		expectedToken := "<auth-token>"
		if requestBody.Token != expectedToken {
			t.Errorf("Unexpected request body field 'token': got %s, expected %s", requestBody.Token, expectedToken)
		}

		expectedPackID := 1
		if requestBody.PackID != expectedPackID {
			t.Errorf("Unexpected request body field 'pid': got %d, expected %d", requestBody.PackID, expectedPackID)
		}

		expectedDocs := []data.Document{
			{
				ID:       123,
				Filename: "doc1.txt",
				Content:  "Document 1 content",
			},
		}
		if len(requestBody.Documents) != len(expectedDocs) {
			t.Errorf("Unexpected number of documents: got %d, expected %d", len(requestBody.Documents), len(expectedDocs))
		}
		for i, doc := range requestBody.Documents {
			if doc.Id != expectedDocs[i].ID {
				t.Errorf("Unexpected document ID: got %d, expected %d", doc.Id, expectedDocs[i].ID)
			}
			if doc.Name != expectedDocs[i].Filename {
				t.Errorf("Unexpected document name: got %s, expected %s", doc.Name, expectedDocs[i].Filename)
			}
			if doc.Content != expectedDocs[i].Content {
				t.Errorf("Unexpected document content: got %s, expected %s", doc.Content, expectedDocs[i].Content)
			}
		}
		// Prepare a sample response
		response := struct {
			TreatmentRequestID int `json:"tr_id"`
		}{
			TreatmentRequestID: 12345,
		}

		// Send the response
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	// Create an instance of the REST class and set the domain to the mocked server URL
	R := &REST{
		Domain: mockServer.URL,
		Token:  "<auth-token>",
	}

	// Prepare the test data
	packID := 1
	docs := []data.Document{
		{
			ID:       123,
			Filename: "doc1.txt",
			Content:  "Document 1 content",
		},
	}

	// Call the SubmitDocuments method
	err := R.SubmitDocuments(packID, docs, "Treat")

	// Check for errors
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestGetHotels(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Verify the request method
		if req.Method != "GET" {
			t.Errorf("Unexpected request method: got %s, expected GET", req.Method)
		}

		// Verify the request body
		expectedBody := map[string]interface{}{
			"package_id": 11,
			"token":      "<auth-token>",
			"city":       "<package-city>",
		}
		var requestBody map[string]interface{}
		err := json.NewDecoder(req.Body).Decode(&requestBody)
		if err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		//if requestBody["package_id"] != expectedBody["package_id"] {
		//	t.Errorf("Unexpected request body field 'package_id'")
		//}
		if requestBody["token"] != expectedBody["token"] {
			t.Errorf("Unexpected request body field 'token': got %v, expected %v", requestBody["token"], expectedBody["token"])
		}
		if requestBody["city"] != expectedBody["city"] {
			t.Errorf("Unexpected request body field 'city': got %v, expected %v", requestBody["city"], expectedBody["city"])
		}

		// Prepare the mock response
		response := []map[string]interface{}{
			{
				"id":          1,
				"hotel_name":  "<hotel-name>",
				"hotel_class": "<golden>",
				"cost":        1200,
				"city":        "<city-name>",
				"address":     "<street, ave>",
			},
			// Add more hotel rooms if needed
		}
		responseBody, err := json.Marshal(response)
		if err != nil {
			t.Fatalf("Failed to marshal response body: %v", err)
		}

		// Set the response status code and body
		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write(responseBody)
	}))
	defer server.Close()

	// Create a REST instance using the mock server's URL
	R := &REST{
		Domain: server.URL,
		Token:  "<auth-token>",
		TreatmentPackage: data.Package{
			ID:   11,
			City: "<package-city>",
		},
	}

	// Call the GetHotels method
	rooms := R.GetHotels()

	// Verify the response
	expectedRooms := []map[string]interface{}{
		{
			"id":          1,
			"hotel_name":  "<hotel-name>",
			"hotel_class": "<golden>",
			"cost":        1200,
			"city":        "<city-name>",
			"address":     "<street, ave>",
		},
		// Add more expected hotel rooms if needed
	}
	if len(rooms) != len(expectedRooms) {
		t.Fatalf("Unexpected number of hotel rooms: got %d, expected %d", len(rooms), len(expectedRooms))
	}
	for i, room := range expectedRooms {
		if room["id"] != rooms[i].ID {
			t.Errorf("Unexpected hotel room ID: got %v, expected %v", rooms[i].ID, room["id"])
		}
		if room["hotel_name"] != rooms[i].HotelName {
			t.Errorf("Unexpected hotel room name: got %v, expected %v", rooms[i].HotelName, room["hotel_name"])
		}
		if room["hotel_class"] != rooms[i].HotelClass {
			t.Errorf("Unexpected hotel room class: got %v, expected %v", rooms[i].HotelClass, room["hotel_class"])
		}
		if room["cost"] != rooms[i].Cost {
			t.Errorf("Unexpected hotel room cost: got %v, expected %v", rooms[i].Cost, room["cost"])
		}
		if room["city"] != rooms[i].City {
			t.Errorf("Unexpected hotel room city: got %v, expected %v", rooms[i].City, room["city"])
		}
		if room["address"] != rooms[i].Address {
			t.Errorf("Unexpected hotel room address: got %v, expected %v", rooms[i].Address, room["address"])
		}
	}
}
