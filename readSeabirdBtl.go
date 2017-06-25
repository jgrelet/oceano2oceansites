// readSeabirdBtl
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var regIsHour = regexp.MustCompile(`^\s+(\d+:\d+:\d+)`)

//var regIsDate = regexp.MustCompile(`^\s+\d+\s+(\w{3}\s+\d{2}\s+\d{4})`)
var regIsMontDayYear = regexp.MustCompile(`^\s+\d+\s+(\w{3})\s+(\d{2})\s+(\d{4})`)

var regIsHeaderBtl = regexp.MustCompile(`^\s*Bottle|^\s+Position|^[*#]`)

// DecodeHeader parse header line from .btl and extract correct information
// use regular expression to parse time with non standard format,
// see: http://golang.org/src/time/format.go
func (nc *Btl) DecodeHeader(str string, profile float64, nbProfile int) {
	//fmt.Println("DecodeHeader for bottle not implemented!")
}

// read .btl files and return dimensions
func (nc *Btl) firstPass(files []string) (int, int) {

	// all locale variale are initialize to 0 by default
	var line int
	var maxLine int
	var bottle float64
	var maxBottle float64
	var maxBottleAll float64

	fmt.Fprintf(echo, "First pass: ")
	fmt.Fprintf(debug, "First pass:\n")
	// loop over each files passed throw command line
	for _, file := range files {
		fid, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		defer fid.Close()

		profile := nc.GetProfileNumber(file)
		scanner := bufio.NewScanner(fid)
		for scanner.Scan() {
			str := scanner.Text() // read the first line for each bottle sample
			match := regIsHeaderBtl.MatchString(str)
			if !match {
				fmt.Fprintf(debug, "First pass: match\n")
				p(str)
				values := strings.Fields(str)
				p("BOTL", mapVar["BOTL"])
				p(values[mapVar["BOTL"]])
				if bottle, err = strconv.ParseFloat(values[mapVar["BOTL"]], 64); err != nil {
					log.Fatal(err)
				}
				fmt.Fprintln(debug, values)
				// read the second line for each bottle sample
				scanner.Scan()
				str = scanner.Text()

			}
			if bottle > maxBottle {
				maxBottle = bottle
				line = line + 1
			}
			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
		}
		fmt.Fprintf(debug, "Read %s size: %d max pres: %4.f\n", file, line, maxBottle)

		if line > maxLine {
			maxLine = line
		}
		// store the maximum pressure value
		nc.ExtraFloat[fmt.Sprintf("BOTL:%d", int(profile))] = maxBottle
		if maxBottle > maxBottleAll {
			maxBottleAll = maxBottle
		}
		// reset value for next loop
		maxBottle = 0
		bottle = 0
		line = 0
	}

	fmt.Fprintf(echo, "%d files read, maximum bottle found: %4.0f db\n", len(files), maxBottle)
	fmt.Fprintf(debug, "%d files read, maximum pressure found: %4.0f db\n", len(files), maxBottle)
	fmt.Fprintf(debug, "size %d x %d\n", len(files), maxLine)
	return len(files), maxLine
}

// read .cnv files and extract data
func (nc *Btl) secondPass(files []string) {

	var month, day, year string

	fmt.Fprintf(echo, "Second pass ...\n")

	// initialize profile and pressure max
	var nbProfile int

	// loop over each files passed throw command line
	for _, file := range files {
		var line int

		fid, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		defer fid.Close()
		// fmt.Printf("Read %s\n", file)

		profile := nc.GetProfileNumber(file)
		scanner := bufio.NewScanner(fid)
		for scanner.Scan() {
			str := scanner.Text()
			match := regIsHeaderBtl.MatchString(str)
			if match {
				nc.DecodeHeader(str, profile, nbProfile)
			} else {
				match = regIsMontDayYear.MatchString(str)
				if match {
					res := regIsMontDayYear.FindStringSubmatch(str)
					month, day, year = res[1], res[2], res[3]
					//date = res[1]
					//f("Date -> %s/%s/%s\n", month, day, year)
				}
				match = regIsHour.MatchString(str)
				if match {
					res := regIsHour.FindStringSubmatch(str)
					time := res[1]
					//f("Time -> %s\n", time)
					// create new Time object, see tools.go
					var t = NewTimeFromString("Jan 02 2006 15:04:05 UTC",
						fmt.Sprintf("%s %s %s %s", month, day, year, time))
					//					v := t.Time2JulianDec()
					//					t1 := NewTimeFromJulian(v)
					y, _ := strconv.ParseFloat(year, 64)
					t2 := NewTimeFromJulianDay(y, t)
					nc.Variables["TIME"] = append(nc.Variables["TIME"].([]float64),
						t2.JulianDayOfYear())
					//p(t2.JulianDayOfYear())
				}
			}
			line++
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		// increment sclice index
		nbProfile++
	}
}

// read .btl files and extract data
func (nc *Btl) Read(files []string) {
	// first pass, return dimensions fron cnv files
	nc.Dimensions["TIME"], nc.Dimensions["DEPTH"] = nc.firstPass(files)

	//	// initialize 2D data
	//	nc.Variables_2D = make(AllData_2D)
	//	for i, _ := range map_var {
	//		nc.Variables_2D.NewData_2D(i, nc.Dimensions["TIME"], nc.Dimensions["DEPTH"])
	//	}

	// second pass, read files again, extract data and fill slices
	nc.secondPass(files)
}
