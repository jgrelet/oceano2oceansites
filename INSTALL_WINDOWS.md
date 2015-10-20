# Windows installation

You have the choice to install the 32-bit or 64-bit toolchain (go and gcc compilers and Netcdf library). The gcc compiler is need to link Netcdf library with the go-netcdf package. The installation of gcc under Windows 64 bit is a little tricky due to exception propagation.

32 bit:
Download golang from https://golang.org/dl and follow the installation instructions
from binary: https://golang.org/doc/install 
or from source : https://golang.org/doc/install/source

Download the mingw on line installer from http://sourceforge.net/projects/mingw/files/latest/download?source=files

Run it and select only msys-base, mingw-developper-tools and mingw32-base

Install Git for Windows https://git-scm.com/download/win

Install Mercurial (hg) from https://mercurial.selenic.com/

Install Netcdf 4.3 NC4-32 from http://www.unidata.ucar.edu/software/netcdf/docs/winbin.html

64 bit:
Download golang from https://golang.org/dl and follow the installation instructions
from binary: https://golang.org/doc/install 
or from source : https://golang.org/doc/install/source

Download the mingw on line installer from http://sourceforge.net/projects/mingw/files/latest/download?source=files

Run it and select only msys-base
Download and run gcc mingw64 installer from http://sourceforge.net/projects/mingw-w64/files/latest/download?source=files.
Select:
version: 4.9.3
Arch: x86_64
threads: win32
exception: sjlj or seh  (both option works)
Build revision: 1 

Install Git for Windows https://git-scm.com/download/win

Install Mercurial (hg) from https://mercurial.selenic.com/

Install Netcdf 4.3 NC4-32 from http://www.unidata.ucar.edu/software/netcdf/docs/winbin.html

Setting environment:
Rename fstab.sample into fstab under C:\MinGW\msys\1.0\etc
Edit fstab and change mount /mingw to your gcc mingw64 directory.
example:
c:/mingw-w64/x86_64-4.9.3-win32-seh-rt_v4-rev1/mingw64		/mingw
c:/ActiveState/perl	/perl
c:/users/<your_home>/go	/go

Update your path env with setx
```
$ setx path "%path%;C:\go\bin;C:\opt\netCDF-4.3.3.1\bin;C:\Program Files (x86)\Git\bin;C:\Program Files\Mercurial\"
```
Run MSYS command tool from C:\MinGW\msys\1.0\msys.bat and check your go and gcc version:

$ go version
go version go1.5.1 windows/amd64

$ gcc --version
gcc.exe (x86_64-win32-seh-rev1, Built by MinGW-W64 project) 4.9.3


Install package go-netcdf from https://github.com/fhs/go-netcdf/
```
$ go get github.com/fhs/go-netcdf/netcdf
```
Installation should be failed during compilation, The pkg-config method currently used to detect C library is not installed under Windows. See http://www.gaia-gis.it/spatialite-3.0.0-BETA/mingw_how_to.html

A faster implementation is to change these cgo directives in `dataset.go` and `attribute.go` files before compilation

Replace :
```
// #cgo pkg-config: nectcdf
```
with:
```
// #cgo windows CFLAGS: -I C:/opt/netCDF-4.3.3.1/include
// #cgo windows LDFLAGS: -lnetcdf -L C:/opt/netCDF-4.3.3.1/lib
```
Restart package installation
```
$ go get github.com/fhs/go-netcdf/netcdf
```
The netcdf.a library sould be installed under `$GOPATH\pkg\windows_amd64\github.com\fhs\go-netcdf`

Install getopt package
```
$ go get github.com/pborman/getopt
```


