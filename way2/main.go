package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/drive/v3"
)

/* TEMPLATE DEFINITION */
var template1 *template.Template

/* Webpage information passing */
type ViewData struct {
	PassedCalendarInfo CalendarPassing `json:"PassedCalendarInfo"`
}

func init() {
	template1 = template.Must(template.ParseGlob("./static/templates/*")) //pass templates
	potentialDates = make(map[string]Appointment)
	getCalendarCreds()                                                    //Get calendar creds
	OAuthGmailService() //Initialize Email 
	OAuthCalendarService() //Initialize Calendar 
	OAuthGoogleDriveService() //Initilaize Google Drive
}

//Handles the Index requests; Ask User if they're legal here
func index(w http.ResponseWriter, r *http.Request) {
	//Build info to pass
	//Rest Calendar passed dates
	calendarPassing.CalendarAllDatesFilled.CalendarDayFilled = getDatesForUse()
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

	fmt.Printf("DEBUG: We are serving files on internet 5000\n")
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
	
	potentialDates = make(map[string]Appointment)
	assembledDateTime := strconv.Itoa(time.Now().Year()) + "-" + strconv.Itoa(int(time.Now().Month())) + "-" +
	strconv.Itoa(time.Now().Day()) + "T09:00:00-06:00"
	startTime, err := time.Parse(time.RFC3339Nano, assembledDateTime)
	if err != nil {
		fmt.Printf("here is our big error: %v\n", err.Error())
	}
	endTime := startTime.AddDate(0,0, 8 * 2)

	fmt.Printf("Starttime is: %v\n EndTime is %v\n", startTime, endTime)

	startTime = startTime.Add(time.Hour * 2)
	fmt.Printf("DEBUG: End time is: %v\n", startTime)

	fillPotentialAppointments()

	//handleRequests() // handle requests
}

/* This fills our dates; called everytime the index page is meant to be loaded */
func fillCalendarDates()CalendarFilledDates{
	theReturnedFilledDates:= CalendarFilledDates{}

	theTimeNow := time.Now().Format(time.RFC3339Nano)

	fmt.Printf("Here is the time now: %v\n", theTimeNow)

	//Time 15 days from now
	theTimeFifDays := time.Now().AddDate(0,0, 8 * 2).Format(time.RFC3339Nano)
	fmt.Printf("Here is the time in two weeks: %v\n", theTimeFifDays)

	return theReturnedFilledDates
}

/* This does all the fun Google Calender reading */
func googleCalendarReadTest() {
	currentTime := time.Now() //Used for debugging
	fmt.Println("Here is the current Google time in YYYY-MM-DD : ", currentTime.Format("2006-01-02"))

	t := time.Now().Format(time.RFC3339)
	events, err := GoogleCalendarService.Events.List("primary").ShowDeleted(false).
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
			//This appears to be for all day items
			if date == "" {
				date = item.Start.Date
			}
			fmt.Printf("Here is Item.Summary: %v\nHere is item.Description: %v\n" +
			"Here is item.Start.DateTime: %v\nHere is item.End.DateTime: %v\n" +
			"Here is item.Start.Date: %v\nHere is item.End.Date: %v\n",
			item.Summary, item.Description, item.Start.DateTime, item.End.DateTime,
			item.Start.Date, item.End.Date)
			fmt.Println()
		}
	}
}

/* This is another test function for Inserting a Google Calendar Event*/
func googleCalendarInsertTestTheSecond() {

	/* Insert File into Google Drive */
	theFileURL, theFileID, theFileTitle, theFileMimeType, anErr, msgs := insertMeetingAttachment()
	if !anErr {
		//Error with making attachment
		fmt.Printf("There was an error creating document for Google Drive:\n")
		for n := 0; n < len(msgs); n++ {
			fmt.Printf("%v\n", msgs[n])
		}
	}
	/* Attachments */
	var theAttachment1 = &calendar.EventAttachment{}
	(*theAttachment1).FileUrl = theFileURL
	(*theAttachment1).FileId = theFileID
	(*theAttachment1).Title = theFileTitle
	(*theAttachment1).MimeType = theFileMimeType

	/* DEBUG PRINT */
	fmt.Printf("DEBUG: The fileURL is: %v\nThe File ID is: %v\n", theFileURL, theFileID)

	var theAttachments []*calendar.EventAttachment
	theAttachments = append(theAttachments, theAttachment1)

	/* Event Attendees */
	var theAttendee = &calendar.EventAttendee{
		Email:       "johnnycowboy39@gmail.com",
		Optional:    true,
		DisplayName: "Some Name",
		Comment:     "This is a comment from an attendee",
	}

	var theAttendees []*calendar.EventAttendee
	theAttendees = append(theAttendees, theAttendee)

	theEvent := &calendar.Event{
		Start: &calendar.EventDateTime{
			DateTime: "2021-11-24T17:06:02.000Z",
			TimeZone: "America/Chicago",
		},
		End: &calendar.EventDateTime{
			DateTime: "2021-11-24T19:06:02.000Z",
			TimeZone: "America/Chicago",
		},
		Summary:     "Test Calendar Creation",
		Location:    "800 Howard St., San Francisco, CA 94103",
		Description: "A test Google Calendar date, created from Golang, with a file...",
		Attachments: theAttachments,
		Attendees:   theAttendees,
	}

	calendarId := "primary"
	event, err2 := GoogleCalendarService.Events.Insert(calendarId, theEvent).SupportsAttachments(true).Do()
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
func insertMeetingAttachment() (string, string, string, string, bool, []string) {

	//Define errors to return and messages
	theBigErr, errMsgss := false, []string{}

	wd, _ := os.Getwd()

	/* Step 1 open file for working with */
	fileDir := filepath.Join(wd, "testFileUploads", "testfile.txt")
	f, err := os.Open(fileDir)

	if err != nil {
		panic(fmt.Sprintf("cannot open file: %v", err))
	}

	defer f.Close()

	// Step 3. Create the directory
	dir, err2 := createDir(GoogleDriveService, "testGoogleDriveFolder", "root")
	if err2 != nil {
		panic(fmt.Sprintf("Uh oh, couldn't create directory: %v\n", err2.Error()))
	}

	// Step 4. Create the file and upload its content
	file, err := createFile(GoogleDriveService, "testfile.txt", "text/plain", f, dir.Id)

	if err != nil {
		panic(fmt.Sprintf("Could not create file: %v\n", err))
	}

	fmt.Printf("File '%s' successfully uploaded in '%s' directory\n", file.Name, dir.Name)

	time.Sleep(2 * time.Second) //DEBUG WAIT

	//Step 5 edit the permissions for this file so others can download/use it
	createPermissionsGoogleAPI(GoogleDriveService, dir.Id, "anyone", "reader")  //For the folder
	createPermissionsGoogleAPI(GoogleDriveService, file.Id, "anyone", "reader") //For the file

	//Get Google Drive File Info
	anErr, theFile := getDriveFileInfo(GoogleDriveService, file.Id)
	if anErr != nil {
		panic(fmt.Sprintf("There was an error getting fileinfo: %v\n", anErr.Error()))
	}

	return theFile.WebViewLink, theFile.Id, file.Name, theFile.MimeType, theBigErr, errMsgss
}

/* This creates a test Google Calendar Event */
func googleCalendarCreateEventTest() {

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
	event, err2 := GoogleCalendarService.Events.Insert(calendarId, theEvent).Do()
	if err2 != nil {
		fmt.Printf("Unable to create event: %v\n", err2)
		log.Fatalf("Unable to create event. %v\n", err2)
	}
	fmt.Printf("Event created: %s\n", event.HtmlLink)
}

/* This reads all files in a certain Google Drive Directory */
func googleDriveList(service *drive.Service) {

	r, err := service.Files.List().PageSize(400).
		Fields("nextPageToken, files(id, name, webViewLink, mimeType)").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}
	fmt.Println("Files:")
	if len(r.Files) == 0 {
		fmt.Println("No files found.")
	} else {
		for _, i := range r.Files {
			fmt.Printf("Here is the file name: %v\nHere is the file ID: %v\nHere is the DriveID: %v\n"+
				"Here is the description: %v\nHere is the WebViewLink: %v\nHere is the webContentLink: %v\n"+
				"Here is owners: %v\nHere is permissions: %v\nHere is permisssionID: %v\n\n",
				i.Name, i.Id, i.DriveId, i.Description, i.WebViewLink, i.WebContentLink, i.Owners, i.Permissions, i.PermissionIds)
		}
	}
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
	//NOTE, you can pass in the parent ID if you want the file placed somewhere
	//I do NOT do that becuase it was a pain to mess with permission IDS with folders in certain heirachies
	f := &drive.File{
		MimeType:                     mimeType,
		Name:                         name,
		Parents:                      []string{parentId},
		Description:                  "A test file created for testing",
		CopyRequiresWriterPermission: false,
		DriveId:                      "12345",
	}
	file, err := service.Files.Create(f).Media(content).SupportsAllDrives(true).Do()

	if err != nil {
		log.Println("Could not create file: " + err.Error())
		return nil, err
	}

	fmt.Printf("DEBUG: The id is:%v\n", file.Id)
	//fmt.Printf("DEBUG: The WebContentLink is: %v\n", file.WebContentLink)
	//fmt.Printf("DEBUG: The WebViewLink is: %v\n", file.WebViewLink)

	return file, nil
}

/* Create permissions for uploaded Google file */
func createPermissionsGoogleAPI(service *drive.Service, theFileID string, permType string, role string) error {
	/* There's other fields here but they cause a write error */
	p := &drive.Permission{
		Type: permType,
		Role: role,
	}

	donePermissions, err := service.Permissions.Create(theFileID, p).Do()
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		return err
	} else {
		fmt.Printf("Permissions added successfully: %v\n", donePermissions.Id)
	}

	return nil
}

/* Gets a test Google Drive file stuff to attach to a meeting */
func getDriveFileInfo(d *drive.Service, fileId string) (error, *drive.File) {
	f, err := d.Files.Get(fileId).Fields("webViewLink, webContentLink, id, mimeType").Do()
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		return err, nil
	}
	fmt.Printf("Description: %v\n", f.Description)
	fmt.Printf("MIME type: %v\n", f.MimeType)
	fmt.Printf("Here is the driveID: %v\n", f.DriveId)
	fmt.Printf("Here is the id: %v\n", f.Id)
	fmt.Printf("Here is the permissionID: %v\n", f.PermissionIds)

	fmt.Printf("Here is the webContentLink: %v\n", f.WebContentLink)
	fmt.Printf("Here is the WebViewLink: %v\n", f.WebViewLink)

	return nil, f
}
