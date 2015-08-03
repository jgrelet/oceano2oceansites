// seawater_test
package main

import (
	//	"fmt"
	"testing"
)

const S1, T1, P1, D1 = 35.0, 30.0, 0.0, 1021.7286394941121      // 1021.729
const S2, T2, P2, D2 = 34.67, 2.48, 10035.0, 1070.1356021881304 //1070.136
const S3, T3, P3, C3 = 35.15542101964218, 26.9900 * 1.00024, 27.000, 5.538891

func TestSw_dens(t *testing.T) {

	v := sw_dens(S1, T1, P1)
	if v != D1 {
		t.Error("Expected 1021.729, got ", v)
	}
	v = sw_dens(S2, T2, P2)
	if v != D2 {
		t.Error("Expected 1070.136, got ", v)
	}
}

func TestSw_sigmat(t *testing.T) {

	v := sw_sigmat(S1, T1, P1)
	if v != 21.728639494112144 { // 21.729
		t.Error("Expected 21.729, got ", v)
	}
	v = sw_sigmat(S2, T2, P2)
	if v != 27.667814966349397 { //  27.668
		t.Error("Expected 27.668 , got ", v)
	}
}

func TestSw_sal(t *testing.T) {

	v := sw_sal(C3, T3, P3)
	if v != S3 {
		t.Error("Expected got ", v)
	}
	//	v = sw_sal(S2, T2, P2)
	//	if v != 27.668 {
	//		t.Error("Expected 1021.729, got ", v)
	//	}
}
