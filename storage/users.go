package storage

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrEmailUsed        = fmt.Errorf("email is already used")
	ErrInvalidSessionId = fmt.Errorf("session id is invalid")
)

type User struct {
	Id        int    `json:"id"`
	SessionId string `json:"sessionId"`
	IsGuest   bool   `json:"isGuest"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Username  string `json:"username"`
}

// AddUser creates a new row with Returns sessionId and error.
func (s storer) AddUser(username string, passwordHash string, email string, sessionId string) (string, error) {
	//Check if email is used.
	emailIsUsed, err := s.emailIsUsed(email)
	if err != nil {
		return "", err
	}

	if emailIsUsed {
		return "", ErrEmailUsed
	}

	//Check if sessionId exists and update old row with new values if user is guest user.
	if sessionId != "" {
		id, err := s.getIdFromSessionId(sessionId)
		if err != nil {
			return "", err
		}

		//If session id is in the database.
		if id != -1 {
			userIsGuest, err := s.userIsGuest(id)
			if err != nil {
				return "", err
			}

			if !userIsGuest {
				return "", fmt.Errorf("User already exist with sessionId given")
			}

			_, err = s.db.Exec(`
				UPDATE users 
					SET password = $1
					SET email = $2
					SET username = $3
					SET isGuest = $4
				WHERE id = $5
			`, passwordHash, email, username, false, id)

			return sessionId, err
		}
	}

	//Create session id and insert new user row.
	sessionId = uuid.New().String()

	_, err = s.db.Exec(`
		INSERT INTO users (sessionId, password, email, username, isGuest) VALUES ($1, $2, $3, $4, $5)
	`, sessionId, passwordHash, email, username, false)

	return sessionId, err
}

// AddGuestUser returns sessionId from new guest user.
func (s storer) AddGuestUser() (string, error) {
	sessionId := uuid.New().String()

	_, err := s.db.Exec(`
		INSERT INTO users (sessionId, isGuest, password, username, email) VALUES ($1, $2, $3, $4, $5)
	`, sessionId, true, "", "", "")

	return sessionId, err
}

// PasswordAndEmailIsCorrect returns the session id of the user with password and email given.
// If the user does not exists and empty string is returned.
func (s storer) PasswordAndEmailIsCorrect(password, email string) (string, error) {
	row := s.db.QueryRow("SELECT sessionId FROM users WHERE password = $1 AND email = $2", password, email)
	var sessionId string
	err := row.Scan(&sessionId)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}

		return "", err
	}

	return sessionId, nil
}

func (s storer) CheckIfSessionIDIsValid(sessionId string) (bool, error) {
	row := s.db.QueryRow("SELECT id FROM users WHERE sessionId = $1", sessionId)
	var id int
	err := row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (s storer) getIdFromSessionId(sessionId string) (int, error) {
	row := s.db.QueryRow("SELECT id FROM users WHERE sessionId = $1", sessionId)
	var id int
	err := row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, ErrInvalidSessionId
		}
		return -1, err
	}

	return id, nil
}

func (s storer) emailIsUsed(email string) (bool, error) {
	row := s.db.QueryRow("SELECT * FROM users WHERE email = $1", email)

	err := row.Scan()
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (s storer) userIsGuest(userId int) (bool, error) {
	row := s.db.QueryRow("SELECT isGuest FROM users WHERE id = $1", userId)

	var isGuest bool
	err := row.Scan(&isGuest)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	return isGuest, nil
}
