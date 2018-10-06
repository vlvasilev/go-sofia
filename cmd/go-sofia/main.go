package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-sofia/internal/diagnostics"
	"src/github.com/gorilla/mux"
)

func main() {
	log.Print("Hello, World")

	blPort := os.Getenv("PORT")
	if len(blPort) == 0 {
		log.Fatal("The applications port should be set")
	}
	diagPort := os.Getenv("DIAG_PORT")
	if len(diagPort) == 0 {
		log.Fatal("The diagnostics port should be set")
	}

	router := mux.NewRouter()
	router.HandleFunc("/", hello)

	possibleErrors := make(chan error, 2)

	go func() {
		log.Print("The application server is preparing tohandle connection...")
		server := http.Server{
			Addr:    ":" + blPort,
			Handler: router,
		}
		err := server.ListenAndServe()
		//server.Shutdown()
		if err != nil {
			possibleErrors <- err
		}
	}()
	go func() {
		log.Print("The diagnostics server is preparing tohandle connection...")
		diagnostics := diagnostics.NewDiagnostics()
		err := http.ListenAndServe(":"+diagPort, diagnostics)
		if err != nil {
			possibleErrors <- err
		}
	}()

	select {
	case err := <-possibleErrors:
		log.Fatal(err)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Print(w, http.StatusText(http.StatusOK))
}
