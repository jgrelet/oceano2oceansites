# Linux installation

## Intel processors

Download and install the latest go version

Check if gcc compiler is installed (debian).

Check if NetCDF library and include files are installed:

    > dpkg --get-selections|grep netcdf

If not, install them with apt-get:

    > sudo apt-get install netcdf-bin libnetcdf-dev gcc

**Notes:**
If you have previously installed the NetCDF library for Python with anaconda, you should not install it with apt-get.
But you have to define the following variables in your Linux session so that the gcc compiler can access the libraries and then perform linking.

    > export PKG_CONFIG_PATH=/home/<user>/miniconda3/lib/pkgconfig
    > export LD_LIBRARY_PATH=/home/<user>/miniconda3/lib

## ARM processors

On a Raspbery PIx, you  must install and compile the go compiler and NetCDF library from source, see:

* <https://golang.org/doc/install/source>
* <http://www.unidata.ucar.edu/software/netcdf/docs/getting_and_building_netcdf.html#building>

## Install packages

* package [oceano2oceansites](https://github.com/jgrelet/oceano2oceansites)

    > go get github.com/jgrelet/oceano2oceansites

This will install automatically these following packages:

* [https://github.com/jgrelet/go-netcdf](https://github.com/fhs/go-netcdf)
* [https://github.com/pborman/getopt](https://github.com/pborman/getopt)
* [https://github.com/BurntSushi/toml](https://github.com/BurntSushi/toml)

The package go-netcdf wrap the netcdf C library with go and need the installation of pkgconfig utility and libnetcdf-dev.
