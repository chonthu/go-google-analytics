package utils

import (
	"fmt"
	"strings"
	"time"
)

func GenWeekRanges(start int, end int, year int) (out []Daterange) {
	store := []string{}
	for week := start; week < end; week++ {
		tmp, tmp1 := checkWeek(year, week)
		store = append(store, tmp)
		store = append(store, tmp1)
	}

	tmp := []string{}
	for i, b := range store[1:] {
		tmp = append(tmp, b)
		if (1+i)%2 == 0 {
			out = append(out, Daterange{strings.Trim(tmp[0], " "),
				strings.Trim(tmp[1], " "),
				"",
				[]*Daterange{},
				false})
			tmp = []string{}
		}
	}
	return
}

func GenDayRanges(start string, end string) (out []Daterange) {
	formatD := "2006-01-02"
	start_date, _ := time.Parse(formatD, start)
	fmt.Println(start_date)
	end_date, _ := time.Parse(formatD, end)
	fmt.Println(end_date)
	for !start_date.Equal(end_date) {
		tmp := start_date.Add(24 * time.Hour)
		fmt.Println(tmp)
		out = append(out, Daterange{start_date.Format(formatD),
			tmp.Format(formatD),
			"",
			[]*Daterange{},
			false})
		start_date = tmp
	}
	return
}

func checkWeek(year int, week int) (string, string) {
	date := firstDayOfISOWeek(year, week, time.UTC)
	// sanity check
	isoYear, isoWeek := date.ISOWeek()
	if year != isoYear || week != isoWeek {
		panic(fmt.Sprintf("Input: year %v, week %v. Result: %v, year %v, week %v\n", year, week, date, isoYear, isoWeek))
	}
	return strings.Split(fmt.Sprintf("%v\n", date.Add(time.Minute*60*24)), "00")[0], strings.Split(fmt.Sprintf("%v\n", date), "00")[0]
}

func firstDayOfISOWeek(year int, week int, timezone *time.Location) time.Time {
	date := time.Date(year, 0, 0, 0, 0, 0, 0, timezone)
	isoYear, isoWeek := date.ISOWeek()
	for date.Weekday() != time.Monday { // iterate back to Monday
		date = date.AddDate(0, 0, -1)
		isoYear, isoWeek = date.ISOWeek()
	}
	for isoYear < year { // iterate forward to the first day of the first week
		date = date.AddDate(0, 0, 1)
		isoYear, isoWeek = date.ISOWeek()
	}
	for isoWeek < week { // iterate forward to the first day of the given week
		date = date.AddDate(0, 0, 1)
		isoYear, isoWeek = date.ISOWeek()
	}
	return date
}
