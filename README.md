# oceano2oceansites

This application read Seabird CTD cnv files, extract data from header files and write result into one ASCII file and NetCDF OceanSITES file.

Binary programs under Windows and Linux are available from the depot under release:

https://github.com/jgrelet/oceano2oceansites/releases

See installation for OS:

- Windows: https://github.com/jgrelet/oceano2oceansites/blob/master/INSTALL_WINDOWS.md
- Linux:  https://github.com/jgrelet/oceano2oceansites/blob/master/INSTALL_LINUX.md

Add some Seabird cnv files under data directory, 
for example data/fr25/*.cnv

Run and test:
```
$ go run *.go --files=data/fr25/fr25*.cnv -e
```
Compile it:
```
$ go build
```
Install binary under go/bin:
```
$ go install
```

Usage:
```
$ oceano2oceansites -h
Usage: oceano2oceansites.exe [-dehv] [-c value] [-f value] [-m value] [parameters ...]
 -c, --config=value
                    Name of the configuration file to use.
 -d                 debug
 -e                 echo
 -f, --files=value  files to process ex: data/dfr25*.cnv
 -h                 help
 -m, --cycle_mesure=value
                    Name of cycle_mesure
 -v, --version      Show version, then exit.

```
The program use by default the configuration files oceano2oceansites.ini
and the code_roscop.csv in the current directory.

You can set a different location or name by setting environment variables OCEANO2OCEANSITES and CODE_ROSCOP.

Edit config file:
```
[global]
author         = your name

[cruise]
cycleMesure    = CRUISE_NAME
plateforme     = SHIP
callsign       = XXXX
institute      = IRD
pi             = PI_NAME
timezone       = GMT
beginDate      = 18/03/2015
endDate        = 13/04/2015
creator        = Firstname.Name@domaine.fr

[ctd]
cruisePrefix   = fr25
stationPrefixLength  = 3
titleSummary  = CTD profiles processed during this cruise
typeInstrument   = SBE911+
instrumentNumber  = 09Pxxxxx-xxx 
split          = PRES,3,DEPTH,4,ETDD,2,TEMP,5,PSAL,18,DENS,20,SVEL,22,DOX2,15,FLU2,13,TUR3,14,NAVG,17

splitAll          = PRES,3,DEPTH,4,ETDD,2,TE01,5,TE02,6,PSA1,18,PSA2,19,CND1,7,CND2,8,DEN1,20,DEN2,21,SVEL,22,DO12,15,DO22,16,DOV1,9,DOV2,10,DVT1,11,DVT2,12,FLU2,13,TUR3,14,NAVG,17
```
split describe the column order of each physical parameter to extract from seabird cnv files. 
The order is used for ASCII file output. 

All the physical parameters definition are decribed inside `code_roscop.csv`. You can update this file with your own definition.
Example:
```
TYPE;string;string;float64;float64;string;float64
TEMP;SEA TEMPERATURE;Celsius degree;0;30;%6.3f;1e36
SSTP;SEA SURFACE TEMPERATURE;Celsius degree;-1.5;38;%6.3f;1e36
PSAL;PRACTICAL SALINITY SCALE 1978;P.S.S.78;33;37;%6.3f;1e36
SSPS;SEA SURFACE PRACTICAL SALINITY;P.S.S.78;0;40;%6.3f;1e36
DOX1;DISSOLVED OXYGEN;ml/l;0;10;%5.2f;1e36
DOX2;DISSOLVED OXYGEN;micromole/kg;0;450;%7.3f;1e36
...
```

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
                PROFILE:format = "%5.0f" ;
                PROFILE:_FillValue = 1.e+36 ;
        double TIME(TIME) ;
                TIME:long_name = "TIME" ;
                TIME:units = "days since 1950-01-01T00:00:00Z" ;
                TIME:valid_min = 0. ;
                TIME:valid_max = 90000. ;
                TIME:format = "%6.6f" ;
                TIME:_FillValue = 1.e+36 ;
        double LATITUDE(TIME) ;
                LATITUDE:long_name = "LATITUDE" ;
                LATITUDE:units = "decimal degree" ;
                LATITUDE:valid_min = -90. ;
                LATITUDE:valid_max = 90. ;
                LATITUDE:format = "%+8.4f" ;
                LATITUDE:_FillValue = 1.e+36 ;
        double LONGITUDE(TIME) ;
                LONGITUDE:long_name = "LONGITUDE" ;
                LONGITUDE:units = "decimal degree" ;
                LONGITUDE:valid_min = -180. ;
                LONGITUDE:valid_max = 180. ;
                LONGITUDE:format = "%+9.4f" ;
                LONGITUDE:_FillValue = 1.e+36 ;
        double BATH(TIME) ;
                BATH:long_name = "BATHYMETRIC DEPTH" ;
                BATH:units = "meter" ;
                BATH:valid_min = 0. ;
                BATH:valid_max = 11000. ;
                BATH:format = "%6.1f" ;
                BATH:_FillValue = 1.e+36 ;
        double DEPH(TIME, DEPTH) ;
                DEPH:long_name = "DEPTH BELOW SEA SURFACE" ;
                DEPH:units = "meter" ;
                DEPH:valid_min = 0. ;
                DEPH:valid_max = 6000. ;
                DEPH:format = "%6.1f" ;
                DEPH:_FillValue = 1.e+36 ;
        double TEMP(TIME, DEPTH) ;
                TEMP:long_name = "SEA TEMPERATURE" ;
                TEMP:units = "Celsius degree" ;
                TEMP:valid_min = 0. ;
                TEMP:valid_max = 30. ;
                TEMP:format = "%6.3f" ;
                TEMP:_FillValue = 1.e+36 ;
        double PRES(TIME, DEPTH) ;
                PRES:long_name = "SEA PRESSURE sea surface=0" ;
                PRES:units = "decibar=10000 pascals" ;
                PRES:valid_min = 0. ;
                PRES:valid_max = 6500. ;
                PRES:format = "%6.1f" ;
                PRES:_FillValue = 1.e+36 ;
        double DOX2(TIME, DEPTH) ;
                DOX2:long_name = "DISSOLVED OXYGEN" ;
                DOX2:units = "micromole/kg" ;
                DOX2:valid_min = 0. ;
                DOX2:valid_max = 450. ;
                DOX2:format = "%7.3f" ;
                DOX2:_FillValue = 1.e+36 ;
        double PSAL(TIME, DEPTH) ;
                PSAL:long_name = "PRACTICAL SALINITY SCALE 1978" ;
                PSAL:units = "P.S.S.78" ;
                PSAL:valid_min = 33. ;
                PSAL:valid_max = 37. ;
                PSAL:format = "%6.3f" ;
                PSAL:_FillValue = 1.e+36 ;

// global attributes:
                :cycle_mesure = "CRUISE" ;
                :plateforme = "SHIP" ;
                ...
data:

 TIME = 23820.5651041667, 23821.5832638889, 23823.3512615741, 23826.5581481481 ;
}
```







