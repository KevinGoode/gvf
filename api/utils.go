package gvf

// Utility functions that were GOSUB routines in original code
import (
	"errors"
	"math"
)

//GetFlowParams is wrapper function around Manning/Darcy Weisbach
func GetFlowParams(in InputParameters, Y float64, A *float64, RH *float64, TW *float64, FS *float64)  {
	if in.FORM == M {
		CalFlowParamsManning(in, Y, A,RH,TW,FS)
	}
	CalFlowParamsDarcyWeisbach(in, Y, A,RH,TW,FS)
}

//CalFlowParamsManning calcs cross sectional area, hydraulic radius ,surface width, and 
//friction slope using manning equation (GOSUB 2040 in original code)
func CalFlowParamsManning(in InputParameters, Y float64, A *float64, RH *float64, TW *float64, FS *float64){
	CalFlowLengths(in, Y, A, RH, TW)
	*FS = math.Pow(in.MN*in.Q/(*A), 2.0) * math.Pow((*RH), -1.33)
}

//CalFlowParamsDarcyWeisbach calcs cross sectional area, hydraulic radius ,surface width, and 
//friction slope using DarcyWeisbach equation (GOSUB 2080 in original code)
func CalFlowParamsDarcyWeisbach(in InputParameters, Y float64, A *float64, RH *float64, TW *float64, FS *float64) {
	CalFlowLengths(in, Y, A, RH, TW)
	V := in.Q / (*A)
	F := CalFrictionFactor(in, *RH, in.KS, V)
	*FS = F * in.Q * in.Q / (8.00 * G * (*A) * (*A) * (*RH))
}

//CalFlowLengths calcs cross sectional area, hydraulic radius and surface width(GOSUB 2130 in original code)
func CalFlowLengths(in InputParameters, Y float64, A *float64, RH *float64, TW *float64) {
	if in.SECT == CIRC {
		HI := PI
		LO := 0.00
		TH := 0.00
		for {
			//Solve equation (R-Y)/Y  = COS(TH) by linear bisection
			//Can't we just use ACOS??
			TH = (HI + LO) / 2.0
			XR := 1.0 - (2*Y)/in.D - math.Cos(TH)
			if XR < 0.0 {
				LO = TH
			} else {
				HI = TH
			}
			Z := (HI + LO) / 2.0
			if math.Abs(Z-TH) < 0.001 {
				break
			}
		}
		//P is wetted perimiter
		P := in.D * TH
		*A = 0.25 * in.D * in.D * (TH - 0.5*math.Sin(2.0*TH))
		*RH = *A / P
		*TW = in.D * math.Sin(TH)

	} else if in.SECT == RECT {
		*A = in.B * Y
		*RH = in.B * Y / (in.B + 2.0*Y)
		*TW = in.B
	} else if in.SECT == TRAP {
		FIRADS := in.FI * PI / 180.0
		*A = Y * (in.B + Y/(math.Tan(FIRADS)))
		*RH = *A / (in.B + (2.0*Y)/math.Sin(FIRADS))
		*TW = in.B + (2.0*Y)/math.Tan(FIRADS)
	}

}

//CalNormalDepth calcs normal depth (GOSUB 2660 in original code)
func CalNormalDepth(in InputParameters) float64 {
	FSH := 0.0
	FSR := math.Pow(in.FO, 0.5)
	TH := 0.00
	HI := 40.00
	LO := 0.001
	Y := 0.00
	if in.SECT == CIRC {
		HI = PI
		LO = 0.001
		//Use linear bisection to solve for angle TH given Q
		for {
			TH = (HI + LO) / 2.0
			A := 0.25 * in.D * in.D * (TH - 0.5*math.Sin(2.0*TH))
			P := in.D * TH
			RH := A / P
			FSH = 0.00
			if in.FORM == M {
				FSH = in.MN * in.Q / (A * math.Pow(RH, 0.67))
			} else if in.FORM == DW {
				V := in.Q / A
				F := CalFrictionFactor(in, RH, in.KS, V)
				FSH = math.Pow(F/(8.0*G*RH), 0.5)*V
			} else {
				panic(errors.New("Unexpected flow type"))
			}
			WW := FSH - FSR
			if WW > 0 {
				LO = TH
			} else {
				HI = TH
			}
			if math.Abs(FSH-FSR)/FSR < 0.001 {
				break
			}
		}
		//Now got TH , set Y
		Y = 0.5 * in.D * (1.0 - math.Cos(TH))

	} else if in.SECT == RECT {
		//Solve for Y by linear bisection
		for {
			Y = (HI + LO) / 2.0
			A := in.B * Y
			RH := A / (in.B + 2.0*Y)
			FSH = 0.00
			if in.FORM == M {
				FSH = in.MN * in.Q / (A * math.Pow(RH, 0.67))
			} else if in.FORM == DW {
				V := in.Q / A
				F := CalFrictionFactor(in, RH, in.KS, V)
				FSH = math.Pow(F/(8.0*G*RH), 0.5*V)
			}
			WW := FSH - FSR
			if WW > 0 {
				LO = Y
			} else {
				HI = Y
			}
			Z := (HI + LO) / 2.0
			if math.Abs(Z-Y) < 0.002 {
				break
			}
		}

	} else if in.SECT == TRAP {
		//Solve for Y by linear bisection
		for {
			FIRADS := in.FI * PI / 180.0
			Y = (HI + LO) / 2.0
			A := Y * (in.B + Y/(math.Tan(FIRADS)))
			RH := A / (in.B + (2.0*Y)/math.Sin(FIRADS))
			FSH = 0.00
			if in.FORM == M {
				FSH = in.MN * in.Q / (A * math.Pow(RH, 0.67))
			} else if in.FORM == DW {
				V := in.Q / A
				F := CalFrictionFactor(in, RH, in.KS, V)
				FSH = math.Pow(F/(8.0*G*RH), 0.5*V)
			}
			WW := FSH - FSR
			if WW > 0 {
				LO = Y
			} else {
				HI = Y
			}
			Z := (HI + LO) / 2.0
			if math.Abs(Z-Y) < 0.002 {
				break
			}
		}
	}
	return Y
}

//CalFrictionFactor calcs fricion factor(GOSUB 2500 in original code)
//TODO K can be removed it is inside in
func CalFrictionFactor(in InputParameters, RH float64, K float64, V float64) float64 {
	KVISC := 1.307E-06
	UPV := 0.5
	LOV := 0.0
	F := (UPV + LOV) / 2.0
	for {
		YY := 1.0 / (math.Sqrt(F))
		X := K/(14.8*RH) + (2.51*KVISC)/(4.0*RH*V*math.Sqrt(F))
		W := YY + 0.88*math.Log(X)
		if W < 0 {
			UPV = F
		} else {
			LOV = F
		}
		Z := (UPV + LOV) / 2.0
		E := math.Abs((Z - F) / F)
		F = Z
		if E < 0.005 {
			break
		}
	}
	return F
}

//CalCriticalDepth calcs critical depth (GOSUB 2290 in original code)
func CalCriticalDepth(in InputParameters, YC *float64) {
	//Critical depth is governed by equation 7.22 in Casey
	//Q^2/(g.A^3)(dA/dy)=1
	//
	//For RECT:
	//Q^2.B/(g.A^3) = 1
	//Q^2.B/(g.B^3.Y^3)=1
	//Y = (Q^2/(B^2.G))^0.3333
	//
	//For TRAP:
	//A = Y.(B+Y/tan(fi))
	//dA/dY = B + 1/tan(fi)
	//Solve non-linear critical depth equation using linear bisection
	//
	//For CIRC:
	// dA/dy = dA/Dtheta / dY/dTHETA
	// A =0.25*D^2*(THETA-0.5*SIN(2.0*THETA))
	// DA/Dtheta = 0.25D^2 (1-COS(2.0*THETA))
	// NOTE TRIG EQUATION:COS(2.0*THETA) = (1-2SIN^2THETA)
	//SO DA.DTheta = 0.5D^2.SIN^2(THETA)
	// Y = D/2.0 (1 - COS(THETA))
	//DY/DTheta = D.SIN(THETA)/2.0
	//So dA/dy = 0.5D^2.SIN^2(THETA)   / D.SIN(THETA)/2.0
	//So dA/dy =  DSIN(THETA) (IE VARIABLE BBB BELOW)
	// Solve non-linear critical depth equation using linear bisection
	if in.SECT == CIRC {
		HI := PI
		LO := 0.00
		for {
			TH := (HI + LO) / 2.00
			A := 0.25 * in.D * in.D * (TH - 0.5*math.Sin(2.0*TH))
			BBB := in.D * math.Sin(TH)
			XR := math.Pow(A, 3.000)/BBB - in.Q*in.Q/G
			if XR < 0 {
				LO = TH
			} else {
				HI = TH
			}
			Z := (HI + LO) / 2.000
			if math.Abs(Z-TH) < 0.001 {
				*YC = 0.5 * in.D * (1.0 - math.Cos(TH))
				break
			}
		}
	} else if in.SECT == RECT {
		*YC = math.Pow(math.Pow(in.Q, 2.00)/(math.Pow(in.B, 2.00)*G), 0.33333)
	} else if in.SECT == TRAP {
		XC := math.Pow(in.Q, 2.00) / G
		TANFI := math.Tan(in.FI * PI / 180.00)
		HI := 20.000
		LO := 0.000
		for {
			YY := (HI + LO) / 2.00
			BB := YY / TANFI
			A := (in.B + BB) * YY
			BBB := in.B + 2.00*BB // IE dA/dy
			XR := math.Pow(A, 3.00)/BBB - XC
			if XR < 0 {
				LO = YY
			} else {
				HI = YY
			}
			Z := (HI + LO) / 2.000
			if math.Abs(Z-YY) < 0.001 {
				*YC = YY
				break
			}
		}
	}
}
