package main

import (
	"fmt"
	"github.com/fhs/go-netcdf/netcdf"
)

// creates the NetCDF file following nc structure.
func CreateNetcdfFile(filename string, nc Nc) error {

	// read roscop definition file for variables attributes
	var roscop = codeRoscopFromCsv("code_roscop.csv")
	//	for k, v := range roscop {
	//		fmt.Printf("%s: ", k)
	//		fmt.Println(v)
	//	}
	//	os.Exit(0)

	// if *optEcho
	fmt.Printf("writing netCDF: %s\n", filename)

	// get variables_1D size
	len_1D := nc.Dimensions["TIME"]
	len_2D := nc.Dimensions["DEPTH"]

	// Create a new NetCDF 3 file. The dataset is returned.
	ds, err := netcdf.CreateFile(filename, netcdf.CLOBBER)
	if err != nil {
		return err
	}
	defer ds.Close()

	// Add the dimensions for our data to the dataset
	dim_1D := make([]netcdf.Dim, 1)
	dim_2D := make([]netcdf.Dim, 2)

	// dimensions for ROSCOP paremeters as DEPTH, PRES, TEMP, PSAL, etc
	dim_2D[0], err = ds.AddDim("TIME", uint64(len_1D))
	if err != nil {
		return err
	}
	dim_2D[1], err = ds.AddDim("DEPTH", uint64(len_2D))
	if err != nil {
		return err
	}
	// dimension for PROFILE, LATITUDE, LONGITUDE and BATH
	dim_1D[0] = dim_2D[0]

	// Add the variable to the dataset that will store our data
	map_1D := make(map[string]netcdf.Var)
	for key, _ := range nc.Variables_1D {
		v, err := ds.AddVar(key, netcdf.DOUBLE, dim_1D)
		if err != nil {
			return err
		}
		map_1D[key] = v

		// define variables attributes, get values from roscop map
		// todos !!! for each type
		a := v.Attr("long_name")
		a.WriteBytes([]byte(roscop[key].long_name))
		a = v.Attr("units")
		a.WriteBytes([]byte(roscop[key].units))
		a = v.Attr("valid_min")
		a.WriteFloat64s([]float64{roscop[key].valid_min})
		a = v.Attr("valid_max")
		a.WriteFloat64s([]float64{roscop[key].valid_max})
		a = v.Attr("format")
		a.WriteBytes([]byte(roscop[key].format))
		a = v.Attr("_FillValue")
		a.WriteFloat64s([]float64{roscop[key]._FillValue})
	}

	map_2D := make(map[string]netcdf.Var)
	for key, _ := range nc.Variables_2D {
		v, err := ds.AddVar(key, netcdf.DOUBLE, dim_2D)
		if err != nil {
			return err
		}
		map_2D[key] = v

		// define variables attributes, get values from roscop map
		// todos !!! for each type
		a := v.Attr("long_name")
		a.WriteBytes([]byte(roscop[key].long_name))
		a = v.Attr("units")
		a.WriteBytes([]byte(roscop[key].units))
		a = v.Attr("valid_min")
		a.WriteFloat64s([]float64{roscop[key].valid_min})
		a = v.Attr("valid_max")
		a.WriteFloat64s([]float64{roscop[key].valid_max})
		a = v.Attr("format")
		a.WriteBytes([]byte(roscop[key].format))
		a = v.Attr("_FillValue")
		a.WriteFloat64s([]float64{roscop[key]._FillValue})
	}

	// defines global attributes
	for key, value := range nc.Attributes {
		a := ds.Attr(key)
		a.WriteBytes([]byte(value))
	}

	// leave define mode in NetCDF3
	ds.EndDef()

	// Create the data with the above dimensions and write it to the file.
	for key, value := range nc.Variables_1D {

		// if *optEcho
		fmt.Printf("writing %s: %d\n", key, len(value))

		err = map_1D[key].WriteFloat64s([]float64(value))
		if err != nil {
			return err
		}
	}

	// write data 2D (value.data) to netcdf variables (key)
	for key, value := range nc.Variables_2D {
		i := 0
		ht := len(value.data)
		wd := len(value.data[0])

		// if *optEcho
		fmt.Printf("writing %s: %d x %d\n", key, ht, wd)

		// Write<type> netcdf methods need []<type>, [][]data has flatten
		gopher := make([]float64, ht*wd)
		for x := 0; x < ht; x++ {
			for y := 0; y < wd; y++ {
				gopher[i] = value.data[x][y]
				i++
			}
		}
		err = map_2D[key].WriteFloat64s(gopher)
		if err != nil {
			return err
		}
	}

	// if *optEcho
	fmt.Printf("writing %s done ...\n", filename)

	return nil
}
