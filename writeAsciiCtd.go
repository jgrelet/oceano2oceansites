// writeAsciiCtd.go
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	codeForProfile = -1
)

// WriteAscii write ascii file
func (nc *Ctd) WriteAscii(mapFormat map[string]string, hdr []string) {
	// define 2 files, profiles header and data
	var filename string

	//	if _, err := os.Stat("/path/to/whatever"); err == nil {
	//		// path/to/whatever exists
	//	}

	// build filenames
	str := nc.Attributes["cycle_mesure"]
	str = strings.Replace(str, "\r", "", -1)
	headerFilename := fmt.Sprintf("%s/ascii/%s.ctd", outputDir, strings.ToLower(str))
	filename = fmt.Sprintf("%s/ascii/%s%s_ctd", outputDir, strings.ToLower(str), prefixAll)
	//fmt.Println(headerFilename)
	//fmt.Println(filename)

	// open header file for writing result
	fidHdr, err := os.Create(headerFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer fidHdr.Close()

	// use buffered mode for writing
	fbufHdr := bufio.NewWriter(fidHdr)

	// open ASCII file for writing result
	fid, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer fid.Close()

	// use buffered mode for writing
	fbuf := bufio.NewWriter(fid)

	// write header to string
	str = fmt.Sprintf("%s  %s  %s  %s %s  %s\n",
		nc.Attributes["cycle_mesure"],
		nc.Attributes["plateforme"],
		nc.Attributes["institute"],
		nc.Attributes["type_instrument"],
		nc.Attributes["instrument_number"],
		nc.Attributes["pi"])

	// write first line header on header file and ascii file
	fmt.Fprintf(fbufHdr, str)
	fmt.Fprintf(fbuf, str)

	// display on screen
	fmt.Printf("%s", str)

	// write physical parameters in second line
	str = ""
	for _, key := range hdr {
		fmt.Fprintf(fbuf, "%s   ", key)
		fmt.Fprintf(debug, "%s   ", key)
	}
	// append new line
	//fmt.Fprintln(fbuf_ascii, "\n")

	// write second line header on ascii file
	fmt.Fprintln(fbuf, str)

	// display on screen
	fmt.Printf("%s", str)

	// get data (slices) from nc struct
	len1D := nc.Dimensions["TIME"]
	len2D := nc.Dimensions["DEPTH"]
	time := nc.Variables.flatten("TIME")
	lat := nc.Variables.flatten("LATITUDE")
	lon := nc.Variables.flatten("LONGITUDE")
	profile := nc.Variables.flatten("PROFILE")
	bath := nc.Variables.flatten("BATH")

	// loop over each profile
	for x := 0; x < len1D; x++ {
		str = ""
		// write profile informations to ASCII data file with DEPTH = -1
		t1 := NewTimeFromJulian(time[x])
		t2 := NewTimeFromJulianDay(nc.ExtraFloat[fmt.Sprintf("ETDD:%d", int(profile[x]))], t1)
		// TODOS: adapt profile format to stationPrefixLength
		fmt.Fprintf(fbuf, "%05.0f %4d %f %f %f %s",
			profile[x],
			codeForProfile,
			t1.JulianDayOfYear(),
			lat[x],
			lon[x],
			t1.Format("20060102150405"))

		// write profile informations to header file, max depth CTD and
		// bathymetrie are in meters
		str = fmt.Sprintf("%05.0f %s %s %s %s %4.4g %4.4g %s %s\n",
			profile[x],
			t1.Format("02/01/2006 15:04:05"),
			t2.Format("02/01/2006 15:04:05"),
			DecimalPosition2String(lat[x], "NS"),
			DecimalPosition2String(lon[x], "EW"),
			nc.ExtraFloat[fmt.Sprintf("DEPTH:%d", int(profile[x]))],
			bath[x],
			nc.ExtraString[fmt.Sprintf("TYPE:%d", int(profile[x]))],
			cfg.Ctd.CruisePrefix+nc.ExtraString[fmt.Sprintf("PRFL_NAME:%d", int(profile[x]))])

		// write profile information to header file
		fmt.Fprintf(fbufHdr, str)

		// display on screen
		fmt.Printf("%s", str)

		// fill last header columns with 1e36
		for i := 0; i < len(hdr)-6; i++ {
			fmt.Fprintf(fbuf, " %g", 1e36)
		}
		fmt.Fprintln(fbuf) // add newline

		// loop over each level
		for y := 0; y < len2D; y++ {
			// goto next profile when max depth reach
			if nc.Variables.get("PRES", x, y).(float64) >=
				nc.ExtraFloat[fmt.Sprintf("PRES:%d", int(profile[x]))] {
				continue
			}
			fmt.Fprintf(fbuf, "%05.0f ", profile[x])
			// loop over each physical parameter (key) in the rigth order
			for _, key := range hdr {
				// if key not in map, goto next key
				if _, ok := nc.Variables[key]; ok {
					// fill 2D slice
					data := nc.Variables.get(key, x, y)
					// print data with it's format, change format for FillValue
					if data == 1e36 {
						fmt.Fprintf(fbuf, "%g ", data)
					} else {
						fmt.Fprintf(fbuf, mapFormat[key]+" ", data)
					}
				}
			}
			fmt.Fprintf(fbuf, "\n")

		}
		fbuf.Flush()
		fbufHdr.Flush()
	}
}
