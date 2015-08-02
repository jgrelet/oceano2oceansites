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

type Process interface {
	Read([]string)
	GetConfig(string)
	//	WriteHeader(map[string]string, []string)
	WriteAscii(map[string]string, []string)
	WriteNetcdf() error
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

var nc Process

type Ctd struct{ Nc }
type Btl struct{ Nc }

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
	//	optCycleMesure := getopt.StringLong("cycle_mesure", 'm', "", "Name of cycle_mesure")
	optFiles := getopt.StringLong("files", 'f', "", "files to process ex: data/fr25*.cnv")

	// parse options line argument
	getopt.Parse()

	if *optFiles == "" {
		files = getopt.Args()
	} else {
		files, _ = filepath.Glob(*optFiles)
	}
	// if no files supplied for arg list, test if files is empty
	if len(files) == 0 {
		getopt.Usage()
		fmt.Println("\nPlease, specify files to process or define --files options")
		os.Exit(0)
	}
	fmt.Fprintln(debug, files)

	// read the first file and try to find the instrument type, return a bit mask
	typeInstrument := AnalyseFirstFile(files)

	// following the instrument type, allocate the rigth receiver based on
	// Process interface
	switch typeInstrument {
	case CTD:
		nc = &Ctd{}
	case BTL:
		nc = &Btl{}
	}

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
	//	if *optCycleMesure != "" {
	//		fmt.Println(*optCycleMesure)
	//		nc.Attributes["cycle_mesure"] = *optCycleMesure
	//	}
	if *optAll {
		prefixAll = "-all"
	}

	// read configuration file, by default, optCfgfile = cfgname
	nc.GetConfig(*optCfgfile)
	// debug
	fmt.Fprintln(debug, map_format)
	// get files list from argument line
	// Args returns the non-option arguments.
	// see https://code.google.com/p/getopt/source/browse/set.go#27

	// read and process all data files
	nc.Read(files)

	// write ASCII file
	nc.WriteAscii(map_format, hdr)

	// write netcdf file
	if err := nc.WriteNetcdf(); err != nil {
		log.Fatal(err)
	}
}
