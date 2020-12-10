package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fhs/go-netcdf/netcdf"
)

// WriteNetcdf creates the NetCDF file following nc structure.
func (nc *Nc) WriteNetcdf(inst instrumentType) {

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
	filename := fmt.Sprintf("%s/netcdf/OS_%s%s_%s.nc",
		outputDir,
		strings.ToUpper(nc.Attributes["cycle_mesure"]),
		strings.ToUpper(prefixAll),
		ncType)

	// get roscop definition file for variables attributes
	var roscop = nc.Roscop
	fmt.Fprintf(echo, "writing file: %s\n", filename)
	fmt.Fprintf(debug, "writing file: %s\n", filename)

	// get variables_1D size
	len1D := nc.Dimensions["TIME"]
	len2D := nc.Dimensions["DEPTH"]

	// Create a new NetCDF 3 file. The dataset is returned.
	ds, err := netcdf.CreateFile(filename, netcdf.CLOBBER)
	if err != nil {
		log.Fatal(err)
	}
	defer ds.Close()

	// Add the dimensions for our data to the dataset
	dim1D := make([]netcdf.Dim, 1)
	dim2D := make([]netcdf.Dim, 2)

	// dimensions for ROSCOP paremeters as DEPTH, PRES, TEMP, PSAL, etc
	dim2D[0], err = ds.AddDim("TIME", uint64(len1D))
	if err != nil {
		log.Fatal(err)
	}
	dim2D[1], err = ds.AddDim("DEPTH", uint64(len2D))
	if err != nil {
		log.Fatal(err)
	}
	// dimension for PROFILE, LATITUDE, LONGITUDE and BATH
	dim1D[0] = dim2D[0]

	// Add the variable to the dataset that will store our data
	map1D := make(map[string]netcdf.Var)

	// create netcdf variables with attributes
	// for key, _ := range nc.Variables {
	for _, key := range nc.GetPhysicalParametersList() {
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
		// add variables
		dims := make([]netcdf.Dim, 1)
		if nc.Variables.isMatrix(key) {
			dims = dim2D
		} else {
			dims = dim1D
		}
		v, err := ds.AddVar(key, netcdfType, dims)
		if err != nil {
			log.Fatal(fmt.Sprintf("%s, %v: %v (%T)", err, key, v, v))
		}
		fmt.Fprintf(debug, "AddVar: %s\n", key)
		map1D[key] = v

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
				err = a.WriteBytes([]byte(value.(string)))
			case int32:
				err = a.WriteInt32s([]int32{value.(int32)})
			case float32:
				err = a.WriteFloat32s([]float32{value.(float32)})
			case float64:
				err = a.WriteFloat64s([]float64{value.(float64)})
			default:
				log.Fatal(fmt.Sprintf("netcdf: create 1D attribute error, %v=%v:%v (%T)",
					key, name, value, value)) // wrong type, check code_roscop file
			}
			if err != nil {
				log.Fatal(fmt.Sprintf("%s, %v: %v (%T)", err, key, v, v))
			}
			fmt.Fprintf(debug, "%s: %s=%v (%T)\n", key, name, value, value)
		}
	}

	// add some global attributes
	nc.Attributes["data_type"] = "OceanSITES profile data"
	nc.Attributes["format_version"] = "1.2"
	nc.Attributes["Conventions"] = "CF-1.4, OceanSITES-1.2"
	nc.Attributes["netcdf_version"] = "3.6"

	// write global attributes
	for key, value := range nc.Attributes {
		a := ds.Attr(key)
		err = a.WriteBytes([]byte(value))
		if err != nil {
			log.Fatal(fmt.Sprintf("%s, %v: %v (%T)", err, key, value, value))
		}
	}

	// leave define mode in NetCDF3
	ds.EndDef()

	// Create the data with the above dimensions and type,
	// write them to the file.
	for _, key := range nc.GetPhysicalParametersList() {
		// remove PRFL from the key list
		if key == "PRFL" {
			continue
		}
		value := nc.Variables[key]
		v := nc.Variables.flatten(key)
		fmt.Fprintf(echo, nc.Variables.printInfo(key))
		fmt.Fprintf(debug, nc.Variables.printInfo(key))

		// convert types from code_roscop structure to native netcdf types
		switch roscop.GetAttributesStringValue(key, "types") {
		case "int32":
			r := Matrix2int32(v)
			if err := map1D[key].WriteInt32s(r); err != nil {
				log.Fatal(fmt.Sprintf("%s, %v: (%T) %v", err, key, value, value))
			}
		case "float32":
			r := Matrix2float32(v)
			if err := map1D[key].WriteFloat32s(r); err != nil {
				log.Fatal(fmt.Sprintf("%s, %v: (%T) %v", err, key, value, value))
			}
		case "float64":
			if err := map1D[key].WriteFloat64s(v); err != nil {
				log.Fatal(fmt.Sprintf("%s, %v: (%T) %v", err, key, v, v))
			}
		default:
			log.Fatal(fmt.Sprintf("%s, %v", err, key))
		}

	}

	fmt.Fprintf(echo, "writing %s done ...\n", filename)
}
