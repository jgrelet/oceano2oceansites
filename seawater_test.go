// seawater_test
package main

import (
	//	"fmt"
	"testing"
)

const S1, T1, P1 = 35.0, 30.0, 0.0
const D1, Sigma_t1, Sigma_teta1, Svel1, Potential_Temp1, Specific_anomaly1, Depth1 = 1021.729, 21.729, 21.729, 1545.595, 30.000, 6.071e-06, 0.000
const Lat = 4.0

const S2, T2, P2 = 34.67, 2.48, 10035.0
const D2, Sigma_t2, Sigma_teta2, Svel2, Potential_Temp2, Specific_anomaly2, Depth2 = 1070.136, 27.668, 27.764, 1633.179, 1.242, 8.352e-07, 9758.558

const S3, T3, P3, C3 = 35.1554, 26.9900 * 1.00024, 27.000, 5.538891
const S4, T4, P4, C4 = 35.7918, 18.1986 * 1.00024, 71.000, 4.705818

func TestSw_dens(t *testing.T) {

	v := sw_dens(S1, T1, P1)
	v = toFixed(v, 3)
	if v != D1 {
		t.Errorf("Expected %f, got %f", D1, v)
	}
	v = sw_dens(S2, T2, P2)
	v = toFixed(v, 3)
	if v != D2 {
		t.Errorf("Expected %f, got %f", D2, v)
	}
}

func TestSw_sal(t *testing.T) {

	v := sw_sal(C3, T3, P3)
	v = toFixed(v, 4)
	if v != S3 {
		t.Errorf("Expected %f, got %f", S3, v)
	}
	v = sw_sal(C4, T4, P4)
	v = toFixed(v, 4)
	if v != S4 {
		t.Errorf("Expected %f, got %f", S4, v)
	}
}

func TestSw_sigmat(t *testing.T) {

	v := sw_sigmat(S1, T1, P1)
	v = toFixed(v, 3)
	if v != Sigma_t1 { // 21.729
		t.Errorf("Expected %f, got %f", Sigma_t1, v)
	}
	v = sw_sigmat(S2, T2, P2)
	v = toFixed(v, 3)
	if v != Sigma_t2 { //  27.668
		t.Errorf("Expected %f, got %f", Sigma_t2, v)
	}
}

func TestSw_sigmateta(t *testing.T) {

	v := sw_sigmateta(S1, T1, P1)
	v = toFixed(v, 3)
	if v != Sigma_teta1 {
		t.Errorf("Expected %6.3f, got %6.3f", Sigma_teta1, v)
	}
	v = sw_sigmateta(S2, T2, P2)
	v = toFixed(v, 3)
	if v != Sigma_teta2 {
		t.Errorf("Expected %6.3f, got %6.3f", Sigma_teta2, v)
	}
}

func TestSw_svel(t *testing.T) {

	v := sw_svel(S1, T1, P1)
	v = toFixed(v, 3)
	if v != Svel1 {
		t.Errorf("Expected %6.3f, got %6.3f", Svel1, v)
	}
	v = sw_svel(S2, T2, P2)
	v = toFixed(v, 3)
	if v != Svel2 {
		t.Errorf("Expected %6.3f, got %6.3f", Svel2, v)
	}
}

func TestSw_ptmp(t *testing.T) {

	v := sw_ptmp(S1, T1, P1, 0)
	v = toFixed(v, 3)
	if v != Potential_Temp1 {
		t.Errorf("Expected %6.3f, got %6.3f", Potential_Temp1, v)
	}
	v = sw_ptmp(S2, T2, P2, 0)
	v = toFixed(v, 3)
	if v != Potential_Temp2 {
		t.Errorf("Expected %6.3f, got %6.3f", Potential_Temp2, v)
	}
}

func TestSw_svan(t *testing.T) {

	v := sw_svan(S1, T1, P1)
	v = toFixed(v, 9)
	if v != Specific_anomaly1 {
		t.Errorf("Expected %6.6g, got %6.6g", Specific_anomaly1, v)
	}
	v = sw_svan(S2, T2, P2)
	v = toFixed(v, 10)
	if v != Specific_anomaly2 {
		t.Errorf("Expected %6.6g, got %6.6g", Specific_anomaly2, v)
	}
}

func TestSw_dpth(t *testing.T) {

	v := sw_dpth(P1, Lat)
	v = toFixed(v, 3)
	if v != Depth1 {
		t.Errorf("Expected %6.3f, got %6.3f", Depth1, v)
	}
	v = sw_dpth(P2, Lat)
	v = toFixed(v, 3)
	if v != Depth2 {
		t.Errorf("Expected %6.3f, got %6.3f", Depth2, v)
	}
}
