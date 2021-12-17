package main

import (
	"time"
	"fmt"
	"log"
	"strconv"
	"strings"
)

/* Google Client/Calendar information */

type CalendarPassing struct {
	CalendarAPIKey     string `json:"CalendarAPIKey"`
	CalendarID         string `json:"CalendarID"`
	GoogleClientID     string `json:"GoogleClientID"`
	GoogleClientSecret string `json:"GoogleClientSecret"`
	GoogleClientCalendarRefreshToken string `json:"GoogleClientCalendarRefreshToken"`
	GoogleClientCalendarAccessToken string `json:"GoogleClientCalendarAccessToken"`
	CurrentEmail string `json:"CurrentEmail"`
	CurrentPWord string `json:"CurrentPWord"`
	EmailClient string `json:"EmailClient"`
	EmailSecret string `json:"EmailSecret"`
	EmailAccess string `json:"EmailAccess"`
	EmailRefresh string `json:"EmailRefresh"`
	GoogleDriveClientID string `json:"GoogleDriveClientID"`
	GoogleDriveClientSecret string `json:"GoogleDriveClientSecret"`
	GoogleDriveRefresh string `json:"GoogleDriveRefresh"`
	GoogleDriveAccess string `json:"GoogleDriveAccess"`
	CurrentTime        string `json:"CurrentTime"`
	PossibleDates []Appointment `json:"PossibleDates"`
	CalendarAllDatesFilled CalendarFilledDates `json:"CalendarAllDatesFilled"`
}

type CalendarFilledDates struct {
	CalendarDayFilled        []CalendarFilledDate `json:"CalendarDayFilled"`
}

type CalendarFilledDate struct {
	AllDay bool `json:"AllDay"`
	DateStart string `json:"DateStart"`
	DateEnd string `json:"DateEnd"`
	DateTimeStart string `json:"DateTimeStart"`
	DateTimeEnd string `json:"DateTimeEnd"`
}

//Date a User can schedule on our app that's available
type DateAvailable struct {
	DateTimeStart string `json:"DateTimeStart"`
	DateTimeEnd string `json:"DateTimeEnd"`
}

//IDK
type Appointment struct {
	IdentifierCode string 	`json:"IdentifierCode"`
	DateTimeStart string `json:"DateTimeStart"`
	DateTimeEnd string `json:"DateTimeEnd"`
	DayNumber int `json:"DayNumber"`
	HourStart string `json:"HourStart"`
	HourEnd string `json:"HourEnd"`
	MonthNum int 	`json:"MonthNum"`
	YearNum int 	`json:"YearNum"`
	FullDateDisplay string `json:"FullDateDisplay"`
	ApptTimeDisplay string `json:"ApptTimeDisplay"`
	ApptDateDisplay string `json:"ApptDateDisplay"`
}

var calendarPassing CalendarPassing
var dateFiller CalendarFilledDates

var potentialDates map[string]Appointment
var potentialDatesSlice []Appointment

const LAYOUTISO = "2006-01-02"
const LAYOUTUS = "January 2, 2006"

/* This function uses our Google Calendar to get dates
within 16 days and 
fill them for later scheudling use */
func getDatesForUse() []CalendarFilledDate{
	datesScheduled := []CalendarFilledDate{}

	startTime := time.Now().Format(time.RFC3339)
	endTime := time.Now().AddDate(0,0, 8 * 2).Format(time.RFC3339)
	events, err := GoogleCalendarService.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(startTime).TimeMax(endTime).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}
	
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		for _, item := range events.Items {
			date := item.Start.DateTime
			//Fill dates based on full day or not
			if date == "" {
				date = item.Start.Date
				gotDate := CalendarFilledDate{
					AllDay: true,
					DateStart: item.Start.Date,
					DateEnd: item.End.Date,
					DateTimeStart: item.Start.Date,
					DateTimeEnd: item.End.Date,
				}
				datesScheduled = append(datesScheduled, gotDate)
			} else {
				gotDate := CalendarFilledDate{
					AllDay: false,
					DateStart: item.Start.Date,
					DateEnd: item.End.Date,
					DateTimeStart: item.Start.DateTime,
					DateTimeEnd: item.End.DateTime,
				}
				datesScheduled = append(datesScheduled, gotDate)
			}
		}
	}

	return datesScheduled
}

/* This fills potential Schedules that will be blocked out later */
func fillPotentialAppointments(){
	potentialDates = make(map[string]Appointment)
	assembledDateTime := strconv.Itoa(time.Now().Year()) + "-" + strconv.Itoa(int(time.Now().Month())) + "-" +
	strconv.Itoa(time.Now().Day()) + "T09:00:00-06:00"
	startTime, err := time.Parse(time.RFC3339Nano, assembledDateTime)
	if err != nil {
		fmt.Printf("here is our big error: %v\n", err.Error())
	}
	endTime := startTime.AddDate(0,0, 8 * 2)

	//While StartTime is less than endTime
	for !startTime.After(endTime) {
		/* Initial check to see if date is Saturday or Sunday */
		if startTime.Weekday().String() == "Sunday" {
			//It's Sunday, add a day onto Monday
			startTime = startTime.AddDate(0,0, 1)
		} else if startTime.Weekday().String() == "Saturday" {
			//It's Saturday, add 2 days onto Monday
			startTime = startTime.AddDate(0,0, 2)
		} else {
			//It is a normal weekday, begin cycling through below
			for startTime.Hour() < 15 {
				endTimeInsert := startTime.Add(time.Hour * 1)
				/* Combination of the following: month num, -,  day, -, startTime Hour*/
				stringAssemble := convertCalendarMonth(startTime.Month().String()) + "-" + 
				convertCalendarDay(startTime.Day()) + 
				"-" + convertCalendarHour(startTime.Hour())
				/* Should be displayed as the following: 1999/08/25 from 09:00-13:00*/
				displayAssemble := strconv.Itoa(startTime.Year()) + "/" + 
				convertCalendarMonth(startTime.Month().String()) + "/" + convertCalendarDay(startTime.Day()) + 
				" from " + convertCalendarHour(startTime.Hour()) + ":00-" +
				convertCalendarHour(endTimeInsert.Hour()) + ":00"

				theAppointment := Appointment{
					IdentifierCode: stringAssemble,
					DateTimeStart: startTime.String(),
					DateTimeEnd: endTimeInsert.String(),
					DayNumber: startTime.Day(),
					HourStart: convertCalendarHour(startTime.Hour()),
					HourEnd: convertCalendarHour(endTimeInsert.Hour()),
					MonthNum: int(startTime.Month()),
					YearNum: startTime.Year(),
					FullDateDisplay: displayAssemble,
					ApptTimeDisplay: "From " + convertCalendarHour(startTime.Hour()) + ":00-" +
					convertCalendarHour(endTimeInsert.Hour()) + ":00",
					ApptDateDisplay: strconv.Itoa(startTime.Year()) + "/" + 
					convertCalendarMonth(startTime.Month().String()) + "/" + convertCalendarDay(startTime.Day()),
				}
				
				potentialDates[stringAssemble] = theAppointment
				potentialDatesSlice = append(potentialDatesSlice, theAppointment)

				startTime = startTime.Add(time.Hour	* 1) //Add an hour for the loop
			}
			startTime = startTime.AddDate(0,0,1)
		}
		
		
		//Set time back to 9am
		calendarMonthStr := convertCalendarMonth(startTime.Month().String())
		calendarDayStr := convertCalendarDay(startTime.Day())
		assembledDateTimeTwo := strconv.Itoa(startTime.Year()) + "-" + calendarMonthStr + "-" +
		calendarDayStr + "T09:00:00-06:00"
		startTime, err = time.Parse(time.RFC3339Nano, assembledDateTimeTwo)
		if err != nil {
			theError := err.Error() //Debug
			fmt.Printf("Big error with start time: %v\n", theError)
		}
	}

	/* DEBUG List our created days */
	/*
	for x := 0; x < len(potentialDatesSlice); x++{
		fmt.Printf("Here is the key: %v\n", potentialDatesSlice[x].IdentifierCode)
	}
	*/
	/*
	for key, _ := range potentialDates{
		fmt.Printf("Here is our key: %v\n", key)
	}
	*/
}

/* Remove Certain Appointments based on dateFiller and potentialDatesSlice*/
func removeCertainAppointments(){
	
	/* Remove all events for today,(Users can't schedule for today */
	for w := 0; w < len(dateFiller.CalendarDayFilled); w++ {
		//Get Current Day
		theCurrentDate := time.Now()
		theMonthTime := convertCalendarMonth(theCurrentDate.Month().String())
		theDayTime := convertCalendarDay(theCurrentDate.Day())

		/* Split each code within the array, see if they have the month
		and the same day; then remove them */
		for l := 0; l < len(potentialDatesSlice); l++ {
			fullStringSplit := strings.Split(potentialDatesSlice[l].IdentifierCode, "-")
			theMonth := fullStringSplit[0]
			theDay := fullStringSplit[1]

			if (strings.ToLower(theDay) == strings.ToLower(theDayTime)) && 
			(strings.ToLower(theMonth) == strings.ToLower(theMonthTime)){
				//The Day and Month match, remove them from the slice
				removePlace := l
				copy(potentialDatesSlice[removePlace:], potentialDatesSlice[removePlace + 1:])
				potentialDatesSlice[len(potentialDatesSlice) - 1] = Appointment{}
				potentialDatesSlice = potentialDatesSlice[:len(potentialDatesSlice) -1]
			}
		}
	}

	for n := 0; n < len(dateFiller.CalendarDayFilled); n++ {
		/* Remove all day events first */
		if dateFiller.CalendarDayFilled[n].AllDay {
			//Get Day and Month of Google Calendar Event
			theGoogleDate, err := time.Parse(LAYOUTISO, dateFiller.CalendarDayFilled[n].DateStart)
			if err != nil {
				fmt.Printf("Bad date error from Google Parse: %v\n", err.Error())
			}
			
			
			theMonthGoogle := convertCalendarMonth(theGoogleDate.Month().String())
			theDayGoogle := convertCalendarDay(theGoogleDate.Day())

			/* Split each code within the array, see if they have the month
			and the same day; then remove them */
			for l := 0; l < len(potentialDatesSlice); l++ {
				fullStringSplit := strings.Split(potentialDatesSlice[l].IdentifierCode, "-")
				theMonth := fullStringSplit[0]
				theDay := fullStringSplit[1]

				if (strings.ToLower(theDay) == strings.ToLower(theDayGoogle)) && 
				(strings.ToLower(theMonth) == strings.ToLower(theMonthGoogle)){
					//The Day and Month match, remove them from the slice
					removePlace := l
					copy(potentialDatesSlice[removePlace:], potentialDatesSlice[removePlace + 1:])
					potentialDatesSlice[len(potentialDatesSlice) - 1] = Appointment{}
					potentialDatesSlice = potentialDatesSlice[:len(potentialDatesSlice) -1]
				}
			}
			
		} else {
			/* Remove Events within a certain day and start time */
			theGoogleDate, err := time.Parse(time.RFC3339Nano, dateFiller.CalendarDayFilled[n].DateTimeStart)
			if err != nil {
				fmt.Printf("Bad date error from Google Parse: %v\n", err.Error())
			}

			theGoogleEndDate, err2 := time.Parse(time.RFC3339Nano, dateFiller.CalendarDayFilled[n].DateTimeEnd)
			if err2 != nil {
				fmt.Printf("Bad date error from Google Parse: %v\n", err2.Error())
			}

			theMonthGoogle := convertCalendarMonth(theGoogleDate.Month().String())
			theDayGoogle := convertCalendarDay(theGoogleDate.Day())
			theStartHour := convertCalendarHour(theGoogleDate.Hour())

			for n := 0; n < len(potentialDatesSlice); n++ {
				fullStringSplit := strings.Split(potentialDatesSlice[n].IdentifierCode, "-")
				theMonth := fullStringSplit[0]
				theDay := fullStringSplit[1]
				theApptStartHour := fullStringSplit[2]
				if (strings.ToLower(theDay) == strings.ToLower(theDayGoogle)) && 
				(strings.ToLower(theMonth) == strings.ToLower(theMonthGoogle)) {
					//Second check to make sure they fall within the same time
					if strings.ToLower(theStartHour) == strings.ToLower(theApptStartHour){
						/* Need to remove EACH potential Appointment that could fall between this time */
						theStartNum := theGoogleDate.Hour()
						theEndNum := theGoogleEndDate.Hour()
						//fmt.Printf("DEBUG: Removing hours between %v and %v on this day: %v\n", theStartNum, theEndNum, theDay)
						//Start with current Appointment
						removePlace := n
						copy(potentialDatesSlice[removePlace:], potentialDatesSlice[removePlace + 1:])
						potentialDatesSlice[len(potentialDatesSlice) - 1] = Appointment{}
						potentialDatesSlice = potentialDatesSlice[:len(potentialDatesSlice) -1]
						//Loop through each hour until every end hour reached
						for theStartNum < theEndNum {
							removePlace := n
							copy(potentialDatesSlice[removePlace:], potentialDatesSlice[removePlace + 1:])
							potentialDatesSlice[len(potentialDatesSlice) - 1] = Appointment{}
							potentialDatesSlice = potentialDatesSlice[:len(potentialDatesSlice) -1]
							theStartNum++ //Increment the StartNum for loop
						}
					}
				}
			}
		}

	}
	
	/* DEBUG PRINT */
	/*
	for v := 0; v < len(potentialDatesSlice); v++ {
		fmt.Printf("Here is a appointment Start Day: %v and End Day: %v\n", potentialDatesSlice[v].DateTimeStart,
		potentialDatesSlice[v].DateTimeEnd)
	}
	*/
}

/* Int conversion for Month */
func convertCalendarMonth(theMonth string) string{
	theStringReturn := "01"
	switch theMonth{
	case "January": 
		theStringReturn = "01"
		break
	case "February":
		theStringReturn = "02"
		break
	case "March":
		theStringReturn = "03"
		break
	case "April":
		theStringReturn = "04"
		break
	case "May":
		theStringReturn = "05"
		break
	case "June":
		theStringReturn = "06"
		break
	case "July":
		theStringReturn = "07"
		break
	case "August":
		theStringReturn = "08"
		break
	case "September":
		theStringReturn = "09"
		break
	case "October":
		theStringReturn = "10"
		break
	case "November":
		theStringReturn = "11"
		break
	case "December":
		theStringReturn = "12"
		break
	default:
		fmt.Printf("DEBUG: Big error problem, here is month: %v\n", theMonth)
		break
	}

	return theStringReturn
}
/* Int conversion for Day */
func convertCalendarDay(theDay int) string{
	returnedStringDay := "01"

	switch theDay{
		case 1:
			returnedStringDay = "01"
			break
		case 2:
			returnedStringDay = "02"
			break
		case 3:
			returnedStringDay = "03"
			break
		case 4:
			returnedStringDay = "04"
			break
		case 5:
			returnedStringDay = "05"
			break
		case 6:
			returnedStringDay = "06"
			break
		case 7:
			returnedStringDay = "07"
			break
		case 8:
			returnedStringDay = "08"
			break
		case 9:
			returnedStringDay = "09"
			break
		default:
			returnedStringDay = strconv.Itoa(theDay)
			break
	}

	return returnedStringDay
}

/* Int Conversion for hour*/
func convertCalendarHour(theHour int) string{
	theReturnHour := "09"
	
	switch theHour {
		case 0:
			theReturnHour = "00"
			break
		case 1:
			theReturnHour = "01"
			break
		case 2:
			theReturnHour = "02"
			break
		case 3:
			theReturnHour = "03"
			break
		case 4:
			theReturnHour = "04"
			break
		case 5:
			theReturnHour = "05"
			break
		case 6:
			theReturnHour = "06"
			break
		case 7:
			theReturnHour = "07"
			break
		case 8:
			theReturnHour = "08"
			break
		case 9:
			theReturnHour = "09"
			break
		case 10:
			theReturnHour = "10"
			break
		case 11:
			theReturnHour = "11"
			break
		case 12:
			theReturnHour = "12"
			break
		case 13:
			theReturnHour = "13"
			break
		case 14:
			theReturnHour = "14"
			break
		case 15:
			theReturnHour = "15"
			break
		case 16:
			theReturnHour = "16"
			break
		case 17:
			theReturnHour = "17"
			break
		case 18:
			theReturnHour = "18"
			break
		case 19:
			theReturnHour = "19"
			break
		case 20:
			theReturnHour = "20"
			break
		case 21:
			theReturnHour = "21"
			break
		case 22:
			theReturnHour = "22"
			break
		case 23:
			theReturnHour = "23"
			break
		case 24:
			theReturnHour = "24"
			break
		default:
			fmt.Printf("Something is messed up: %v\n", theHour)
			break
	}
	return theReturnHour
}

/* This evaluates Dates that are available based on 16 days out
and returns an array of date times a User can schedule */
func makeScheduleDates(dateSchedules []CalendarFilledDate)[]DateAvailable{
	datesAvailable := []DateAvailable{}

	//Make the day line up for today, at 9AM
	assembledDateTime := strconv.Itoa(time.Now().Year()) + "-" + strconv.Itoa(int(time.Now().Month())) + "-" +
	strconv.Itoa(time.Now().Day()) + "T09:00:00-06:00"
	finalTime, err := time.Parse(time.RFC3339Nano, assembledDateTime)
	if err != nil {
		fmt.Printf("here is our big error: %v\n", err.Error())
	}
	dayMover := finalTime

	/* Initial check to see if date is Saturday or Sunday */
	if dayMover.Weekday().String() == "Sunday" {
		//It's Sunday, add a day onto Monday
		dayMover = time.Now().AddDate(0,0, 1)
	} else if dayMover.Weekday().String() == "Saturday" {
		//It's Saturday, add 2 days onto Monday
		dayMover = time.Now().AddDate(0,0, 2)
	} else {
		//It is a normal weekday, begin cycling through below
	}

	//Create days we can't work with, blocked out, whatever and days that are scheduled
	arrayOBadDays := make(map[int]int) //Blocked Out Days
	arrayODaysScheduled := make(map[int]int) //Days with appointments on them
	mapOApptsInDays := make(map[int][]CalendarFilledDate)
	for j := 0; j < len(dateSchedules); j++ {

		if dateSchedules[j].AllDay == true {
			fmt.Printf("Found an all day event on this day: %v\n", dayMover.Day())
			arrayOBadDays[dayMover.Day()] = dayMover.Day()
			break
		}
		//Set the dateSchedule to a variable to work with
		dateWorkingWith, err := time.Parse(time.RFC3339Nano, dateSchedules[j].DateTimeStart)
		if err != nil {
			fmt.Printf("Error parsing time: %v\n", err.Error())
			break
		}
		//Add Date to list of days that have time schedules
		if _, ok := arrayODaysScheduled[dateWorkingWith.Day()]; ok {
			//Date already in map, don't worry about it
		} else {
			//Add Date Integer to map
			arrayODaysScheduled[dateWorkingWith.Day()] = dateWorkingWith.Day()
			//Add this specific dateScheudle to our array to work with
			mapOApptsInDays[dateWorkingWith.Day()] = append(mapOApptsInDays[dateWorkingWith.Day()], dateSchedules[j])
		}
	}

	//Loop through Days to create Schedule
	for dayMover.Day() >= time.Now().AddDate(0,0,8 * 2).Day(){
		//Makes sure day dosen't fall on bad day
		if _, ok := arrayOBadDays[dayMover.Day()]; ok {
			//Day is bad, skip it
			if dayMover.Weekday().String() == "Friday" {
				dayMover.AddDate(0,0, 3)
			} else {
				//Done with day not on friday, add one plus day
				dayMover.AddDate(0,0,1)
			}
			break
		} else {
			//Day not bad, go on scheduling
			//If Day has appointments in it, process schedules with care
			if _, okay := arrayODaysScheduled[dayMover.Day()]; okay{
				//This day has scheduled appointments in here, be careful
				//currentHour := dayMover.Hour() //Should be 9
				//Loop through all appointments for this day
				

			} else {
				//This day has no scheduled apppointments
				datesAvailable = append(datesAvailable, addWholeDayScheduling(dayMover)...)
			}

			//Done scheduling for the day; make sure we go past Saturday and Sunday
			if dayMover.Weekday().String() == "Friday" {
				dayMover.AddDate(0,0, 3)
				newDateTime := strconv.Itoa(dayMover.Year()) + "-" + strconv.Itoa(int(dayMover.Month())) + "-" +
				strconv.Itoa(dayMover.Day()) + "T09:00:00-06:00"
				thefinalTime, err := time.Parse(time.RFC3339Nano, newDateTime)
				if err != nil {
					fmt.Printf("here is our big error: %v\n", err.Error())
				}
				dayMover = thefinalTime
			} else {
				//Done with day not on friday, add one plus day
				dayMover.AddDate(0,0,1)
				newDateTime := strconv.Itoa(dayMover.Year()) + "-" + strconv.Itoa(int(dayMover.Month())) + "-" +
				strconv.Itoa(dayMover.Day()) + "T09:00:00-06:00"
				thefinalTime, err := time.Parse(time.RFC3339Nano, newDateTime)
				if err != nil {
					fmt.Printf("here is our big error: %v\n", err.Error())
				}
				dayMover = thefinalTime
			}
		}
		
	}
	

	return datesAvailable
}

/* Adds a whole day of open scheduling */
func addWholeDayScheduling(theTime time.Time)[]DateAvailable{
	theDayAvailable := []DateAvailable{}

	//Make the day line up for today, at 9AM
	//9AM - 10AM
	assembledTime := "T09:00:00-06:00"
	assembledEndTime := "T10:00:00-06:00"
	assembledDateTime := strconv.Itoa(time.Now().Year()) + "-" + strconv.Itoa(int(time.Now().Month())) + "-" +
	strconv.Itoa(time.Now().Day())
	dateAvailable := DateAvailable{DateTimeStart: assembledDateTime + assembledTime, 
		DateTimeEnd: assembledDateTime + assembledEndTime}
	theDayAvailable = append(theDayAvailable, dateAvailable)
	//10AM-11AM
	assembledTime = "T10:00:00-06:00"
	assembledEndTime = "T11:00:00-06:00"
	dateAvailable = DateAvailable{DateTimeStart: assembledDateTime + assembledTime, 
		DateTimeEnd: assembledDateTime + assembledEndTime}
	theDayAvailable = append(theDayAvailable, dateAvailable)
	//11AM-12PM
	assembledTime = "T11:00:00-06:00"
	assembledEndTime = "T12:00:00-06:00"
	dateAvailable = DateAvailable{DateTimeStart: assembledDateTime + assembledTime, 
		DateTimeEnd: assembledDateTime + assembledEndTime}
	theDayAvailable = append(theDayAvailable, dateAvailable)
	//12PM-1PM
	assembledTime = "T12:00:00-06:00"
	assembledEndTime = "T13:00:00-06:00"
	dateAvailable = DateAvailable{DateTimeStart: assembledDateTime + assembledTime, 
		DateTimeEnd: assembledDateTime + assembledEndTime}
	theDayAvailable = append(theDayAvailable, dateAvailable)
	//1PM-2PM
	assembledTime = "T13:00:00-06:00"
	assembledEndTime = "T14:00:00-06:00"
	dateAvailable = DateAvailable{DateTimeStart: assembledDateTime + assembledTime, 
		DateTimeEnd: assembledDateTime + assembledEndTime}
	theDayAvailable = append(theDayAvailable, dateAvailable)
	//2PM-3PM
	assembledTime = "T14:00:00-06:00"
	assembledEndTime = "T15:00:00-06:00"
	dateAvailable = DateAvailable{DateTimeStart: assembledDateTime + assembledTime, 
		DateTimeEnd: assembledDateTime + assembledEndTime}
	theDayAvailable = append(theDayAvailable, dateAvailable)
	//3PM-4PM
	assembledTime = "T15:00:00-06:00"
	assembledEndTime = "T16:00:00-06:00"
	dateAvailable = DateAvailable{DateTimeStart: assembledDateTime + assembledTime, 
		DateTimeEnd: assembledDateTime + assembledEndTime}
	theDayAvailable = append(theDayAvailable, dateAvailable)
	//4PM-5PM
	assembledTime = "T16:00:00-06:00"
	assembledEndTime = "T17:00:00-06:00"
	dateAvailable = DateAvailable{DateTimeStart: assembledDateTime + assembledTime, 
		DateTimeEnd: assembledDateTime + assembledEndTime}
	theDayAvailable = append(theDayAvailable, dateAvailable)

	return theDayAvailable
}