package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Day struct {
	Date string `json:"Date"`
	Day  string `json:"Day"`
}

var days = make([]Day, 0)

func main() {
	handleRequests()
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", defaultPage)
	router.HandleFunc("/date/", useCurrentDate).Methods("GET")
	router.HandleFunc("/date/{date}", useSpecificDate).Methods("GET") //date format should be YYYYMMDD

	log.Fatal(http.ListenAndServe(":9090", router))
}

//This page will be the landing page
func defaultPage(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w,"Instructions:\n")
	fmt.Fprintf(w,"To get the next seven days from CURRENT DATE, enter http://localhost:9090/date/.\n")
	fmt.Fprintf(w,"To get the next seven days from a SPECIFIC DATE, enter http://localhost:9090/date/{date} where {date} should be in format YYYYMMDD.")
}

//This page will display the response schema if a user didn't specify a date in the request
func useCurrentDate(w http.ResponseWriter, r *http.Request){
	retrieveNextSevenDays(time.Now())
	prettyPrintJson(days, w)

	//Reset or clear the list
	days = nil
}

//This page will display the response schema if a user specified a date in the request
func useSpecificDate(w http.ResponseWriter, r *http.Request){
	//Validate if the parameter is a valid date
	date, err := time.Parse("20060102", mux.Vars(r)["date"])
	if (err == nil){
		retrieveNextSevenDays(date)
		prettyPrintJson(days, w)
		
		//Reset or clear the list
		days = nil
	} else {
		w.WriteHeader(400)
		fmt.Fprintf(w,  "Invalid date.")
	}
	
}

//This will retrieve the next seven days from the start date specified.
func retrieveNextSevenDays(startingDate time.Time){
	index := 0
	
	for{
		index +=1
		if(index > 7){
			break
		}
		
		//Add number of days from the starting date
		nextDay := startingDate.AddDate(0, 0, index)
		
		nextDate := string(nextDay.Format("January 2, 2006"))

		//Get weekday of the date
		weekDay := nextDay.Weekday().String()

		var nextDays = Day{
				Date: nextDate,
				Day: weekDay,
			}

		//Add the days to the list of days
		days = append(days, nextDays)
	}
}

//This will format the json data to be more readable
func prettyPrintJson(data interface{}, w http.ResponseWriter){
	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	enc.Encode(data)
}