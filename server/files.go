package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"smartread/storage"

	"github.com/gorilla/mux"
)

type addFileResponse struct {
	FileId   string `json:"fileId"`
	FileName string `json:"filename"`
}

func (s Server) addFileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authenticate request.
		sessionId, err := s.getSessionId(w, r)
		if err != nil {
			handleError(w, err, http.StatusInternalServerError, "")
			return
		}

		// Get file data.
		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			handleError(w, err, http.StatusUnprocessableEntity, "missing file")
			return
		}
		defer file.Close()

		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, file)
		if err != nil {
			handleError(w, err, http.StatusInternalServerError, "")
			return
		}

		fileText := buf.String()

		// Add file data to local database and vector database.
		fileId, err := s.storer.AddFile(sessionId, fileText, fileHeader.Filename)
		if err != nil {
			handleError(w, err, http.StatusInternalServerError, "")
			return
		}

		// Send response to the client with file name and file id.
		body, err := json.Marshal(addFileResponse{FileId: fileId, FileName: fileHeader.Filename})
		if err != nil {
			handleError(w, err, http.StatusInternalServerError, "")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(body)
		if err != nil {
			handleError(w, err, http.StatusInternalServerError, "")
			return
		}
	}
}

type queryFileResponse struct {
	Response string `json:"response"`
}

func (s Server) queryFileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authenticate request.
		sessionId, err := s.getSessionId(w, r)
		if err != nil {
			handleError(w, err, http.StatusInternalServerError, "")
			return
		}

		// Get query string and file id from form values and params.
		fileId, ok := mux.Vars(r)["fileId"]
		if !ok {
			handleError(w, nil, http.StatusUnprocessableEntity, "missing id")
			return
		}

		err = r.ParseForm()
		if err != nil {
			handleError(w, err, http.StatusUnprocessableEntity, "malformed form data")
		}
		query := r.PostFormValue("query")
		if query == "" {
			handleError(w, nil, http.StatusUnprocessableEntity, "missing query")
			return
		}

		// Do query.
		response, err := s.storer.QueryFile(sessionId, fileId, query)
		if err != nil {
			status := http.StatusInternalServerError
			if errors.Is(err, storage.ErrFileNotExist) {
				status = http.StatusNotFound
			} else if errors.Is(err, storage.ErrFileNotOwnedByUser) {
				status = http.StatusUnauthorized
			}

			handleError(w, err, status, "")
			return
		}

		// Send response to client.
		body, err := json.Marshal(queryFileResponse{Response: response})
		if err != nil {
			handleError(w, err, http.StatusInternalServerError, "")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(body)
		if err != nil {
			handleError(w, err, http.StatusInternalServerError, "")
		}
	}
}

type getFilesResponse struct {
	Files []storage.File `json:"files"`
}

func (s Server) getFilesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authenticate request.
		sessionId, err := s.getSessionId(w, r)
		if err != nil {
			handleError(w, err, http.StatusInternalServerError, "")
			return
		}

		// Get files from database.
		files, err := s.storer.GetAllUserFiles(sessionId)
		if err != nil {
			handleError(w, err, http.StatusInternalServerError, "")
			return
		}

		// Send response to client.
		body, err := json.Marshal(getFilesResponse{Files: files})
		if err != nil {
			handleError(w, err, http.StatusInternalServerError, "")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(body)
		if err != nil {
			handleError(w, err, http.StatusInternalServerError, "")
			return
		}
	}
}
