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

type Time struct {
	time.Time
}

func New(format, value string) *Time {
	t, _ := time.Parse(format, value)
	return &Time{t}
}

// convert date/time to a decimal julian day with origin 1950/1/1
func (t *Time) Date2JulianDec() float64 {
	const DIFF_ORIGIN = 2433283.0 // diff between UNIX DATE and 1950/1/1 00:00:00
	a := int(14-t.Month()) / 12
	y := t.Year() + 4800 - a
	m := int(t.Month()) + 12*a - 3
	julianDay := int(t.Day()) + (153*m+2)/5 + 365*y + y/4
	julianDay = julianDay - y/100 + y/400 - 32045.0 - DIFF_ORIGIN
	fmt.Println("Julian day:", julianDay)
	return float64(julianDay) + float64(t.Hour())/24 + float64(t.Minute())/1440 + float64(t.Second())/86400
}

func (t *Time) JulianDay() float64 {
	julianDay := t.YearDay()
	return float64(julianDay) + float64(t.Hour())/24 + float64(t.Minute())/1440 + float64(t.Second())/86400
}

// convert position "DD MM.SS S" to decimal position
func PositionDeci(pos string) (float64, error) {

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
