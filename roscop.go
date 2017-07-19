// Roscop.go
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type m map[string]map[string]string

// Roscop represents the attributes associated with each netCDF variable
type Roscop struct {
	m
	physicalParametersOrderedList []string
	attributesOrderedList         map[string][]string
	attributesType                map[string]string
}

// NewRoscop read Roscop definition from fileName and return Roscop object
// read documentation for csv is at http://golang.org/pkg/encoding/csv/
func NewRoscop(filename string) Roscop {

	// use a map of map to store for each physical parameter a map where keys are
	// attributes
	var r = Roscop{
		m: make(m),
		physicalParametersOrderedList: []string{},
		attributesOrderedList:         make(map[string][]string),
		attributesType:                make(map[string]string),
	}

	// open csv file
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

	// init new reader
	reader := csv.NewReader(file)
	// options are available at:
	// http://golang.org/src/pkg/encoding/csv/reader.go?s=3213:3671#L94
	reader.Comma = ';'

	// read first header line
	fields, err := reader.Read()
	//fmt.Println(fields)

	// read second line
	types, err := reader.Read()
	//fmt.Println(types)

	// fill map of attribute type: string or numeric
	for i := 0; i < len(fields); i++ {
		r.attributesType[fields[i]] = types[i]
	}
	//fmt.Println(r.attributesType)

	// read next lines and fill Roscop object
	for {

		// initialize a new empty map to store attributes variable
		// with pair of name and value
		mfields := map[string]string{}

		// read just one record, but we could ReadAll() as well
		record, err := reader.Read()
		// end-of-file is fitted into err
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
		}
		// key is physical parameter
		key := record[0]

		for i := 1; i < len(fields); i++ {
			if record[i] != "" {
				// store attributes position
				// the iteration order is not specified and is not guaranteed to be
				// the same from one iteration to the next in golang
				r.attributesOrderedList[key] = append(r.attributesOrderedList[key], fields[i])
				// fill map with non empty values
				mfields[fields[i]] = record[i]
			}
		}

		// put the new map to Roscop map with the rigth physical parameter
		r.m[key] = mfields

		// the iteration order is not specified and is not guaranteed to be
		// the same from one iteration to the next in golang
		r.physicalParametersOrderedList = append(r.physicalParametersOrderedList, key)
		fmt.Fprintf(debug, "%s: %v\n", record[0], mfields)
	}
	//f("%#v", r)
	return r
}

// GetPhysicalParameters returm an ordered list of all physical parameters
func (r Roscop) GetPhysicalParameters() []string {
	return r.physicalParametersOrderedList
}

// GetAttributes returm an ordered list of attributes for an physical parameter
func (r Roscop) GetAttributes(physicalParameter string) []string {
	// remove first key "types"
	return r.attributesOrderedList[physicalParameter][1:]
}

// GetAttributesStringValue returm the attribute value as a string for a physical parameter
func (r Roscop) GetAttributesStringValue(physicalParameter string, attributeName string) string {
	return r.m[physicalParameter][attributeName]
}

// GetAttributesValue returm the attribute value with the right type for a physical parameter
func (r Roscop) GetAttributesValue(physicalParameter string, attributeName string) interface{} {
	switch r.attributesType[attributeName] {
	case "string":
		return r.m[physicalParameter][attributeName]
	case "numeric":
		// convert generic numeric type with the rigth type of physical parameter
		switch r.m[physicalParameter]["types"] {
		case "char", "byte":
			return byte(r.m[physicalParameter][attributeName][0])
		case "int", "int32":
			value, _ := strconv.ParseInt(r.m[physicalParameter][attributeName], 10, 32)
			return int32(value)
		case "float32", "float":
			value, _ := strconv.ParseFloat(r.m[physicalParameter][attributeName], 32)
			return float32(value)
		case "float64", "double":
			value, _ := strconv.ParseFloat(r.m[physicalParameter][attributeName], 64)
			return value
		default:
			log.Fatalf("Error: check the column types  of your Roscop file,"+
				" valid values are: char, byte, int, int32, float, float32, double and float64\n"+
				"find \"%s\" for %s, %s = %s", r.m[physicalParameter]["types"], physicalParameter,
				attributeName, r.m[physicalParameter][attributeName])
		}
	default:
		log.Fatal("Error: check the second line of your Roscop file," +
			" values sould be string or numeric")
	}

	return r.m[physicalParameter][attributeName]
}

//func main() {

//	// initialize new Roscop object from file
//	Roscop := NewRoscop("code_Roscop.csv")

//	// loop over each physicalParameter and display attribute values
//	for _, physicalParameter := range Roscop.GetPhysicalParameters() {
//		fmt.Printf("%s => ", physicalParameter)

//		for _, attributeName := range Roscop.GetAttributes(physicalParameter) {
//			value := Roscop.GetAttributesValue(physicalParameter, attributeName)
//			fmt.Printf("%s: %v (%T), ", attributeName, value, value)
//		}
//		fmt.Println()
//	}
//}
