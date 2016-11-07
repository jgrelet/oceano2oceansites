package main

import (
	"fmt"
	_ "log"
	_ "os"

	"github.com/jgrelet/oceano2oceansites/fileExtractor"
)

// usefull macro
var p = fmt.Println
var pf = fmt.Printf

// example main
func main() {

	// initialize options
	//	opts := fileExtractor.NewFileExtractOptions().SetFilename("test.gps")
	//	opts.SetVarsList("TIME,1,LATITUDE,2,LONGITUDE,3,TEMP,4")
	//	opts.SetSkipLine(2)

	// pirata-FR23
	opts := fileExtractor.NewFileExtractOptions("ctd/dfr26001.cnv")
	opts.SetVarsList("PRES,3,DEPTH,4,ETDD,2,TEMP,5,PSAL,17,DENS,20,SVEL,22,DOX2,15,FLU2,13,TUR3,14,NAVG,23")
	//opts.SetSkipLine(354)
	opts.SetHdrDelimiter("*END*")

	// print options
	p(opts)

	// initialize fileExtractor from options
	//ext := fileExtractor.NewFileExtractor(opts)

	fe := fileExtractor.NewFilesExtractor("ctd/dfr2600?.cnv")
	p(fe)
	/*
		// read thes files
		_, err := ext.Read()
		if err != nil {
			log.Fatal(err)
			//os.Exit(1)
		}

		// display the value
		pres := ext.Data()["PRES"]
		temp := ext.Data()["TEMP"]
		psal := ext.Data()["PSAL"]
		for i := 0; i < ext.Size()[1]; i++ {
			p := pres[i]
			t := temp[i]
			s := psal[i]
			pf("%f\t%f\t%f\n", p, t, s)
		}
	*/
}
