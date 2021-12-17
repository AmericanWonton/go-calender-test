var allDates = [];

function dateSetter(pIdentifierCode, pDateTimeStart, pDateTimeEnd, pDayNumber, 
    pHourStart, pHourEnd, pMonthNum, pYearNum, pFullDateDisplay,
    pApptTimeDisplay, pApptDateDisplay){
    var theDate = {
        IdentifierCode: String(pIdentifierCode),
        DateTimeStart: String(pDateTimeStart),
        DateTimeEnd: String(pDateTimeEnd),
        DayNumber: Number(pDayNumber),
        HourStart: String(pHourStart),
        HourEnd: String(pHourEnd),
        MonthNum: Number(pMonthNum),
        YearNum: Number(pYearNum),
        FullDateDisplay: String(pFullDateDisplay),
        ApptTimeDisplay: String(pApptTimeDisplay),
        ApptDateDisplay: String(pApptDateDisplay),
    }

    console.log("Here is this date: " + theDate.DateTimeStart + " " + theDate.DateTimeEnd);
    allDates.push(theDate);
}