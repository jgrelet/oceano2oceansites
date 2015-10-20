// GetOptions
package main

import (
	"github.com/pborman/getopt"
	"fmt"
	"os"
	"path/filepath"
)

// global arg list options
var optDebug *bool
var optEcho *bool
var optAll *bool

func GetOptions() ([]string, string) {

	var files []string

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

	// process bloc when option is set
	if *optHelp {
		getopt.Usage()
		os.Exit(0)
	}
	if *optVersion {
		fmt.Println(PROGNAME + ": v" + PROGVERSION)
		fmt.Printf("Environnement variable OCEANO2OCEANSITES: %s\n", os.Getenv("OCEANO2OCEANSITES"))
		fmt.Printf("Environnement variable ROSCOP: %s\n", os.Getenv("ROSCOP"))
		fmt.Printf("Configuration file: %s\n", cfgname)
		fmt.Printf("Code ROSCOP file: %s\n", code_roscop)
		fmt.Printf("GOPATH: %s\n", os.Getenv("GOPATH"))
		fmt.Printf("GOBIN: %s\n", os.Getenv("GOBIN"))
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
		fmt.Println("\nPlease, specify files to process or define --files options")
		os.Exit(0)
	}
	//	if *optCycleMesure != "" {
	//		fmt.Println(*optCycleMesure)
	//		nc.Attributes["cycle_mesure"] = *optCycleMesure
	//	}
	fmt.Fprintln(debug, "Arg files list: ", files)
	return files, *optCfgfile
}
