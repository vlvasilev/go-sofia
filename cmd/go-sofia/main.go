package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-sofia/internal/diagnostics"
	"src/github.com/gorilla/mux"
)

type serverConf struct {
	port   string
	router http.Handler
	name   string
}

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
	diagnostics := diagnostics.NewDiagnostics()

	possibleErrors := make(chan error, 2)

	configurations := []serverConf{
		{
			port:   blPort,
			router: router,
			name:   "application server",
		},
		{
			port:   diagPort,
			router: diagnostics,
			name:   "diagnostics server",
		},
	}

	servers := make([]*http.Server, 2)
	for i, sc := range configurations {
		go func(conf serverConf, i int) {
			log.Print("The %s server is preparing tohandle connection...", conf.name)
			servers[i] = &http.Server{
				Addr:    ":" + conf.port,
				Handler: conf.router,
			}
			err := servers[i].ListenAndServe()
			//server.Shutdown()
			if err != nil {
				possibleErrors <- err
			}
		}(sc, i)
	}

	select {
	case err := <-possibleErrors:
		for _, s := range servers {
			context
			s.Shutdown(context.Background())
		}
		log.Fatal(err)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Print(w, http.StatusText(http.StatusOK))
}
