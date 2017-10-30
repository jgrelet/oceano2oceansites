// GetOptions
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fhs/go-netcdf/netcdf"
	"github.com/pborman/getopt"
)

// global arg list options
var optDebug *bool
var optEcho *bool
var optAll *bool

// GetOptions parse arguments on command line and return files list to process
func GetOptions() ([]string, string) {

	var files []string
	var err error

	// parse options, moved outside main
	optDebug = getopt.Bool('d', "debug", "Display debug info")
	optEcho = getopt.Bool('e', "echo", "Display processing in stdout")
	optHelp := getopt.Bool('h', "help", "Help")
	optAll = getopt.Bool('a', "all", "Process all parameters")
	optVersion := getopt.BoolLong("version", 'v', "Show version, then exit.")
	optCfgfile := getopt.StringLong("config", 'c', cfgname, "use specific configuration .toml file", "oceano2oceansites.toml")
	//	optCycleMesure := getopt.StringLong("cycle_mesure", 'm', "", "Name of cycle_mesure")
	optFiles := getopt.StringLong("files", 'f', "", "files to process ex: data/fr25*.cnv", "files")
	optRoscop := getopt.StringLong("roscop", 'r', codeRoscop, "use a specific .csv file for physical parameter ", "code_roscop.csv")
	optOutput := getopt.StringLong("output", 'o', outputDir, "write files into", "output directory")

	// parse options line argument
	getopt.Parse()

	// process bloc when option is set
	if *optHelp {
		getopt.Usage()
		os.Exit(0)
	}
	if *optDebug {
		debug = os.Stdout
	}
	if *optEcho {
		echo = os.Stdout
	}
	if *optAll {
		prefixAll = "-all"
	}
	// get files list from argument line
	// Args() returns the non-option arguments
	if *optFiles == "" {
		files = getopt.Args()
	} else {
		files, err = filepath.Glob(*optFiles)
	}
	fmt.Fprintln(debug, files)

	if err != nil {
		p(err)
	}
	if *optCfgfile != "" {
		cfgname = *optCfgfile
	}
	if *optRoscop != "" {
		codeRoscop = *optRoscop
	}
	if *optOutput != "" {
		outputDir = *optOutput
	}
	// show version and env
	if *optVersion {
		fmt.Println(Binary + ": v" + Version + ", build date: " + BuildTime)
		fmt.Printf("Netcdf version: %s\n", netcdf.Version())
		fmt.Printf("Environnement variables:\n")
		v := os.Getenv("OCEANO2OCEANSITES_CFG")
		if v == "" {
			v = "not defined"
		}
		fmt.Printf(" - OCEANO2OCEANSITES_CFG: %s\n", v)
		r := os.Getenv("ROSCOP_CSV")
		if r == "" {
			r = "not defined"
		}
		fmt.Printf(" - ROSCOP_CSV: %s\n", r)
		fmt.Printf("Configuration file: %s\n", cfgname)
		fmt.Printf("Code ROSCOP file: %s\n", codeRoscop)
		fmt.Printf("GOPATH: %s\n", os.Getenv("GOPATH"))
		fmt.Printf("GOBIN: %s\n", os.Getenv("GOBIN"))
		os.Exit(0)
	}
	// if no files supplied for arg list, test if files is empty
	if len(files) == 0 {
		getopt.Usage()
		fmt.Println("\nError: please, specify files to process or define --files options")
		os.Exit(0)
	}
	//	if *optCycleMesure != "" {
	//		fmt.Println(*optCycleMesure)
	//		nc.Attributes["cycle_mesure"] = *optCycleMesure
	//	}
	fmt.Fprintln(debug, "Arg files list: ", files)

	return files, *optCfgfile
}
