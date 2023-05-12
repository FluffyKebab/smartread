package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type newUserResponse struct {
	SessionId string `json:"sessionId"`
}

func (s Server) newUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authenticate user.
		sessionId, err := s.getSessionId(w, r)
		if err != nil {
			handleError(w, err, http.StatusInternalServerError, "")
			return
		}

		// Get username, password and email from form vales.
		r.ParseForm()
		username := r.Form.Get("username")
		password := r.Form.Get("password")
		email := r.Form.Get("email")

		if username == "" || password == "" || email == "" {
			handleError(w, nil, http.StatusUnprocessableEntity, "missing from field")
			return
		}

		// Add user data to database
		sessionId, err = s.storer.AddUser(username, password, email, sessionId)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Set cookie
		setSessionIdCookie(w, sessionId)

		// Send response
		body, err := json.Marshal(newUserResponse{SessionId: sessionId})
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(body)
		if err != nil {
			fmt.Println(err)
		}
	}
}

type loginResponse struct {
	Success   bool   `json:"success"`
	SessionId string `json:"sessionId"`
}

func (s Server) loginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authenticate user
		sessionId, err := s.getSessionId(w, r)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Getting form values
		r.ParseForm()
		password := r.Form.Get("password")
		email := r.Form.Get("email")

		if password == "" || email == "" {
			fmt.Println("Missing form field")
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("Missing form field"))
			return
		}

		// Check if login is correct
		sessionId, err = s.storer.PasswordAndEmailIsCorrect(password, email)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Send response
		success := true
		if sessionId == "" {
			success = false
		}

		body, err := json.Marshal(loginResponse{Success: success, SessionId: sessionId})
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(body)
		if err != nil {
			fmt.Println(err)
		}
	}
}
