package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"smartread/text"

	_ "github.com/lib/pq"
)

var ErrEnvNotSet = errors.New("environment variable not set")

const (
	_postgresPasswordEnvVar = "POSTGRES_PASSWORD"
	_postgresNameEnvVar     = "POSTGRES_DB"
	_postgresUserEnvVar     = "POSTGRES_USER"
	_postgresHostEnvVar     = "POSTGRES_HOST"

	_port = 5432
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

func New() (Storer, error) {
	password := os.Getenv(_postgresPasswordEnvVar)
	if password == "" {
		return nil, fmt.Errorf("%w: %s", ErrEnvNotSet, _postgresPasswordEnvVar)
	}

	name := os.Getenv(_postgresNameEnvVar)
	if name == "" {
		return nil, fmt.Errorf("%w: %s", ErrEnvNotSet, _postgresNameEnvVar)
	}

	user := os.Getenv(_postgresUserEnvVar)
	if user == "" {
		return nil, fmt.Errorf("%w: %s", ErrEnvNotSet, _postgresUserEnvVar)
	}

	host := os.Getenv(_postgresHostEnvVar)
	if host == "" {
		return nil, fmt.Errorf("%w: %s", ErrEnvNotSet, _postgresHostEnvVar)
	}

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, _port, user, password, name,
	)

	db, err := sql.Open("postgres", psqlInfo)
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
			number SERIAL PRIMARY KEY, 
			id TEXT NOT NULL, 
			fileName TEXT NOT NULL,
			ownerId INTEGER NOT NULL
		);`,
	)

	return err
}
