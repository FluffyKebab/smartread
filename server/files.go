package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"smartread/storage"
)

type addFileResponse struct {
	FileId string `json:"fileId"`
}

func (s Server) addFileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authenticate request
		sessionId, err := s.getSessionId(w, r)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Get file data
		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer file.Close()

		fileData, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fileText := string(fileData)

		// Add file data to local database and vector database
		fileId, err := s.storer.AddFile(sessionId, fileText, fileHeader.Filename)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Send response
		body, err := json.Marshal(addFileResponse{FileId: fileId})
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(body)
		if err != nil {
			fmt.Println(err)
		}
	}
}

type queryFileResponse struct {
	Response string `json:"response"`
}

func (s Server) queryFileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authenticate request
		sessionId, err := s.getSessionId(w, r)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Get form values
		r.ParseForm()
		query := r.Form.Get("query")
		fileId := r.Form.Get("fileId")
		if query == "" || fileId == "" {
			fmt.Println("Missing form field")
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("Missing form field"))
			return
		}

		// Do query
		response, err := s.storer.QueryFile(sessionId, fileId, query)
		if err != nil {
			if err == storage.ErrFileNotExist {
				w.WriteHeader(http.StatusNotFound)
			} else if err == storage.ErrFileNotOwnedByUser {
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}

			fmt.Println(err.Error())
			return
		}

		// Send response in JSON
		body, err := json.Marshal(queryFileResponse{Response: response})
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(body)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

type getFilesResponse struct {
	Files []storage.File `json:"files"`
}

func (s Server) getFilesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authenticate request
		sessionId, err := s.getSessionId(w, r)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Get files database
		files, err := s.storer.GetAllUserFiles(sessionId)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}

		// Write response
		body, err := json.Marshal(getFilesResponse{Files: files})
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(body)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
