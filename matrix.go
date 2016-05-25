package main

import (
	"fmt"
	"log"
)

// map used a matrix for each parameters
type matrix map[string]interface{}

func (m matrix) NewMatrix(name string, fillValue float64, size ...int) *matrix {
	switch len(size) {
	case 1:
		s := make([]float64, size[0])
		for i := 0; i < size[0]; i++ {
			s[i] = fillValue
		}
		m[name] = s
	case 2:
		s := make([][]float64, size[0])
		for i := 0; i < size[0]; i++ {
			s[i] = make([]float64, size[1])
			for j := 0; j < size[1]; j++ {
				s[i][j] = fillValue
			}
		}
		m[name] = s
	default:
		log.Fatal("matrix.NewMatrix: Too many dimension")
	}

	return &m
}

func (m matrix) set(key string, value float64, size ...int) {
	switch len(size) {
	case 1:
		s := m[key].([]float64)
		if size[0] < len(s) {
			s[size[0]] = value
		} else {
			log.Fatalf("matrix.set: index out of range: %d", size[0])
		}
	case 2:
		s := m[key].([][]float64)
		s[size[0]][size[1]] = value
	default:
		log.Fatal("matrix.set: Too many dimension")
	}
}

func (m matrix) get(key string, size ...int) (value interface{}) {
	switch m[key].(type) {
	case []float64:
		switch len(size) {
		case 0:
			value = m[key].([]float64)
		case 1:
			s := m[key].([]float64)
			if size[0] < len(s) {
				value = s[size[0]]
			} else {
				log.Fatalf("matrix.set: index out of range: %d", size[0])
			}
		case 2:
			s := m[key].([][]float64)
			value = s[size[0]][size[1]]
		default:
			log.Fatal("matrix.get: Too many dimension")
		}
	case [][]float64:
		switch len(size) {
		case 0:
			value = m[key].([][]float64)
		case 1:
			s := m[key].([][]float64)
			if size[0] < len(s) {
				value = s[size[0]]
			} else {
				log.Fatalf("matrix.set: index out of range: %d", size[0])
			}
		case 2:
			s := m[key].([][]float64)
			value = s[size[0]][size[1]]
		default:
			log.Fatal("matrix.get: Too many dimension")
		}
	default:
		log.Fatalf("matrix.get: unknow type %T", m[key])
	}
	return value
}

func (m matrix) getDim(key string) (x int, y int) {
	switch m[key].(type) {
	case []float64:
		x = len(m[key].([]float64))
		y = 0
	case [][]float64:
		r := m[key].([][]float64)
		x = len(r)
		y = len(r[0])
	default:
		log.Fatalf("matrix.getToArray: unknow type %T", m[key])
	}
	return int(x), int(y)
}

func (m matrix) printInfo(key string) (s string) {
	switch m[key].(type) {
	case []float64:
		x := len(m[key].([]float64))
		s = fmt.Sprintf("writing %s: %d\n", key, x)
	case [][]float64:
		r := m[key].([][]float64)
		x := len(r)
		y := len(r[0])
		s = fmt.Sprintf("writing %s: %d x %d\n", key, x, y)
	default:
		log.Fatalf("matrix.getInfo: unknow type %T", m[key])
	}
	return s
}

func (m matrix) isMatrix(key string) bool {
	switch m[key].(type) {
	case [][]float64:
		return true
	default:
		return false
	}
}

func (m matrix) flatten(key string) (value []float64) {
	switch m[key].(type) {
	case []float64:
		value = m[key].([]float64)
	case [][]float64:
		r := m[key].([][]float64)
		ht := len(r)
		wd := len(r[0])
		v := make([]float64, 0, ht*wd)
		for _, row := range r {
			v = append(v, row...)
		}
		value = v
	default:
		log.Fatalf("matrix.getToArray: unknow type %T", m[key])
	}
	return value
}

func Matrix2int32(ar []float64) []int32 {
	newar := make([]int32, len(ar))
	for i, v := range ar {
		newar[i] = int32(v)
	}
	return newar
}

func Matrix2float32(ar []float64) []float32 {
	newar := make([]float32, len(ar))
	for i, v := range ar {
		newar[i] = float32(v)
	}
	return newar
}
