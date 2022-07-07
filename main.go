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
	router.HandleFunc("/date/", useCurrentDate)
	router.HandleFunc("/date/{date}", useSpecificDate) //date format should be YYYYMMDD

	log.Fatal(http.ListenAndServe(":9090", router))
}

func defaultPage(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w,"Instructions:\n")
	fmt.Fprintf(w,"To get the next seven days from CURRENT DATE, enter http://localhost:9090/date/\n")
	fmt.Fprintf(w,"To get the next seven days from a SPECIFIC DATE, enter http://localhost:9090/date/{date} where {date} should be in format YYYYMMDD")
}

func useCurrentDate(w http.ResponseWriter, r *http.Request){
	retrieveNextSevenDays(time.Now())
	prettyPrintJson(days, w)
	days = nil
}

func useSpecificDate(w http.ResponseWriter, r *http.Request){
	date, err := time.Parse("20060102", mux.Vars(r)["date"])
	if (err == nil){
		retrieveNextSevenDays(date)
		prettyPrintJson(days, w)
		days = nil
	} else {
		fmt.Fprintf(w,  "Invalid date.")
	}
	
}

func retrieveNextSevenDays(startingDate time.Time){
	index := 0
	
	for{
		index +=1
		if(index > 7){
			break
		}
		
		nextDay := startingDate.AddDate(0, 0, index)
		
		nextDate := string(nextDay.Format("January 2, 2006"))
		weekDay := nextDay.Weekday().String()

		var nextDays = Day{
				Date: nextDate,
				Day: weekDay,
			}

		days = append(days, nextDays)
	}
}

func prettyPrintJson(data interface{}, w http.ResponseWriter){
	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	enc.Encode(data)
}