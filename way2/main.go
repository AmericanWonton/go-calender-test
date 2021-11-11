package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

/* TEMPLATE DEFINITION */
var template1 *template.Template

/* Google Client/Calendar information */

type CalendarPassing struct {
	CalendarAPIKey     string `json:"CalendarAPIKey"`
	CalendarID         string `json:"CalendarID"`
	GoogleClientID     string `json:"GoogleClientID"`
	GoogleClientSecret string `json:"GoogleClientSecret"`
	CurrentTime        string `json:"CurrentTime"`
}

var calendarPassing CalendarPassing

/* Webpage information passing */
type ViewData struct {
	PassedCalendarInfo CalendarPassing `json:"PassedCalendarInfo"`
}

func init() {
	template1 = template.Must(template.ParseGlob("./static/templates/*")) //pass templates
	getCalendarCreds()                                                    //Get calendar creds
}

//Handles the Index requests; Ask User if they're legal here
func index(w http.ResponseWriter, r *http.Request) {
	//Build info to pass
	currentTime := time.Now()
	calendarPassing.CurrentTime = currentTime.Format("2006-01-02")
	vd := ViewData{
		PassedCalendarInfo: calendarPassing,
	}
	/* Execute template, handle error */
	err1 := template1.ExecuteTemplate(w, "index.gohtml", vd)
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

	fmt.Printf("DEBUG: We are serving files on internet\n")
	http.Handle("/favicon.ico", http.NotFoundHandler()) //For missing FavIcon
	//Serve our pages
	myRouter.HandleFunc("/", index)
	//Serve Google Calendar stuff

	//Serve our static files
	myRouter.Handle("/", http.FileServer(http.Dir("./static")))
	myRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	log.Fatal(http.ListenAndServe(":5000", myRouter))
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano()) //Randomly Seed

	googleCalendarInsertTestTheSecond()
	//googleCalendarCreateEventTest()
	//googleCalendarReadTest()
	handleRequests() // handle requests
}

/* This gets our enviornment varialbles to create our google calendar information */
func getCalendarCreds() {
	_, ok := os.LookupEnv("GDESS_GOOGLE_CLIENT_ID")
	if !ok {
		message := "This ENV Variable is not present: " + "GDESS_GOOGLE_CLIENT_ID"
		panic(message)
	}

	_, ok2 := os.LookupEnv("GDESS_GOOGLE_CLIENT_SECRET")
	if !ok2 {
		message := "This ENV Variable is not present: " + "GDESS_GOOGLE_CLIENT_SECRET"
		panic(message)
	}

	_, ok3 := os.LookupEnv("GDESS_GOOGLE_CALENDAR_APIKEY")
	if !ok3 {
		message := "This ENV Variable is not present: " + "GDESS_GOOGLE_CALENDAR_APIKEY"
		panic(message)
	}

	_, ok4 := os.LookupEnv("GDESS_CALENDAR_ID")
	if !ok4 {
		message := "This ENV Variable is not present: " + "GDESS_CALENDAR_ID"
		panic(message)
	}

	calendarPassing.GoogleClientID = os.Getenv("GDESS_GOOGLE_CLIENT_ID")
	calendarPassing.GoogleClientSecret = os.Getenv("GDESS_GOOGLE_CLIENT_SECRET")
	calendarPassing.CalendarAPIKey = os.Getenv("GDESS_GOOGLE_CALENDAR_APIKEY")
	calendarPassing.CalendarID = os.Getenv("GDESS_CALENDAR_ID")
}

/* This does all the fun Google Calender reading */
func googleCalendarReadTest() {
	currentTime := time.Now() //Used for debugging
	fmt.Println("Here is the current Google time in YYYY-MM-DD : ", currentTime.Format("2006-01-02"))

	ctx := context.Background()
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		fmt.Printf("Unable to parse client secret file to config: %v", err)
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		fmt.Printf("Unable to retrieve Calendar client: %v", err)
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	t := time.Now().Format(time.RFC3339)
	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}
	fmt.Println("Upcoming events:")
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		for _, item := range events.Items {
			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}
			fmt.Printf("%v (%v)\n", item.Summary, date)
			fmt.Printf("Here's a description: %v\n", item.Description)
		}
	}

}

/* This is another test function for Inserting a Google Calendar Event*/

func googleCalendarInsertTestTheSecond() {
	ctx := context.Background()

	b, err := ioutil.ReadFile("credentials-insert.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarEventsScope)
	if err != nil {
		fmt.Printf("Unable to parse client secret file to config: %v", err)
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	tokFile := "insertToken.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}

	calendarService, err := calendar.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, tok)))

	if err != nil {
		fmt.Printf("Client Setup failed: %v", err)
		log.Fatalf("Client Setup failed: %v", err)
	}

	theEvent := &calendar.Event{
		Start: &calendar.EventDateTime{
			DateTime: "2021-11-021T09:00:00-07:00",
			TimeZone: "America/Los_Angeles",
		},
		End: &calendar.EventDateTime{
			DateTime: "2021-11-21T17:00:00-07:00",
			TimeZone: "America/Saint_Louis",
		},
		Summary:     "Test Calendar Creation",
		Location:    "800 Howard St., San Francisco, CA 94103",
		Description: "A test Google Calendar date, created from Golang",
	}

	calendarId := "primary"
	event, err2 := calendarService.Events.Insert(calendarId, theEvent).Do()
	if err2 != nil {
		fmt.Printf("Unable to create event: %v\n", err2.Error)
		//fmt.Printf("Here is the header: %v\n", event.Header)
		//fmt.Printf("Here is the HTTPStatusCode: %v\n", event.HTTPStatusCode)
		fmt.Printf("Here is the MarshalJSON: %v\n", event.MarshalJSON)
		log.Fatalf("Unable to create event. %v\n", err2)
	}
	fmt.Printf("Event created: %s\n", event.HtmlLink)
}

/* This creates a test Google Calendar Event */
func googleCalendarCreateEventTest() {

	ctx := context.Background()
	b, err := ioutil.ReadFile("credentials-insert.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		fmt.Printf("Unable to parse client secret file to config: %v", err)
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		fmt.Printf("Unable to retrieve Calendar client: %v", err)
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	theEvent := &calendar.Event{
		Start: &calendar.EventDateTime{
			DateTime: "2021-11-021T09:00:00-07:00",
			TimeZone: "America/Los_Angeles",
		},
		End: &calendar.EventDateTime{
			DateTime: "2021-11-21T17:00:00-07:00",
			TimeZone: "America/Saint_Louis",
		},
		Summary:     "Test Calendar Creation",
		Location:    "800 Howard St., San Francisco, CA 94103",
		Description: "A test Google Calendar date, created from Golang",
	}

	calendarId := "primary"
	event, err2 := srv.Events.Insert(calendarId, theEvent).Do()
	if err2 != nil {
		fmt.Printf("Unable to create event: %v\n", err2)
		log.Fatalf("Unable to create event. %v\n", err2)
	}
	fmt.Printf("Event created: %s\n", event.HtmlLink)
}
