package main


import (
	"os"
	"context"
	"log"
	"fmt"
	"time"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/gmail/v1"
)

/* Services Declaration */
var GmailService *gmail.Service //This gets initialized in init
var GoogleCalendarService *calendar.Service // This gets initilized in init
var GoogleDriveService *drive.Service //This gets initialized in init

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

	_, ok5 := os.LookupEnv("GDESS_CALENDAR_REFRESHTOKEN")
	if !ok5 {
		message := "This ENV Variable is not present: " + "GDESS_CALENDAR_REFRESHTOKEN"
		panic(message)
	}

	_, ok6 := os.LookupEnv("GDESS_CALENDAR_ACCESSTOKEN")
	if !ok6 {
		message := "This ENV Variable is not present: " + "GDESS_CALENDAR_ACCESSTOKEN"
		panic(message)
	}

	_, ok7 := os.LookupEnv("GDESS_GOOGLE_EMAIL_CLIENT_ID")
	if !ok7 {
		message := "This ENV Variable is not present: " + "GDESS_GOOGLE_EMAIL_CLIENT_ID"
		panic(message)
	}

	_, ok8 := os.LookupEnv("GDESS_GOOGLE_EMAIL_CLIENT_SECRET")
	if !ok8 {
		message := "This ENV Variable is not present: " + "GDESS_GOOGLE_EMAIL_CLIENT_SECRET"
		panic(message)
	}

	_, ok9 := os.LookupEnv("GDESS_GOOGLE_EMAIL_ACCESS_TOKEN")
	if !ok9 {
		message := "This ENV Variable is not present: " + "GDESS_GOOGLE_EMAIL_ACCESS_TOKEN"
		panic(message)
	}

	_, ok10 := os.LookupEnv("GDESS_GOOGLE_EMAIL_REFRESH_TOKEN")
	if !ok10 {
		message := "This ENV Variable is not present: " + "GDESS_GOOGLE_EMAIL_REFRESH_TOKEN"
		panic(message)
	}

	_, ok11 := os.LookupEnv("GDESS_EMAIL")
	if !ok11 {
		message := "This ENV Variable is not present: " + "GDESS_EMAIL"
		panic(message)
	}

	_, ok12 := os.LookupEnv("GDESS_PASSWORD")
	if !ok12 {
		message := "This ENV Variable is not present: " + "GDESS_PASSWORD"
		panic(message)
	}

	_, ok13 := os.LookupEnv("GDESS_GOOGLE_DRIVE_CLIENT_ID")
	if !ok13 {
		message := "This ENV Variable is not present: " + "GDESS_GOOGLE_DRIVE_CLIENT_ID"
		panic(message)
	}

	_, ok14 := os.LookupEnv("GDESS_GOOGLE_DRIVE_CLIENT_SECRET")
	if !ok14 {
		message := "This ENV Variable is not present: " + "GDESS_GOOGLE_DRIVE_CLIENT_SECRET"
		panic(message)
	}

	_, ok15 := os.LookupEnv("GDESS_GOOGLEDRIVE_REFRESH")
	if !ok15 {
		message := "This ENV Variable is not present: " + "GDESS_GOOGLEDRIVE_REFRESH"
		panic(message)
	}

	_, ok16 := os.LookupEnv("GDESS_GOOGLEDRIVE_ACCESS")
	if !ok16 {
		message := "This ENV Variable is not present: " + "GDESS_GOOGLEDRIVE_ACCESS"
		panic(message)
	}

	calendarPassing.GoogleClientID = os.Getenv("GDESS_GOOGLE_CLIENT_ID")
	calendarPassing.GoogleClientSecret = os.Getenv("GDESS_GOOGLE_CLIENT_SECRET")
	calendarPassing.CalendarAPIKey = os.Getenv("GDESS_GOOGLE_CALENDAR_APIKEY")
	calendarPassing.CalendarID = os.Getenv("GDESS_CALENDAR_ID")
	calendarPassing.GoogleClientCalendarRefreshToken = os.Getenv("GDESS_CALENDAR_REFRESHTOKEN")
	calendarPassing.GoogleClientCalendarAccessToken = os.Getenv("GDESS_CALENDAR_ACCESSTOKEN")
	calendarPassing.EmailClient = os.Getenv("GDESS_GOOGLE_EMAIL_CLIENT_ID")
	calendarPassing.EmailSecret = os.Getenv("GDESS_GOOGLE_EMAIL_CLIENT_SECRET")
	calendarPassing.EmailAccess = os.Getenv("GDESS_GOOGLE_EMAIL_ACCESS_TOKEN")
	calendarPassing.EmailRefresh = os.Getenv("GDESS_GOOGLE_EMAIL_REFRESH_TOKEN")
	calendarPassing.CurrentEmail = os.Getenv("GDESS_EMAIL")
	calendarPassing.CurrentPWord = os.Getenv("GDESS_PASSWORD")
	calendarPassing.GoogleDriveClientID = os.Getenv("GDESS_GOOGLE_DRIVE_CLIENT_ID")
	calendarPassing.GoogleDriveClientSecret = os.Getenv("GDESS_GOOGLE_DRIVE_CLIENT_SECRET")
	calendarPassing.GoogleDriveRefresh = os.Getenv("GDESS_GOOGLEDRIVE_REFRESH")
	calendarPassing.GoogleDriveAccess = os.Getenv("GDESS_GOOGLEDRIVE_ACCESS")
}

/* This creates a service for our email */
func OAuthGmailService() {
	config := oauth2.Config{
		ClientID:     calendarPassing.EmailClient,
		ClientSecret: calendarPassing.EmailSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost",
	}

	token := oauth2.Token{
		AccessToken:  calendarPassing.EmailAccess,
		RefreshToken: calendarPassing.EmailRefresh,
		TokenType:    "Bearer",
		Expiry:       time.Now(),
	}

	//Create a context to use for our gmail services

	var tokenSource = config.TokenSource(context.Background(), &token)

	srv, err := gmail.NewService(context.Background(), option.WithTokenSource(tokenSource))
	if err != nil {
		errMsg := "Unable to retrieve Gmail client: " + err.Error()
		fmt.Println(errMsg)
	}

	GmailService = srv
	if GmailService != nil {
		succMsg := "Email service is initialized"
		fmt.Printf("%v\n", succMsg)
	}
}

/* This creates the services for our calendar */
func OAuthCalendarService(){
	// If modifying these scopes, delete your previously saved token.json.
	config := oauth2.Config{
		ClientID: calendarPassing.GoogleClientID,
		ClientSecret: calendarPassing.GoogleClientSecret,
		Endpoint: google.Endpoint,
		RedirectURL: "http://localhost",
	}
	token := oauth2.Token{
		AccessToken: calendarPassing.GoogleClientCalendarAccessToken,
		RefreshToken: calendarPassing.GoogleClientCalendarRefreshToken,
		TokenType: "Bearer",
		Expiry: time.Now(),
	}
	var tokenSource = config.TokenSource(context.Background(), &token)

	srv, err := calendar.NewService(context.Background(), option.WithTokenSource(tokenSource))
	if err != nil {
		fmt.Printf("Unable to retrieve Calendar client: %v", err)
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	GoogleCalendarService = srv
	if GoogleCalendarService != nil {
		fmt.Printf("Google Calendar Initialized\n")
	}
}

/* This initializes Google Drive Service */
func OAuthGoogleDriveService(){
	// If modifying these scopes, delete your previously saved token.json.
	config := oauth2.Config{
		ClientID: calendarPassing.GoogleDriveClientID,
		ClientSecret:	calendarPassing.GoogleDriveClientSecret,
		Endpoint: google.Endpoint,
		RedirectURL: "http://localhost",
	}
	token := oauth2.Token{
		AccessToken: calendarPassing.GoogleDriveAccess,
		RefreshToken: calendarPassing.GoogleDriveRefresh,
		TokenType: "Bearer",
		Expiry: time.Now(),
	}
	var tokenSource = config.TokenSource(context.Background(), &token)

	srv, err := drive.NewService(context.Background(), option.WithTokenSource(tokenSource))
	if err != nil {
		fmt.Printf("Unable to retrieve drive client: %v", err)
		log.Fatalf("Unable to retrieve drive client: %v", err)
	}

	GoogleDriveService = srv
	if GoogleDriveService != nil {
		fmt.Printf("GoogleDrive Service initialized.\n")
	}
}