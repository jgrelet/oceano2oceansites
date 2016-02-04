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

type m map[string]map[string]string

type roscop struct {
	m
	physicalParametersOrderedList []string
	attributesOrderedList         map[string][]string
	attributesType                map[string]string
}

// documentation for csv is at http://golang.org/pkg/encoding/csv/
func NewRoscop(filename string) roscop {

	// use a map of map to store for each physical parameter a map where keys are
	// attributes
	var r = roscop{
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

	// fill map of attribute type
	for i := 0; i < len(fields); i++ {
		r.attributesType[fields[i]] = types[i]
	}
	fmt.Println(r.attributesType)

	// read next lines
	for {

		// initialize a new empty map to store variable attributes
		// name and the value
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

		// put the new map to roscop map with the rigth physical parameter
		r.m[key] = mfields

		// the iteration order is not specified and is not guaranteed to be
		// the same from one iteration to the next in golang
		r.physicalParametersOrderedList = append(r.physicalParametersOrderedList, key)
		//fmt.Printf("%s: %v\n", record[0], mfields)
	}
	//f("%#v", r)
	return r
}

// returm an ordered list of all physical parameters
func (r roscop) GetPhysicalParameters() []string {
	return r.physicalParametersOrderedList
}

// returm an ordered list of attributes for an physical parameter
func (r roscop) GetAttributes(physicalParameter string) []string {
	// remove first key "types"
	return r.attributesOrderedList[physicalParameter][1:]
}

// returm the attribute value as a string for a physical parameter
func (r roscop) GetAttributesStringValue(physicalParameter string, attributeName string) string {
	return r.m[physicalParameter][attributeName]
}

// returm the attribute value with the right type for a physical parameter
func (r roscop) GetAttributesValue(physicalParameter string, attributeName string) interface{} {
	switch r.attributesType[attributeName] {
	case "string":
		return r.m[physicalParameter][attributeName]
	case "char", "byte":
		return byte(r.m[physicalParameter][attributeName][0])
	case "int", "int32":
		value, _ := strconv.ParseInt(r.m[physicalParameter][attributeName], 32, 32)
		return value
	case "float32", "float":
		value, _ := strconv.ParseFloat(r.m[physicalParameter][attributeName], 32)
		return value
	case "float64", "double":
		value, _ := strconv.ParseFloat(r.m[physicalParameter][attributeName], 64)
		return value
	default:
		log.Fatal("bad type")

	}

	return r.m[physicalParameter][attributeName]
}

//func main() {

//	// initialize new roscop object from file
//	roscop := NewRoscop("code_roscop.csv")

//	// loop over each physicalParameter and display attribute values
//	for _, physicalParameter := range roscop.GetPhysicalParameters() {
//		fmt.Printf("%s => ", physicalParameter)

//		for _, attributeName := range roscop.GetAttributes(physicalParameter) {
//			value := roscop.GetAttributesValue(physicalParameter, attributeName)
//			fmt.Printf("%s: %v (%T), ", attributeName, value, value)
//		}
//		fmt.Println()
//	}
//}
