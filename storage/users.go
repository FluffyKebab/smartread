package storage

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

var ErrEmailUsed = fmt.Errorf("email is already used")

type User struct {
	Id        int    `json:"id"`
	SessionId string `json:"sessionId"`
	IsGuest   bool   `json:"isGuest"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Username  string `json:"username"`
}

// Returns sessionId and error
func (s storer) AddUser(username string, passwordHash string, email string, sessionId string) (string, error) {
	//Check if email is already used
	emailIsUsed, err := s.emailIsUsed(email)
	if err != nil {
		return "", err
	}

	if emailIsUsed {
		return "", ErrEmailUsed
	}

	//Check if sessionId exists and update old row with new values if user is guest user
	if sessionId != "" {
		id, err := s.getIdFromSessionId(sessionId)
		if err != nil {
			return "", err
		}

		//If session id has row
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
					SET password = ?,
					SET email = ?
					SET username = ?
					SET isGuest = ?
				WHERE id = ?
			`, passwordHash, email, username, false, id)

			return sessionId, err
		}
	}

	//Create session id and insert new user row
	sessionId = uuid.New().String()

	_, err = s.db.Exec(`
		INSERT INTO users (sessionId, password, email, username isGuest) VALUES (?, ?, ?, ?)
	`, sessionId, passwordHash, email, username, false)

	return sessionId, err
}

// Returns sessionId and error
func (s storer) AddGuestUser() (string, error) {
	sessionId := uuid.New().String()

	_, err := s.db.Exec(`
		INSERT INTO users (sessionId, isGuest) VALUES (?, ?)
	`, sessionId, true)

	return sessionId, err
}

func (s storer) PasswordAndEmailIsCorrect(password, email string) (string, error) {
	row := s.db.QueryRow("SELECT sessionId FROM users WHERE password = ? AND email = ?", password, email)
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

func (s storer) getIdFromSessionId(sessionId string) (int, error) {
	row := s.db.QueryRow("SELECT id FROM users WHERE sessionId = ?", sessionId)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (s storer) emailIsUsed(email string) (bool, error) {
	row := s.db.QueryRow("SELECT * FROM users WHERE email = ?", email)

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
	row := s.db.QueryRow("SELECT isGuest FROM users WHERE id = ?", userId)

	var isGuest bool
	err := row.Scan(&isGuest)
	if err != nil {
		return false, err
	}

	return isGuest, nil
}
