package main

import (
	"fmt"
	"log"
	"os"
	//	"strconv"
	"strings"

	"github.com/fhs/go-netcdf/netcdf"
)

// creates the NetCDF file following nc structure.
//func WriteNetcdf(any interface{}) error {
func (nc *Nc) WriteNetcdf(inst InstrumentType) {

	var ncType string

	switch inst {
	case CTD:
		ncType = "CTD"
	case BTL:
		ncType = "BTL"
		//	case XBT:
		//		ncType = "XBT"
	default:
		fmt.Printf("WriteNetcdf: invalide InstrumentType -> %d\n", inst)
		fmt.Println("Exiting...")
		os.Exit(0)
	}

	// build filename
	filename := fmt.Sprintf("netcdf/OS_%s%s_%s.nc",
		strings.ToUpper(nc.Attributes["cycle_mesure"]),
		strings.ToUpper(prefixAll),
		ncType)
	//fmt.Println(filename)

	// get roscop definition file for variables attributes
	var roscop = nc.Roscop
	//	for k, v := range roscop {
	//		fmt.Printf("%s: ", k)
	//		fmt.Println(v)
	//	}
	//	os.Exit(0)

	fmt.Fprintf(echo, "writing netCDF: %s\n", filename)

	// get variables_1D size
	len_1D := nc.Dimensions["TIME"]
	len_2D := nc.Dimensions["DEPTH"]

	// Create a new NetCDF 3 file. The dataset is returned.
	ds, err := netcdf.CreateFile(filename, netcdf.CLOBBER)
	if err != nil {
		log.Fatal(err)
	}
	defer ds.Close()

	// Add the dimensions for our data to the dataset
	dim_1D := make([]netcdf.Dim, 1)
	dim_2D := make([]netcdf.Dim, 2)

	// dimensions for ROSCOP paremeters as DEPTH, PRES, TEMP, PSAL, etc
	dim_2D[0], err = ds.AddDim("TIME", uint64(len_1D))
	if err != nil {
		log.Fatal(err)
	}
	dim_2D[1], err = ds.AddDim("DEPTH", uint64(len_2D))
	if err != nil {
		log.Fatal(err)
	}
	// dimension for PROFILE, LATITUDE, LONGITUDE and BATH
	dim_1D[0] = dim_2D[0]

	// Add the variable to the dataset that will store our data
	map_1D := make(map[string]netcdf.Var)

	// create netcdf variables with attributes
	for key, _ := range nc.Variables_1D {

		// convert types from code_roscop structure to native netcdf types
		var netcdfType netcdf.Type

		pa := roscop.GetAttributesStringValue(key, "types")
		switch pa {
		case "int32":
			netcdfType = netcdf.INT
		case "float32":
			netcdfType = netcdf.FLOAT
		case "float64":
			netcdfType = netcdf.DOUBLE
		default:
			log.Fatal(fmt.Sprintf("Error: key: %s, Value: [%s], check roscop file\n", key, pa)) // wrong type, check code_roscop file
		}
		// add variables
		v, err := ds.AddVar(key, netcdfType, dim_1D)
		if err != nil {
			log.Fatal(err)
		}
		map_1D[key] = v

		// define variable attributes with the right type
		// for an physical parameter, get a slice of attributes name
		for _, name := range roscop.GetAttributes(key) {
			// for each attribute, get the value
			value := roscop.GetAttributesValue(key, name)
			// add new attribute to the variable v
			a := v.Attr(name)
			// value is an interface{}, need type assertion
			switch value.(type) {
			case string:
				a.WriteBytes([]byte(value.(string)))
			case int32:
				a.WriteInt32s([]int32{value.(int32)})
			case float32:
				a.WriteFloat32s([]float32{value.(float32)})
			case float64:
				a.WriteFloat64s([]float64{value.(float64)})
			default:
				log.Fatal("netcdf: create 1D attribute error") // wrong type, check code_roscop file
			}
		}
	}

	// Add the variable to the dataset that will store our data
	map_2D := make(map[string]netcdf.Var)

	// use the order list gave by split or splitAll (config file) because
	// the iteration order is not specified and is not guaranteed to be
	// the same from one iteration to the next in golang
	// for key, _ := range nc.Variables_2D {
	for _, key := range hdr {
		// remove PRFL from the key list
		if key == "PRFL" {
			continue
		}
		// convert types from code_roscop structure to native netcdf types
		var netcdfType netcdf.Type

		pa := roscop.GetAttributesStringValue(key, "types")
		switch pa {
		case "int32":
			netcdfType = netcdf.INT
		case "float32":
			netcdfType = netcdf.FLOAT
		case "float64":
			netcdfType = netcdf.DOUBLE
		default:
			log.Fatal(fmt.Sprintf("Error: key: %s, Value: [%s], check roscop file\n", key, pa)) // wrong type, check code_roscop file
		}
		v, err := ds.AddVar(key, netcdfType, dim_2D)
		if err != nil {
			log.Fatal(err)
		}
		map_2D[key] = v

		// define variable attributes with the right type
		// for an physical parameter, get a slice of attributes name
		for _, name := range roscop.GetAttributes(key) {
			// for each attribute, get the value
			value := roscop.GetAttributesValue(key, name)
			// add new attribute to the variable v
			a := v.Attr(name)
			// value is an interface{}, need type assertion
			switch value.(type) {
			case string:
				a.WriteBytes([]byte(value.(string)))
			case int32:
				a.WriteInt32s([]int32{value.(int32)})
			case float32:
				a.WriteFloat32s([]float32{value.(float32)})
			case float64:
				a.WriteFloat64s([]float64{value.(float64)})
			default:
				log.Fatal(fmt.Sprintf("netcdf: create 1D attribute error: key: %s, Value: [%s], \n", key, value))
			}
		}
	}

	// defines global attributes
	for key, value := range nc.Attributes {
		a := ds.Attr(key)
		a.WriteBytes([]byte(value))
	}

	// leave define mode in NetCDF3
	ds.EndDef()

	// Create the data with the above dimensions and type,
	// write them to the file.
	for key, value := range nc.Variables_1D {

		// convert types from code_roscop structure to native netcdf types
		switch roscop.GetAttributesStringValue(key, "types") {
		case "int32":
			//tmp := value.([]int32)
			length := len(value.([]float64))
			v := make([]int32, length)
			for i := 0; i < length; i++ {
				v[i] = int32(value.([]float64)[i])
			}
			err = map_1D[key].WriteInt32s(v)
			fmt.Fprintf(echo, "writing %s: %d\n", key, len(v))
			if err != nil {
				log.Fatal(err)
			}
		case "float32":
			//v := value.([]float32)
			length := len(value.([]float64))
			v := make([]float32, length)
			for i := 0; i < length; i++ {
				v[i] = float32(value.([]float64)[i])
			}
			err = map_1D[key].WriteFloat32s(v)
			fmt.Fprintf(echo, "writing %s: %d\n", key, len(v))
			if err != nil {
				log.Fatal(err)
			}
		case "float64":
			//v := value.([]float64)
			length := len(value.([]float64))
			v := make([]float64, length)
			for i := 0; i < length; i++ {
				v[i] = float64(value.([]float64)[i])
			}
			err = map_1D[key].WriteFloat64s(v)
			fmt.Fprintf(echo, "writing %s: %d\n", key, len(v))
			if err != nil {
				log.Fatal(err)
			}
		default:
			log.Fatal(err) // wrong type, check code_roscop file
		}

	}

	// write data 2D (value.data) to netcdf variables
	// for key, value := range nc.Variables_2D {
	for _, key := range hdr {
		// remove PRFL from the key list
		if key == "PRFL" {
			continue
		}
		value := nc.Variables_2D[key]
		i := 0
		ht := len(value.data)
		wd := len(value.data[0])
		fmt.Fprintf(echo, "writing %s: %d x %d\n", key, ht, wd)
		fmt.Fprintf(debug, "writing %s: %d x %d\n", key, ht, wd)
		// Write<type> netcdf methods need []<type>, [][]data will be flatten
		gopher := make([]float64, ht*wd)
		for x := 0; x < ht; x++ {
			for y := 0; y < wd; y++ {
				gopher[i] = value.data[x][y]
				i++
			}
		}
		switch roscop.GetAttributesStringValue(key, "types") {
		case "int32":
			v := make([]int32, ht*wd)
			for i := 0; i < ht*wd; i++ {
				v[i] = int32(gopher[i])
			}
			err = map_2D[key].WriteInt32s(v)
			if err != nil {
				log.Fatal(err)
			}
		case "float32":
			v := make([]float32, ht*wd)
			for i := 0; i < ht*wd; i++ {
				v[i] = float32(gopher[i])
			}
			err = map_2D[key].WriteFloat32s(v)
			if err != nil {
				log.Fatal(err)
			}
		case "float64":
			err = map_2D[key].WriteFloat64s(gopher)
			if err != nil {
				log.Fatal(err)
			}
		default:
			log.Fatal(err) // wrong type, check code_roscop file
		}
	}
	fmt.Fprintf(echo, "writing %s done ...\n", filename)
	//return nil
}
