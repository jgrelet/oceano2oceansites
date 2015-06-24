// writeAscii.go
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func WriteAsciiFiles(nc Nc, map_format map[string]string, hdr []string) {
	// define 2 files, profiles header and data
	headerFilename := fmt.Sprintf("%s.ctd", nc.Attributes["cycle_mesure"])
	asciiFilename := fmt.Sprintf("%s_ctd", nc.Attributes["cycle_mesure"])

	// open header file for writing result
	fhdr, err := os.Create(headerFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer fhdr.Close()
	// use buffered mode for writing
	hdr_id := bufio.NewWriter(fhdr)

	// open ascii file for writing result
	fascii, err := os.Create(asciiFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer fascii.Close()
	// use buffered mode for writing
	ascii := bufio.NewWriter(fascii)

	// write header
	fmt.Fprintf(ascii, "%s  %s  %s  %s  %s  %s\n", nc.Attributes["cycle_mesure"],
		nc.Attributes["plateforme"], nc.Attributes["institute"],
		nc.Attributes["type"], nc.Attributes["sn"], nc.Attributes["pi"])
	// TODOS:
	for _, key := range hdr {
		fmt.Fprintf(ascii, "%s   ", key)
	}
	fmt.Fprintln(ascii) // add new line

	// dims
	len_1D := nc.Dimensions["TIME"]
	len_2D := nc.Dimensions["DEPTH"]
	time := nc.Variables_1D["TIME"]
	lat := nc.Variables_1D["LATITUDE"]
	lon := nc.Variables_1D["LONGITUDE"]
	profile := nc.Variables_1D["PROFILE"]
	date := nc.Extras_s
	// loop over each profile
	for x := 0; x < len_1D; x++ {
		// write header profile, level = -1
		key := fmt.Sprintf("DATE:%d", int(profile[x]))
		fmt.Fprintf(ascii, "%4d %4d %f %f %f %s", int(profile[x]), -1, time[x], lat[x], lon[x], date[key])
		fmt.Fprintf(hdr_id, "%4d %4d %f %f %f %s\n", int(profile[x]), -1, time[x], lat[x], lon[x], date[key])
		// fill last header columns with 1e36
		for i := 0; i < len(hdr)-6; i++ {
			fmt.Fprintf(ascii, " %g", 1e36)
		}
		fmt.Fprintln(ascii) // add newline
		// loop over each level
		for y := 0; y < len_2D; y++ {
			fmt.Fprintf(ascii, "%4d ", int(profile[x]))
			// loop over each physical parameter (key)
			for _, key := range hdr {
				// if key not in map, goto next key
				if _, ok := nc.Variables_2D[key]; ok {
					// fill 2D slice
					data := nc.Variables_2D[key].data[x][y]
					// print data with it's format, change format for FillValue
					if data == 1e36 {
						fmt.Fprintf(ascii, "%g ", data)
					} else {
						fmt.Fprintf(ascii, map_format[key]+" ", data)
					}
				}
			}
			fmt.Fprintf(ascii, "\n")
			// goto next profile when max depth reach
			if nc.Variables_2D["DEPTH"].data[x][y] >= nc.Extras_f[fmt.Sprintf("DEPTH:%d", int(profile[x]))] {
				break
			}
		}
		ascii.Flush()
		hdr_id.Flush()
	}
}
