package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ddddami/bindle/jsonx"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type UserList struct {
	Users []User `json:"users"`
	Total int    `json:"total"`
	Page  int    `json:"page"`
}

type Pagination struct {
	CurrentPage int `json:"current_page"`
	TotalPages  int `json:"total_pages"`
	PerPage     int `json:"per_page"`
	Total       int `json:"total"`
}

var users = []User{
	{ID: 1, Name: "Dami O.", Email: "dami@dami.me", CreatedAt: time.Now().Add(-72 * time.Hour)},
	{ID: 2, Name: "Ngozi A.", Email: "b@gm.com", CreatedAt: time.Now().Add(-48 * time.Hour)},
	{ID: 3, Name: "Feranmi", Email: "ff@ff", CreatedAt: time.Now().Add(-24 * time.Hour)},
	{ID: 4, Name: "Alfred", Email: "dico@mail.com", CreatedAt: time.Now()},
}

func main() {
	http.HandleFunc("/api/users", handleUsers)
	http.HandleFunc("/api/users/create", handleCreateUser)
	http.HandleFunc("/api/error", handleError)
	http.HandleFunc("/api/custom-error", handleCustomError)
	http.HandleFunc("/api/complex", handleComplexResponse)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// GET /api/users
func handleUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonx.SendError(w, errors.New("method not allowed"))
		return
	}

	jsonx.Send(w, users)
	// jsonx.SendSuccess(w, users)
}

// POST /api/users/create
func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonx.SendError(w, errors.New("method not allowed"))
		return
	}

	var input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := jsonx.DecodeJSONFromRequest(r, &input); err != nil {
		jsonx.SendError(w, err)
		return
	}

	if input.Name == "" || input.Email == "" {
		jsonx.SendError(w, errors.New("name and email are required"))
		return
	}

	newUser := User{
		ID:        len(users) + 1,
		Name:      input.Name,
		Email:     input.Email,
		CreatedAt: time.Now(),
	}

	users = append(users, newUser)

	jsonx.RespondWithJSON(w, newUser, jsonx.Options{
		SuccessStatus: http.StatusCreated,
		Headers: map[string]string{
			"X-Resource-ID": fmt.Sprintf("%d", newUser.ID),
		},
	})
}

// GET /api/error - Demonstrate error handling
func handleError(w http.ResponseWriter, r *http.Request) {
	err := errors.New("something went wrong")
	jsonx.SendError(w, err)
}

// GET /api/custom-error
func handleCustomError(w http.ResponseWriter, r *http.Request) {
	customErr := jsonx.ErrorDetail{
		Code:    "RESOURCE_NOT_FOUND",
		Message: "The requested resource could not be found",
	}
	// customErr := "customerror"

	jsonx.RespondWithError(w, customErr, jsonx.Options{
		ErrorStatus:    http.StatusNotFound,
		IndentResponse: true,
		ContentType:    "application/vnd.api+json",
	})
}

// GET /api/complex
func handleComplexResponse(w http.ResponseWriter, r *http.Request) {
	userList := UserList{
		Users: users,
		Total: len(users),
		Page:  1,
	}

	meta := Pagination{
		CurrentPage: 1,
		TotalPages:  1,
		PerPage:     10,
		Total:       len(users),
	}

	jsonx.RespondWithSuccess(w, userList, meta, jsonx.Options{
		IndentResponse: true,
	})
}
