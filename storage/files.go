package storage

import (
	"database/sql"
	"fmt"
)

type File struct {
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
		INSERT INTO files (id, fileName, ownerId ) VALUES (?, ?, ?)
	`, fileId, fileName, ownerId)

	return fileId, nil
}

func (s storer) GetAllUserFiles(userSessionId string) ([]File, error) {
	ownerId, err := s.getIdFromSessionId(userSessionId)
	if err != nil {
		return []File{}, err
	}

	rows, err := s.db.Query("SELECT * FROM files WHERE ownerId = ?", ownerId)
	if err != nil {
		return []File{}, err
	}

	files := make([]File, 0)
	for rows.Next() {
		var f File
		err := rows.Scan(&f.Id, &f.FileName, &f.OwnerId)
		if err != nil {
			return []File{}, err
		}

		files = append(files, f)
	}

	return files, nil
}

var ErrFileNotExist = fmt.Errorf("File with given id does not exist")
var ErrFileNotOwnedByUser = fmt.Errorf("File is not owned by user")

func (s storer) QueryFile(userSessionId, fileId, query string) (string, error) {
	//Check if file exists and if user owns file
	ownerId, err := s.getIdFromSessionId(userSessionId)
	if err != nil {
		return "", err
	}

	row := s.db.QueryRow(`SELECT * FROM files WHERE ownerId = ? AND id = ?`, ownerId, fileId)
	err = row.Scan()
	if err != nil {
		if err != sql.ErrNoRows {
			return "", err
		}

		// Check if file exist and is owned by other user or if file doesn't exist
		row = s.db.QueryRow(`SELECT * FROM files WHERE AND id = ?`, fileId)
		err := row.Scan()
		if err != nil {
			if err != sql.ErrNoRows {
				return "", err
			}

			return "", ErrFileNotExist // No file with fileId given exists
		}

		return "", ErrFileNotOwnedByUser // File exists but with other ownerId
	}

	return s.textHandler.QueryFile(fileId, query)
}
