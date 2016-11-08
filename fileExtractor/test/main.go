package main

import (
	"fmt"

	"github.com/jgrelet/oceano2oceansites/fileExtractor"
)

// usefull macro
var p = fmt.Println
var pf = fmt.Printf

// example main
func main() {

	// CTD
	opts := fileExtractor.NewFileExtractOptions()
	opts.SetVarsList("PRES,3,DEPTH,4,ETDD,2,TEMP,5,PSAL,17,DENS,20,SVEL,22,DOX2,15,FLU2,13,TUR3,14,NAVG,23")
	//opts.SetSkipLine(354)
	opts.SetHdrDelimiter("*END*")
	// print options
	p(opts)
	prf := fileExtractor.NewProfilesExtractor("ctd/dfr2600?.cnv", *fileExtractor.NewFileExtractor(opts))
	prf.Read()
	p(prf)

	//TSG
	opts = fileExtractor.NewFileExtractOptions()
	opts.SetVarsList("DATE,3,TIME,4,LATITUDE,6,LONGITUDE,9,SSTP,19,SSJT,20,COND,21,SSPS,22")
	opts.SetSkipLine(6)
	opts.SetSeparator(",")
	p(opts)
	traj := fileExtractor.NewTrajectoriesExtractor("tsg/*-TS_COLCOR.COLCOR", *fileExtractor.NewFileExtractor(opts))
	traj.Read()
	p(traj)

}
