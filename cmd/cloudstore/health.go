package main

import (
	"fmt"
	"net/http"

	"github.com/rusenask/cloudstore"

	log "github.com/sirupsen/logrus"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func startHealthHandler() {
	http.HandleFunc("/health", handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cloudstore.DefaultHTTPHealthCheckPort), nil))
}
