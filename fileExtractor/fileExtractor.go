package fileExtractor

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// define new slice type for trajectory, profile and timeSerie
type data map[string][]float64
type profile map[string][][]float64
type trajectory data
type timeSerie data

// FileExtractOptions contains configurable options used to read an ASCII file.
type FileExtractOptions struct {

	// map with physical parameter as key and the column number position as value
	// the order of keys is used to defined the header
	// example: PRES,3,DEPTH,4,ETDD,2,TEMP,5,PSAL,17,DENS,20
	varsList map[string]int

	// slice from varList key map, ex: PRES  DEPTH  ETDD  TEMP  PSAL  DENS
	hdr []string

	// the separator of data in ascii file, default is space
	separator string

	// number of line to skip before read data
	skipLine int

	// header delimiter
	hdrDelimiter string
}

// NewFileExtractOptions will create a new FileExtractOptions type with some
// empty default values.
func NewFileExtractOptions() *FileExtractOptions {
	o := &FileExtractOptions{
		hdr:          []string{},
		varsList:     map[string]int{},
		separator:    "",
		skipLine:     -1,
		hdrDelimiter: "",
	}
	return o
}

// SetVars will set the parameters and their columns to extract from file
func (o *FileExtractOptions) SetVarsList(split string) *FileExtractOptions {
	// create empty map and header list
	m := map[string]int{}
	h := []string{}

	// construct map from split
	fields := strings.Split(split, ",")
	for i := 0; i < len(fields); i += 2 {
		if v, err := strconv.Atoi(fields[i+1]); err == nil {
			m[fields[i]] = v
			h = append(h, fields[i])
		} else {
			log.Fatalf("Check the input of SetVars: %v\n", err)
		}
	}
	// copy map and header list to FileExtractOptions object
	o.varsList = m
	o.hdr = h
	return o
}

// VarsList getter
func (o *FileExtractOptions) VarsList() map[string]int {
	return o.varsList
}

// SetSeparator will override the default separator (space)
func (o *FileExtractOptions) SetSeparator(sep string) *FileExtractOptions {
	o.separator = sep
	return o
}

// SetSkipLine will set to skip header line
func (o *FileExtractOptions) SetSkipLine(line int) *FileExtractOptions {
	o.skipLine = line
	return o
}

// header delimiter
func (o *FileExtractOptions) SetHdrDelimiter(delimiter string) *FileExtractOptions {
	o.hdrDelimiter = delimiter
	return o
}

// display FileExtractOptions object
func (o FileExtractOptions) String() string {
	return fmt.Sprintf("FileExtractOptions:\nFields:%s\nVars: %v\nSkipLine: %d\nHdr delimiter: %s",
		o.hdr, o.varsList, o.skipLine, o.hdrDelimiter)
}

// FileExtractor contains FileExtractOptions object and map data extracted from ASCII file.
type FileExtractor struct {

	// options read from toml configuration file
	*FileExtractOptions

	// extracted data
	data

	// the size of extracted data
	length int
}

// NewFileExtracter will create a new FileExtractor type with some values from
// configuration (not implemented)
func NewFileExtractor(o *FileExtractOptions) *FileExtractor {

	// in this constructor, we use composition (or embedding) vs inheritance
	fe := &FileExtractor{
		FileExtractOptions: o,
		data:               make(map[string][]float64),
		length:             0,
	}
	// initialize map for each key to a slice
	for _, name := range fe.hdr {
		fe.data[name] = []float64{}
	}
	return fe
}

// get the the size of map data
func (fe FileExtractor) Length() int {
	return fe.length
}

// Read an ASCII file and extract data and save then to map data
func (ext *FileExtractor) Read(filename string) (err error) {

	fid, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fid.Close()

	// open bufio for file
	scanner := bufio.NewScanner(fid)

	// skip hdr lines if skipLiner defined
	if ext.skipLine != -1 {
		for i := 0; i < ext.skipLine; i++ {
			scanner.Scan()
		}
	}
	if ext.hdrDelimiter != "" {
		for scanner.Scan() {
			str := scanner.Text()
			if str == ext.hdrDelimiter {
				break
			}
		}
	}

	// read file
	for scanner.Scan() {
		var values []string

		// parse each line to string
		str := scanner.Text()

		// split the string str with defined separator
		if ext.separator != "" {
			values = strings.Split(str, ext.separator)
		} else {
			// split the string str with one or more space
			values = strings.Fields(str)
		}

		// fill map data
		for key, column := range ext.varsList {
			// slice index start at 0
			ind := column - 1

			if ind < len(values) {
				if v, err := strconv.ParseFloat(values[ind], 64); err == nil {
					// column start at zero
					ext.data[key] = append(ext.data[key], v)
				}
			} else {
				continue
			}
		}

		ext.length += 1
	}
	return nil
}

// get the the map data
func (fe *FileExtractor) Data() data {
	return fe.data
}

// print the result
func (ext FileExtractor) String() string {
	var s []string

	s = append(s, "FileExtractor:\n")
	for key, _ := range ext.varsList {
		s = append(s, fmt.Sprintf("\n%s: %7.3f", key, ext.data[key]))
	}
	return strings.Join(s, "")
}

type ProfilesExtractor struct {

	// may contain wildcards if multiple files to read
	fileNames []string

	//FileExtractor
	fileExtractor FileExtractor

	values profile
}

func NewProfilesExtractor(files string, fe FileExtractor) *ProfilesExtractor {
	fs, _ := filepath.Glob(files)
	return &ProfilesExtractor{
		fileNames:     fs,
		fileExtractor: fe,
		values:        make(map[string][][]float64),
	}
}

// Len returns the amount of element in the slice.
func (fse *ProfilesExtractor) Len() int {
	return len((fse.values))
}

// Read one or multiple ASCII files and extract data
// add return error
func (fse *ProfilesExtractor) Read() {

	for _, fileName := range fse.fileNames {

		err := fse.fileExtractor.Read(fileName)
		if err != nil {
			log.Fatalf("ProfilesExtractor.Read(): %s", err)
		}
		for key, data := range fse.fileExtractor.data {
			fse.values[key] = append(fse.values[key], data)
		}
	}
}

func (fse *ProfilesExtractor) String() string {
	var s []string

	s = append(s, fmt.Sprintf("ProfilesExtractor:\nFiles: %v\n", fse.fileNames))
	for key, _ := range fse.values {
		s = append(s, fmt.Sprintf("%s:%v\n", key, fse.values[key]))
	}
	return strings.Join(s, "")
}

type TrajectoriesExtractor struct {

	// may contain wildcards if multiple files to read
	fileNames []string

	//FileExtractor
	fileExtractor FileExtractor

	values trajectory
}

func NewTrajectoriesExtractor(files string, fe FileExtractor) *TrajectoriesExtractor {
	fs, _ := filepath.Glob(files)
	return &TrajectoriesExtractor{
		fileNames:     fs,
		fileExtractor: fe,
		values:        make(map[string][]float64),
	}
}

// Read one or multiple ASCII files and extract data
// add return error
func (fse *TrajectoriesExtractor) Read() {

	for _, fileName := range fse.fileNames {

		err := fse.fileExtractor.Read(fileName)
		if err != nil {
			log.Fatalf("TrajectoriesExtractor.Read(): %s", err)
		}
		for key, _ := range fse.fileExtractor.data {
			// append a slice with the ...
			fse.values[key] = append(fse.values[key], fse.fileExtractor.data[key]...)
		}
	}
}

func (fse *TrajectoriesExtractor) String() string {
	var s []string

	s = append(s, fmt.Sprintf("TrajectoriesExtractor:\nFiles: %v\n", fse.fileNames))
	for key, _ := range fse.values {
		s = append(s, fmt.Sprintf("%s:%v\n", key, fse.values[key]))
	}
	return strings.Join(s, "")
}
