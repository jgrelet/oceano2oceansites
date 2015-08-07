package main

import (
	"code.google.com/p/gcfg"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

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
		Institute   string
		Pi          string
		Timezone    string
		BeginDate   string
		EndDate     string
		Creator     string
	}
	Ctd struct {
		CruisePrefix        string
		StationPrefixLength string
		Split               string
		SplitAll            string
		TypeInstrument      string
		InstrumentNumber    string
		TitleSummary        string
	}
	Btl struct {
		CruisePrefix        string
		StationPrefixLength string
		Split               string
		TypeInstrument      string
		InstrumentNumber    string
		TitleSummary        string
		Comment             string
	}
}

func (nc *Nc) GetConfig(configFile string) {

	//	var split, header, format string
	var split, splitAll string

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

	// add some global attributes for profile, change in future
	nc.Attributes["data_type"] = "OceanSITES profile data"

	err := gcfg.ReadFileInto(&cfg, configFile)
	if err == nil {
		switch typeInstrument {
		case CTD:
			split = cfg.Ctd.Split
			splitAll = cfg.Ctd.SplitAll
		case BTL:
			split = cfg.Btl.Split
		default:
			fmt.Printf("main: invalide option typeInstrument -> %d\n", typeInstrument)
			fmt.Println("Exiting...")
			os.Exit(0)

		}

		//		cruisePrefix = cfg.Ctd.CruisePrefix
		//		stationPrefixLength = cfg.Ctd.StationPrefixLength
		// TODOS: complete
		nc.Attributes["cycle_mesure"] = cfg.Cruise.CycleMesure
		nc.Attributes["plateforme"] = cfg.Cruise.Plateforme
		nc.Attributes["institute"] = cfg.Cruise.Institute
		nc.Attributes["pi"] = cfg.Cruise.Pi
		nc.Attributes["timezone"] = cfg.Cruise.Timezone
		nc.Attributes["begin_date"] = cfg.Cruise.BeginDate
		nc.Attributes["end_date"] = cfg.Cruise.EndDate
		nc.Attributes["creator"] = cfg.Cruise.Creator
		nc.Attributes["type_instrument"] = cfg.Ctd.TypeInstrument
		nc.Attributes["instrument_number"] = cfg.Ctd.InstrumentNumber

	} else {
		fmt.Println("function GetConfig error:")
		fmt.Printf("Please, check location for %s file\n", configFile)
		log.Fatal(err)
	}

	switch typeInstrument {
	case CTD:
		// First column should be PRFL
		hdr = append(hdr, "PRFL")
	case BTL:
		hdr = append(hdr, "PRFL")
		hdr = append(hdr, "ETDD")
	default:
	}

	// fill map_var from split
	// store the position (column) of each physical parameter
	var fields []string
	if *optAll {
		fields = strings.Split(splitAll, ",")
	} else {
		fields = strings.Split(split, ",")
	}
	fmt.Fprintln(debug, "getConfig: ", fields)

	// construct header slice from split
	for i := 0; i < len(fields); i += 2 {
		if v, err := strconv.Atoi(fields[i+1]); err == nil {
			map_var[fields[i]] = v - 1
			hdr = append(hdr, fields[i])
		}
	}
	fmt.Fprintln(debug, "getConfig: ", hdr)

	// fill map_format from code_roscop
	for _, key := range hdr {
		map_format[key] = nc.Roscop[key].format
	}
	//return nc
}
