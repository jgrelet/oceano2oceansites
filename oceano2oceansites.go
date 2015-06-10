package main

import (
	"bufio"
	"code.google.com/p/getopt"
	"fmt"
	"log"
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
type Nc struct {
	Dimensions   map[string]int
	Variables_1D map[string][]float64
	Variables_2D map[string][]float64
	Attributes   map[string]string
}

type MyType struct {
	someslice [][]int
}

type MyData map[string]MyType

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
	day, _ := time.ParseInLocation("20060102150405", date, time.UTC)
	const julian = 2453738.4195
	// Easiest way to get the time.Time of the Unix time.
	// (See comments for the UnixDate in package Time.)
	unix := time.Unix(1136239445, 0)
	const oneDay = float64(86400. * time.Second)
	return julian + float64(day.Sub(unix))/oneDay
}

/*
sub positionDeci (deg, min, sign string) float64 {

  #controle du format des parametres
  $hemi = uc $hemi;
  if ( $hemi !~ /[NSEW]/ ) {
    die "Mauvais parametre '$hemi' : format hemi [NSEW]\n";
  }

  my $sign = 1;
  if ( $hemi eq "S" || $hemi eq "W" ) {
    $sign = -1;
  }
  my $tmp = $min;
  $min = abs $tmp;
  my $sec = ( $tmp - $min ) * 100;
  return ( ( $deg + ( $min + $sec / 100 ) / 60 ) * $sign );
}
*/

func DecodeHeader(str string) {
	// decode Systeme Upload Time
	match := regSystemTime.MatchString(str)
	if match {
		res := regSystemTime.FindStringSubmatch(str)
		value := res[1]
		if *optDebug {
			fmt.Println(value)
		}
		// parse time with non standard format, see:
		// http://golang.org/src/time/format.go
		if t, err := time.Parse("Jan 02 2006 15:04:05", value); err == nil {
			if *optDebug {
				fmt.Println(t.Format("20060102150405"))
				fmt.Printf("Julian: %f\n", Julian(t.Format("20060102150405")))
			}
		} else {
			fmt.Println("Failed to decode time:", err)
		}
	}
	match = regNmeaLatitude.MatchString(str)
	if match {
		//		res := regNmeaLatitude.FindStringSubmatch(str)
		//		deg := res[1]; min:=res[2]; sign:=res[3]
		//		fmt.Printf("%s %s %s", deg, min, sign)
	}
	match = regNmeaLongitude.MatchString(str)
	if match {
		res := regNmeaLongitude.FindStringSubmatch(str)
		value := res[1]
		if *optDebug {
			fmt.Println(value)
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
	//	fmt.Fprintf(w, "%s\n", header.Time)
}

func GetProfileNumber(path string) float64 {
	reg := fmt.Sprintf("%s(\\d{%s})", cfg.Ctd.CruisePrefix, cfg.Ctd.StationPrefixLength)
	r := regexp.MustCompile(reg)
	match := r.FindStringSubmatch(strings.ToLower(path))
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

func (mp MyData) NewMyType(name string, width, height int) *MyData {
	mt := new(MyType)
	mt.someslice = make([][]int, width)
	for i := range mt.someslice {
		mt.someslice[i] = make([]int, height)
		//		for j := range mt.someslice[i] {
		//			mt.someslice[i][j] = j
		//		}
	}
	mp[name] = *mt
	return &mp
}

func firstPass(files []string) {
	// loop over each files passed throw command line
	for _, file := range files {
		fid, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		defer fid.Close()

		fmt.Printf("Read %s\n", file)
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
					fmt.Fprintf(w, map_format[hdr[i]]+" ", data[hdr[i]])
				}
				// add new line
				if *optEcho {
					fmt.Printf("\n")
				}
				fmt.Fprintf(w, "\n")
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
	}
	if *optDebug {
		fmt.Println(nc.Variables_1D["PROFILE"])
	}
}

func secondPass() {
	// loop over each files passed throw command line
	for _, file := range files {
		fid, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		defer fid.Close()

		fmt.Printf("Read %s\n", file)
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
					fmt.Fprintf(w, map_format[hdr[i]]+" ", data[hdr[i]])
				}
				// add new line
				if *optEcho {
					fmt.Printf("\n")
				}
				fmt.Fprintf(w, "\n")
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

	// init the data 2D for test
	data := make(MyData)

	// initialize map from netcdf structure
	nc.Dimensions = make(map[string]int)
	nc.Variables_1D = make(map[string][]float64)
	nc.Variables_2D = make(map[string][]float64)
	nc.Attributes = make(map[string]string)
	nc.Variables_1D["PROFILE"] = []float64{}

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

	// open output file for writing result
	fout, err := os.Create(outputFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer fout.Close()

	// use buffered mode for writing
	w := bufio.NewWriter(fout)

	firstPass()

//		CreateNetcdfFile(outputFilename+".nc", nc)
}
