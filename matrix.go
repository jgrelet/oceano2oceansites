package main

import (
	"fmt"
	"log"
)

// map used a matrix for each parameters
type Matrix map[string]interface{}

func (mp Matrix) NewMatrix(name string, fillValue float64, size ...int) *Matrix {
	switch len(size) {
	case 1:
		mt := make([]float64, size[0])
		for i := 0; i < size[0]; i++ {
			mt[i] = fillValue
		}
		mp[name] = mt
	case 2:
		mt := make([][]float64, size[0])
		for i := 0; i < size[0]; i++ {
			mt[i] = make([]float64, size[1])
			for j := 0; j < size[1]; j++ {
				mt[i][j] = fillValue
			}
		}
		mp[name] = mt
	default:
		log.Fatal("Matrix.NewMatrix: Too many dimension")
	}

	return &mp
}

func (mp Matrix) set(key string, value float64, size ...int) {
	switch len(size) {
	case 1:
		m := mp[key].([]float64)
		if size[0] < len(m) {
			m[size[0]] = value
		} else {
			log.Fatalf("Matrix.set: index out of range: %d", size[0])
		}
	case 2:
		m := mp[key].([][]float64)
		m[size[0]][size[1]] = value
	default:
		log.Fatal("Matrix.set: Too many dimension")
	}
}

func (mp Matrix) get(key string, size ...int) (value interface{}) {
	switch mp[key].(type) {
	case []float64:
		switch len(size) {
		case 0:
			value = mp[key].([]float64)
		case 1:
			m := mp[key].([]float64)
			if size[0] < len(m) {
				value = m[size[0]]
			} else {
				log.Fatalf("Matrix.set: index out of range: %d", size[0])
			}
		case 2:
			m := mp[key].([][]float64)
			value = m[size[0]][size[1]]
		default:
			log.Fatal("Matrix.get: Too many dimension")
		}
	case [][]float64:
		switch len(size) {
		case 0:
			value = mp[key].([][]float64)
		case 1:
			m := mp[key].([][]float64)
			if size[0] < len(m) {
				value = m[size[0]]
			} else {
				log.Fatalf("Matrix.set: index out of range: %d", size[0])
			}
		case 2:
			m := mp[key].([][]float64)
			value = m[size[0]][size[1]]
		default:
			log.Fatal("Matrix.get: Too many dimension")
		}
	default:
		log.Fatalf("Matrix.get: unknow type %T", mp[key])
	}
	return value
}

func (mp Matrix) getDim(key string) (x int, y int) {
	switch mp[key].(type) {
	case []float64:
		x = len(mp[key].([]float64))
		y = 0
	case [][]float64:
		r := mp[key].([][]float64)
		x = len(r)
		y = len(r[0])
	default:
		log.Fatalf("Matrix.getToArray: unknow type %T", mp[key])
	}
	return int(x), int(y)
}

func (mp Matrix) printInfo(key string) (s string) {
	switch mp[key].(type) {
	case []float64:
		x := len(mp[key].([]float64))
		s = fmt.Sprintf("writing %s: %d\n", key, x)
	case [][]float64:
		r := mp[key].([][]float64)
		x := len(r)
		y := len(r[0])
		s = fmt.Sprintf("writing %s: %d x %d\n", key, x, y)
	default:
		log.Fatalf("Matrix.getInfo: unknow type %T", mp[key])
	}
	return s
}

func (mp Matrix) isMatrix(key string) bool {
	switch mp[key].(type) {
	case [][]float64:
		return true
	default:
		return false
	}
}

func (mp Matrix) flatten(key string) (value []float64) {
	switch mp[key].(type) {
	case []float64:
		value = mp[key].([]float64)
	case [][]float64:
		r := mp[key].([][]float64)
		ht := len(r)
		wd := len(r[0])
		v := make([]float64, 0, ht*wd)
		for _, row := range r {
			v = append(v, row...)
		}
		value = v
	default:
		log.Fatalf("Matrix.getToArray: unknow type %T", mp[key])
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
