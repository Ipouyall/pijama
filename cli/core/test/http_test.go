package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"saaj/core"
	"testing"
)

func TestAuthenticate(t *testing.T) {
	// Create a mocked server for testing
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		success := true
		// Check the request URL
		if r.URL.Path != "/api/v1/login" {
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
			json.NewEncoder(w).Encode(response)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			response := map[string]string{
				"error": "Invalid credentials",
			}
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer mockServer.Close()

	// Create an instance of the REST class and set the domain to the mocked server URL
	rest := core.NewREST(mockServer.URL)

	// Test case 1: Successful authentication
	success, prompt := rest.Authenticate("<test-username>", "<test-password>")
	if !success {
		t.Errorf("Expected success to be true, got false")
	}
	if prompt != "" {
		t.Errorf("Expected prompt to be empty, got '%s'", prompt)
	}
	if rest.Token != "<auth-token>" {
		t.Errorf("Unexpected token value: %s", rest.Token)
	}

	rest = core.NewREST(mockServer.URL)
	// Test case 2: Authentication failure
	success, prompt = rest.Authenticate("<test-username-wrong>", "<test-password-wrong>")
	if success {
		t.Errorf("Expected success to be false, got true")
	}
	expectedPrompt := "Invalid credentials"
	if prompt != expectedPrompt {
		t.Errorf("Expected prompt to be '%s', got '%s'", expectedPrompt, prompt)
	}
	if rest.Token != "" {
		t.Errorf("Expected token to be empty, got '%s'", rest.Token)
	}
}
