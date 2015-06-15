package main

import (
	"bufio"
	"code.google.com/p/getopt"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	//	"runtime"  cf runtime.GOSS get osname
)

const PROGNAME string = "oceano2oceansites"
const PROGVERSION string = "0.1.0"

type Header struct {
	Time      string
	Julian    float64
	Latitude  float64
	Longitude float64
}

type Config struct {
	Global struct {
		Author string
		Debug  bool
		Echo   bool
	}
	Cruise struct {
		CycleMesure string
		Plateforme  string
		Callsign    string
	}
	Ctd struct {
		CruisePrefix        string
		StationPrefixLength string
		Header              string
		Split               string
		Format              string
	}
	Ctdall struct {
		Header string
		Split  string
		Format string
	}
}

type Data_2D struct {
	data [][]float64
}

type AllData_2D map[string]Data_2D

type Nc struct {
	Dimensions   map[string]int
	Variables_1D map[string][]float64
	Variables_2D AllData_2D
	Attributes   map[string]string
}

// configuration file
const cfgname string = "test.ini"

// Create an empty map.
var map_var = map[string]int{}
var map_format = map[string]string{}
var data = map[string]float64{}
var hdr []string
var cfg Config
var header Header
var optDebug *bool
var optEcho *bool

// define regexp
var regIsHeader = regexp.MustCompile(`^[*#]`)
var regEndOfHeader = regexp.MustCompile(`\*END\*`)
var regCruise = regexp.MustCompile(`Cruise\s*:\s*(.*)`)
var regShip = regexp.MustCompile(`Ship\s*:\s*(.*)`)
var regStation = regexp.MustCompile(`Station\s*:\s*(.*)`)
var regType = regexp.MustCompile(`Type\s*:\s*(.*)`)
var regOperator = regexp.MustCompile(`Operator\s*:\s*(.*)`)
var regBottomDepth = regexp.MustCompile(`Bottom Depth\s*:\s*(\d*\.?\d+?)\s*\S*`)
var regDate = regexp.MustCompile(`Date\s*:\s*(\d+)/(\d+)/(\d+)`)
var regHour = regexp.MustCompile(`[Heure|Hour]\s*:\s*(\d+)[:hH](\d+):(\d+)`)
var regLatitude = regexp.MustCompile(`Latitude\s*:\s*(\d+)\s+(\d+.\d+)\s+(\w)`)
var regLongitude = regexp.MustCompile(`Longitude\s*:\s*(\d+)\s+(\d+.\d+)\s+(\w)`)
var regSystemTime = regexp.MustCompile(`System UpLoad Time =\s+(.*)`)
var regNmeaLatitude = regexp.MustCompile(`NMEA Latitude\s*=\s*(\d+\s+\d+.\d+\s+\w)`)
var regNmeaLongitude = regexp.MustCompile(`NMEA Longitude\s*=\s*(\d+\s+\d+.\d+\s+\w)`)

var outputFilename string = "test_ctd"

//var nc map[string]map[string]string
var nc Nc

func Julian(date string) float64 {
	const DIFF_ORIGIN = 2433283.0 // diff between UNIX DATE and 1950/1/1 00:00:00

	t, _ := time.Parse("20060102150405", date)
	a := int(14-t.Month()) / 12
	y := t.Year() + 4800 - a
	m := int(t.Month()) + 12*a - 3
	julianDay := int(t.Day()) + (153*m+2)/5 + 365*y + y/4
	julianDay = julianDay - y/100 + y/400 - 32045.0 - DIFF_ORIGIN //+
	//float64(t.Hour())/24 + float64(t.Minute())/1440 + float64(t.Second())/86400
	if *optDebug {
		fmt.Println("Julian day:", date, " -> ", julianDay)
	}
	return float64(julianDay) + float64(t.Hour())/24 + float64(t.Minute())/1440 + float64(t.Second())/86400
}

func positionDeci(pos string) (float64, error) {

	var multiplier float64 = 1
	var value float64

	var regNmeaPos = regexp.MustCompile(`(\d+)\s+(\d+.\d+)\s+(\w)`)

	if strings.Contains(pos, "S") || strings.Contains(pos, "W") {
		multiplier = -1.0
	}
	match := regNmeaPos.MatchString(pos)
	if match {
		res := regNmeaPos.FindStringSubmatch(pos)
		deg, _ := strconv.ParseFloat(res[1], 64)
		min, _ := strconv.ParseFloat(res[2], 64)
		tmp := math.Abs(min)
		sec := (tmp - min) * 100.0
		value = (deg + (min+sec/100.0)/60.0) * multiplier
		if *optDebug {
			fmt.Println("positionDeci:", pos, " -> ", value)
		}
	} else {
		return 1e36, errors.New("positionDeci: failed to decode position")
	}
	return value, nil
}

func DecodeHeader(str string) {
	// decode Systeme Upload Time
	var v float64 = 1e36

	match := regSystemTime.MatchString(str)
	if match {
		res := regSystemTime.FindStringSubmatch(str)
		value := res[1]
		if *optDebug {
			fmt.Printf("%s -> ", value)
		}
		// parse time with non standard format, see:
		// http://golang.org/src/time/format.go
		if t, err := time.Parse("Jan 02 2006 15:04:05", value); err == nil {
			if *optDebug {
				v = Julian(t.Format("20060102150405"))
				nc.Variables_1D["TIME"] = append(nc.Variables_1D["TIME"], v)
			}
		} else {
			fmt.Println("Failed to decode time:", err)
			nc.Variables_1D["TIME"] = append(nc.Variables_1D["TIME"], 1e36)
		}
	}
	match = regNmeaLatitude.MatchString(str)
	if match {
		if v, err := positionDeci(str); err == nil {
			nc.Variables_1D["LATITUDE"] = append(nc.Variables_1D["LATITUDE"], v)
		} else {
			nc.Variables_1D["LATITUDE"] = append(nc.Variables_1D["LATITUDE"], 1e36)
		}
	}
	match = regNmeaLongitude.MatchString(str)
	if match {
		if v, err := positionDeci(str); err == nil {
			nc.Variables_1D["LONGITUDE"] = append(nc.Variables_1D["LONGITUDE"], v)
		} else {
			nc.Variables_1D["LATITUDE"] = append(nc.Variables_1D["LATITUDE"], 1e36)
		}
	}
	match = regCruise.MatchString(str)
	if match {
		res := regCruise.FindStringSubmatch(str)
		value := res[1]
		if *optDebug {
			fmt.Println(value)
		}
		nc.Attributes["cycle_mesure"] = value
	}
	match = regShip.MatchString(str)
	if match {
		res := regShip.FindStringSubmatch(str)
		value := res[1]
		if *optDebug {
			fmt.Println(value)
		}
		nc.Attributes["plateforme"] = value
	}
	match = regStation.MatchString(str)
	if match {
		res := regStation.FindStringSubmatch(str)
		value := res[1]
		if v, err := strconv.ParseFloat(value, 64); err == nil {
			if *optDebug {
				fmt.Println(v)
			}
			nc.Variables_1D["PROFILE"] = append(nc.Variables_1D["PROFILE"], v)
		} else {
			nc.Variables_1D["PROFILE"] = append(nc.Variables_1D["PROFILE"], 1e36)
		}
	}
	match = regBottomDepth.MatchString(str)
	if match {
		res := regBottomDepth.FindStringSubmatch(str)
		value := res[1]
		if v, err := strconv.ParseFloat(value, 64); err == nil {
			if *optDebug {
				fmt.Println(v)
			}
			nc.Variables_1D["BATH"] = append(nc.Variables_1D["BATH"], v)
		} else {
			nc.Variables_1D["BATH"] = append(nc.Variables_1D["BATH"], 1e36)
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
}

func WriteHeader() {
	if *optDebug {
		fmt.Printf("%s\n", header.Time)
	}
}

func GetProfileNumber(path string) float64 {
	reg := fmt.Sprintf("%s(\\d{%s})", cfg.Ctd.CruisePrefix, cfg.Ctd.StationPrefixLength)
	r := regexp.MustCompile(reg)
	match := r.FindStringSubmatch(strings.ToLower(path))
	if *optDebug {
		fmt.Println("Get Profile number: ", path, "-> ", match)
	}
	value, _ := strconv.ParseFloat(match[1], 64)
	return value
}

func DecodeData(str string, profile float64) {
	fields := strings.Fields(str)
	for key, value := range map_var {
		if v, err := strconv.ParseFloat(fields[value], 64); err == nil {
			data[key] = v
		}
	}
	data["PRFL"] = profile
}

func (mp AllData_2D) NewData_2D(name string, width, height int) *AllData_2D {
	mt := new(Data_2D)
	mt.data = make([][]float64, width)
	for i := range mt.data {
		mt.data[i] = make([]float64, height)
		for j := range mt.data[i] {
			mt.data[i][j] = 1e36
		}
	}
	mp[name] = *mt
	return &mp
}

// read all cnv files and return dimensions
func firstPass(files []string) (int, int) {

	var line int = 0
	var depth int = 0

	fmt.Printf("First pass:\n")

	// loop over each files passed throw command line
	for _, file := range files {
		fid, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		defer fid.Close()

		//	profile := GetProfileNumber(file)  // ?
		scanner := bufio.NewScanner(fid)
		for scanner.Scan() {
			str := scanner.Text()
			match := regIsHeader.MatchString(str)
			if !match {
				//	DecodeData(str, profile) // ?
				line += 1
			}
		}
		if *optDebug {
			fmt.Printf("Read %s size: %d\n", file, line)
		}
		if line > depth {
			depth = line
		}
		line = 0

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
	return len(files), depth
}

func secondPass(files []string) {

	var nbProfile int = 0

	fmt.Printf("Second pass:\n")

	// open output file for writing result
	fout, err := os.Create(outputFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer fout.Close()

	// use buffered mode for writing
	w := bufio.NewWriter(fout)

	// loop over each files passed throw command line
	for _, file := range files {
		var line int = 0
		fid, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		defer fid.Close()

		// fmt.Printf("Read %s\n", file)

		profile := GetProfileNumber(file)
		scanner := bufio.NewScanner(fid)
		for scanner.Scan() {
			str := scanner.Text()
			match := regIsHeader.MatchString(str)
			if match {
				DecodeHeader(str)
			} else {
				DecodeData(str, profile)
				for i := range hdr {
					if *optEcho {
						fmt.Printf(map_format[hdr[i]]+" ", data[hdr[i]])
					}
					// write data to ascii file
					fmt.Fprintf(w, map_format[hdr[i]]+" ", data[hdr[i]])

					// fill 2D slice for netcdf
					if hdr[i] != "PRFL" {
						if *optEcho {
							//fmt.Printf("%1d %2d %s %6.3f\n", nbProfile, line, hdr[i], data[hdr[i]])
						}
						nc.Variables_2D[hdr[i]].data[nbProfile][line] = data[hdr[i]]
					}
				}
				// add new line
				if *optEcho {
					fmt.Printf("\n")
				}
				fmt.Fprintf(w, "\n")
				line++
			}
			match = regEndOfHeader.MatchString(str)
			if match {
				WriteHeader()
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		w.Flush()

		// increment index in sclice
		nbProfile += 1
	}
	if *optDebug {
		fmt.Println(nc.Variables_1D["PROFILE"])
	}
}

func main() {

	// to change the flags on the default logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// parse option, move outside main
	optDebug = getopt.Bool('d', "debug", "Display debug info")
	optEcho = getopt.Bool('e', "debug", "Display processing in stdout")
	optHelp := getopt.Bool('h', "help", "Help")
	optVersion := getopt.BoolLong("version", 'v', "Show version, then exit.")
	optCfgfile := getopt.StringLong("config", 'c', cfgname, "Name of the configuration file to use.")
	optCycleMesure := getopt.StringLong("cycle_mesure", 'm', "", "Name of cycle_mesure")
	optFiles := getopt.StringLong("files", 'f', "", "files to process ex: data/fr25*.cnv")

	// parse options line argument
	getopt.Parse()

	// process bloc when option is set
	if *optHelp {
		getopt.Usage()
		os.Exit(0)
	}
	if *optVersion {
		fmt.Println(PROGNAME + ": v" + PROGVERSION)
		os.Exit(0)
	}
	if *optCycleMesure != "" {
		fmt.Println(*optCycleMesure)
		nc.Attributes["cycle_mesure"] = *optCycleMesure
	}

	// read configuration file, by default, optCfgfile = cfgname
	GetConfig(*optCfgfile)

	if *optDebug {
		fmt.Println(map_format)
	}

	// initialize map from netcdf structure
	nc.Dimensions = make(map[string]int)
	nc.Attributes = make(map[string]string)
	nc.Variables_1D = make(map[string][]float64)
	nc.Variables_1D["PROFILE"] = []float64{}
	nc.Variables_1D["TIME"] = []float64{}
	nc.Variables_1D["LATITUDE"] = []float64{}
	nc.Variables_1D["LONGITUDE"] = []float64{}

	// get files list from argument line
	// Args returns the non-option arguments.
	// see https://code.google.com/p/getopt/source/browse/set.go#27
	var files []string
	if *optFiles == "" {
		files = getopt.Args()
	} else {
		files, _ = filepath.Glob(*optFiles)
	}
	// if no files supplied for arg list, test if files is empty
	if len(files) == 0 {
		getopt.Usage()
		os.Exit(0)
	}
	if *optDebug {
		fmt.Println(files)
	}

	// first pass, return dimensions fron cnv files
	nc.Dimensions["TIME"], nc.Dimensions["DEPTH"] = firstPass(files)

	// initialize 2D data
	nc.Variables_2D = make(AllData_2D)
	for i, _ := range map_var {
		nc.Variables_2D.NewData_2D(i, nc.Dimensions["TIME"], nc.Dimensions["DEPTH"])
	}

	// second pass, read files again, extract data and fill slices
	secondPass(files)

	// write netcdf file
	CreateNetcdfFile(outputFilename+".nc", nc)
}
