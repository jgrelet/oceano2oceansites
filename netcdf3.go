package main

import (
	"fmt"
	"github.com/fhs/go-netcdf/netcdf"
)

// CreateExampleFile creates an example NetCDF file containing only one variable.
func CreateNetcdfFile(filename string, nc Nc) error {

	// if *optEcho
	fmt.Printf("wrinting netCDF: %s\n", filename)

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
	}
	map_2D := make(map[string]netcdf.Var)
	for key, _ := range nc.Variables_2D {
		v, err := ds.AddVar(key, netcdf.DOUBLE, dim_2D)
		if err != nil {
			return err
		}
		map_2D[key] = v

		// define attribbute, modify it !!!!
		a := v.Attr("_FillValue")
		a.WriteFloat64s([]float64{1e36})
	}
	/*
		// defined variable attributes
		a := v.Attr("long_name")
		a.WriteBytes([]byte("TEMPERATURE"))
		a = v.Attr("max_value")
		a.WriteInt32s([]int32{32})
		a = v.Attr("min_value")
		a.WriteInt32s([]int32{0})
	*/

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
		fmt.Printf("wrinting %s: %d\n", key, len(value))

		err = map_1D[key].WriteFloat64s([]float64(value))
		if err != nil {
			return err
		}
	}

	for key, value := range nc.Variables_2D {

		i := 0
		ht := len(value.data)
		wd := len(value.data[0])

		// if *optEcho
		fmt.Printf("wrinting %s: %d x %d\n", key, ht, wd)

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
	fmt.Printf("wrinting %s done ...", filename)

	return nil
}
