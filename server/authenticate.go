package server

import (
	"fmt"
	"net/http"
	"time"
)

// Gets sessionId from cookie. If no sessionId exist, a new guest user is created
func (s Server) getSessionId(w http.ResponseWriter, r *http.Request) (string, error) {
	sessionCookie, err := r.Cookie("sessionId")
	if err == nil {
		return sessionCookie.Value, nil
	}

	if err != http.ErrNoCookie {
		return "", err
	}

	//If no cookie is found the function creates a guest user
	sessionId, err := s.storer.AddGuestUser()
	if err != nil {
		return "", fmt.Errorf("Unable to create guest user: %s", err.Error())
	}

	setSessionIdCookie(w, sessionId)
	return sessionId, nil
}

func setSessionIdCookie(w http.ResponseWriter, sessionId string) {
	http.SetCookie(w, &http.Cookie{
		Name:    "sessionId",
		Path:    "/",
		Value:   sessionId,
		Expires: time.Now().AddDate(5, 0, 0), // Expires in 5 years
	})
}
