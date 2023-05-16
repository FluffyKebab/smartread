package server

import (
	"fmt"
	"net/http"
	"time"
)

// getSessionId gets the sessionId from a cookie. If no sessionId
// exist, a new guest user will be created.
func (s Server) getSessionId(w http.ResponseWriter, r *http.Request) (string, error) {
	sessionCookie, err := r.Cookie("sessionId")
	if err != nil {
		if err != http.ErrNoCookie {
			return "", err
		}
	}

	if err == nil {
		id := sessionCookie.Value
		isValid, err := s.storer.CheckIfSessionIDIsValid(id)
		if err != nil {
			return "", err
		}

		if isValid {
			return id, nil
		}
	}

	//If no cookie is found or the session id is invalid a guest
	// user is created.
	sessionId, err := s.storer.AddGuestUser()
	if err != nil {
		return "", fmt.Errorf("creating guest user: %w", err)
	}

	setSessionIdCookie(w, sessionId)
	return sessionId, nil
}

func setSessionIdCookie(w http.ResponseWriter, sessionId string) {
	http.SetCookie(w, &http.Cookie{
		Name:    "sessionId",
		Path:    "/",
		Value:   sessionId,
		Expires: time.Now().AddDate(5, 0, 0), // Expires in 5 years.
	})
}
