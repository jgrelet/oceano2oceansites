# Linux installation

###Intel processors:

Download and install the latest go version

Check if gcc compiler is install (debian).

Check if NetCDF library and include files are installed:
````
dpkg --get-selections|grep netcdf
````

If not, install them with:
````
sudo apt-get install netcdf-bin libnetcdf-dev gcc
````

###ARM processors:

On a Raspbery PI2, you  must install and compile the go compiler and NetCDF library from source, see:
* https://golang.org/doc/install/source
* http://www.unidata.ucar.edu/software/netcdf/docs/getting_and_building_netcdf.html#building

###Install packages 
* package [go-netcdf](https://github.com/fhs/go-netcdf)
````
$ go get github.com/fhs/go-netcdf/netcdf
````
* package getopt 
````
$ go get github.com/pborman/getopt
````
* package [gcfg](https://gopkg.in/gcfg.v1)
````
go get gopkg.in/gcfg.v1
````

