/* Variables for Google Calendar */
let googleClientID = "";
let googleClientSecret = "";
let thegoogleCalendarAPIKey = "";
let calendarID = "";

/* Sets variables for Calendar */
function calendarVariableSet(gClientID, gClientSecretID, gCalendarAPIKey, idCalendar){
    googleClientID = gClientID;
    googleClientSecret = gClientSecretID;
    thegoogleCalendarAPIKey = gCalendarAPIKey;
    calendarID = idCalendar;
}


/* Initialize the 'FullCalendar' plugin into the 'calendar' div*/
document.addEventListener('DOMContentLoaded', function(){
    console.log("DEBUG: We can recognize fullcalendar");
    /* Documentation Refferences:
    initialView: How the calendar displays initially

    events: An array of events that will displayed on the calendar.
    Theoritically, you don't need to use Google Calendar as a plugin for this,
    just load as an array, then pass that in. NOTE, DO NOT 
    PUT COMMA AFTER LAST EVENT IN ARRAY, OR IE WILL HAVE ISSUES
    
    */
    var calendarEl = document.getElementById('calendar');
    var calendar = new FullCalendar.Calendar(calendarEl, {
        initialView: 'dayGridMonth',
        eventSources: [

            // your event source
            {
              events: [ // put the array in the `events` property
                {
                  title  : 'event1',
                  start  : '2021-10-01'
                },
                {
                  title  : 'event2',
                  start  : '2021-10-05',
                  end    : '2021-10-07'
                },
                {
                  title  : 'event3',
                  start  : '2021-10-09T12:30:00',
                },
                {
                    title  : 'event4',
                    start  : '2021-10-29T12:30:00',
                    end : '2021-10-30T12:30:00',
                    allDay : false // will make the time show
                }
              ],
              color: 'black',     // an option!
              textColor: 'yellow' // an option!
            },
        
            // any other event sources...
            {
                events: [ // put the array in the `events` property
                    {
                    title  : 'event5',
                    start  : '2021-10-15'
                    }
                ],
                color: 'black',     // an option!
                textColor: 'yellow' // an option!
            }
          ]
    });
    calendar.render();
    console.log("Initialize Calender");
});

