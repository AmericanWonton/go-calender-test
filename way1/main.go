package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

/* TEMPLATE DEFINITION */
var template1 *template.Template

//Google Auth config
var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:3000/GoogleCallback",
		ClientID:     os.Getenv("GDESS_GOOGLE_CLIENT_ID"),     // from https://console.developers.google.com/project/<your-project-id>/apiui/credential
		ClientSecret: os.Getenv("GDESS_GOOGLE_CLIENT_SECRET"), // from https://console.developers.google.com/project/<your-project-id>/apiui/credential
		Scopes:       []string{},
		Endpoint:     google.Endpoint,
	}
	// Some random string, random for each request
	oauthStateString = "random"
)

func init() {
	template1 = template.Must(template.ParseGlob("./static/templates/*"))
}

//Handles the Index requests; Ask User if they're legal here
func index(w http.ResponseWriter, r *http.Request) {
	/* Execute template, handle error */
	err1 := template1.ExecuteTemplate(w, "index.gohtml", nil)
	HandleError(w, err1)
}

// Handle Errors passing templates
func HandleError(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalln(err)
	}
}

//Handles all requests coming in
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	http.Handle("/favicon.ico", http.NotFoundHandler()) //For missing FavIcon
	//Serve our pages
	myRouter.HandleFunc("/", index)
	//Serve Google Calendar stuff
	http.HandleFunc("/handleGoogleLogin", handleGoogleLogin)
	http.HandleFunc("/handleGoogleCallback", handleGoogleCallback)
	//Serve our static files
	myRouter.Handle("/", http.FileServer(http.Dir("./static")))
	myRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	log.Fatal(http.ListenAndServe(":5000", myRouter))
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano()) //Randomly Seed

	handleRequests() // handle requests
}

func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))
}
