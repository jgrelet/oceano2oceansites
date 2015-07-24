// writeAscii.go
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

func WriteAsciiFiles(nc Nc, map_format map[string]string, hdr []string) {
	// define 2 files, profiles header and data
	var asciiFilename string

	// build filenames
	headerFilename := fmt.Sprintf("%s.ctd",
		strings.ToLower(nc.Attributes["cycle_mesure"]))
	asciiFilename = fmt.Sprintf("%s%s_ctd",
		strings.ToLower(nc.Attributes["cycle_mesure"]), prefixAll)
	//	fmt.Println(headerFilename)
	//	fmt.Println(asciiFilename)

	// open header file for writing result
	fid_hdr, err := os.Create(headerFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer fid_hdr.Close()

	// use buffered mode for writing
	fbuf_hdr := bufio.NewWriter(fid_hdr)

	// open ASCII file for writing result
	fid_ascii, err := os.Create(asciiFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer fid_ascii.Close()

	// use buffered mode for writing
	fbuf_ascii := bufio.NewWriter(fid_ascii)

	// write header to header file
	fmt.Fprintf(fbuf_hdr, "%s  %s  %s  %s %s  %s\n",
		nc.Attributes["cycle_mesure"],
		nc.Attributes["plateforme"],
		nc.Attributes["institute"],
		nc.Attributes["type_instrument"],
		nc.Attributes["instrument_number"],
		nc.Attributes["pi"])

	// write header to ascii file, first line
	fmt.Fprintf(fbuf_ascii, "%s  %s  %s  %s %s  %s\n",
		nc.Attributes["cycle_mesure"],
		nc.Attributes["plateforme"],
		nc.Attributes["institute"],
		nc.Attributes["type_instrument"],
		nc.Attributes["instrument_number"],
		nc.Attributes["pi"])

	// write physical parameters in second line
	for _, key := range hdr {
		fmt.Fprintf(fbuf_ascii, "%s   ", key)
		fmt.Fprintf(debug, "%s   ", key)
	}
	fmt.Fprintln(fbuf_ascii) // add new line

	// get data (slices) from nc struct
	len_1D := nc.Dimensions["TIME"]
	len_2D := nc.Dimensions["DEPTH"]
	time := nc.Variables_1D["TIME"]
	lat := nc.Variables_1D["LATITUDE"]
	lon := nc.Variables_1D["LONGITUDE"]
	bath := nc.Variables_1D["BATH"]
	profile := nc.Variables_1D["PROFILE"]

	// loop over each profile
	for x := 0; x < len_1D; x++ {
		// write profile informations to ASCII data file with DEPTH = -1
		t := NewTimeFromJulian(time[x])
		// TODOS: adapt profile format to stationPrefixLength
		fmt.Fprintf(fbuf_ascii, "%05.0f %4d %f %f %f %s",
			profile[x],
			codeForProfile,
			t.JulianDayOfYear(),
			lat[x],
			lon[x],
			t.Format("20060102150405"))

		// write profile informations to header file
		fmt.Fprintf(fbuf_hdr, "%05.0f %s %s %s %4.4g %4.4g\n",
			profile[x],
			t.Format("01/02/2006 15:04:05"),
			DecimalPosition2String(lat[x], 0),
			DecimalPosition2String(lon[x], 0),
			nc.Extras_f[fmt.Sprintf("DEPTH:%d", int(profile[x]))],
			bath[x])

		// fill last header columns with 1e36
		for i := 0; i < len(hdr)-6; i++ {
			fmt.Fprintf(fbuf_ascii, " %g", 1e36)
		}
		fmt.Fprintln(fbuf_ascii) // add newline

		// loop over each level
		for y := 0; y < len_2D; y++ {
			fmt.Fprintf(fbuf_ascii, "%05.0f ", profile[x])
			// loop over each physical parameter (key) in the rigth order
			for _, key := range hdr {
				// if key not in map, goto next key
				if _, ok := nc.Variables_2D[key]; ok {
					// fill 2D slice
					data := nc.Variables_2D[key].data[x][y]
					// print data with it's format, change format for FillValue
					if data == 1e36 {
						fmt.Fprintf(fbuf_ascii, "%g ", data)
					} else {
						fmt.Fprintf(fbuf_ascii, map_format[key]+" ", data)
					}
				}
			}
			fmt.Fprintf(fbuf_ascii, "\n")
			// goto next profile when max depth reach
			if nc.Variables_2D["DEPTH"].data[x][y] >=
				nc.Extras_f[fmt.Sprintf("DEPTH:%d", int(profile[x]))] {
				break
			}
		}
		fbuf_ascii.Flush()
		fbuf_hdr.Flush()
	}
}
