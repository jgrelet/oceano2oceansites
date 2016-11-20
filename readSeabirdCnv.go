// readSeabirdCnv.go
package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// define regexp
var regIsHeader = regexp.MustCompile(`^[*#]`)
var regEndOfHeader = regexp.MustCompile(`\*END\*`)
var regCruise = regexp.MustCompile(`Cruise\s*:\s*(.*)`)
var regShip = regexp.MustCompile(`Ship\s*:\s*(.*)`)
var regStation = regexp.MustCompile(`Station\s*:\s*\D*(\d*)`)
var regType = regexp.MustCompile(`Type\s*:\s*(.*)`)
var regOperator = regexp.MustCompile(`Operator\s*:\s*(.*)`)
var regBottomDepth = regexp.MustCompile(`Bottom Depth\s*:\s*(\d*\.?\d+?)\s*\S*`)
var regDummyBottomDepth = regexp.MustCompile(`Bottom Depth\s*:\s*$`)
var regDate = regexp.MustCompile(`Date\s*:\s*(\d+)/(\d+)/(\d+)`)
var regHour = regexp.MustCompile(`[Heure|Hour]\s*:\s*(\d+)[:hH](\d+):(\d+)`)
var regLatitude = regexp.MustCompile(`Latitude\s*:\s*(\d+)\s+(\d+.\d+)\s+(\w)`)
var regLongitude = regexp.MustCompile(`Longitude\s*:\s*(\d+)\s+(\d+.\d+)\s+(\w)`)
var regSystemTime = regexp.MustCompile(`System UpLoad Time =\s+(.*)`)
var regNmeaLatitude = regexp.MustCompile(`NMEA Latitude\s*=\s*(\d+\s+\d+.\d+\s+\w)`)
var regNmeaLongitude = regexp.MustCompile(`NMEA Longitude\s*=\s*(\d+\s+\d+.\d+\s+\w)`)

// DecodeHeader parse header line from .cnv and extract correct information
// use regular expression to parse time with non standard format, 
// see: http://golang.org/src/time/format.go
func (nc *Nc) DecodeHeader(str string, profile float64, nbProfile int) {

	// decode Systeme Upload Time
	match := regSystemTime.MatchString(str)
	if match {
		res := regSystemTime.FindStringSubmatch(str)
		value := res[1]
		fmt.Fprintf(debug, "%s -> ", value)
		// create new Time object, see tools.go
		var t = NewTimeFromString("Jan 02 2006 15:04:05", value)
		v := t.Time2JulianDec()
		nc.Variables.set("TIME", v, nbProfile)
	}
	match = regNmeaLatitude.MatchString(str)
	if match {
		if v, err := Position2Decimal(str); err == nil {
			//nc.Variables_1D["LATITUDE"].([]float64)[nbProfile] = v
			nc.Variables.set("LATITUDE", v, nbProfile)
		}
	}
	match = regNmeaLongitude.MatchString(str)
	if match {
		if v, err := Position2Decimal(str); err == nil {
			//nc.Variables_1D["LONGITUDE"].([]float64)[nbProfile] = v
			nc.Variables.set("LONGITUDE", v, nbProfile)
		}
	}
	match = regCruise.MatchString(str)
	if match {
		res := regCruise.FindStringSubmatch(str)
		value := res[1]
		fmt.Fprintln(debug, value)
		nc.Attributes["cycle_mesure"] = value
	}
	match = regShip.MatchString(str)
	if match {
		res := regShip.FindStringSubmatch(str)
		value := res[1]
		fmt.Fprintln(debug, value)

		nc.Attributes["plateforme"] = value
	}
	match = regStation.MatchString(str)
	if match {
		res := regStation.FindStringSubmatch(str)
		value := res[1]
		if v, err := strconv.ParseFloat(value, 64); err == nil {
			fmt.Fprintln(debug, v)
			// ch
			//			format := "%0" + cfg.Ctd.StationPrefixLength + ".0f"
			//			p := fmt.Sprintf(format, profile)
			//			//s := fmt.Sprintf(format, v)
			//			fmt.Println(p, v)
			//			if p != v {
			//				fmt.Printf("Warning: profile for header differ from file name: %s <=> %s\n", p, v)
			//			}
			//nc.Variables_1D["PROFILE"].([]float64)[nbProfile] = profile
			nc.Variables.set("PROFILE", profile, nbProfile)

		}
	}
	match = regBottomDepth.MatchString(str)
	if match {
		res := regBottomDepth.FindStringSubmatch(str)
		value := res[1]
		if v, err := strconv.ParseFloat(value, 64); err == nil {
			fmt.Fprintf(debug, "Bath: %f\n", v)
			// nc.Variables_1D["BATH"].([]float64)[nbProfile] = v
			nc.Variables.set("BATH", v, nbProfile)
		}
	}
	match = regOperator.MatchString(str)
	if match {
		res := regOperator.FindStringSubmatch(str)
		value := res[1]
		if *optDebug {
			fmt.Println(value)
		}
	}
	// TODOS: uncomment, add optionnal value from seabird header
	match = regType.MatchString(str)
	if match {
		res := regType.FindStringSubmatch(str)
		value := strings.ToUpper(res[1]) // convert to upper case
		var v float64
		switch value {
		case "PHY":
			v = float64(PHY)
		case "GEO":
			v = float64(GEO)
		case "BIO":
			v = float64(BIO)
		default:
			v = float64(UNKNOW)
		}
		//f("Type: %f\n", v)
		//nc.Variables_1D["TYPECAST"] = append(nc.Variables_1D["TYPECAST"].([]float64), v)
		nc.Variables.set("TYPECAST", v, nbProfile)

		if *optDebug {
			fmt.Println(value)
		}
		nc.ExtraString[fmt.Sprintf("TYPE:%d", int(profile))] = value
	}
}

// GetProfileNumber return the profile number from filename. 
// Use CruisePrefix and StationPrefixLength defined in configuration file
// TODOS:  the prefix could be extract from filename
func (nc *Nc) GetProfileNumber(str string) float64 {
	var value float64
	var err error

	fmt.Fprintln(debug, "GetProfileNumber():\n--------------")
	reg := fmt.Sprintf("%s(\\d{%d})", cfg.Ctd.CruisePrefix, cfg.Ctd.StationPrefixLength)
	res := regexp.MustCompile(reg)
	match := res.MatchString(str)
	if match {
		t := res.FindStringSubmatch(strings.ToLower(str))
		fmt.Fprintf(debug, "Get profile number: %s -> %s\n", str, t[1])
		if value, err = strconv.ParseFloat(t[1], 64); err == nil {
			// get profile name, eg: csp00101
			nc.ExtraString[fmt.Sprintf("PRFL_NAME:%d", int(value))] = t[1]
		} else {
			log.Fatal(fmt.Sprintf("Error func GetProfileNumber: %s, var: %v", err, t[1]))
		}

	} else {
		log.Fatal(fmt.Sprintf("Error func GetProfileNumber:  %s -> %s\n", reg, str))
	}
	return value

}

// DecodeData extract data from the line read in str with order gave by hash mapVar
// values:  1318 81.583900 3.000 2.983 29.5431 29.5464 5 ...
// mapVar: PRES:2 DEPTH:3 PSAL:21 DOX2:18 ...
func (nc *Ctd) DecodeData(str string, profile float64, file string, line int) {

	// split the string str using whitespace characters
	values := strings.Fields(str)
	nbValue := len(values)

	// for each physical parameter, extract its data from the rigth column
	// and save it in map data
	for key, column := range mapVar {
		if column > nbValue {
			log.Fatal(fmt.Sprintf("Error in func DecodeData() "+
				"configuration mismatch\nFound %d values, and we try to use column %d",
				nbValue, column))
		}
		if v, err := strconv.ParseFloat(values[column], 64); err == nil {
			data[key] = v
		} else {
			log.Printf("Can't parse line: %d in file: %s\n", line, file)
			log.Fatal(err)
		}
	}
	data["PRFL"] = profile
}

// read .cnv files and return dimensions
func (nc *Ctd) firstPass(files []string) (int, int) {

	var line int
	var maxLine int 
	var pres float64 
	var depth float64 
	var maxDepth float64 
	var maxPres float64 
	var maxPresAll float64 

	fmt.Fprintf(echo, "First pass: ")
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
			str := scanner.Text()
			match := regIsHeader.MatchString(str)
			if !match {
				values := strings.Fields(str)
				// read the pressure
				if pres, err = strconv.ParseFloat(values[mapVar["PRES"]], 64); err != nil {
					log.Fatal(err)
				}
				// read the depth
				if depth, err = strconv.ParseFloat(values[mapVar["DEPTH"]], 64); err != nil {
					log.Fatal(err)
				} else {
					//p(math.Floor(depth))
				}
			}
			if pres > maxPres {
				maxPres = pres
				maxDepth = depth
				line = line + 1
			}
			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
		}
		fmt.Fprintf(debug, "Read %s size: %d max pres: %4.f\n", file, line, maxPres)

		if line > maxLine {
			maxLine = line
		}
		// store the maximum pressure and maximum depth value per cast
		nc.ExtraFloat[fmt.Sprintf("PRES:%d", int(profile))] = maxPres
		nc.ExtraFloat[fmt.Sprintf("DEPTH:%d", int(profile))] = math.Floor(maxDepth)
		if maxPres > maxPresAll {
			maxPresAll = maxPres
		}
		// reset value for next loop
		maxPres = 0
		maxDepth = 0
		pres = 0
		line = 0
	}

	fmt.Fprintf(echo, "%d files read, maximum pressure found: %4.0f db\n", len(files), maxPresAll)
	fmt.Fprintf(debug, "%d files read, maximum pressure found: %4.0f db\n", len(files), maxPresAll)
	fmt.Fprintf(debug, "size %d x %d\n", len(files), maxLine)
	return len(files), maxLine
}

// read .cnv files and extract data
func (nc *Ctd) secondPass(files []string) {

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
		downcast := true
		for scanner.Scan() {
			str := scanner.Text()
			match := regIsHeader.MatchString(str)
			if match {
				nc.DecodeHeader(str, profile, nbProfile)
			} else {
				// fill map data with information contain in read line str
				nc.DecodeData(str, profile, file, line)

				if downcast {
					// fill 2D slice
					for _, key := range hdr {
						if key != "PRFL" {
							//fmt.Println("Line: ", line, "key: ", key, " data: ", data[key])
							//nc.Variables_2D[key].data[nbProfile][line] = data[key].(float64)
							nc.Variables.set(key, data[key].(float64), nbProfile, line)
						}
					}
					// exit loop if reach maximum pressure for the profile
					if data["PRES"] == nc.ExtraFloat[fmt.Sprintf("PRES:%d", int(profile))] {
						downcast = false
					}
				} else {
					// store last julian day for end profile
					nc.ExtraFloat[fmt.Sprintf("ETDD:%d", int(profile))] = data["ETDD"].(float64)
					//fmt.Println(presMax)
				}
				line++
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		// increment sclice index
		nbProfile ++

		// store last julian day for end profile
		nc.ExtraFloat[fmt.Sprintf("ETDD:%d", int(profile))] = data["ETDD"].(float64)
		//fmt.Println(presMax)
	}
	fmt.Fprintln(debug, nc.Variables.get("PROFILE"))
}

// read cnv files in two pass, the first pass to get dimensions
// the second pass to get data and fill slice
func (nc *Ctd) Read(files []string) {

	// first pass, return dimensions fron cnv files
	dimx, dimy := nc.firstPass(files)
	nc.Dimensions["TIME"] = dimx
	nc.Dimensions["DEPTH"] = dimy

	//
	nc.InitVariables(dimx, dimy)

	// second pass, read files again, extract data and fill slices
	nc.secondPass(files)
}
