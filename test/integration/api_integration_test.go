package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GradiyantoS/go-dealls-test-app/repositories"
	"github.com/GradiyantoS/go-dealls-test-app/routes"
	"github.com/stretchr/testify/assert"
)

func TestAPIIntegration(t *testing.T) {
	// Initialize the repository and test wrapper
	baseRepo := repositories.NewUserRepository()
	testRepo := NewResettableTestRepository(baseRepo)

	// Seed test data initially
	testRepo.SeedTestData()

	// Set up the router with the test repository
	router := routes.SetupRouterWithRepo(testRepo.GetRepository())

	// Define test cases
	testCases := []struct {
		name               string
		method             string
		url                string
		body               interface{}
		headers            map[string]string
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:               "SignUp - Success",
			method:             "POST",
			url:                "/signup",
			body:               map[string]string{"email": "newuser@example.com", "password": "newpassword", "phone": "1112223333", "name": "New User", "gender": "male"},
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:               "Login - Success",
			method:             "POST",
			url:                "/login",
			body:               map[string]string{"identifier": "test1@example.com", "password": "password1"},
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reseed test data before each test
			testRepo.SeedTestData()

			// Create the request
			var req *http.Request
			if tc.body != nil {
				body, _ := json.Marshal(tc.body)
				req = httptest.NewRequest(tc.method, tc.url, bytes.NewReader(body))
			} else {
				req = httptest.NewRequest(tc.method, tc.url, nil)
			}

			// Set headers
			for key, value := range tc.headers {
				req.Header.Set(key, value)
			}

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(rr, req)

			// Assert status code
			assert.Equal(t, tc.expectedStatusCode, rr.Code)

			// Optional: Assert response body
			if tc.expectedResponse != "" {
				assert.Contains(t, rr.Body.String(), tc.expectedResponse)
			}
		})
	}
}
