package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	samples := []struct {
		inputJSON  string
		statusCode int
		email      string
	}{
		{
			inputJSON:  `{"email": "djodi@example.com"}`,
			statusCode: 201,
			email:      "djodi@example.com",
		},
		{
			inputJSON:  `{}`,
			statusCode: 422,
		},
		{
			inputJSON:  `{"email": "djodi@example.com"}`,
			statusCode: 500,
		},
	}

	for _, v := range samples {
		r := gin.Default()
		r.POST("/api/user", server.CreateUser)
		req, err := http.NewRequest(http.MethodPost, "/api/user", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("Error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		responseInterface := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseInterface)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 201 {
			responseMap := responseInterface["response"].(map[string]interface{})
			assert.Equal(t, responseMap["email"], v.email)
		}

		if v.statusCode == 422 || v.statusCode == 500 {
			responseMap := responseInterface["error"].(map[string]interface{})

			if responseMap["email"] != nil {
				assert.Equal(t, responseMap["email"], responseMap["email"])
			}
		}
	}
}

func TestGetUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	_, err = seedUsers()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.GET("/users", server.GetUsers)

	req, err := http.NewRequest(http.MethodGet, "/users", nil)
	if err != nil {
		t.Errorf("Error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	usersMap := make(map[string]interface{})

	err = json.Unmarshal([]byte(rr.Body.String()), &usersMap)
	if err != nil {
		log.Fatalf("Cannot convert to json: %v\n", err)
	}
	// This is so that we can get the length of the users:
	theUsers := usersMap["response"].([]interface{})
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(theUsers), 2)
}

func TestGetUserByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	user, err := seedUser()
	if err != nil {
		log.Fatal(err)
	}
	userSample := []struct {
		id         string
		statusCode int
		username   string
		email      string
	}{
		{
			id:         strconv.Itoa(int(user.ID)),
			statusCode: 200,
			email:      user.Email,
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
		{
			id:         strconv.Itoa(12322), //an id that does not exist
			statusCode: 404,
		},
	}
	for _, v := range userSample {
		req, _ := http.NewRequest("GET", "/user/"+v.id, nil)
		rr := httptest.NewRecorder()

		r := gin.Default()
		r.GET("/user/:id", server.GetUserByID)
		r.ServeHTTP(rr, req)

		responseInterface := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseInterface)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			responseMap := responseInterface["response"].(map[string]interface{})
			assert.Equal(t, responseMap["email"], v.email)
		}

		if v.statusCode == 400 || v.statusCode == 404 {
			responseMap := responseInterface["error"].(map[string]interface{})

			if responseMap["invalid_request"] != nil {
				assert.Equal(t, responseMap["invalid_request"], "Invalid request")
			}
			if responseMap["no_user"] != nil {
				assert.Equal(t, responseMap["no_user"], "No user found")
			}
		}
	}
}

func TestDeleteUser(t *testing.T) {

	gin.SetMode(gin.TestMode)

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	user, err := seedUser()
	if err != nil {
		log.Fatal(err)
	}

	userSample := []struct {
		id         string
		statusCode int
	}{
		{
			// Convert int32 to int first before converting to string
			id:         strconv.Itoa(int(user.ID)),
			statusCode: 200,
		},
		{
			// When bad request data is given:
			id:         strconv.Itoa(int(4)),
			statusCode: 404,
		},
	}
	for _, v := range userSample {

		r := gin.Default()
		r.DELETE("/user/:id", server.DeleteUser)
		req, _ := http.NewRequest(http.MethodDelete, "/user/"+v.id, nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		responseInterface := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseInterface)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, responseInterface["response"], "Selected user has been deleted successfully.")
		}

		if v.statusCode == 400 || v.statusCode == 401 {
			responseMap := responseInterface["error"].(map[string]interface{})

			if responseMap["Invalid_request"] != nil {
				assert.Equal(t, responseMap["invalid_request"], "Invalid request")
			}
		}
	}
}
