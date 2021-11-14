package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/drive/v3"
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
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano()) //Randomly Seed

	//googleCalendarInsertTestTheSecond()
	//googleCalendarCreateEventTest()
	//googleCalendarReadTest()
	//insertMeetingAttachment()
	/*
		getService, err := getService()
		if err != nil {
			panic(fmt.Sprintf("Could not get service: %v\n", err.Error()))
		}
	*/
	//googleDriveList(getService)
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

	wd, _ := os.Getwd()
	credDir := filepath.Join(wd, "creds", "credentials-insert.json")
	b, err := ioutil.ReadFile(credDir)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarEventsScope)
	if err != nil {
		fmt.Printf("Unable to parse client secret file to config: %v", err)
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	tokFile := filepath.Join(wd, "creds", "insertToken.json")
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	/*
		token := oauth2.Token{
			AccessToken:  "ya29.a0ARrdaM_vWvmg9vU8bujpanzVEiQ3GXWBhP5CyA7sHDsAzCTqcsAbpktrMObQZ1W6-QxmCoj5ZK00zfeiT0Bpmf8pN45UWrTqpUDPBCGqG4cjHUDXDjiwmwS4wBnWUYMr2YAbEfFHbDEgRJCJR5-1UvL_wq_y",
			RefreshToken: "1//04pFlAHNvTESECgYIARAAGAQSNwF-L9Ir2v2T0JPHk_H1fxAWJA7JH7eGbbGCHJ0R5WqJNaV3WR9MAbriryVPuDRW86sGsay-5G4",
			TokenType:    "Bearer",
			Expiry:       time.Now(),
		}

		var tokenSource = config.TokenSource(context.Background(), &token)
	*/

	//calendarService, err := calendar.NewService(ctx, option.WithTokenSource(tokenSource))
	calendarService, err := calendar.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, tok)))

	if err != nil {
		fmt.Printf("Client Setup failed: %v", err)
		log.Fatalf("Client Setup failed: %v", err)
	}

	theEvent := &calendar.Event{
		Start: &calendar.EventDateTime{
			DateTime: "2021-11-12T17:06:02.000Z",
			TimeZone: "America/Chicago",
		},
		End: &calendar.EventDateTime{
			DateTime: "2021-11-12T19:06:02.000Z",
			TimeZone: "America/Chicago",
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

/* This is a test fnction for inserting a Google Drive Attachment to a meeting */
func insertMeetingAttachment() {
	wd, _ := os.Getwd()

	/* Step 1 open file for working with */
	fileDir := filepath.Join(wd, "testFileUploads", "testfile.txt")
	f, err := os.Open(fileDir)

	if err != nil {
		panic(fmt.Sprintf("cannot open file: %v", err))
	}

	defer f.Close()

	/* Step 2 get Google Service */
	service, err := getService()

	if err != nil {
		panic(fmt.Sprintf("Uh oh, couldn't create service: %v\n", err.Error()))
	}

	// Step 3. Create the directory
	dir, err2 := createDir(service, "testGoogleDriveFolder", "root")
	if err2 != nil {
		panic(fmt.Sprintf("Uh oh, couldn't create directory: %v\n", err2.Error()))
	}

	// Step 4. Create the file and upload its content
	file, err := createFile(service, "testfile.txt", "text/plain", f, dir.Id)

	if err != nil {
		panic(fmt.Sprintf("Could not create file: %v\n", err))
	}

	fmt.Printf("File '%s' successfully uploaded in '%s' directory\n", file.Name, dir.Name)
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

/* This reads all files in a certain Google Drive Directory */
func googleDriveList(service *drive.Service) {

	r, err := service.Files.List().PageSize(10).
		Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}
	fmt.Println("Files:")
	if len(r.Files) == 0 {
		fmt.Println("No files found.")
	} else {
		for _, i := range r.Files {
			fmt.Printf("%s (%s)\n", i.Name, i.Id)
		}
	}
}

/* Get Google Drive Service */
func getService() (*drive.Service, error) {
	wd, _ := os.Getwd()
	credDir := filepath.Join(wd, "creds", "google-drive-credentials.json")

	b, err := ioutil.ReadFile(credDir)
	if err != nil {
		fmt.Printf("Unable to read credentials.json file. Err: %v\n", err)
		return nil, err
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveFileScope)

	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	tokFile := filepath.Join(wd, "creds", "googleDriveToken.json")
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}

	service, err := drive.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, tok)))

	if err != nil {
		fmt.Printf("Cannot create the Google Drive service: %v\n", err)
		return nil, err
	}

	return service, err
}

/* Create Google Drive Directory */
func createDir(service *drive.Service, name string, parentId string) (*drive.File, error) {
	d := &drive.File{
		Name:     name,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []string{parentId},
	}

	file, err := service.Files.Create(d).Do()

	if err != nil {
		log.Println("Could not create dir: " + err.Error())
		return nil, err
	}

	return file, nil
}

/* Create Google Drive file in a specific directory */
func createFile(service *drive.Service, name string, mimeType string, content io.Reader, parentId string) (*drive.File, error) {
	f := &drive.File{
		MimeType: mimeType,
		Name:     name,
		Parents:  []string{parentId},
	}
	file, err := service.Files.Create(f).Media(content).Do()

	if err != nil {
		log.Println("Could not create file: " + err.Error())
		return nil, err
	}

	return file, nil
}
