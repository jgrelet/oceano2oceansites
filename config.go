package main

import (
	"code.google.com/p/gcfg"
	"fmt"
	"log"
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

func GetConfig(nc Nc, configFile string) Nc {

	//	var split, header, format string
	var split, splitAll, header string

	err := gcfg.ReadFileInto(&cfg, configFile)
	if err == nil {
		split = cfg.Ctd.Split
		splitAll = cfg.Ctd.SplitAll

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

	// First column should be PRFL
	hdr = append(hdr, "PRFL")

	// fill map_var from split
	// store the position (column) of each physical parameter
	var fields []string
	if *optAll {
		fields = strings.Split(splitAll, ",")
	} else {
		fields = strings.Split(split, ",")
	}
	fmt.Fprintln(debug, fields)

	for i := 0; i < len(fields); i += 2 {
		if v, err := strconv.Atoi(fields[i+1]); err == nil {
			map_var[fields[i]] = v - 1
			hdr = append(hdr, fields[i])
		}
	}

	// fill map_format from code_roscop
	for _, key := range hdr {
		map_format[key] = nc.Roscop[key].format
	}
	if *optDebug {
		fmt.Println(header)
	}
	return nc
}
