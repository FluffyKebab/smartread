package storage

// To start the database in a docker run:
// docker run --rm --name database -p 5432:5432 -e POSTGRES_PASSWORD=123456789 -e POSTGRES_DB=database postgres

import (
	"database/sql"
	"fmt"
	"smartread/text"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type Storer interface {
	AddUser(username, passwordHash, email, sessionId string) (string, error)
	AddGuestUser() (string, error)
	CheckIfSessionIDIsValid(string) (bool, error)
	PasswordAndEmailIsCorrect(password, email string) (string, error)

	AddFile(userSessionId, fileData, fileName string) (string, error)
	GetAllUserFiles(userSessionId string) ([]File, error)
	QueryFile(userSessionId, fileId, query string) (string, error)
}

type storer struct {
	db          *sql.DB
	textHandler text.Handler
}

var _ storer = storer{}

func New(name, password string) (Storer, error) {
	db, err := sql.Open("pgx", fmt.Sprintf("postgres://postgres:%s@localhost:5432/%s", password, name))
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	err = configDatabase(db)
	if err != nil {
		return nil, err
	}

	textHandler, err := text.NewHandler()
	if err != nil {
		return nil, err
	}

	return storer{
		db:          db,
		textHandler: textHandler,
	}, nil
}

func configDatabase(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY, 
			sessionId TEXT NOT NULL,
			isGuest BOOLEAN NOT NULL,
			password TEXT NOT NULL,
			email TEXT NOT NULL, 
			username TEXT NOT NULL 
		);`,
	)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS files (
			id TEXT NOT NULL, 
			fileName TEXT NOT NULL,
			ownerId INTEGER NOT NULL
		);`,
	)

	return err
}
