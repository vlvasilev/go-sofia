package diagnostics

import (
	"fmt"
	"net/http"

	"src/github.com/gorilla/mux"
)

func NewDiagnostics() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/healthz", healthz)
	router.HandleFunc("/ready", ready)
	return router
}

func healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Print(w, http.StatusText(http.StatusOK))
}

func ready(w http.ResponseWriter, r *http.Request) {
	fmt.Print(w, http.StatusText(http.StatusOK))
}
