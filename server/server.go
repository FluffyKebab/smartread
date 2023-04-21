package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"smartread/storage"

	"github.com/gorilla/mux"
)

type Server struct {
	storer storage.Storer
	router *mux.Router
}

func New() (Server, error) {
	r := mux.NewRouter()
	s, err := storage.NewStorer()
	if err != nil {
		return Server{}, err
	}

	server := Server{
		router: r,
		storer: s,
	}
	server.config()

	return server, nil
}

func (s Server) config() {
	// User handlers
	s.router.HandleFunc("/api/new_user", s.newUserHandler()).Methods("POST")
	s.router.HandleFunc("/api/login", s.loginHandler()).Methods("POST")

	// File and query handlers
	s.router.HandleFunc("/api/new_file", s.addFileHandler()).Methods("POST")
	s.router.HandleFunc("/api/query_file/{fileId}", s.queryFileHandler()).Methods("GET")
	s.router.HandleFunc("/api/get_files", s.getFilesHandler()).Methods("GET")

	// Static handlers
	s.router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./client/css"))))
	s.router.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("./client/js"))))

	// All other routing is done in the client
	s.router.NotFoundHandler = s.indexHandler()
}

func (s Server) ListenAndServer() error {
	fmt.Println("Running server...")
	return http.ListenAndServe(":8080", s.router)
}

func (s Server) indexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := s.getSessionId(w, r)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		staticHandler("./client/index.html", "text/html")(w, r)
	}
}

func staticHandler(filepath string, contentType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile(filepath)
		if err != nil {
			fmt.Println("Unable to read file data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", contentType)
		w.Write(data)
	}
}
