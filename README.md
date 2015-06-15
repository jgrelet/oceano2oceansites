# oceano2oceansites

This application read Seabird CTD cnv files, extract data from header files and write result into one ASCII file and NetCDF OceanSITES file.

Installation:

Install golang from https://golang.org

Install tdm-gcc from http://tdm-gcc.tdragon.net/

Install Git for Windows

Install Mercurial (hg) from https://mercurial.selenic.com/

Insall Netcdf 4.3 NC4-64 from http://www.unidata.ucar.edu/software/netcdf/docs/winbin.html

Install package go-netcdf from https://github.com/fhs/go-netcdf/
```
$ go get github.com/fhs/go-netcdf/netcdf
```
Install getopt package
```
go get code.google.com/p/getopt
```
Add some Seabird cnv files under data directory, 
for example data/fr24/*.cnv

Run and test (linux, cygwin):
```
$ go run *.go data/fr24/dfr24*.cnv 
```
on Windows console, use :
```
$ go run *.go --files=data/fr24/dfr24*.cnv 
```
Compile it:
```
$ go build
```

Usage:
```
> oceano2oceansites -h
Usage: oceano2oceansites.exe [-dehv] [-c value] [-f value] [-m value] [parameters ...]
 -c, --config=value
                    Name of the configuration file to use.
 -d                 debug
 -e                 debug
 -f, --files=value  files to process ex: data/fr25*.cnv
 -h                 help
 -m, --cycle_mesure=value
                    Name of cycle_mesure
 -v, --version      Show version, then exit.

```
Edit config file:
```
[global]
author         = your name

[cruise]
cycleMesure    = CRUISE_NAME
plateforme     = SHIP
callsign       = XXXX

[ctd]
cruisePrefix   = fr25
stationPrefixLength  = 3
header         = PRFL  PRES   DEPH   ETDD      TEMP    PSAL     DENS   SVEL    DOX2    FLU2    TUR3   NAVG
format         = %05.0f  %4.0f    %6.1f %10.6f    %7.4f    %7.4f   %7.5f   %6.3f  %7.2f   %6.3f   %7.3f   %4.0f
split          = ETDD,2,PRES,3,DEPH,4,TEMP,5,CNDC,7,FLU2,13,TUR3,14,DOX2,15,NAVG,17,PSAL,18,DENS,20,SVEL,22

#header         = PRFL  PRES   DEPH   ETDD      TE01    TE02     PSA1    PSA2    CND1    CND2   DEN1   DEN2   SVEL    FLU2   TUR3   DO12   DO22    DOV1  DOV2     DVT1     DVT2   NAVG
#split          = ETDD,2,PRES,3,DEPH,4,TE01,5,TE02,6,CND1,7,CND2,8,DOV1,9,DOV2,10,DVT1,11,DVT2,12,FLU2,13,TUR3,14,DO12,15,DO22,16,NAVG,17,PSA1,18,PSA2,19,DEN1,20,DEN2,21,SVEL,22
#format         = %05.0f  %4.0f    %6.1f %10.6f    %7.4f   %7.4f     %7.4f   %7.4f   %7.5f   %7.5f  %6.3f  %6.3f  %7.2f   %6.3f  %6.2f  %7.3f  %7.3f   %6.4f %6.4f   %+7.5f   %+7.5f   %4.0f
```
header and format are used for output ASCII file. header give also the parameter list for Netcdf file
split is used to extract parameter from seabird cnv files.

Check data with ncdump:
```
$ ncdump -v TIME  test_ctd.nc
netcdf test_ctd {
dimensions:
        TIME = 4 ;
        DEPTH = 13 ;
variables:
        double PROFILE(TIME) ;
                PROFILE:long_name = "PROFILE NUMBER" ;
                PROFILE:units = "N/A" ;
                PROFILE:valid_min = 1. ;
                PROFILE:valid_max = 99999. ;
                PROFILE:format = "%4.0f" ;
                PROFILE:_FillValue = 1.e+36 ;
        double TIME(TIME) ;
                TIME:long_name = "TIME" ;
                TIME:units = "days since 1950-01-01T00:00:00Z" ;
                TIME:valid_min = 0. ;
                TIME:valid_max = 90000. ;
                TIME:format = "%6.6d" ;
                TIME:_FillValue = 1.e+36 ;
        double LATITUDE(TIME) ;
                LATITUDE:long_name = "LATITUDE" ;
                LATITUDE:units = "decimal degree" ;
                LATITUDE:valid_min = -90. ;
                LATITUDE:valid_max = 90. ;
                LATITUDE:format = "%+8.4lf" ;
                LATITUDE:_FillValue = 1.e+36 ;
        double LONGITUDE(TIME) ;
                LONGITUDE:long_name = "LONGITUDE" ;
                LONGITUDE:units = "decimal degree" ;
                LONGITUDE:valid_min = -180. ;
                LONGITUDE:valid_max = 180. ;
                LONGITUDE:format = "%+9.4lf" ;
                LONGITUDE:_FillValue = 1.e+36 ;
        double BATH(TIME) ;
                BATH:long_name = "BATHYMETRIC DEPTH" ;
                BATH:units = "meter" ;
                BATH:valid_min = 0. ;
                BATH:valid_max = 11000. ;
                BATH:format = "%6.1lf" ;
                BATH:_FillValue = 1.e+36 ;
        double DEPH(TIME, DEPTH) ;
                DEPH:long_name = "DEPTH BELOW SEA SURFACE" ;
                DEPH:units = "meter" ;
                DEPH:valid_min = 0. ;
                DEPH:valid_max = 6000. ;
                DEPH:format = "%6.1lf" ;
                DEPH:_FillValue = 1.e+36 ;
        double TEMP(TIME, DEPTH) ;
                TEMP:long_name = "SEA TEMPERATURE" ;
                TEMP:units = "Celsius degree" ;
                TEMP:valid_min = 0. ;
                TEMP:valid_max = 30. ;
                TEMP:format = "%6.3lf" ;
                TEMP:_FillValue = 1.e+36 ;
        double PRES(TIME, DEPTH) ;
                PRES:long_name = "SEA PRESSURE sea surface=0" ;
                PRES:units = "decibar=10000 pascals" ;
                PRES:valid_min = 0. ;
                PRES:valid_max = 6500. ;
                PRES:format = "%6.1lf" ;
                PRES:_FillValue = 1.e+36 ;
        double DOX2(TIME, DEPTH) ;
                DOX2:long_name = "DISSOLVED OXYGEN" ;
                DOX2:units = "micromole/kg" ;
                DOX2:valid_min = 0. ;
                DOX2:valid_max = 450. ;
                DOX2:format = "%7.3lf" ;
                DOX2:_FillValue = 1.e+36 ;
        double PSAL(TIME, DEPTH) ;
                PSAL:long_name = "PRACTICAL SALINITY SCALE 1978" ;
                PSAL:units = "P.S.S.78" ;
                PSAL:valid_min = 33. ;
                PSAL:valid_max = 37. ;
                PSAL:format = "%6.3lf" ;
                PSAL:_FillValue = 1.e+36 ;

// global attributes:
                :cycle_mesure = "PIRATA-FR25" ;
                :plateforme = "THALASSA" ;
data:

 TIME = 23820.5651041667, 23821.5832638889, 23823.3512615741, 23826.5581481481 ;
}
```






