package main

// #cgo pkg-config: netcdf
// #include <stdlib.h>
// #include <netcdf.h>
import "C"

// NetcdfVersion returns a string identifying the version of the netCDF library, and when it was built.
// This function takes no arguments, and thus no errors are possible in its invocation.
func NetcdfVersion() string {
	return C.GoString(C.nc_inq_libvers())
}
