var allDates = [];

function dateSetter(pAllDay, pDateStart, pDateEnd, pDateTimeStart, pDateTimeEnd ){
    var theDate = {
        AllDay: Boolean(pAllDay),
        DateStart: String(pDateStart),
        DateEnd: String(pDateEnd),
        DateTimeStart: String(pDateTimeStart),
        DateTimeEnd: String(pDateTimeEnd)
    }

    console.log("Here is this date: " + theDate.DateEnd + " " + theDate.DateTimeStart);
    allDates.push(theDate);
}