package main

import (
	"code.google.com/p/getopt"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const PROGNAME string = "oceano2oceansites"
const PROGVERSION string = "0.2.0"

type Data_2D struct {
	data [][]float64
}

type AllData_2D map[string]Data_2D

type Nc struct {
	Dimensions   map[string]int
	Variables_1D map[string][]float64
	Variables_2D AllData_2D
	Attributes   map[string]string
	Extras_f     map[string]float64 // used to store max of profiles value
	Extras_s     map[string]string  // used to store max of profiles type
	Roscop       map[string]RoscopAttribute
}

// configuration file
var cfgname string = "oceano2oceansites.ini"
var code_roscop string = "code_roscop.csv"

// file prefix for --all option: "-all" for all parameters, "" empty by default
var prefixAll = ""

// Create an empty map.
var map_var = map[string]int{}
var map_format = map[string]string{}
var data = map[string]float64{}
var hdr []string
var cfg Config

// global arg list options
var optDebug *bool
var optEcho *bool
var optAll *bool

// use for debug mode
var debug io.Writer = ioutil.Discard

// use for echo mode
var echo io.Writer = ioutil.Discard

var nc Nc

func main() {

	var files []string
	// to change the flags on the default logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if os.Getenv("OCEANO2OCEANSITES") != "" {
		cfgname = os.Getenv("OCEANO2OCEANSITES")
	}
	if os.Getenv("ROSCOP") != "" {
		code_roscop = os.Getenv("ROSCOP")
	}
	fmt.Fprintf(debug, "Configuration file:", os.Getenv("OCEANO2OCEANSITES"))
	fmt.Fprintf(debug, "Code ROSCOP file:", os.Getenv("ROSCOP"))
	fmt.Fprintf(debug, "GOPATH:", os.Getenv("GOPATH"))
	fmt.Fprintf(debug, "GOBIN:", os.Getenv("GOBIN"))

	// parse option, move outside main
	optDebug = getopt.Bool('d', "debug", "Display debug info")
	optEcho = getopt.Bool('e', "echo", "Display processing in stdout")
	optHelp := getopt.Bool('h', "help", "Help")
	optAll = getopt.Bool('a', "all", "Process all parameters")
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
	if *optDebug {
		debug = os.Stdout
	}
	if *optEcho {
		echo = os.Stdout
	}
	if *optCycleMesure != "" {
		fmt.Println(*optCycleMesure)
		nc.Attributes["cycle_mesure"] = *optCycleMesure
	}
	if *optAll {
		prefixAll = "-all"
	}

	// initialize map from netcdf structure
	nc.Dimensions = make(map[string]int)
	nc.Attributes = make(map[string]string)
	nc.Extras_f = make(map[string]float64)
	nc.Extras_s = make(map[string]string)
	nc.Variables_1D = make(map[string][]float64)
	nc.Variables_1D["PROFILE"] = []float64{}
	nc.Variables_1D["TIME"] = []float64{}
	nc.Variables_1D["LATITUDE"] = []float64{}
	nc.Variables_1D["LONGITUDE"] = []float64{}
	nc.Roscop = codeRoscopFromCsv(code_roscop)

	// read configuration file, by default, optCfgfile = cfgname
	GetConfig(*optCfgfile)
	// debug
	fmt.Fprintln(debug, map_format)
	// get files list from argument line
	// Args returns the non-option arguments.
	// see https://code.google.com/p/getopt/source/browse/set.go#27
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
	fmt.Fprintln(debug, files)

	readSeabirdCnv(files)

	// add some global attributes for profile, change in future
	nc.Attributes["data_type"] = "OceanSITES profile data"

	// write ASCII file
	WriteAsciiFiles(nc, map_format, hdr)

	// write netcdf file
	if err := WriteNetcdfFile(nc); err != nil {
		log.Fatal(err)
	}
}
