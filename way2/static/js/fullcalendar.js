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
/*
document.addEventListener('DOMContentLoaded', function(){
    console.log("DEBUG: We can recognize fullcalendar");
    /* Documentation Refferences:
    initialView: How the calendar displays initially

    events: An array of events that will displayed on the calendar.
    Theoritically, you don't need to use Google Calendar as a plugin for this,
    just load as an array, then pass that in. NOTE, DO NOT 
    PUT COMMA AFTER LAST EVENT IN ARRAY, OR IE WILL HAVE ISSUES
    
    */
   /*
    var calendarEl = document.getElementById('calendar');
    var calendar = new FullCalendar.Calendar(calendarEl, {
        initialView: 'dayGridMonth',
        eventSources: [

            // your event source
            {
              events: [ // put the array in the `events` property
                {
                  title  : 'event1',
                  start  : '2021-12-01'
                },
                {
                  title  : 'event2',
                  start  : '2021-12-05',
                  end    : '2021-12-07'
                },
                {
                  title  : 'event3',
                  start  : '2021-12-09T12:30:00',
                },
                {
                  title  : 'event7',
                  start  : '2021-12-29T08:30:00',
                  end : '2021-12-29T12:00:00',
                  allDay : false // will make the time show
                },
                {
                    title  : 'event4',
                    start  : '2021-12-29T12:30:00',
                    end : '2021-12-29T13:30:00',
                    allDay : false // will make the time show
                },
                {
                  title  : 'event6',
                  start  : '2021-12-29T14:30:00',
                  end : '2021-12-29T15:00:00',
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
                    start  : '2021-12-15'
                    }
                ],
                color: 'black',     // an option!
                textColor: 'yellow' // an option!
            }
          ],
        businessHours: [ // specify an array instead
          {
            daysOfWeek: [ 1, 2, 3 ], // Monday, Tuesday, Wednesday
            startTime: '08:00', // 8am
            endTime: '18:00' // 6pm
          },
          {
            daysOfWeek: [ 4, 5 ], // Thursday, Friday
            startTime: '10:00', // 10am
            endTime: '16:00' // 4pm
          }
        ],
        displayEventEnd: true,
        displayEventTime: true
    });
    //Handler for clicking dates
    calendar.on('dateClick', function(info){
      console.log("Clicked on: " + info.dateStr);
    });
    //Set Height Initially
    //calendar.setOption('height', 650);
    //Sets aspect ratio initially
    calendar.setOption('aspectRatio', 4.0);

    calendar.render(); //Render the calendar
    console.log("Initialize Calender");
});
*/
