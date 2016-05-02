// GetOptions
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pborman/getopt"
)

// global arg list options
var optDebug *bool
var optEcho *bool
var optAll *bool

func GetOptions() ([]string, string) {

	var files []string
	var err error

	// parse options, moved outside main
	optDebug = getopt.Bool('d', "debug", "Display debug info")
	optEcho = getopt.Bool('e', "echo", "Display processing in stdout")
	optHelp := getopt.Bool('h', "help", "Help")
	optAll = getopt.Bool('a', "all", "Process all parameters")
	optVersion := getopt.BoolLong("version", 'v', "Show version, then exit.")
	optCfgfile := getopt.StringLong("config", 'c', cfgname, "use specific configuration .ini file", "oceano2oceansites.ini")
	//	optCycleMesure := getopt.StringLong("cycle_mesure", 'm', "", "Name of cycle_mesure")
	optFiles := getopt.StringLong("files", 'f', "", "files to process ex: data/fr25*.cnv", "files")
	optRoscop := getopt.StringLong("roscop", 'r', code_roscop, "use a specific .csv file for physical parameter ", "code_roscop.csv")

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
		code_roscop = *optRoscop
	}
	// show version and env
	if *optVersion {
		fmt.Println(PROGNAME + ": v" + PROGVERSION)
		fmt.Printf("Environnement variable OCEANO2OCEANSITES_INI: %s\n", os.Getenv("OCEANO2OCEANSITES_INI"))
		fmt.Printf("Environnement variable ROSCOP_CSV: %s\n", os.Getenv("ROSCOP_CSV"))
		fmt.Printf("Configuration file: %s\n", cfgname)
		fmt.Printf("Code ROSCOP file: %s\n", code_roscop)
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
