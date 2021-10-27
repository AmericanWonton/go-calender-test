package main

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
)

func createCalendar() {
	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))

	calendarService, err := calendar.New(client)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
}
