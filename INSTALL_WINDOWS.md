# Windows installation

You have the choice to install the 32-bit or 64-bit toolchain (go and gcc compilers and Netcdf library). The gcc compiler is needed to link Netcdf library with the go-netcdf package. The installation of gcc under 32 bit was tested under Windows 7 and 64 bit installation under Windows 7, 10 and 11.

It's better to use a git bash terminal window under VSC because I have a link editing problem with the command.com.

## 32 bit

* Download golang from [https://golang.org/dl](https://golang.org/dl) and follow the installation instructions
  * from binary: [https://golang.org/doc/install](https://golang.org/doc/install)
  * from source : [https://golang.org/doc/install/source](https://golang.org/doc/install/source)
* Download the mingw online installer from [sourceforge](http://sourceforge.net/projects/mingw/files/latest/download?source=files)
* Run it and select only msys-base, mingw-developper-tools,mingw32-base and pkg-config
* Install [Git for Windows](https://git-scm.com/download/win)
* Install [NetCDF library and tools](http://www.unidata.ucar.edu/software/netcdf/docs/winbin.html) under c:\opt\netCDF directory for example

## 64 bit

* Download golang from [https://golang.org/dl](https://golang.org/dl) and follow the installation instructions
  * from binary: [https://golang.org/doc/install](https://golang.org/doc/install)
  * from source : [https://golang.org/doc/install/source](https://golang.org/doc/install/source)
* Download the mingw-64 on line installer from [sourceforge](https://sourceforge.net/projects/mingw-w64/files/mingw-w64/)
  or install it with chocolatey which is an excellent alternative
* Test gcc installation with :

```bash
$ g++ -v
  Using built-in specs.
  ...
  Target: x86_64-w64-mingw32
  ...
```

* Install [Git for Windows](https://git-scm.com/download/win)
* Install [Netcdf library and tools](http://www.unidata.ucar.edu/software/netcdf/docs/winbin.html) under c:\opt\netCDF directory for example (don't install it in directory with space)

## Install pkg-config

* Install [pkg-config-lite](https://sourceforge.net/projects/pkgconfiglite/) inside gcc binary directory or with chocolatey.

* Define PKG_CONFIG_PATH which as an environment variable that specifies additional paths in which pkg-config will search for its .pc files.

```bash
$ export PKG_CONFIG_PATH=/c/opt/netCDF/lib/pkgconfig:$PKG_CONFIG_PATH

$ echo $PKG_CONFIG_PATH
/c/opt/netCDF/lib/pkgconfig:/mingw64/lib/pkgconfig:/mingw64/share/pkgconfig
```

Edit the C:\opt\netCDF\lib\pkgconfig\netcdf.pc file with the correct path.

```bash
prefix=C:/opt/netCDF
exec_prefix=C:/opt/netCDF
libdir=C:/opt/netCDF/lib
includedir=C:/opt/netCDF/include
#ccompiler=C:/opt/mingw-w64/x86_64-8.1.0-posix-seh-rt_v6-rev0/mingw64/bin/gcc.exe
ccompiler=C:/ProgramData/chocolatey/lib/mingw/tools/install/mingw64/bin/gcc.exe

Name: netCDF
Description: NetCDF Client Library for C
URL: http://www.unidata.ucar.edu/netcdf
Version: 4.x.x
Libs: -L${libdir} -lnetcdf -lhdf5 -lzlib
Cflags: -I${includedir}
```

Check it:

```bash
$ pkg-config --cflags netcdf 
-IC:/opt/netCDF/include
$ pkg-config --libs netcdf
-LC:/opt/netCDF/lib -lnetcdf -lhdf5 -lzlib
```
  
* ## Install packages

* package [go-netcdf](https://github.com/fhs/go-netcdf)

```bash
go get github.com/fhs/go-netcdf
```

Build manually the package go-netcdf (optional):

```bash
cd github.com/fhs/go-netcdf/netcdf
go build -a -v
```

The netcdf.a library sould be installed under       `$GOPATH\pkg\windows_amd64\github.com\fhs\go-netcdf`

* package [oceano2oceansites](https://github.com/jgrelet/oceano2oceansites)

```bash
go get github.com/jgrelet/oceano2oceansites
```

This will install automatically these following packages:

* [https://github.com/pborman/getopt](https://github.com/pborman/getopt)
* [https://github.com/BurntSushi/toml](https://github.com/BurntSushi/toml)
