package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"smartread/storage"

	"github.com/gorilla/mux"
)

var verbose = true

type Server struct {
	storer storage.Storer
	router *mux.Router
}

func New() (Server, error) {
	r := mux.NewRouter()
	s, err := storage.New()
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
	// User handlers:
	s.router.HandleFunc("/api/new_user", s.newUserHandler).Methods("POST")
	s.router.HandleFunc("/api/login", s.loginHandler).Methods("POST")

	// File and query handlers:
	s.router.HandleFunc("/api/new_file", s.addFileHandler).Methods("POST")
	s.router.HandleFunc("/api/query_file/{fileId}", s.queryFileHandler).Methods("POST")
	s.router.HandleFunc("/api/get_files", s.getFilesHandler).Methods("GET")

	// Static handlers:
	s.router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./client/css"))))
	s.router.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("./client/js"))))
	s.router.PathPrefix("/img/").Handler(http.StripPrefix("/img/", http.FileServer(http.Dir("./client/img"))))

	// All other routing is done in the client.
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
			handleError(w, err, http.StatusInternalServerError, "failed to create guest user")
			return
		}

		staticHandler("./client/index.html", "text/html")(w, r)
	}
}

func staticHandler(filepath string, contentType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile(filepath)
		if err != nil {
			handleError(w, err, http.StatusInternalServerError, "")
			return
		}

		w.Header().Set("Content-Type", contentType)
		w.Write(data)
	}
}

func handleError(w http.ResponseWriter, err error, status int, userMessage string) {
	if verbose && err != nil {
		fmt.Printf("%s: %s \n", userMessage, err.Error())
	} else if verbose {
		fmt.Println(userMessage)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	w.Write([]byte(userMessage))
}
