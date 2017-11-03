# Windows installation

You have the choice to install the 32-bit or 64-bit toolchain (go and gcc compilers and Netcdf library). The gcc compiler is needed to link Netcdf library with the go-netcdf package. The installation of gcc under Windows 64 bit is a little tricky due to exceptions propagation.

## 32 bit

* Download golang from [https://golang.org/dl](https://golang.org/dl) and follow the installation instructions
  * from binary: [https://golang.org/doc/install](https://golang.org/doc/install)
  * from source : [https://golang.org/doc/install/source](https://golang.org/doc/install/source)
* Download the mingw online installer from [sourceforge](http://sourceforge.net/projects/mingw/files/latest/download?source=files)
* Run it and select only msys-base, mingw-developper-tools,mingw32-base and pkg-config
* Install [Git for Windows](https://git-scm.com/download/win)
* Install [Mercurial](https://mercurial.selenic.com/) (hg)
* Install [Netcdf 4.5.0 NC4-32](http://www.unidata.ucar.edu/software/netcdf/docs/winbin.html) under c:\opt\netCDF directory for example

## 64 bit

* Download golang from [https://golang.org/dl](https://golang.org/dl) and follow the installation instructions
  * from binary: [https://golang.org/doc/install](https://golang.org/doc/install)
  * from source : [https://golang.org/doc/install/source](https://golang.org/doc/install/source)
* Download the mingw on line installer from [sourceforge](http://sourceforge.net/projects/mingw/files/latest/download?source=files)
* Run it and select only msys-base
* Download and run gcc mingw64 installer from: <http://sourceforge.net/projects/mingw-w64/files/latest/download?source=files>
* Run the installer and select:

    version: 4.9.3
    Arch: x86_64
    threads: win32
    exception: sjlj or seh  (both option works)
    Build revision: 1

* Install [Git for Windows](https://git-scm.com/download/win)
* Install [Mercurial](https://mercurial.selenic.com/)
* Install [Netcdf 4.5.0 NC4-64] (http://www.unidata.ucar.edu/software/netcdf/docs/winbin.html) under c:\opt\netCDF directory for example

## Setting environment

* Rename fstab.sample to fstab, edit fstab and change mount /mingw to your gcc mingw64 directory.


   vi C:\MinGW\msys\1.0\etc\fstab


_example:_

    c:/mingw-w64/x86_64-4.9.3-win32-seh-rt_v4-rev1/mingw64  /mingw
    c:/ActiveState/perl /perl
    c:/users/<your_home>/go /go

* Update your path env with setx

    > setx path "%path%;C:\go\bin;C:\opt\netCDF\bin;C:\Program Files (x86)\Git\bin;C:\Program Files\Mercurial\"

* Run MSYS command tool from C:\MinGW\msys\1.0\msys.bat and check your go and gcc version:

    > go version
    go version go1.9.2 windows/amd64

    > gcc --version
    gcc.exe (x86_64-win32-seh-rev1, Built by MinGW-W64 project) 4.9.3

## Install packages

* package [go-netcdf](https://github.com/jgrelet/go-netcdf)

    > go get github.com/jgrelet/go-netcdf/netcdf

_Installation should be failed during compilation, the pkg-config method currently used to detect C library is not installed under Windows. See <http://www.gaia-gis.it/spatialite-3.0.0-BETA/mingw_how_to.html#pkg-config>_

Define PKG_CONFIG_PATH which is a environment variable that specifies additional paths in which pkg-config will search for its .pc files.

    > echo $PKG_CONFIG_PATH
    C:\opt\netCDF\lib\pkgconfig

Edit the C:\opt\netCDF\lib\pkgconfig\netcdf.pc file with the correct path.

    prefix=C:/opt/netCDF
    exec_prefix=C:/opt/netCDF
    libdir=C:/opt/netCDF/lib
    includedir=C:/opt/netCDF/include
    ccompiler=C:/MinGW/msys/1.0/bin/gcc.exe

    Name: netCDF
    Description: NetCDF Client Library for C
    URL: http://www.unidata.ucar.edu/netcdf
    Version: 4.5.0
    Libs: -L${libdir} -lnetcdf -lhdf5 -hdf5_hl -lzlib
    Cflags: -I${includedir}

Check it:

    > pkg-config --cflags netcdf
    -IC:/opt/netCDF/include
    > pkg-config --libs netcdf
    -LC:/opt/netCDF/lib -lnetcdf -lhdf5 -lhdf5_hl -lzlib

A faster implementation but not standard and portable is to change directly these cgo directives in `dataset.go` and `attribute.go` files before compilation.
cgo allows Go programs to interoperate with C libraries.

Replace the line :

    // #cgo pkg-config: nectcdf

 with the path where NetCDF is installed, c:/opt in this case for pseudo #cgo directives CFLAGS and LDFLAGS:

    // #cgo windows CFLAGS: -I C:/opt/netCDF/include
    // #cgo windows LDFLAGS: -lnetcdf -L C:/opt/netCDF/lib

Build manually the package go-netcdf:

    > cd github.com/jgrelet/go-netcdf/netcdf
    > go build -a -v

The netcdf.a library sould be installed under       `$GOPATH\pkg\windows_amd64\github.com\jgrelet\go-netcdf`

* package [oceano2oceansites](https://github.com/jgrelet/oceano2oceansites)

    > go get github.com/jgrelet/oceano2oceansites

This will install automatically these following packages:

* [https://github.com/pborman/getopt](https://github.com/pborman/getopt)
* [https://github.com/BurntSushi/toml](https://github.com/BurntSushi/toml)