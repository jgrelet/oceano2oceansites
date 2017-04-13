package main

import (
	"errors"
	"fmt"
	_ "log"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	secondsPerMinute           = 60
	secondsPerHour             = 60 * 60
	secondsPerDay              = 24 * secondsPerHour
	secondsPerWeek             = 7 * secondsPerDay
	daysPer400Years            = 365*400 + 97
	daysPer100Years            = 365*100 + 24
	daysPer4Years              = 365*4 + 1
	unixToInternal       int64 = (1969*365 + 1969/4 - 1969/100 + 1969/400) * secondsPerDay
	internalToUnix       int64 = -unixToInternal
	oceanSitesToInternal int64 = (1949*365 + 1949/4 - 1949/100 + 1949/400) * secondsPerDay
	internalToOceanSites int64 = -oceanSitesToInternal
)

// ProfileType is integer
type ProfileType int

// ProfileType define constante for profile
const (
	UNKNOW ProfileType = iota // UNKNOW == 0
	PHY                       // PHY = 1
	BIO                       // BIO == 2
	GEO                       // GEO == 3

)

// Time represent an object time.Time with the number of seconds
// elapsed since January 1, 1970 UTC
type Time struct {
	time.Time
	nsec int64
}

// NewTimeFromString construct time object from a string date
func NewTimeFromString(format, value string) *Time {
	t, _ := time.Parse(format, value)
	return &Time{t, t.Unix()}
}

// NewTimeFromSec construct time object from number of second since 01/01/1970
func NewTimeFromSec(nsec int64) *Time {
	t := time.Unix(nsec, 0).UTC()
	return &Time{t, t.Unix()}
}

// NewTimeFromJulian construct time object from decimal julian day since 01/01/1950
func NewTimeFromJulian(julian float64) *Time {
	t := time.Unix(int64(julian*86400.0)+oceanSitesToInternal+internalToUnix, 0).UTC()
	return &Time{t, t.Unix()}
}

// NewTimeFromJulianDay construct time object from decimal julian day from current year
func NewTimeFromJulianDay(julian float64, c *Time) *Time {
	nsec := time.Date(c.Year(), time.January, 0, 0, 0, 0, 0, time.UTC)
	t := time.Unix(int64(julian*86400.0)+nsec.Unix(), 0).UTC()
	return &Time{t, t.Unix()}
}

// Time2JulianDec compute from time object a decimal julian day from 1950
func (t *Time) Time2JulianDec() float64 {
	const diffOrigin = 2433283.0 // diff between UNIX DATE and 1950/1/1 00:00:00
	a := int(14-t.Month()) / 12
	y := t.Year() + 4800 - a
	m := int(t.Month()) + 12*a - 3
	julianDay := int(t.Day()) + (153*m+2)/5 + 365*y + y/4
	julianDay = julianDay - y/100 + y/400 - 32045.0 - diffOrigin
	d := float64(julianDay) + float64(t.Hour())/24 +
		float64(t.Minute())/1440 + float64(t.Second())/86400
	fmt.Fprintf(debug, "Julian day: %v ", d)
	return d
}

// JulianDayOfYear compute from time object a decimal julian day from the current year
func (t *Time) JulianDayOfYear() float64 {
	julianDay := t.YearDay()
	return float64(julianDay) + float64(t.Hour())/24 +
		float64(t.Minute())/1440 + float64(t.Second())/86400
}

// Position2Decimal convert position "DD MM.SS S" to decimal position
func Position2Decimal(pos string) (float64, error) {

	var multiplier float64 = 1
	var value float64

	var regNmeaPos = regexp.MustCompile(`(\d+)\s+(\d+.\d+)\s+(\w)`)

	if strings.Contains(pos, "S") || strings.Contains(pos, "W") {
		multiplier = -1.0
	}
	match := regNmeaPos.MatchString(pos)
	if match {
		res := regNmeaPos.FindStringSubmatch(pos)
		deg, _ := strconv.ParseFloat(res[1], 64)
		min, _ := strconv.ParseFloat(res[2], 64)
		tmp := math.Abs(min)
		sec := (tmp - min) * 100.0
		value = (deg + (min+sec/100.0)/60.0) * multiplier
		fmt.Fprintln(debug, "positionDeci:", pos, " -> ", value)
	} else {
		return 1e36, errors.New("positionDeci: failed to decode position")
	}
	return value, nil
}

// DecimalPosition2String convert  decimal position to string,
// hemi = 0 for latitude, 1 for longitude
func DecimalPosition2String(position float64, hemi string) string {
	var neg, pos, geo rune
	var str string

	if strings.Contains(hemi, "EW") {
		neg = 'W'
		pos = 'E'
	} else {
		neg = 'S'
		pos = 'N'
	}
	if position < 0 {
		geo = neg
	} else {
		geo = pos
	}
	tmp := math.Abs(position)
	deg := int(tmp)
	tmp = (tmp - float64(deg)) * 60
	min := tmp

	if strings.Contains(hemi, "EW") {
		str = fmt.Sprintf("%03d%s%06.3f %c", deg, "\u00B0", min, geo)
	} else {
		str = fmt.Sprintf("%02d%s%06.3f %c", deg, "\u00B0", min, geo)
	}
	return str
}

func isArray(a interface{}) bool {
	var v reflect.Value
	v = reflect.ValueOf(a)

	var k reflect.Kind
	k = v.Kind()

	if k == reflect.Array {
		return true
	}
	return false
}

func zeros(size int) []float64 {
	return make([]float64, size)
}
func ones(size int) []float64 {
	m := make([]float64, size)
	for i := range m {
		m[i] = 1
	}
	return m
}

/*
func fillSlice(size int, value float64) []float64 {
	m := make([]float64, size)
	for i, _ := range m {
		m[i] = value
	}
	return m
}
func fillSliceInt32(size int, value int32) []float64 {
	m := make([]float64, size)
	for i, _ := range m {
		m[i] = float64(value)
	}
	return m
}
*/

// I'm just starting in Go and found it surprising that it has neither a
// "toFixed" function (as in JavaScript), which would accomplish what you want,
// nor even a "round" function.
// I picked up a one-liner round function from elsewhere, and also made
// toFixed() which depends on round():
// from http://stackoverflow.com/
// How can we truncate float64 type to a particular precision in golang?
// Usage:
// fmt.Println(toFixed(1.2345678, 0))  // 1.0
// fmt.Println(toFixed(1.2345678, 1))  // 1.2
// fmt.Println(toFixed(1.2345678, 2))  // 1.23
// fmt.Println(toFixed(1.2345678, 3))  // 1.235 (rounded up)

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func mkOutputDir() {
	if _, err := os.Stat("ascii"); os.IsNotExist(err) {
		os.Mkdir(outputDir+string(filepath.Separator)+"ascii", 0777)
	}
	if _, err := os.Stat("netcdf"); os.IsNotExist(err) {
		os.Mkdir(outputDir+string(filepath.Separator)+"netcdf", 0777)
	}
	if _, err := os.Stat("odv"); os.IsNotExist(err) {
		os.Mkdir(outputDir+string(filepath.Separator)+"odv", 0777)
	}
}
