/*
Copyright © 2025 ganlinden@gmail.com
*/
package cmd

import (
	"encoding/json"
	"path/filepath"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/mattn/go-runewidth"
)

var PERIODS_PATH = filepath.Join(os.Getenv("HOME"), "go/bin/data/periods.json")
const MAX_DAYS = 100
const MAX_FUTURE_DAYS = 10

type JsonStruct struct {
	Item string `json:"item,omitempty` // unique if not empty
	ItemAlias string `json:"itemAlias,omitempty` // unique if not empty
	Period []string `json:"period,omitempty` // [12/25, 12/30, 1/2, 1/6], null if branch node
	Members []*JsonStruct `json:"members,omitempty` // members of this category, null if leaf node
}

type CliStruct struct {
	Item string
	ItemAlias string
	Period string // ...#....#..#...#.
}

/*
member/memberAlias -> period
tree of jsonstructconst
a list of names member/memberAlias/group/groupAlias
*/
func Unmarshal() (*JsonStruct, []string ) {
	var jsonTree JsonStruct
	var allString []string

	data, err := os.ReadFile(PERIODS_PATH)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(data, &jsonTree); err != nil {
		panic(err)
	}

	var helper func(*JsonStruct)
	helper = func(s *JsonStruct) {
		allString = append(allString, s.Item)
		allString = append(allString, s.ItemAlias)
		if s.Members == nil {
			return
		}
		for _, m := range s.Members {
			helper(m)
		}
	}
	helper(&jsonTree)
	return &jsonTree, allString
}

func Marshal(input *JsonStruct) {
	data, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(PERIODS_PATH, data, 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println("Water success!")
}

/*
Take in a set of strings as query
Return a list of pointers of JsonStruct as items
If the query is empty, return empty slice
*/
func SelectItems(query map[string]struct{}, jsonTree *JsonStruct) []*JsonStruct {
	var result []*JsonStruct

	var helper func(*JsonStruct, bool)
	helper = func(s *JsonStruct, choose bool) {
		_, ok1 := query[s.Item]
		_, ok2 := query[s.ItemAlias]
		choose = ok1 || ok2 || choose
		if s.Members == nil && choose {
			result = append(result, s)
			return
		}
		for _, m := range s.Members {
			helper(m, choose)
		}
	}
	helper(jsonTree, false)
	
	return result
}

func date2String(date time.Time) string {
	return date.Format("1/2")
}

/*
Possible input str: 8/3, 1/17, 12/9, 04/24, 12/31
*/
func string2Date(str string, year int) time.Time {
	res, err := time.Parse("1/2", str)
	if err != nil {
		panic(err)
	}
	res = res.AddDate(year, 0, 0)
	return res
}

func WaterItems(items []*JsonStruct) {
	var today string = date2String(time.Now())
	var year int = time.Now().Year()
	for _, i := range items {
		// Check if we already watered today
		lastTime := ""
		if len(i.Period) > 0 {
			lastTime = date2String(string2Date(i.Period[len(i.Period) - 1], year))
		}
		if lastTime == today {
			fmt.Println("Already watered", i.Item, "today.")
		} else {
			i.Period = append(i.Period, today)
			fmt.Println("Watering", i.Item)
		}
	}
}

/*
Map a slice of mmdd strings to a slice of time.Time.
Example input: [12/25, 12/30, 1/2, 1/6] (say today is 12/31)

Should be able to deduce year by itself based on the assumption that:
The last date in the input slice is in the range of
today plus minus half a year.
Example 1: if today is 2025/8/1, then we assume 1/1 as 2026/1/1,
and assume 3/1 as 2025/3/1.
Example 2: if today is 2025/2/1, then we assume 11/1 as 2024/11/1,
and assume 7/1 as 2025/7/1.
*/
const halfYear time.Duration = 182 * 24 * time.Hour
func strings2Dates(input []string) []time.Time {
	var res = make([]time.Time, len(input))
	if len(input) == 0 {
		return res
	}
	// Deduce year
	var year int = time.Now().Year()
	tentative := string2Date(input[len(input) - 1], year)
    if time.Since(tentative) > halfYear {
		year++
	} else if time.Until(tentative) > halfYear {
		year--
	}
    // Reversely traverse input and update year if needed
	var i = len(input) - 1
	res[i] = string2Date(input[i], year)
	for i = len(input) - 2; i >= 0; i-- {
		currDate := string2Date(input[i], year)
		laterDate := res[i + 1]
		// If currDate is 12/31, laterDate is 1/1, then we know the year
		// becomes the previous year.
		if currDate.After(laterDate) {
			year--
			currDate = string2Date(input[i], year)
		}
		res[i] = currDate
	}
	return res
}

func Json2CliStructs(input []*JsonStruct ) ([]CliStruct, time.Time) {
	// Convert Period of each JsonStruct to a slice of dates as Time.time,
	// and keep track of the latest date among all dates.
	var today = string2Date(date2String(time.Now()), time.Now().Year())
	var latestDate = today
	var dates2D [][]time.Time
	for _, js := range input {
		var dates = strings2Dates(js.Period)
		dates2D = append(dates2D, dates)
		if len(dates) > 0 && dates[len(dates) - 1].After(latestDate) {
			latestDate = dates[len(dates) - 1]
		}
	}
	// Compute latest and earliest dates to determine the range of days.
	// Both latest and earliest dates are inclusive.
	if latestDate.After(today.AddDate(0, 0, MAX_FUTURE_DAYS)) {
		latestDate = today.AddDate(0, 0, MAX_FUTURE_DAYS)
	}
	var earliestDate = latestDate.AddDate(0, 0, -1 * (MAX_DAYS - 1))
	// Map JsonStruct to CliStruct, and convert slices of dates to strings
	var res []CliStruct
	for i, js := range input {
		// Convert to set
		dates := make(map[time.Time]struct{})
		for _, date := range dates2D[i] {
			dates[date] = struct{}{}
		}
		var sb strings.Builder
		for date := earliestDate; !date.After(latestDate); date = date.AddDate(0, 0, 1) {
			if _, ok := dates[date]; ok {
				sb.WriteByte('#')
			} else {
				sb.WriteByte('.')
			}
		}
		var cs = CliStruct{Item: js.Item, ItemAlias: js.ItemAlias, Period: sb.String()}
		res = append(res, cs)
	}
	return res, earliestDate
}

func countDays(a, b time.Time) int {
	a = string2Date(date2String(a), a.Year())
	b = string2Date(date2String(b), b.Year())
	if a.After(b) {
		a, b = b, a
	}
	i := 0
	for ; b.After(a); a = a.AddDate(0, 0, 1) {
		i++
	}
	return i
}

func Print(input []CliStruct, earliestDate time.Time) {
	var tw = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	// Print the header like --*------*------|------*--
	var sb strings.Builder
	days := countDays(earliestDate, time.Now())
	sb.WriteString(strings.Repeat("-", days % 7))
	sb.WriteString(strings.Repeat("*------", days / 7))
	sb.WriteByte('|')
	days = countDays(time.Now(), earliestDate.AddDate(0, 0, MAX_DAYS - 1))
	sb.WriteString(strings.Repeat("------*", days / 7))
	sb.WriteString(strings.Repeat("-", days % 7))
	fmt.Fprintln(tw, sb.String() + "\tItem(Alias)")
	for _, cs := range input {
		name := cs.Item
		if cs.ItemAlias != "" {
			name = fmt.Sprintf("%s(%s)", cs.Item, cs.ItemAlias)
		}
		fmt.Fprintf(tw, "%s\t%s\n", cs.Period, runewidth.FillRight(name, 20))
	}
	tw.Flush()
}