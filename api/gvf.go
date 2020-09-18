package gvf

import (
	"fmt"
	"math"
)

//FlowType manning / dw
type FlowType string

const (
	//M uses Manning flow equation
	M FlowType = "M"
	//DW uses Darcy Weisbach flow equation
	DW FlowType = "Darcy-Weisbach"
)

const (
	//G is gravity m3/s
	G float64 = 9.8100001
)
const (
	//PI is 3.142
	PI float64 = 3.412
)

//CompDir (proceeding upstream or downstream)
type CompDir string

const (
	//UP is computation proceeding upstream
	UP CompDir = "UP"
	//DN is computation proceeding downstream
	DN CompDir = "DN"
)

//ChannelShape is enum for CIRC/RECT/TRAP
type ChannelShape string

const (
	//CIRC is circualt cross section
	CIRC ChannelShape = "CIRC"
	//RECT is rectangular cross section
	RECT ChannelShape = "RECT"
	// TRAP is trapezoidal cross section
	TRAP ChannelShape = "TRAP"
)

//InputParameters are those entered via stdin by user or file input file
type InputParameters struct {
	//NUMB is 1 for manning and 2 for Darcy-Weisbach
	NUMB int
	//FORM is formula type
	FORM FlowType
	// MN is Manning value (only used when FORM=M)
	MN float64
	//KS is wall roughness in m
	KS float64
	// NO is channel cross section type: 1 CIRCULAR 2 RECTANGULAR 3 TRAPEZOIDAL
	NO int
	//SECT is channel ross section type (As string)
	SECT ChannelShape
	//D is channel diameter (WHEN NO=1 (circular only))
	D float64
	//B is channel width (rectangular channel) or bottom width (trapezoidal channel)
	B float64
	// FI Angle of side to horizontal (trapezoidal channel)
	FI float64
	//YO is depth at Control section
	YO float64
	//FO is channel bedslope
	FO float64
	//DELX is channel step computational length (m)
	DELX float64
	//Q is discharge (Used in Noflow case only)
	Q float64
	//ANS is computational direction (Used in Noflow case only)
	ANS CompDir
	//NS is number of computaional steps (Used in Noflow case only)
	NS int
}

//StateVariables are gvf state parameters
type StateVariables struct {
}

//Gvf is main Gvf server
type Gvf struct {
	Input InputParameters
	State StateVariables
}

// RunNoFlow runs GVF with no flow
func (s *Gvf) RunNoFlow() {
	var Y float64    //Depth(m)
	var A float64    //Cross sectional area(m2)
	var RH float64   //Hyrdaualic Radius(m)
	var TW float64   //Top Width(m)
	var FS float64   //Friction slope
	var DIST float64 //Distance  (m)
	var YC float64   //Critical depth (m)
	NNS := 1         //Number of steps
	YO := s.Input.YO
	if s.Input.ANS == UP {
		s.Input.DELX = -s.Input.DELX
	}
	//**** Computation of Critical Depth *****
	CalCriticalDepth(s.Input, &YC)
	fmt.Printf("CRITICAL DEPTH %f (mm)\n", YC*1000.0)
	//**** Computation of Normal Depth *****
	YN := CalNormalDepth(s.Input)
	fmt.Printf("NORMAL DEPTH %f (mm)\n", YN*1000.0)
	Y = YO
	fmt.Printf("DISTANCE (m)    DEPTH(mm)\n")
	fmt.Printf("%.1f        %.1f\n", DIST, Y*1000.0)
	for {
		GetFlowParams(s.Input, Y, &A, &RH, &TW, &FS)
		A1 := getDYDX(s.Input, A, FS, TW)
		Y = YO + 0.5*A1*s.Input.DELX
		GetFlowParams(s.Input, Y, &A, &RH, &TW, &FS)
		A2 := getDYDX(s.Input, A, FS, TW)
		Y = YO + 0.5*A2*s.Input.DELX
		GetFlowParams(s.Input, Y, &A, &RH, &TW, &FS)
		A3 := getDYDX(s.Input, A, FS, TW)
		Y = YO + A3*s.Input.DELX
		GetFlowParams(s.Input, Y, &A, &RH, &TW, &FS)
		A4 := getDYDX(s.Input, A, FS, TW)
		DELY := (s.Input.DELX / 6.00) * (A1 + 2.0*A2 + 2.0*A3 + A4)
		Y = YO + DELY
		YO = Y
		DIST = DIST + s.Input.DELX
		fmt.Printf("%.1f        %.1f\n", DIST, Y*1000.0)
		NNS++
		if NNS > s.Input.NS {
			break
		}
	}
}

// RunInFlow runs GVF with in flow
func (s *Gvf) RunInFlow() {

}

// RunOutFlow runs GVF with out flow
func (s *Gvf) RunOutFlow() {

}

func getDYDX(in InputParameters, A float64, FS float64, TW float64) float64 {
	DYDX := (in.FO - FS) / (1.0 - in.Q*in.Q*TW/(G*math.Pow(A, 3.0)))
	return DYDX
}
