# Linux installation

## Intel processors

Download and install the latest go version

Check if gcc compiler is installed (debian).

Check if NetCDF library and include files are installed:

    dpkg --get-selections|grep netcdf

If not, install them with apt-get:

    sudo apt-get install netcdf-bin libnetcdf-dev gcc

## ARM processors

On a Raspbery PIx, you  must install and compile the go compiler and NetCDF library from source, see:

* <https://golang.org/doc/install/source>
* <http://www.unidata.ucar.edu/software/netcdf/docs/getting_and_building_netcdf.html#building>

## Install packages

* package [oceano2oceansites](https://github.com/jgrelet/oceano2oceansites)

     go get github.com/jgrelet/oceano2oceansites

This will install automatically these following packages:

* [https://github.com/fhs/go-netcdf](https://github.com/fhs/go-netcdf)
* [https://github.com/pborman/getopt](https://github.com/pborman/getopt)
* [https://github.com/BurntSushi/toml](https://github.com/BurntSushi/toml)

The package go-netcdf wrap the netcdf C library with go and need the installation
of pkgconfig utility and libnetcdf-dev.