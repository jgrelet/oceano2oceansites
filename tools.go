package main

import (
	"errors"
	"fmt"
	"math"
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

type Time struct {
	time.Time
	nsec int64 // the number of seconds elapsed since January 1, 1970 UTC
}

// construct time object from a string date
func NewTimeFromString(format, value string) *Time {
	t, _ := time.Parse(format, value)
	return &Time{t, t.Unix()}
}

// construct time object from number of second since 01/01/1970
func NewTimeFromSec(nsec int64) *Time {
	t := time.Unix(nsec, 0)
	return &Time{t, t.Unix()}
}

// construct time object from decimal julian day since 01/01/1950
func NewTimeFromJulian(julian float64) *Time {
	t := time.Unix(int64(julian*86400.0)+oceanSitesToInternal+internalToUnix, 0)
	return &Time{t, t.Unix()}
}

// compute from time object a decimal julian day from 1950
func (t *Time) Time2JulianDec() float64 {
	const DIFF_ORIGIN = 2433283.0 // diff between UNIX DATE and 1950/1/1 00:00:00
	a := int(14-t.Month()) / 12
	y := t.Year() + 4800 - a
	m := int(t.Month()) + 12*a - 3
	julianDay := int(t.Day()) + (153*m+2)/5 + 365*y + y/4
	julianDay = julianDay - y/100 + y/400 - 32045.0 - DIFF_ORIGIN
	//fmt.Println("Julian day:", julianDay)
	return float64(julianDay) + float64(t.Hour())/24 +
		float64(t.Minute())/1440 + float64(t.Second())/86400
}

// compute from time object a decimal julian day from the current year
func (t *Time) JulianDayOfYear() float64 {
	julianDay := t.YearDay()
	return float64(julianDay) + float64(t.Hour())/24 +
		float64(t.Minute())/1440 + float64(t.Second())/86400
}

// convert position "DD MM.SS S" to decimal position
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

// convert  decimal position to string, hemi = 0 for latitude, 1 for longitude
func DecimalPosition2String(position float64, hemi int) string {
	var neg, pos, geo rune

	if hemi == 1 {
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
	h := math.Floor(position)
	m := (position - h) * 60
	if math.Abs(m) > 59 {
		if position > 0 {
			h = h + 1
		} else {
			h = h - 1
		}
		m = 0
	}
	return fmt.Sprintf("%02.0f°%06.3f %c", math.Abs(h), math.Abs(m), geo)
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
