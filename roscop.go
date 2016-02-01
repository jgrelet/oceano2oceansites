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
	types      string
}

// documentation for csv is at http://golang.org/pkg/encoding/csv/
// TODO: could not find
func codeRoscopFromCsv(filename string) map[string]RoscopAttribute {

	// if ROSCOP env is define, read file
	if code_roscop != "" {
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
		//f("%#v", roscop)
		return roscop
	} else {
		// roscop map is hard define here
		roscop := map[string]RoscopAttribute{
			"CUPW":      {long_name: "14C PRODUCTION UNKNOWN FILTER", units: "milligram carbon/(m3.day)", valid_min: 0, valid_max: 200, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"C1UW":      {long_name: "14C UPTAKE 0.2-1 MICRON", units: "milligram carbon/(m3.day)", valid_min: 0, valid_max: 100, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"BFUP":      {long_name: "19'BUTANOYLOXYFUCOXANTHINE", units: "milligram/m3", valid_min: 0, valid_max: 5, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"DENS":      {long_name: "DENSITY (Sigma-theta)", units: "kg/m3", valid_min: 10, valid_max: 35, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"TH4W":      {long_name: "DISSOLVED 234TH", units: "Bq/m3", valid_min: 0, valid_max: 90, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"DRYT":      {long_name: "DRY BULB TEMPERATURE", units: "Celsius degree", valid_min: 0, valid_max: 90, format: "%5.1f", _FillValue: 1e+36, types: "float32"},
			"ETDD":      {long_name: "ELLAPSED TIME", units: "1", valid_min: 0, valid_max: 999, format: "%9.5f", _FillValue: 1e+36, types: "float64"},
			"MPBS":      {long_name: "Pb IN THE SEDIMENT", units: "ppm", valid_min: 0, valid_max: 900, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"REDS":      {long_name: "REDOX POTENTIAL", units: "millivolt", valid_min: -110, valid_max: 200, format: "%+4.0f", _FillValue: 1e+36, types: "float32"},
			"MSIS":      {long_name: "Si IN THE SEDIMENT", units: "%", valid_min: 0, valid_max: 99, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"MBAP":      {long_name: "Ba IN SUSPENDED MATTER", units: "milligram/m3", valid_min: 0, valid_max: 99.9999, format: "%7.4f", _FillValue: 1e+36, types: "float32"},
			"WSPN":      {long_name: "WIND SPEED NORTHWARD COMPONENT", units: "meter/second", valid_min: 0, valid_max: 100, format: "%7.3f", _FillValue: 1e+36, types: "float32"},
			"NETR":      {long_name: "NET RADIATION", units: "watt/m2", valid_min: -500, valid_max: 500, format: "%+5.1f", _FillValue: 1e+36, types: "float32"},
			"NT1P":      {long_name: "TOTAL PARTICULATE NITROGEN", units: "micromole/kg", valid_min: 0, valid_max: 10, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"LCAW":      {long_name: "LIGHT CARBON ABSORPTION", units: "milligram/m3", valid_min: 0, valid_max: 40, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"CHLT":      {long_name: "CHLOROPHYLL-TOTAL", units: "milligram/m3", valid_min: 0, valid_max: 99, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"NSCS":      {long_name: "CURRENT NORTH STD. DEVIATION", units: "meter/second", valid_min: 0, valid_max: 20, format: "%+7.3f", _FillValue: 1e+36, types: "float32"},
			"PTNP":      {long_name: "TOTAL PARTICULATE NITROGEN", units: "milligram/m3", valid_min: 0, valid_max: 5, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"PSA1":      {long_name: "PRACTICAL SALINITY PRIMARY SENSOR", units: "P.S.S.78", valid_min: 33, valid_max: 37, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"TYPE":      {long_name: "string", units: "string", valid_min: 0, valid_max: 0, format: "string", _FillValue: 0, types: "float32"},
			"MCAP":      {long_name: "Ca IN SUSPENDED MATTER", units: "milligram/m3", valid_min: 0, valid_max: 99.999, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"FLU1":      {long_name: "FLUORESCENCE", units: "volt", valid_min: -1, valid_max: 3, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"PC1P":      {long_name: "PARTICULATE ORGANIC CARBON/POC", units: "millimole/m3", valid_min: 0, valid_max: 100, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"ATMS":      {long_name: "ATMOSPHERIC PRESSURE - SEA LEV", units: "hectopascal", valid_min: 1000, valid_max: 1030, format: "%8.3f", _FillValue: 1e+36, types: "float32"},
			"MSZW":      {long_name: "MESOZOOPLANCTON DRY WEIGHT", units: "milligram/m3", valid_min: 0, valid_max: 90, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"SIOS":      {long_name: "SEDIMENT BIOGENIC SiO2", units: "%", valid_min: 0, valid_max: 100, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"SILT":      {long_name: "SILT IN THE SEDIMENT", units: "%", valid_min: 0, valid_max: 100, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"NOTT":      {long_name: "TOTAL ORGANIC NITROGEN (D+P)", units: "micromole/kg", valid_min: 0, valid_max: 10, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"NEOW":      {long_name: "DISSOLVED NEON", units: "nanomole/kg", valid_min: 4, valid_max: 8, format: "%7.4f", _FillValue: 1e+36, types: "float32"},
			"XCO2":      {long_name: "CO2 MOLE FRACTION IN DRY GAS", units: "ppm", valid_min: 0, valid_max: 5000, format: "%8.3f", _FillValue: 1e+36, types: "float32"},
			"CO2W":      {long_name: "DISSOLVED CARBON DIOXYD (CO2)", units: "millimole/m3", valid_min: 0, valid_max: 4800, format: "%6.1f", _FillValue: 1e+36, types: "float32"},
			"TPHS":      {long_name: "TOTAL PHOSPHORUS (P) CONTENT", units: "millimole/m3", valid_min: 0, valid_max: 10, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"NHRW":      {long_name: "AMMONIUM REGENERATION", units: "micromole nitrogen/(m3.day)", valid_min: 0, valid_max: 99, format: "%5.0f", _FillValue: 1e+36, types: "float32"},
			"OSAT":      {long_name: "OXYGEN SATURATION", units: "%", valid_min: 0, valid_max: 10, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"MPPP":      {long_name: "P IN SUSPENDED MATTER", units: "milligram/m3", valid_min: 0, valid_max: 99.999, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"TE10":      {long_name: "SEA TEMPERATURE", units: "Celsius degree", valid_min: -2, valid_max: 32, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"GPMP":      {long_name: "TOTAL SUSP. PART. MATTER/GLASS", units: "gram/m3", valid_min: 0, valid_max: 99, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"TUR4":      {long_name: "TURBIDITY", units: "N.T.U Nephelo Turb. Unit", valid_min: 0, valid_max: 100, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"DIPF":      {long_name: "DIATOMS PPC FLUX", units: "milligram C/(m2.day)", valid_min: 0, valid_max: 99.999, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"DPSF":      {long_name: "DEPTH BELOW SEA FLOOR", units: "meter", valid_min: 0, valid_max: 11000, format: "%6.1f", _FillValue: 1e+36, types: "float32"},
			"LTHF":      {long_name: "LITHOGENIC FRACTION FLUX", units: "milligram/(m2.day)", valid_min: 0, valid_max: 3400, format: "%7.2f", _FillValue: 1e+36, types: "float32"},
			"POCP":      {long_name: "PARTICULATE ORGANIC CARBON/POC", units: "milligram/m3", valid_min: 0, valid_max: 999, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"HEAD":      {long_name: "PLATFORM HEADING REL. NORTH", units: "degree", valid_min: -360, valid_max: 360, format: "%+5.1f", _FillValue: 1e+36, types: "float32"},
			"MZNS":      {long_name: "Zn IN THE SEDIMENT", units: "ppm", valid_min: 0, valid_max: 900, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"ALKW":      {long_name: "ALKALINITY", units: "micromole/kg", valid_min: 0, valid_max: 9000, format: "%6.1f", _FillValue: 1e+36, types: "float32"},
			"TPAW":      {long_name: "DISS. TRIPHOSPHATE ADENOSINE", units: "milligram/m3", valid_min: 0, valid_max: 1, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"CHAD":      {long_name: "DIVINYL CHLOROPHYLL-A", units: "milligram/m3", valid_min: 0, valid_max: 9, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"FLUO":      {long_name: "FLUORESCENCE", units: "relative unit", valid_min: -0.1, valid_max: 10, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"OXIR":      {long_name: "ISOTOPIC RATIO O18/O16", units: "per thousand", valid_min: -10, valid_max: 10, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"PHOW":      {long_name: "PHOSPHATE (PO4-P) CONTENT", units: "micromole/kg", valid_min: 0, valid_max: 10, format: "%5.3f", _FillValue: 1e+36, types: "float32"},
			"TCO2":      {long_name: "TOTAL CARBON DIOXYD (CO2)", units: "mole/m3", valid_min: 0, valid_max: 5000, format: "%8.3f", _FillValue: 1e+36, types: "float32"},
			"HCDT":      {long_name: "DIRECTION REL. TRUE NORTH", units: "degree", valid_min: 0, valid_max: 360, format: "%5.1f", _FillValue: 1e+36, types: "float32"},
			"CO3S":      {long_name: "SEDIMENT CARBONATES", units: "%", valid_min: 0, valid_max: 100, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"TCCS":      {long_name: "SEDIMENT TOTAL CARBON", units: "%", valid_min: 0, valid_max: 100, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"PTPP":      {long_name: "TOTAL PARTICULATE PHOSPHORUS", units: "milligram/m3", valid_min: 0, valid_max: 5, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"TUR6":      {long_name: "TURBIDITY", units: "milliF.T.U Formaz Turb Unit", valid_min: 0, valid_max: 5000, format: "%6.1f", _FillValue: 1e+36, types: "float32"},
			"MTHS":      {long_name: "Th IN THE SEDIMENT", units: "ppm", valid_min: 0, valid_max: 900, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"DEN2":      {long_name: "DENSITY (Sigma-theta) SECONDARY SENSORS", units: "kg/m3", valid_min: 10, valid_max: 35, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"CPH1":      {long_name: "CHLOROPHYLL-A TOTAL", units: "milligram/m3", valid_min: 0, valid_max: 99, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"LSIC":      {long_name: "LITHOGENIC CONTENT", units: "%", valid_min: 0, valid_max: 99.99, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"OMMF":      {long_name: "PART. ORGANIC MATTER FLUX", units: "milligram/(m2.day)", valid_min: 0, valid_max: 600, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"SWDR":      {long_name: "SWELL DIRECTION  REL TRUE N.", units: "degree", valid_min: 0, valid_max: 360, format: "%5.1f", _FillValue: 1e+36, types: "float32"},
			"MALF":      {long_name: "AL FLUX IN SETTLING PARTICLES", units: "milligram/(m2.day)", valid_min: 0, valid_max: 180, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"ETHW":      {long_name: "DISSOLVED 234 TH ACT. ERROR", units: "Bq/m3", valid_min: 0, valid_max: 90, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"CORG":      {long_name: "DISSOLVED ORGANIC CARBON", units: "millimole/m3", valid_min: 0, valid_max: 1000, format: "%5.0f", _FillValue: 1e+36, types: "float32"},
			"CNDC":      {long_name: "ELECTRICAL CONDUCTIVITY", units: "mho/meter", valid_min: 3, valid_max: 7, format: "%5.3f", _FillValue: 1e+36, types: "float32"},
			"PHTF":      {long_name: "PHAEOPIGMENTS VERTICAL FLUX", units: "milligram/(m2.day)", valid_min: 0, valid_max: 10, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"VTDH":      {long_name: "SIGNIFICANT WAVE HEIGHT", units: "meter", valid_min: 0, valid_max: 99.999, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"TOCW":      {long_name: "TOTAL ORGANIC CARBON", units: "millimole/m3", valid_min: 0, valid_max: 999.999, format: "%7.3f", _FillValue: 1e+36, types: "float32"},
			"PHTP":      {long_name: "TOTAL PHEOPHYTINE", units: "milligram/m3", valid_min: 0, valid_max: 99, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"VTZA":      {long_name: "AVER ZERO CROSSING WAVE PERIOD", units: "second", valid_min: 0, valid_max: 500, format: "%3.0f", _FillValue: 1e+36, types: "float32"},
			"DVT1":      {long_name: "DISSOLVED OXYGEN PRIMARY SENSOR dV/dt", units: "dv/dt", valid_min: -1, valid_max: 1, format: "%+7.5f", _FillValue: 1e+36, types: "float32"},
			"VERT":      {long_name: "VERTICAL DISPLACEMENT", units: "meter", valid_min: 0, valid_max: 9999, format: "%8.3f", _FillValue: 1e+36, types: "float32"},
			"FCO2":      {long_name: "CO2 FUGACITY", units: "microatmosphere", valid_min: 0, valid_max: 1000, format: "%6.1f", _FillValue: 1e+36, types: "float32"},
			"RDIN":      {long_name: "INCIDENT RADIATION", units: "watt/m2", valid_min: -500, valid_max: 500, format: "%+5.1f", _FillValue: 1e+36, types: "float32"},
			"IODI":      {long_name: "IODINE", units: "millimole/m3", valid_min: 0, valid_max: 10, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"NTRI":      {long_name: "NITRITE (NO2-N) CONTENT", units: "millimole/m3", valid_min: 0, valid_max: 10, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"TICW":      {long_name: "TOTAL INORGANIC CARBON", units: "micromole/kg", valid_min: 0, valid_max: 9000, format: "%7.2f", _FillValue: 1e+36, types: "float32"},
			"TUR5":      {long_name: "TURBIDITY", units: "relative unit", valid_min: 0, valid_max: 10, format: "%7.4f", _FillValue: 1e+36, types: "float32"},
			"NAVG":      {long_name: "AVERAGED DATA CYCLE NUMBER", units: "1", valid_min: 0, valid_max: 999, format: "%3.0f", _FillValue: 1e+36, types: "float32"},
			"BCCS":      {long_name: "BACTERIA NUMBER SEDIMENT", units: "10+9 cell/dm3", valid_min: 0, valid_max: 9999.9, format: "%6.1f", _FillValue: 1e+36, types: "float32"},
			"LGH3":      {long_name: "LIGHT IRRADIANCE CORRECTED PAR", units: "micromole photon/(m2.s)", valid_min: 0, valid_max: 3000, format: "%8.3f", _FillValue: 1e+36, types: "float32"},
			"CICW":      {long_name: "NUMBER OF SW CILIATES", units: "10+3 cell/m3", valid_min: 0, valid_max: 90000, format: "%5.0f", _FillValue: 1e+36, types: "float32"},
			"TE07":      {long_name: "SEA TEMPERATURE", units: "Celsius degree", valid_min: -2, valid_max: 32, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"LIPS":      {long_name: "SEDIMENT LIPIDS", units: "microgram/g", valid_min: 0, valid_max: 5000, format: "%4.0f", _FillValue: 1e+36, types: "float32"},
			"WMSP":      {long_name: "WIND SPEED - MAX AVER PER 2 MN", units: "meter/second", valid_min: 0, valid_max: 200, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"CND1":      {long_name: "ELECTRICAL CONDUCTIVITY PRIMARY SENSOR", units: "mho/meter", valid_min: 3, valid_max: 7, format: "%5.3f", _FillValue: 1e+36, types: "float32"},
			"MFES":      {long_name: "Fe IN THE SEDIMENT", units: "%", valid_min: 0, valid_max: 99, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"AMIS":      {long_name: "SEDIMENT AMINO-ACIDS", units: "microgram/g", valid_min: 0, valid_max: 7000, format: "%4.0f", _FillValue: 1e+36, types: "float32"},
			"WDIR":      {long_name: "WIND DIRECTION REL. TRUE NORTH", units: "degree", valid_min: 0, valid_max: 360, format: "%+5.1f", _FillValue: 1e+36, types: "float32"},
			"ASDW":      {long_name: "ABSORPTION STANDARD DEVIATION", units: "milligram/m3", valid_min: 0, valid_max: 10, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"MMNP":      {long_name: "Mn IN SUSPENDED MATTER", units: "milligram/m3", valid_min: 0, valid_max: 99.999, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"CE1W":      {long_name: "DISSOLVED CFC11 ERROR", units: "picomole/kg", valid_min: 0, valid_max: 0.5, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"MSDP":      {long_name: "MEAN SPHERIC DIAM. OF PARTICLE", units: "millimeter", valid_min: 0, valid_max: 99, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"MNTH":      {long_name: "MONTH", units: "mm", valid_min: 1, valid_max: 12, format: "%2.2d", _FillValue: 1e+36, types: "float32"},
			"PCEW":      {long_name: "NUMBER OF SW PICOEUCARYOTES", units: "10+6 cell/m3", valid_min: 0, valid_max: 90000, format: "%5.0f", _FillValue: 1e+36, types: "float32"},
			"SSPS":      {long_name: "SEA SURFACE PRACTICAL SALINITY", units: "1", valid_min: 0, valid_max: 40, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"MZNF":      {long_name: "ZN FLUX IN SETTLING PARTICLES", units: "microgram/(m2.day)", valid_min: 0, valid_max: 7600, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"DOV1":      {long_name: "DISSOLVED OXYGEN PRIMARY SENSOR VOLTAGE", units: "V", valid_min: 0, valid_max: 10, format: "%6.4f", _FillValue: 1e+36, types: "float32"},
			"CODW":      {long_name: "CHEMICAL OXYGEN DEMAND", units: "millimole/m3", valid_min: 0, valid_max: 650, format: "%7.3f", _FillValue: 1e+36, types: "float32"},
			"NFAW":      {long_name: "AUTOTROPHIC NANOFLAGELLATES", units: "10+3 cell/m3", valid_min: 0, valid_max: 900000, format: "%6.0f", _FillValue: 1e+36, types: "float32"},
			"MALP":      {long_name: "Al IN SUSPENDED MATTER", units: "milligram/m3", valid_min: 0, valid_max: 99.999, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"C13D":      {long_name: "DELTA 13C (13C/12C)", units: "per thousand", valid_min: -50, valid_max: 100, format: "%5.1f", _FillValue: 1e+36, types: "float32"},
			"CHAE":      {long_name: "EPIMERE CHLOROPHYLL-A", units: "milligram/m3", valid_min: 0, valid_max: 5, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"PHCW":      {long_name: "NUM. OF SW PHYTOPLANKTON CELLS", units: "10+3 cell/m3", valid_min: 0, valid_max: 900000, format: "%6.0f", _FillValue: 1e+36, types: "float32"},
			"SACC":      {long_name: "SALINITY ACCURACY", units: ".", valid_min: 0, valid_max: 4, format: "%2d", _FillValue: 1e+36, types: "float32"},
			"PHNS":      {long_name: "SEDIMENT PHENOLS", units: "microgram/g", valid_min: 0, valid_max: 200, format: "%5.1f", _FillValue: 1e+36, types: "float32"},
			"TNNS":      {long_name: "SEDIMENT TOTAL NITROGEN", units: "%", valid_min: 0, valid_max: 100, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"MCES":      {long_name: "Ce IN THE SEDIMENT", units: "ppm", valid_min: 0, valid_max: 900, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"DATE":      {long_name: "DATE", units: "mmdd", valid_min: 101, valid_max: 1231, format: "%4.4d", _FillValue: 1e+36, types: "float32"},
			"DOX2":      {long_name: "DISSOLVED OXYGEN", units: "micromole/kg", valid_min: 0, valid_max: 450, format: "%7.3f", _FillValue: 1e+36, types: "float32"},
			"NFCW":      {long_name: "NUMBER OF SW NANOFLAGELLATES", units: "10+3 cell/m3", valid_min: 0, valid_max: 900000, format: "%6.0f", _FillValue: 1e+36, types: "float32"},
			"MNIS":      {long_name: "Ni IN THE SEDIMENT", units: "ppm", valid_min: 0, valid_max: 900, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"MBAS":      {long_name: "Ba IN THE SEDIMENT", units: "ppm", valid_min: 0, valid_max: 5000, format: "%7.2f", _FillValue: 1e+36, types: "float32"},
			"POCF":      {long_name: "PART. ORGANIC CARBON FLUX", units: "milligram/(m2.day)", valid_min: 0, valid_max: 99, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"PN1P":      {long_name: "PARTICULATE ORGANIC NITROGEN", units: "millimole/m3", valid_min: 0, valid_max: 2, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"PERP":      {long_name: "PERIDININE", units: "milligram/m3", valid_min: 0, valid_max: 5, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"TE05":      {long_name: "SEA TEMPERATURE", units: "Celcius degree", valid_min: -2, valid_max: 32, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"CPH2":      {long_name: "CHLOROPHYLL-A TOTAL", units: "milligram/m3", valid_min: 0, valid_max: 99, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"PRCW":      {long_name: "NUMBER OF SW PROCHLOROCOCCUS", units: "10+6 cell/m3", valid_min: 0, valid_max: 900000, format: "%6.0f", _FillValue: 1e+36, types: "float32"},
			"PHEC":      {long_name: "PHEOPHYTIN-C", units: "milligram/m3", valid_min: 0, valid_max: 99, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"HEDW":      {long_name: "HELIUM DEV. OF ISOTOPIC RATIO", units: "%", valid_min: -1.5, valid_max: 100, format: "%8.4f", _FillValue: 1e+36, types: "float32"},
			"CPHL":      {long_name: "CHLOROPHYLL-A TOTAL", units: "milligram/m3", valid_min: 0, valid_max: 99, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"HBAW":      {long_name: "NB OF HETEROTROPHIC BACTERIA", units: "10+6 cell/m3", valid_min: 0, valid_max: 900000, format: "%6.0f", _FillValue: 1e+36, types: "float32"},
			"BCCW":      {long_name: "NUMBER OF SW BACTERIA", units: "10+9 cell/m3", valid_min: 0, valid_max: 9000, format: "%7.2f", _FillValue: 1e+36, types: "float32"},
			"SCDT":      {long_name: "SEA SURF CURRENT DIR. REL T. N", units: "degree", valid_min: 0, valid_max: 360, format: "%5.1f", _FillValue: 1e+36, types: "float32"},
			"SLCW":      {long_name: "SILICATE (SIO4-SI) CONTENT", units: "micromole/kg", valid_min: 0, valid_max: 195, format: "%5.1f", _FillValue: 1e+36, types: "float32"},
			"UREA":      {long_name: "UREA", units: "millimole/m3", valid_min: 0, valid_max: 5, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"ATMP":      {long_name: "ATMOSPHERIC PRESSURE", units: "hectopascal", valid_min: 1000, valid_max: 1030, format: "%8.3f", _FillValue: 1e+36, types: "float32"},
			"TUR2":      {long_name: "LIGHT ATTENUATION COEFFICIENT", units: "m-1", valid_min: 0, valid_max: 10, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"PP1P":      {long_name: "PART. ORGANIC PHOSPHORUS (P)", units: "millimole/m3", valid_min: 0, valid_max: 10, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"TSMP":      {long_name: "TOTAL SUSPENDED MATTER", units: "gram/m3", valid_min: 0, valid_max: 100, format: "%7.3f", _FillValue: 1e+36, types: "float32"},
			"TINW":      {long_name: "DISSOLVED INORGANIC NITROGEN", units: "millimole/m3", valid_min: 0, valid_max: 100, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"TCCF":      {long_name: "PART. TOTAL CARBON FLUX", units: "milligram/(m2.day)", valid_min: 0, valid_max: 999, format: "%7.3f", _FillValue: 1e+36, types: "float32"},
			"FLPF":      {long_name: "FLAGELLATES PPC FLUX", units: "milligram C/(m2.day)", valid_min: 0, valid_max: 99.999, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"CHC3":      {long_name: "CHLOROPHYLL-C3", units: "milligram/m3", valid_min: 0, valid_max: 5, format: "%6.3f", _FillValue: -99.9999, types: "float32"},
			"NTIW":      {long_name: "NITRITE (NO2-N) CONTENT", units: "micromole/kg", valid_min: 0, valid_max: 100, format: "%7.3f", _FillValue: 1e+36, types: "float32"},
			"SLCP":      {long_name: "PARTICULATE ORGANIC SILICA(SI)", units: "millimole/m3", valid_min: 0, valid_max: 1, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"TIME":      {long_name: "time of measurement", units: "days since 1950-01-01T00:00:00Z", valid_min: 0, valid_max: 90000, format: "%6.6d", _FillValue: -1e+36, types: "float64"},
			"BATH":      {long_name: "BATHYMETRIC DEPTH", units: "meter", valid_min: 0, valid_max: 11000, format: "%6.1f", _FillValue: 1e+36, types: "float32"},
			"EWCT":      {long_name: "CURRENT EAST COMPONENT", units: "meter/second", valid_min: -100, valid_max: 100, format: "%+7.3f", _FillValue: 1e+36, types: "float32"},
			"TSMF":      {long_name: "TOTAL MASS FLUX", units: "milligram/(m2.day)", valid_min: 0, valid_max: 9999.9, format: "%6.1f", _FillValue: 1e+36, types: "float32"},
			"CH1T":      {long_name: "CHLOROPHYLL TOTAL", units: "microgram/kg", valid_min: 0, valid_max: 10, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"EF2W":      {long_name: "DISSOLVED CFC12 ERROR", units: "%", valid_min: 0, valid_max: 100, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"NEEW":      {long_name: "DISSOLVED NEON ERROR", units: "nanomole/kg", valid_min: 0, valid_max: 1, format: "%7.4f", _FillValue: 1e+36, types: "float32"},
			"NORG":      {long_name: "DISSOLVED ORGANIC NITROGEN", units: "millimole/m3", valid_min: 0, valid_max: 50, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"VZMX":      {long_name: "MAXI ZERO CROSSING WAVE HEIGHT", units: "metre", valid_min: 0, valid_max: 10, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"DTCW":      {long_name: "NUMBER OF SW DIATOMS", units: "10+3 cell/m3", valid_min: 0, valid_max: 900000, format: "%6.0f", _FillValue: 1e+36, types: "float32"},
			"OPAP":      {long_name: "OPAL CONTENT", units: "%", valid_min: 0, valid_max: 99.99, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"NTOT":      {long_name: "TOTAL NITROGEN (N) CONTENT", units: "millimole/m3", valid_min: 0, valid_max: 90, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"EWCS":      {long_name: "CURRENT EAST  STD. DEVIATION", units: "meter/second", valid_min: 0, valid_max: 20, format: "%+6.3f", _FillValue: 1e+36, types: "float32"},
			"DVT2":      {long_name: "DISSOLVED OXYGEN SECONDARY SENSOR dV/dt", units: "dv/dt", valid_min: -1, valid_max: 1, format: "%+7.5f", _FillValue: 1e+36, types: "float32"},
			"LONX":      {long_name: "LONGITUDE", units: "degree_east", valid_min: -180, valid_max: 180, format: "%+9.4f", _FillValue: 1e+36, types: "float32"},
			"SAND":      {long_name: "SAND IN THE SEDIMENT", units: "%", valid_min: 0, valid_max: 100, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"ZXAP":      {long_name: "ZEAXANTHINE", units: "milligram/m3", valid_min: 0, valid_max: 5, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"DOV2":      {long_name: "DISSOLVED OXYGEN SECONDARY SENSOR VOLTAGE", units: "V", valid_min: 0, valid_max: 10, format: "%6.4f", _FillValue: 1e+36, types: "float32"},
			"EF1W":      {long_name: "DISSOLVED CFC11 ERROR", units: "%", valid_min: 0, valid_max: 100, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"TXAP":      {long_name: "DIATOXANTHINE", units: "milligram/m3", valid_min: 0, valid_max: 5, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"WETT":      {long_name: "WET BULB TEMPERATURE", units: "Celsius degree", valid_min: 0, valid_max: 90, format: "%5.1f", _FillValue: 1e+36, types: "float32"},
			"CLAY":      {long_name: "CLAY IN THE SEDIMENT", units: "%", valid_min: 0, valid_max: 100, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"HEEW":      {long_name: "DISSOLVED HELIUM ERROR", units: "nanomole/kg", valid_min: 0, valid_max: 1, format: "%7.4f", _FillValue: 1e+36, types: "float32"},
			"WSPD":      {long_name: "HORIZONTAL WIND SPEED", units: "meter/second", valid_min: 0, valid_max: 300, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"LGH5":      {long_name: "IMMERGED/SURF IRRADIANCE RATIO", units: "%", valid_min: 0, valid_max: 100, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"TOCS":      {long_name: "SEDIMENT TOTAL ORGANIC CARBON", units: "%", valid_min: 0, valid_max: 100, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"TPHP":      {long_name: "TOTAL PHAEOPIGMENTS", units: "milligram/m3", valid_min: 0, valid_max: 100, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"TE02":      {long_name: "SEA TEMPERATURE SECONDARY SENSOR", units: "Celsius degree", valid_min: -2, valid_max: 32, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"DEN1":      {long_name: "DENSITY (Sigma-theta) PRIMARY SENSORS", units: "kg/m3", valid_min: 10, valid_max: 35, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"MMGP":      {long_name: "Mg IN SUSPENDED MATTER", units: "milligram/m3", valid_min: 0, valid_max: 99.999, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"COPP":      {long_name: "NUMBER OF COPEPODS", units: "number/m3", valid_min: 0, valid_max: 99999, format: "%5.0f", _FillValue: 1e+36, types: "float32"},
			"PONP":      {long_name: "PARTICULATE ORGANIC NITROGEN", units: "milligram/m3", valid_min: 0, valid_max: 100, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"MSCS":      {long_name: "Sc IN THE SEDIMENT", units: "ppm", valid_min: 0, valid_max: 900, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"CO3P":      {long_name: "CARBONATES CONTENT", units: "%", valid_min: 0, valid_max: 99.99, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"TE12":      {long_name: "SEA TEMPERATURE", units: "Celsius degree", valid_min: -2, valid_max: 32, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"SWHT":      {long_name: "SWELL HEIGHT", units: "meter", valid_min: 0, valid_max: 30, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"TDCW":      {long_name: "TOTAL DISSOLVED CARBON", units: "millimole/m3", valid_min: 0, valid_max: 9999, format: "%4.0f", _FillValue: 1e+36, types: "float32"},
			"MMNS":      {long_name: "Mn IN THE SEDIMENT", units: "%", valid_min: 0, valid_max: 99, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"CF3W":      {long_name: "DISSOLVED CFC113", units: "picomole/kg", valid_min: 0, valid_max: 5, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"CE2W":      {long_name: "DISSOLVED CFC12 ERROR", units: "picomole/kg", valid_min: 0, valid_max: 0.5, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"HCSP":      {long_name: "HORIZONTAL CURRENT SPEED", units: "meter/second", valid_min: 0, valid_max: 9, format: "%5.3f", _FillValue: 1e+36, types: "float32"},
			"MKKS":      {long_name: "K IN THE SEDIMENT", units: "%", valid_min: 0, valid_max: 99, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"TUR1":      {long_name: "LIGHT DIFFUSION COEFFICIENT", units: "m-1", valid_min: 0, valid_max: 10, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"NTZW":      {long_name: "NITRATE + NITRITE CONTENT", units: "micromole/kg", valid_min: 0, valid_max: 100, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"MPPS":      {long_name: "P IN THE SEDIMENT", units: "%", valid_min: 0, valid_max: 99, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"ABCP":      {long_name: "ALPHA BETA CAROTENES", units: "milligram/m3", valid_min: 0, valid_max: 5, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"TH4P":      {long_name: "PARTICULATE 234TH ACTIVTY", units: "Bq/m3", valid_min: 0, valid_max: 99, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"TDPW":      {long_name: "TOTAL DISSOLVED PHOSPHORUS", units: "millimole/m3", valid_min: 0, valid_max: 40, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"CO3F":      {long_name: "PART. CaCO3 FLUX", units: "milligram/(m2.day)", valid_min: 0, valid_max: 600, format: "%7.3f", _FillValue: 1e+36, types: "float32"},
			"DAYS":      {long_name: "DAY WITHIN YEAR", units: "decimal day", valid_min: 1, valid_max: 366, format: "%9.5f", _FillValue: 1e+36, types: "float32"},
			"MFEF":      {long_name: "FE FLUX IN SETTLING PARTICLES", units: "milligram/(m2.day)", valid_min: 0, valid_max: 100, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"NOUW":      {long_name: "NITRATE UPTAKE", units: "micromole nitrogen/(m3.day)", valid_min: 0, valid_max: 900, format: "%5.0f", _FillValue: 1e+36, types: "float32"},
			"CHTS":      {long_name: "SEDIMENT CARBOHYDRATES", units: "microgram/g", valid_min: 0, valid_max: 7000, format: "%4.0f", _FillValue: 1e+36, types: "float32"},
			"EPMP":      {long_name: "TOTAL SUSP. PART. MATTER/ESTER", units: "gram/m3", valid_min: 0, valid_max: 9, format: "%5.3f", _FillValue: 1e+36, types: "float32"},
			"AMON":      {long_name: "AMMONIUM (NH4-N) CONTENT", units: "millimole/m3", valid_min: 0, valid_max: 10, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"TE08":      {long_name: "SEA TEMPERATURE", units: "Celsius degree", valid_min: -2, valid_max: 32, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"BOTL":      {long_name: "SEA WATER SAMPLE BOTTLE NUMBER", units: "N/A", valid_min: 1, valid_max: 36, format: "%3d", _FillValue: 255, types: "uint8"},
			"TE03":      {long_name: "SEA TEMPERATURE", units: "Celsius degree", valid_min: -2, valid_max: 32, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"N15D":      {long_name: "DELTA 15N (15N/14N)", units: "per thousand", valid_min: -2, valid_max: 10, format: "%5.1f", _FillValue: 1e+36, types: "float32"},
			"MKKP":      {long_name: "K IN SUSPENDED MATTER", units: "milligram/m3", valid_min: 0, valid_max: 99.999, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"TUR3":      {long_name: "LIGHT TRANSMISSION", units: "%", valid_min: 0, valid_max: 100, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"PRFL":      {long_name: "PROFILE NUMBER", units: "N/A", valid_min: 1, valid_max: 99999, format: "%5.0f", _FillValue: 1e+36, types: "float32"},
			"CHLC":      {long_name: "CHLOROPHYLL-C TOTAL", units: "milligram/m3", valid_min: 0, valid_max: 99, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"PHEA":      {long_name: "PHEOPHYTIN-A", units: "milligram/m3", valid_min: 0, valid_max: 99, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"TNNF":      {long_name: "PART. TOTAL NITROGEN FLUX", units: "milligram/(m2.day)", valid_min: 0, valid_max: 99, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"DEPH":      {long_name: "DEPTH BELOW SEA SURFACE", units: "meter", valid_min: 0, valid_max: 6000, format: "%6.1f", _FillValue: 1e+36, types: "float32"},
			"MNAS":      {long_name: "Na IN THE SEDIMENT", units: "%", valid_min: 0, valid_max: 99, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"CPH3":      {long_name: "CHLOROPHYLL-A/2 MICRON FILTER", units: "milligram/m3", valid_min: 0, valid_max: 99, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"NTRZ":      {long_name: "NITRATE + NITRITE CONTENT", units: "millimole/m3", valid_min: 0, valid_max: 100, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"SNCW":      {long_name: "NUMBER OF SW SYNECHOCOCCUS", units: "10+6 cell/m3", valid_min: 0, valid_max: 900000, format: "%6.0f", _FillValue: 1e+36, types: "float32"},
			"PHEB":      {long_name: "PHEOPHYTIN-B", units: "milligram/m3", valid_min: 0, valid_max: 99, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"MSSP":      {long_name: "S IN SUSPENDED MATTER", units: "milligram/m3", valid_min: 0, valid_max: 99.999, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"TE11":      {long_name: "SEA TEMPERATURE", units: "Celsius degree", valid_min: -2, valid_max: 32, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"TE09":      {long_name: "SEA TEMPERATURE", units: "Celsius degree", valid_min: -2, valid_max: 32, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"SECS":      {long_name: "SECONDS WITHIN MINUTE", units: "ss", valid_min: 0, valid_max: 59, format: "%2.2d", _FillValue: 1e+36, types: "float32"},
			"CDFW":      {long_name: "DARK FIXATION", units: "milligram carbon/(m3.day)", valid_min: 0, valid_max: 2, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"H2OS":      {long_name: "WATER CONTENT", units: "%", valid_min: 0, valid_max: 100, format: "%5.1f", _FillValue: 1e+36, types: "float32"},
			"MCAS":      {long_name: "Ca IN THE SEDIMENT", units: "%", valid_min: 0, valid_max: 99, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"DAYX":      {long_name: "DAY WITHIN MONTH", units: "dd", valid_min: 1, valid_max: 31, format: "%2.2d", _FillValue: 1e+36, types: "float32"},
			"DPS1":      {long_name: "DEPTH BELOW BOTTOM-LOWER LIMIT", units: "meter", valid_min: 0, valid_max: 11000, format: "%6.1f", _FillValue: 1e+36, types: "float32"},
			"TOMP":      {long_name: "ORGANIC MATTER CONTENT", units: "%", valid_min: 0, valid_max: 99.99, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"RELH":      {long_name: "RELATIVE HUMIDITY", units: "%", valid_min: 0, valid_max: 100, format: "%5.1f", _FillValue: 1e+36, types: "float32"},
			"PSA2":      {long_name: "PRACTICAL SALINITY SECONDARY SENSOR", units: "P.S.S.78", valid_min: 33, valid_max: 37, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"CH1P":      {long_name: "CHL-A(LESS DIVINYLCHL-A)", units: "milligram/m3", valid_min: 0, valid_max: 3, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"MCLP":      {long_name: "Cl IN SUSPENDED MATTER", units: "milligram/m3", valid_min: 0, valid_max: 99.999, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"CF1W":      {long_name: "DISSOLVED CFC11", units: "picomole/kg", valid_min: -0.01, valid_max: 90, format: "%7.4f", _FillValue: 1e+36, types: "float32"},
			"LSCT":      {long_name: "LIGHT SCATTERING", units: "%", valid_min: 0, valid_max: 100, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"LIPI":      {long_name: "LIPIDS IN THE WATER COLUMN", units: "milligram/m3", valid_min: 0, valid_max: 200, format: "%5.0f", _FillValue: 1e+36, types: "float32"},
			"NTAW":      {long_name: "NITRATE (NO3-N) CONTENT", units: "micromole/kg", valid_min: 0, valid_max: 90, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"TEMP":      {long_name: "SEA TEMPERATURE", units: "Celsius degree", valid_min: 0, valid_max: 30, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"PT1P":      {long_name: "TOTAL PARTICULATE PHOSPHORUS", units: "micromole/kg", valid_min: 0, valid_max: 5, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"CH2P":      {long_name: "CHL-B(LESS DIVINYLCHL-B)", units: "milligram/m3", valid_min: 0, valid_max: 1, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"HSUL":      {long_name: "HYDROGEN SULPHIDE (H2S)", units: "millimole/m3", valid_min: 0, valid_max: 500, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"MIIS":      {long_name: "I IN THE SEDIMENT", units: "ppm", valid_min: 0, valid_max: 900, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"MMOS":      {long_name: "Mo IN THE SEDIMENT", units: "ppm", valid_min: 0, valid_max: 900, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"NTRA":      {long_name: "NITRATE (NO3-N) CONTENT", units: "millimole/m3", valid_min: 0, valid_max: 56, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"ETHP":      {long_name: "PARTICULATE 234 TH ACT. ERROR", units: "Bq/m3", valid_min: 0, valid_max: 99, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"MTIP":      {long_name: "Ti IN SUSPENDED MATTER", units: "milligram/m3", valid_min: 0, valid_max: 99.999, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"PREC":      {long_name: "CORRECTED SEA PRESSURE", units: "decibar=10000 pascals", valid_min: 0, valid_max: 6500, format: "%6.1f", _FillValue: 1e+36, types: "float32"},
			"EF3W":      {long_name: "DISSOLVED CFC113 ERROR", units: "%", valid_min: 0, valid_max: 100, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"DOX1":      {long_name: "DISSOLVED OXYGEN", units: "ml/l", valid_min: 0, valid_max: 10, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"HEDE":      {long_name: "HELIUM ISOTOPIC RATIO ERROR", units: "%", valid_min: 0, valid_max: 100, format: "%8.4f", _FillValue: 1e+36, types: "float32"},
			"SSTM":      {long_name: "MODEL SEA SURFACE TEMPERATURE", units: "Celsius degree", valid_min: -1.5, valid_max: 38, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"NUMP":      {long_name: "NUMBER OF PARTICLES", units: "number/m3", valid_min: 0, valid_max: 999999, format: "%9.3e", _FillValue: 1e+36, types: "float32"},
			"SLEV":      {long_name: "OBSERVED SEA LEVEL", units: "meter", valid_min: 0, valid_max: 6000, format: "%8.3f", _FillValue: 1e+36, types: "float32"},
			"PHPH":      {long_name: "PH", units: "pH unit", valid_min: 7.4, valid_max: 8.4, format: "%5.3f", _FillValue: 1e+36, types: "float32"},
			"DXAP":      {long_name: "DIADINOXANTHINE", units: "milligram/m3", valid_min: 0, valid_max: 5, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"SVEL":      {long_name: "SOUND VELOCITY", units: "meter/second", valid_min: 1350, valid_max: 1600, format: "%7.2f", _FillValue: 1e+36, types: "float32"},
			"SLCA":      {long_name: "SILICATE (SIO4-SI) CONTENT", units: "millimole/m3", valid_min: 0, valid_max: 200, format: "%7.3f", _FillValue: 1e+36, types: "float32"},
			"MMGS":      {long_name: "Mg IN THE SEDIMENT", units: "%", valid_min: 0, valid_max: 99, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"MTIS":      {long_name: "Ti IN THE SEDIMENT", units: "%", valid_min: 0, valid_max: 99, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"MCRS":      {long_name: "Cr IN THE SEDIMENT", units: "ppm", valid_min: 0, valid_max: 900, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"LATD":      {long_name: "latitude of measurement", units: "degrees_north", valid_min: -90, valid_max: 90, format: "%+3.0f", _FillValue: 1e+36, types: "float32"},
			"TUR0":      {long_name: "LIGHT TRANSMISSION -  NOT USED", units: "%", valid_min: 0, valid_max: 100, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"DFCW":      {long_name: "NUMBER OF SW DINOFLAGELLATES", units: "10+3 cell/m3", valid_min: 0, valid_max: 900000, format: "%6.0f", _FillValue: 1e+36, types: "float32"},
			"WSPE":      {long_name: "WIND SPEED EASTWARD COMPONENT", units: "meter/second", valid_min: 0, valid_max: 100, format: "%7.3f", _FillValue: 1e+36, types: "float32"},
			"DO21":      {long_name: "DISSOLVED OXYGEN SECONDARY SENSOR", units: "ml/l", valid_min: 0, valid_max: 10, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"HEIG":      {long_name: "ATMOSPHERIC HEIGHT", units: "meter", valid_min: 0, valid_max: 40000, format: "%8.2f", _FillValue: 0, types: "float32"},
			"MFEP":      {long_name: "Fe IN SUSPENDED MATTER", units: "milligram/m3", valid_min: 0, valid_max: 99.999, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"CHCZ":      {long_name: "CHLOROPHYLL-C1+C2", units: "milligram/m3", valid_min: 0, valid_max: 900, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"PCO2":      {long_name: "CO2 PART. PRES IN DRY/WET GAS", units: "microatmosphere", valid_min: 100, valid_max: 700, format: "%7.3f", _FillValue: 1e+36, types: "float32"},
			"DFPF":      {long_name: "DINOFLAGELLATES PPC FLUX", units: "milligram C/(m2.day)", valid_min: 0, valid_max: 99.999, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"DOXY":      {long_name: "DISSOLVED OXYGEN", units: "millimole/m3", valid_min: 0, valid_max: 650, format: "%7.3f", _FillValue: 1e+36, types: "float32"},
			"TRIW":      {long_name: "DISSOLVED TRITIUM", units: "TU", valid_min: 0, valid_max: 2, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"TE06":      {long_name: "SEA TEMPERATURE", units: "Celsius degree", valid_min: -2, valid_max: 32, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"VAVH":      {long_name: "AVER. HEIGHT HIGHEST 1/3 WAVE", units: "metre", valid_min: 0, valid_max: 30, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"NORW":      {long_name: "NITRIFICATION", units: "micromole nitrogen/(m3.day)", valid_min: 0, valid_max: 100, format: "%5.0f", _FillValue: 1e+36, types: "float32"},
			"SIOF":      {long_name: "PART. BIOGENIC Si FLUX", units: "milligram/(m2.day)", valid_min: 0, valid_max: 600, format: "%7.3f", _FillValue: 1e+36, types: "float32"},
			"SCSP":      {long_name: "SEA SURFACE CURRENT SPEED", units: "meter/second", valid_min: 0, valid_max: 10, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"TE04":      {long_name: "SEA TEMPERATURE", units: "Celsius degree", valid_min: -2, valid_max: 32, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"MSIP":      {long_name: "Si IN SUSPENDED MATTER", units: "milligram/m3", valid_min: 0, valid_max: 999.999, format: "%7.3f", _FillValue: 1e+36, types: "float32"},
			"ALKY":      {long_name: "ALKALINITY", units: "millimole/m3", valid_min: 1500, valid_max: 2500, format: "%4.0f", _FillValue: 1e+36, types: "float32"},
			"AXAP":      {long_name: "ALLOXANTHINE", units: "milligram/m3", valid_min: 0, valid_max: 5, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"HELW":      {long_name: "DISSOLVED HELIUM", units: "nanomole/kg", valid_min: 1, valid_max: 2, format: "%7.4f", _FillValue: 1e+36, types: "float32"},
			"LINC":      {long_name: "LONG-WAVE INCOMING RADIATION", units: "watt/m2", valid_min: -500, valid_max: 500, format: "%+5.1f", _FillValue: 1e+36, types: "float32"},
			"MNBS":      {long_name: "Nb IN THE SEDIMENT", units: "ppm", valid_min: 0, valid_max: 900, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"PTZF":      {long_name: "PROTOZOA PPC FLUX", units: "milligram C/(m2.day)", valid_min: 0, valid_max: 99.999, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"DO12":      {long_name: "DISSOLVED OXYGEN PRIMARY SENSOR", units: "micromole/kg", valid_min: 0, valid_max: 450, format: "%7.3f", _FillValue: 1e+36, types: "float32"},
			"HFUP":      {long_name: "19'HEXANOYLOXYFUCOXANTHINE", units: "milligram/m3", valid_min: 0, valid_max: 5, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"PSAL":      {long_name: "sea water salinity", units: "1", valid_min: 0, valid_max: 40, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"TACC":      {long_name: "TEMPERATURE ACCURACY", units: ".", valid_min: 0, valid_max: 4, format: "%2d", _FillValue: 1e+36, types: "float32"},
			"GLUC":      {long_name: "GLUCIDE", units: "milligram/m3", valid_min: 0, valid_max: 200, format: "%5.0f", _FillValue: 1e+36, types: "float32"},
			"CHBD":      {long_name: "DIVINYL-CHLOROPHYLL-B", units: "milligram/m3", valid_min: 0, valid_max: 5, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"DRDD":      {long_name: "DURATION (DAYS)", units: "ddd", valid_min: 0, valid_max: 999, format: "%3.0f", _FillValue: 1e+36, types: "float32"},
			"NFHW":      {long_name: "HETEROTROPHIC NANOFLAGELLATES", units: "10+3 cell/m3", valid_min: 0, valid_max: 900000, format: "%6.0f", _FillValue: 1e+36, types: "float32"},
			"LOND":      {long_name: "latitude of measurement", units: "degree_east", valid_min: -179, valid_max: 180, format: "%+4.0f", _FillValue: 1e+36, types: "float32"},
			"MSMP":      {long_name: "MEAN SPHERIC DIAM. MEDIAN", units: "millimeter", valid_min: 0, valid_max: 99, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"COCW":      {long_name: "NUMBER OF SW COCCOLITHOPHORIDS", units: "10+3 cell/m3", valid_min: 0, valid_max: 900000, format: "%6.0f", _FillValue: 1e+36, types: "float32"},
			"PRES":      {long_name: "sea water pressure", units: "decibar", valid_min: 0, valid_max: 6500, format: "%6.1f", _FillValue: 1e+36, types: "float32"},
			"NHUW":      {long_name: "AMMONIUM UPTAKE", units: "micromole nitrogen/(m3.day)", valid_min: 0, valid_max: 900, format: "%5.0f", _FillValue: 1e+36, types: "float32"},
			"POTT":      {long_name: "TOTAL ORGANIC PHOSPHORUS (D+P)", units: "micromole/kg", valid_min: 0, valid_max: 5, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"CHTF":      {long_name: "PART. TOTAL CARBOHYDRATES FLUX", units: "milligram/(m2.day)", valid_min: 0, valid_max: 99, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"PHOS":      {long_name: "PHOSPHATE (PO4-P) CONTENT", units: "millimole/m3", valid_min: 0, valid_max: 4, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"MYYS":      {long_name: "Y IN THE SEDIMENT", units: "ppm", valid_min: 0, valid_max: 900, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"COPF":      {long_name: "COCCOLITHOPHORIDAE PPC FLUX", units: "milligram C/(m2.day)", valid_min: 0, valid_max: 99.999, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"CE3W":      {long_name: "DISSOLVED CFC113 ERROR", units: "picomole/kg", valid_min: 0, valid_max: 0.5, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"PODW":      {long_name: "DISSOLVED ORGANIC PHOSPHORUS", units: "micromole/kg", valid_min: 0, valid_max: 5, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"FUCP":      {long_name: "FUCOXANTHINE", units: "milligram/m3", valid_min: 0, valid_max: 5, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"LATX":      {long_name: "LATITUDE", units: "degree_north", valid_min: -90, valid_max: 90, format: "%+8.4f", _FillValue: 1e+36, types: "float32"},
			"VERR":      {long_name: "VELOCITY ERROR", units: "meter/second", valid_min: 0, valid_max: 90, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"SSTP":      {long_name: "SEA SURFACE TEMPERATURE", units: "Celsius degree", valid_min: -1.5, valid_max: 38, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"DAYD":      {long_name: "DECIMAL JULIAN DAY TIME ORIGIN 0", units: "decimal day", valid_min: 0, valid_max: 3660, format: "%9.5f", _FillValue: 1e+36, types: "float64"},
			"SPDI":      {long_name: "INDICATED PLATFORM SPEED-SHIP", units: "meter/second", valid_min: 0, valid_max: 90, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"MCUS":      {long_name: "Cu IN THE SEDIMENT", units: "ppm", valid_min: 0, valid_max: 900, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"LATM":      {long_name: "LATITUDE MINUTES", units: "minute", valid_min: 0, valid_max: 59.999, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"LGH4":      {long_name: "LIGHT IRRADIANCE SURFACE PAR", units: "micromole photon/(m2.s)", valid_min: 0, valid_max: 3000, format: "%8.3f", _FillValue: 1e+36, types: "float32"},
			"PXAP":      {long_name: "PRASINOXANTHINE", units: "milligram/m3", valid_min: 0, valid_max: 5, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"YEAR":      {long_name: "YEAR", units: "yyyy", valid_min: 1900, valid_max: 2020, format: "%4.4d", _FillValue: 1e+36, types: "float32"},
			"DO22":      {long_name: "DISSOLVED OXYGEN SECONDARY SENSOR", units: "micromole/kg", valid_min: 0, valid_max: 450, format: "%7.3f", _FillValue: 1e+36, types: "float32"},
			"CHLB":      {long_name: "CHLOROPHYLL-B TOTAL", units: "milligram/m3", valid_min: 0, valid_max: 99, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"PRRT":      {long_name: "PRECIPITATION RATE", units: "millimeter/hour", valid_min: 0, valid_max: 900, format: "%7.3f", _FillValue: 1e+36, types: "float32"},
			"PROT":      {long_name: "PROTEIN", units: "milligram/m3", valid_min: 0, valid_max: 500, format: "%5.0f", _FillValue: 1e+36, types: "float32"},
			"DCAW":      {long_name: "DARK CARBON ABSORPTION", units: "milligram/m3", valid_min: -0.05, valid_max: 10, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"NSCT":      {long_name: "CURRENT NORTH COMPONENT", units: "meter/second", valid_min: -100, valid_max: 100, format: "%+7.3f", _FillValue: 1e+36, types: "float32"},
			"HELD":      {long_name: "DELTA HELIUM 3", units: "%", valid_min: -99, valid_max: 99, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"MMNF":      {long_name: "MN FLUX IN SETTLING PARTICLES", units: "microgram/(m2.day)", valid_min: 0, valid_max: 3400, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"MUUS":      {long_name: "U IN THE SEDIMENT", units: "ppm", valid_min: 0, valid_max: 900, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"VCSP":      {long_name: "VERTICAL CURRENT SPEED", units: "meter/second", valid_min: 0, valid_max: 9, format: "%5.3f", _FillValue: 1e+36, types: "float32"},
			"MBRS":      {long_name: "Br IN THE SEDIMENT", units: "ppm", valid_min: -50, valid_max: 50, format: "%+6.2f", _FillValue: 1e+36, types: "float32"},
			"CHAF":      {long_name: "CHLOROPHYLL-A VERTICAL FLUX", units: "milligram/(m2.day)", valid_min: 0, valid_max: 99, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"FLU2":      {long_name: "FLUORESCENCE", units: "milligram/m3", valid_min: 0, valid_max: 1, format: "%6.4f", _FillValue: 1e+36, types: "float32"},
			"CND2":      {long_name: "ELECTRICAL CONDUCTIVITY SECONDARY SENSOR", units: "mho/meter", valid_min: 3, valid_max: 7, format: "%5.3f", _FillValue: 1e+36, types: "float32"},
			"AMOW":      {long_name: "AMMONIUM (NH4-N) CONTENT", units: "micromole/kg", valid_min: 0, valid_max: 1, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"LONM":      {long_name: "LONGITUDE MINUTES", units: "minute", valid_min: 0, valid_max: 59.999, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"MPBF":      {long_name: "PB FLUX IN SETTLING PARTICLES", units: "microgram/(m2.day)", valid_min: 0, valid_max: 200, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"TE01":      {long_name: "SEA TEMPERATURE PRIMARY SENSOR", units: "Celsius degree", valid_min: -2, valid_max: 32, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"TREW":      {long_name: "DISSOLVED TRITIUM ERROR", units: "TU", valid_min: 0, valid_max: 5, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"NODW":      {long_name: "DISSOLVED ORGANIC NITROGEN", units: "micromole/kg", valid_min: 0, valid_max: 99, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"MNDS":      {long_name: "Nd IN THE SEDIMENT", units: "ppm", valid_min: 0, valid_max: 900, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"OSMP":      {long_name: "ORGANIC SUSPENDED MATTER", units: "gram/m3", valid_min: 0, valid_max: 10, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"MSRS":      {long_name: "Sr IN THE SEDIMENT", units: "ppm", valid_min: 0, valid_max: 900, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"MVVS":      {long_name: "V IN THE SEDIMENT", units: "ppm", valid_min: -50, valid_max: 900, format: "%+6.2f", _FillValue: 1e+36, types: "float32"},
			"SSJT":      {long_name: "SEA SURFACE WATER JACKET TEMPERATURE", units: "Celsius degree", valid_min: -1.5, valid_max: 38, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"CL4W":      {long_name: "DISSOLVED C-TETRACHLORIDE", units: "picomole/kg", valid_min: 0, valid_max: 10, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"GSPD":      {long_name: "GUST WIND SPEED", units: "meter/second", valid_min: 0, valid_max: 99, format: "%2.0f", _FillValue: 1e+36, types: "float32"},
			"MRBS":      {long_name: "Rb IN THE SEDIMENT", units: "ppm", valid_min: 0, valid_max: 900, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"SSAL":      {long_name: "SALINITY (PRE-1978 DEFN)", units: "P.S.U.", valid_min: 0, valid_max: 40, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"VOCP":      {long_name: "VOLUME CONC. OF PARTICLES", units: "p.p.m.", valid_min: 0, valid_max: 99, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"MZRS":      {long_name: "Zr IN THE SEDIMENT", units: "ppm", valid_min: 0, valid_max: 900, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"CPH4":      {long_name: "CHLOROPHYLL-A/20 MICRON FILTER", units: "milligram/m3", valid_min: 0, valid_max: 99, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"BCMW":      {long_name: "BACTERIAL BIOMASS IN SEA WATER", units: "milligram C/m3", valid_min: 0, valid_max: 16, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"D13C":      {long_name: "DELTA 13 C SIGNATURE", units: "%", valid_min: -50, valid_max: 100, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"DOPW":      {long_name: "DISSOLVED ORGANIC PHOSPHORUS", units: "millimole/m3", valid_min: 0, valid_max: 50, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"ISMP":      {long_name: "INORGANIC SUSPENDED MATTER", units: "gram/m3", valid_min: 0, valid_max: 10, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"MALS":      {long_name: "Al IN THE SEDIMENT", units: "%", valid_min: 0, valid_max: 99, format: "%6.3f", _FillValue: 1e+36, types: "float32"},
			"CF2W":      {long_name: "DISSOLVED CFC12", units: "picomole/kg", valid_min: -0.01, valid_max: 90, format: "%7.4f", _FillValue: 1e+36, types: "float32"},
			"LTUW":      {long_name: "LEUCINE UPTAKE RATE", units: "microgram carbon/(m3.h)", valid_min: 0, valid_max: 90, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"LGHT":      {long_name: "LIGHT IRRADIANCE IMMERGED PAR", units: "micromole photon/(m2.s)", valid_min: 0, valid_max: 4000, format: "%8.3f", _FillValue: 1e+36, types: "float32"},
			"MLAS":      {long_name: "La IN THE SEDIMENT", units: "ppm", valid_min: 0, valid_max: 900, format: "%6.2f", _FillValue: 1e+36, types: "float32"},
			"DO11":      {long_name: "DISSOLVED OXYGEN PRIMARY SENSOR", units: "ml/l", valid_min: 0, valid_max: 10, format: "%5.2f", _FillValue: 1e+36, types: "float32"},
			"MCUF":      {long_name: "CU FLUX IN SETTLING PARTICLES", units: "microgram/(m2.day)", valid_min: 0, valid_max: 120, format: "%6.2f", _FillValue: 1e+36},
			"DEPTH":     {long_name: "Depth of measurement", units: "meters", valid_min: 0, valid_max: 6000, format: "%6.1f", _FillValue: 1e+36, types: "float32"},
			"TYPECAST":  {long_name: "TYPE OF CAST: UNKNOW=0, PHY=1, BIO=2, GEO=3", units: "N/A", valid_min: 0, valid_max: 9, format: "%1d", _FillValue: 255, types: "uint8"},
			"PROFILE":   {long_name: "PROFILE NUMBER", units: "N/A", valid_min: 1, valid_max: 99999, format: "%5d", _FillValue: 9999999, types: "int32"},
			"LATITUDE":  {long_name: "LATITUDE", units: "degree_north", valid_min: -90, valid_max: 90, format: "%+8.4f", _FillValue: 1e+36, types: "float32"},
			"LONGITUDE": {long_name: "LONGITUDE", units: "degree_east", valid_min: -180, valid_max: 180, format: "%+9.4f", _FillValue: 1e+36, types: "float32"}}

		//f("%#v", roscop)
		//f("Code roscop: %v\n", roscop["TIME"])
		return roscop
	}
}
