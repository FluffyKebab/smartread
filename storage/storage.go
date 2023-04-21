package storage

import (
	"database/sql"
	"fmt"
	"smartread/text"

	_ "modernc.org/sqlite"
)

type Storer interface {
	AddUser(username, passwordHash, email, sessionId string) (string, error)
	AddGuestUser() (string, error)
	PasswordAndEmailIsCorrect(password, email string) (string, error)

	AddFile(userSessionId, fileData, fileName string) (string, error)
	GetAllUserFiles(userSessionId string) ([]File, error)
	QueryFile(userSessionId, fileId, query string) (string, error)
}

type storer struct {
	db          *sql.DB
	textHandler text.Handler
}

func NewStorer() (Storer, error) {
	db, err := sql.Open("sqlite", "./database.db")
	if err != nil {
		return storer{}, fmt.Errorf("Unable to open database: %s", err.Error())
	}

	err = configDatabase(db)
	if err != nil {
		return storer{}, err
	}

	textHandler, err := text.NewHandler()
	if err != nil {
		return storer{}, err
	}

	return storer{
		db:          db,
		textHandler: textHandler,
	}, nil
}

func configDatabase(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			sessionId string,
			isGuest boolean,
			password string,
			email string, 
			username string 
		);`,
	)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS files (
			id string, 
			fileName string,
			ownerId INTEGER
		);`,
	)
	if err != nil {
		return err
	}

	return nil
}
