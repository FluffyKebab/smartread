package storage

import (
	"database/sql"
	"errors"
)

var (
	ErrFileNotExist       = errors.New("file with given id does not exist")
	ErrFileNotOwnedByUser = errors.New("file is not owned by user")
)

type File struct {
	Number   int    `json:"number"`
	Id       string `json:"id"`
	FileName string `json:"filename"`
	OwnerId  int    `json:"ownerId"`
}

func (s storer) AddFile(userSessionId string, fileData string, fileName string) (string, error) {
	ownerId, err := s.getIdFromSessionId(userSessionId)
	if err != nil {
		return "", err
	}

	fileId, err := s.textHandler.AddFile(fileData)
	if err != nil {
		return "", err
	}

	_, err = s.db.Exec(`
		INSERT INTO files (id, fileName, ownerId ) VALUES ($1, $2,$3)
	`, fileId, fileName, ownerId)

	return fileId, err
}

func (s storer) GetAllUserFiles(userSessionId string) ([]File, error) {
	ownerId, err := s.getIdFromSessionId(userSessionId)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query("SELECT * FROM files WHERE ownerId = $1 ORDER BY number DESC", ownerId)
	if err != nil {
		if err == sql.ErrNoRows {
			return []File{}, nil
		}

		return nil, err
	}

	files := make([]File, 0)
	for rows.Next() {
		var f File
		err := rows.Scan(&f.Number, &f.Id, &f.FileName, &f.OwnerId)
		if err != nil {
			return nil, err
		}

		files = append(files, f)
	}

	return files, nil
}

func (s storer) QueryFile(userSessionId, fileId, query string) (string, error) {
	//Check if file exists and if user owns file.
	ownerId, err := s.getIdFromSessionId(userSessionId)
	if err != nil {
		return "", err
	}

	rows, err := s.db.Query(`SELECT * FROM files WHERE ownerId = $1 AND id = $2`, ownerId, fileId)
	if err != nil {
		return "", err
	}

	if !rows.Next() {
		if err != sql.ErrNoRows {
			return "", err
		}

		// Check if file exist and is owned by other user or if file doesn't exist.
		rows, err = s.db.Query(`SELECT * FROM files WHERE AND id = $1`, fileId)
		if err != nil {
			return "", err
		}
		if !rows.Next() {
			return "", ErrFileNotExist
		}
		return "", ErrFileNotOwnedByUser
	}

	return s.textHandler.QueryFile(fileId, query)
}
