package main

import (
	"context"
	"fmt"

	"net/http"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"

	//"net/http/httputil"
	"github.com/gorilla/mux"
	"golang.org/x/net/http2"
)

var ()

const ()

type contextKey string

const contextEventKey contextKey = "event"

type parsedData struct {
	Msg string `json:"msg"`
}

func eventsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		headers := ""

		for name, values := range r.Header {
			for _, value := range values {
				fmt.Println(name, value)
				headers = headers + name + ":  " + value + "\n"
			}
		}
		event := &parsedData{
			Msg: headers,
		}
		ctx := context.WithValue(r.Context(), contextEventKey, *event)

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func httpError(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(code)
	fmt.Fprintln(w, error)
}

func gethandler(w http.ResponseWriter, r *http.Request) {
	//val := r.Context().Value(contextKey("event")).(parsedData)
	cloudPlatformScope := "https://www.googleapis.com/auth/cloud-platform"
	ts, err := google.DefaultTokenSource(context.Background(), cloudPlatformScope)
	if err != nil {
		httpError(w, fmt.Sprintf("Error getting default tokensource %v", err), http.StatusForbidden)
		return
	}

	tok, err := ts.Token()
	if err != nil {
		httpError(w, fmt.Sprintf("Error getting token %v", err), http.StatusForbidden)
		return
	}
	ctx := context.Background()
	oauth2Service, err := oauth2.NewService(ctx, option.WithScopes(oauth2.OpenIDScope))
	if err != nil {
		fmt.Printf("Error %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ti, err := oauth2Service.Tokeninfo().AccessToken(tok.AccessToken).Do()
	if err != nil {
		httpError(w, fmt.Sprintf("Error getting default tokeninfo %v", err), http.StatusForbidden)
		return
	}
	fmt.Printf("Token Issued to : %s", ti.Email)

	fmt.Fprint(w, ti.Email)
}

func main() {

	router := mux.NewRouter()
	router.Methods(http.MethodGet).Path("/get").HandlerFunc(gethandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: eventsMiddleware(router),
	}
	http2.ConfigureServer(server, &http2.Server{})
	fmt.Println("Starting Server..")
	err := server.ListenAndServe()

	fmt.Printf("Unable to start Server %v", err)

}
