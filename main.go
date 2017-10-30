package main

// build arg: -i -ldflags "-linkmode external -extldflags -Wl,-Bstatic"

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var (
	// Version set in Makefile
	// see: http://stackoverflow.com/questions/11354518/golang-application-auto-build-versioning
	Version = "undefined"
	// Binary set in Makefile
	Binary = "oceano2oceansites"
	// BuildTime set in Makefile
	BuildTime = "undefined"
)

// use for echo mode
// Discard is an io.Writer on which all Write calls succeed
// without doing anything. Discard = devNull(0) = int(0)
// when echo is define in program argument list, echo = os.Stdout
var echo = ioutil.Discard

// use for debug mode
var debug = ioutil.Discard

// usefull shortcut macros
var p = fmt.Println
var f = fmt.Printf

// default configuration file
var cfgname = "oceano2oceansites.toml"

// default physical parameters file definition is embeded in code_roscop.go
var codeRoscop = "roscop/code_roscop.csv"

// file prefix for --all option: "-all" for all parameters, "" empty by default
var prefixAll = ""

// default output directory (current)
var outputDir = "."

// Create an empty map.
var mapVar = map[string]int{}
var mapFormat = map[string]string{}
var data = make(map[string]interface{})

var hdr []string
var cfg Config

// Nc is the representation in memory of a data set is similar to
// that of a netcdf file
type Nc struct {
	// store dimensions
	Dimensions map[string]int

	// store one dimension variables (eg: TIME, LATITUDE, ...)
	//Variables_1D map[string]interface{}
	// store two dimensions variables (eg: PRES, DEPTH, TEMP, ...)
	//Variables_2D AllData_2D
	Variables matrix
	// store global attributes
	Attributes map[string]string

	// used to store max of profiles value
	ExtraFloat map[string]float64
	// used to store max of profiles type
	ExtraString map[string]string
	// give access to physical parameters
	Roscop Roscop

	// store header
	//hdr []string
}

// Process is an interface common for all data sets like profiles,
// trajectories and time-series
type Process interface {
	Read([]string)
	GetConfig(string)
	//	WriteHeader(map[string]string, []string)
	WriteASCII(map[string]string, []string)
	WriteNetcdf(instrumentType)
}

// nc implement interface Process
var nc Process

// Ctd define new receiver type based on netcdf equivalent structure
type Ctd struct{ Nc }

// Btl define new receiver type based on netcdf equivalent structure
type Btl struct{ Nc }

// NewCtd return an interface type based on Netcdf representation
func NewCtd() *Ctd { return &Ctd{} }

// NewBtl return an interface type based on Netcdf representation
func NewBtl() *Btl { return &Btl{} }

// main body
func main() {

	// slice of filename to read and extract data
	var files []string

	// to change the flags on the default logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if os.Getenv("OCEANO2OCEANSITES_CFG") != "" {
		cfgname = os.Getenv("OCEANO2OCEANSITES_CFG")
	}
	if os.Getenv("ROSCOP_CSV") != "" {
		codeRoscop = os.Getenv("ROSCOP_CSV")
	}

	// get options parse args list and return all given files to read
	// and configuration file name
	files, optCfgfile := GetOptions()

	// test if output directories exists and create them if not
	mkOutputDir()

	// read the first file and try to find the instrument type, return a bit mask
	typeInstrument = analyseFirstFile(files)

	// following the instrument type, allocate the rigth receiver based on
	// Process interface
	switch typeInstrument {
	case CTD:
		//nc = &Ctd{}
		nc = NewCtd()
	case BTL:
		//nc = &Btl{}
		nc = NewBtl()
	default:
		f("main: invalide option typeInstrument -> %d\n", typeInstrument)
		p("Exiting...")
		os.Exit(0)
	}

	// read configuration file, by default, optCfgfile = cfgname
	nc.GetConfig(optCfgfile)
	// debug
	fmt.Fprintln(debug, mapFormat)

	// read and process all data files
	nc.Read(files)

	// write ASCII file
	nc.WriteASCII(mapFormat, hdr)

	// write netcdf file
	nc.WriteNetcdf(typeInstrument)
}
