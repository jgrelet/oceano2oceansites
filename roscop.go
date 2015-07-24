// roscop.go
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type RoscopAttribute struct {
	long_name  string
	units      string
	valid_min  float64
	valid_max  float64
	format     string
	_FillValue float64
}

// documentation for csv is at http://golang.org/pkg/encoding/csv/
// TODO: could not find
func codeRoscopFromCsv(filename string) map[string]RoscopAttribute {

	var roscop = make(map[string]RoscopAttribute)

	file, err := os.Open(filename)
	if err != nil {
		// err is printable
		// elements passed are separated by space automatically
		fmt.Println("function codeRoscopFromCsv error")
		fmt.Printf("Please, check location for %s file\n", filename)
		//os.Exit(1)
		log.Fatal(err)
	}
	// automatically call Close() at the end of current method
	defer file.Close()
	//
	reader := csv.NewReader(file)
	// options are available at:
	// http://golang.org/src/pkg/encoding/csv/reader.go?s=3213:3671#L94
	reader.Comma = ';'

	for {
		r := RoscopAttribute{}

		// read just one record, but we could ReadAll() as well
		record, err := reader.Read()
		// end-of-file is fitted into err
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)

		}
		// record is an array of string so is directly printable
		//fmt.Println("Record ", record, "and has", len(record), "fields")
		// and we can iterate on top of that
		r.long_name = record[1]
		r.units = record[2]
		if v, err := strconv.ParseFloat(record[3], 64); err == nil {
			r.valid_min = v
		}
		if v, err := strconv.ParseFloat(record[4], 64); err == nil {
			r.valid_max = v
		}
		r.format = record[5]
		if v, err := strconv.ParseFloat(record[6], 64); err == nil {
			r._FillValue = v
		}
		roscop[record[0]] = r

	}
	return roscop
}
