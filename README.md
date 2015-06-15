# oceano2oceansites

This application read Seabird CTD cnv files, extract data from header files and write result into one ASCII file and NetCDF OceanSITES file.

Installation:

Install golang from https://golang.org

Install tdm-gcc from http://tdm-gcc.tdragon.net/

Install Git for Windows

Install Mercurial (hg) from https://mercurial.selenic.com/

Insall Netcdf 4.3 NC4-64 from http://www.unidata.ucar.edu/software/netcdf/docs/winbin.html

Install package go-netcdf from https://github.com/fhs/go-netcdf/
> go get github.com/fhs/go-netcdf/netcdf

Install getopt package
> go get code.google.com/p/getopt

Add some Seabird cnv files under data directory, 
for example data/fr24/*.cnv

Run and test (linux, cygwin):
> go run *.go data/fr24/dfr24*.cnv 

on Windows console, use :
> go run *.go --files=data/fr24/dfr24*.cnv 

Compile it:
> go build

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








