package main

import (
	"fmt"
	"net/http"

	"github.com/alecholmes/xfccparser"
	"github.com/gorilla/mux"
	"golang.org/x/net/http2"
)

var ()

const ()

func getTokenHandler(w http.ResponseWriter, r *http.Request) {
	xfccHeader := r.Header.Get(xfccparser.ForwardedClientCertHeader)
	clientCerts, err := xfccparser.ParseXFCCHeader(xfccHeader)
	if err != nil {
		http.Error(w, "error parsing xfccparser", http.StatusInternalServerError)
		return
	}
	if len(clientCerts) == 0 {
		http.Error(w, " clientCerts is empty", http.StatusInternalServerError)
		return
	}
	u := clientCerts[0].URI
	if len(u) == 0 {
		http.Error(w, " clientCerts uri is empty", http.StatusInternalServerError)
		return
	}

	// get an access_token for the service u
	token := fmt.Sprintf("token_for_service_%s", u[0])
	fmt.Fprint(w, token)
}

func main() {
	router := mux.NewRouter()
	router.Methods(http.MethodGet).Path("/computeMetadata/v1/instance/service-accounts/default/token").HandlerFunc(getTokenHandler)
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	http2.ConfigureServer(server, &http2.Server{})
	err := server.ListenAndServe()
	fmt.Printf("Unable to start Server %v", err)
}
