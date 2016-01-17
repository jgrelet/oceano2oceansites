package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

const PROGNAME string = "oceano2oceansites"
const PROGVERSION string = "0.2.0"

type Data_2D struct {
	data [][]float64
}

type AllData_2D map[string]Data_2D

type Nc struct {
	Dimensions   map[string]int
	Variables_1D map[string]interface{}
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
	WriteNetcdf(InstrumentType)
}

// default configuration file
var cfgname string = "oceano2oceansites.ini"

// default physical parameters file definition is embeded in code_roscop.go
var code_roscop string = ""

// file prefix for --all option: "-all" for all parameters, "" empty by default
var prefixAll = ""

// default output directory
var outputDir = "out"

// Create an empty map.
var map_var = map[string]int{}
var map_format = map[string]string{}
var data = make(map[string]interface{})
var hdr []string
var cfg Config

// use for debug mode
var debug io.Writer = ioutil.Discard

// use for echo mode
var echo io.Writer = ioutil.Discard

// usefull macro
var p = fmt.Println
var f = fmt.Printf

// nc implement interface Process
var nc Process

// define new receiver type based on netcdf equivalent structure
type Ctd struct{ Nc }
type Btl struct{ Nc }

// main body
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

	// get options on argument line
	files, optCfgfile := GetOptions()

	// read the first file and try to find the instrument type, return a bit mask
	typeInstrument = AnalyseFirstFile(files)

	// following the instrument type, allocate the rigth receiver based on
	// Process interface
	switch typeInstrument {
	case CTD:
		nc = &Ctd{}
	case BTL:
		nc = &Btl{}
	default:
		f("main: invalide option typeInstrument -> %d\n", typeInstrument)
		p("Exiting...")
		os.Exit(0)
	}

	// read configuration file, by default, optCfgfile = cfgname
	nc.GetConfig(optCfgfile)
	// debug
	fmt.Fprintln(debug, map_format)

	// read and process all data files
	nc.Read(files)

	// write ASCII file
	nc.WriteAscii(map_format, hdr)

	// write netcdf file
	nc.WriteNetcdf(typeInstrument)
}
