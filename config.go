package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
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
		StationPrefixLength int
		Split               string
		SplitAll            string
		TypeInstrument      string
		InstrumentNumber    string
		TitleSummary        string
		Comment             string
	}
	Btl struct {
		CruisePrefix        string
		StationPrefixLength int
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

	// define map from netcdf structure
	nc.Dimensions = make(map[string]int)
	nc.Attributes = make(map[string]string)
	nc.Variables = make(matrix)
	nc.Extras_f = make(map[string]float64)
	nc.Extras_s = make(map[string]string)

	nc.Roscop = NewRoscop(code_roscop)

	//  read config file
	if _, err := toml.DecodeFile(configFile, &cfg); err != nil {
		log.Fatal(fmt.Sprintf("Error func GetConfig: file= %s -> %s\n", configFile, err))
	}
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

	// TODOS: complete
	nc.Attributes["cycle_mesure"] = cfg.Cruise.CycleMesure
	nc.Attributes["plateforme"] = cfg.Cruise.Plateforme
	nc.Attributes["callsign"] = cfg.Cruise.Callsign
	nc.Attributes["institute"] = cfg.Cruise.Institute
	nc.Attributes["pi"] = cfg.Cruise.Pi
	nc.Attributes["timezone"] = cfg.Cruise.Timezone
	nc.Attributes["begin_date"] = cfg.Cruise.BeginDate
	nc.Attributes["end_date"] = cfg.Cruise.EndDate
	nc.Attributes["creator"] = cfg.Cruise.Creator
	nc.Attributes["type_instrument"] = cfg.Ctd.TypeInstrument
	nc.Attributes["instrument_number"] = cfg.Ctd.InstrumentNumber

	// add specific column(s) to the first header line in ascii file
	switch typeInstrument {
	case CTD:
		// First column should be PRFL
		hdr = append(hdr, "PRFL")
	case BTL:
		hdr = append(hdr, "PRFL")
		hdr = append(hdr, "ETDD")
	default:
	}

	// fill map_var from split (read in .ini configuration file)
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
		// Change this call in next version !!!!
		map_format[key] = nc.Roscop.m[key]["format"]
	}
	//return nc
}

func (nc *Nc) InitVariables(dimx int, dimy int) {

	// initialize map entry from nil interface with _FillValue
	v := nc.Roscop.GetAttributesValue("PROFILE", "_FillValue").(int32)
	fmt.Fprintln(debug, "\nInitVariables():\n--------------")
	fmt.Fprintf(debug, "PROFILE with: %v (%T)\n", reflect.ValueOf(v), v)
	//nc.Variables_1D["PROFILE"] = fillSliceInt32(dimx, v)
	nc.Variables.NewMatrix("PROFILE", float64(v), dimx)

	y := nc.Roscop.GetAttributesValue("TIME", "valid_min").(float64)
	fmt.Fprintf(debug, "TIME with: %v (%T)\n", reflect.ValueOf(y), y)
	//nc.Variables_1D["TIME"] = fillSlice(dimx, y)
	nc.Variables.NewMatrix("TIME", float64(y), dimx)

	x := nc.Roscop.GetAttributesValue("LATITUDE", "_FillValue").(float32)
	//nc.Variables_1D["LATITUDE"] = fillSlice(dimx, 0) // need to be verify
	nc.Variables.NewMatrix("LATITUDE", float64(y), dimx)

	x = nc.Roscop.GetAttributesValue("LONGITUDE", "_FillValue").(float32)
	//nc.Variables_1D["LONGITUDE"] = fillSlice(dimx, 0)
	nc.Variables.NewMatrix("LONGITUDE", float64(y), dimx)

	x = nc.Roscop.GetAttributesValue("BATH", "_FillValue").(float32)
	fmt.Fprintf(debug, "BATH with: %v (%T)\n", reflect.ValueOf(x), x)
	//nc.Variables_1D["BATH"] = fillSlice(dimx, float64(x))
	nc.Variables.NewMatrix("BATH", float64(x), dimx)

	v = nc.Roscop.GetAttributesValue("TYPECAST", "_FillValue").(int32)
	fmt.Fprintf(debug, "TYPECAST with: %v (%T)\n", reflect.ValueOf(v), v)
	//nc.Variables_1D["TYPECAST"] = fillSliceInt32(dimx, v)
	nc.Variables.NewMatrix("TYPECAST", float64(v), dimx)

	// initialize 2D data
	//nc.Variables_2D = make(AllData_2D)
	for physicalParameter, _ := range map_var {
		// fmt.Printf("Initialize 2D var: %v\n", physicalParameter)
		fv := nc.Roscop.GetAttributesValue(physicalParameter, "_FillValue")
		switch fv.(type) {
		case int32:
			nc.Variables.NewMatrix(physicalParameter, float64(fv.(int32)), dimx, dimy)
		case float32:
			nc.Variables.NewMatrix(physicalParameter, float64(fv.(float32)), dimx, dimy)
		case float64:
			nc.Variables.NewMatrix(physicalParameter, float64(fv.(float64)), dimx, dimy)
		default:
			// todos !!!
		}

	}
}

// return an ordered list of parameters
func (nc *Nc) GetPhysicalParametersList() []string {
	r := append([]string{"PROFILE", "TIME", "LATITUDE", "LONGITUDE", "BATH", "TYPECAST"}, hdr...)
	return r
}
