package main

import (
	"errors"
	"fmt"
	gvf "gvf/api"
	"os"
)

// PROGRAM GVF: Analyses steady gradually varied flow
func main() {
	print("\nProgram GVF\n\n")
	print("THIS PROGRAM ANALYSES THREE CATEGORIES OF STEADY GRADUALLY\n")
	print("VARIED OPEN CHANNEL FLOW; IT CATERS FOR RECTANGULAR, TRAPEZOIDAL\n")
	print("AND CIRCULAR CHANNEL SECTIONS AND OFFERS A CHOICE BETWEEN THE\n")
	print("MANNING AND DARCY-WEISBACH FLOW EQUATIONS. THE ANALYSES USES A\n")
	print("FOURTH ORDER RUNGE-KUTTA NUMERICAL COMPUTATIONAL SCHEMA IN THE\n")
	print("SOLUTION OF THE RELEVANT WATER SURFACE SLOPE EQUATION.\n")
	print("\nTHE THREE FLOW CATEGORIES ARE:\n")
	print("    (1) GVF WITHOUT LATERAL INFLOW OR OUTFLOW\n")
	print("    (2) GVF WITH LATERAL INFLOW\n")
	print("    (3) GVF WITH LATERAL OUTFLOW\n")
	print("\n")
	file := openInputFile()
	defer closeInputFile(file)
	sim := gvf.Gvf{}
	simType := 0
	readInput("ENTER NUMBER OF YOUR CHOICE", file, &simType)
	if simType == 1 {
		readNoFlowParams(file, &sim.Input)
		sim.RunNoFlow()
	} else if simType == 2 {
		readInFlowParams(file, &sim.Input)
		sim.RunInFlow()
	} else if simType == 3 {
		readOutFlowParams(file, &sim.Input)
		sim.RunOutFlow()
	} else {
		check(errors.New("Unrecognised sim type"))
	}

}
func readNoFlowParams(file *os.File, in *gvf.InputParameters) {
	//**********FLOW CATEGORY 1 **********
	print("ANALYSIS OF GVF IN CHANNELS WITHOUT LATERAL INFLOW/OUTFLOW\n")
	print("The program computes the flow depth at specified intervals along\n")
	print("the channel, starting from a control point at which the depth\n")
	print("is specified. The program outputs distance from the control\n")
	print("point and corresponding flow depth. Note the distances measured\n")
	print("upstream from the control point are printed as negative values.\n")
	print("\nDATA ENTRY:\n")
	readFlowType(file, in)
	readChannelData(file, in)
	readInput("DEPTH AT CONTROL SECTION", file, &in.YO)
	in.YO = mmToM(in.YO)
	readInput("ENTER CHANNEL BED SLOPE", file, &in.FO)
	readInput("ENTER DISCHARGE(m**3/s)", file, &in.Q)
	readInput("IS COMPUTATION PROCEEDING UPSTREAM (UP) OR DOWNSTREAM (DN) FROM CONTROL SECTION?", file, &in.ANS)
	if in.ANS != gvf.UP && in.ANS != gvf.DN {
		check(errors.New("Invalid entry. Should be 'UP' or 'DN'"))
	}
	readInput("ENTER CHANNEL STEP COMPUTATIONAL LENGTH (m)", file, &in.DELX)
	readInput("ENTER NUMBER OF COMPUTATION STEPS", file, &in.NS)
	print("...data input complete;computation now in progress...\n")
}
func readInFlowParams(file *os.File, in *gvf.InputParameters) {

}
func readOutFlowParams(file *os.File, in *gvf.InputParameters) {

}
func readFlowType(file *os.File, in *gvf.InputParameters) {
	readInput("DO YOU WISH TO USE 1 MANNING OR 2 DARCY-WEISBACH ?", file, &in.NUMB)
	if in.NUMB == 1 {
		in.FORM = gvf.M
		readInput("MANNING N-VALUE", file, &in.MN)
	} else if in.NUMB == 2 {
		in.FORM = gvf.DW
		readInput("Wall Roughness(mm)", file, &in.KS)
		//Convert to m
		in.KS = mmToM(in.KS)
	} else {
		check(errors.New("Error entering formula type"))
	}
}
func readChannelData(file *os.File, in *gvf.InputParameters) {
	print("ENTER CHANNEL DATA:\n")
	readInput("IS SECTION 1 CIRCULAR 2 RECTANGULAR 3 TRAPEZOIDAL ?", file, &in.NO)
	if in.NO == 1 {
		in.SECT = gvf.CIRC
		readInput("DIAMETER (mm)", file, &in.D)
		//Convert to m
		in.D = mmToM(in.D)
	} else if in.NO == 2 {
		in.SECT = gvf.RECT
		readInput("CHANNEL WIDTH (mm)", file, &in.B)
		//Convert to m
		in.B = mmToM(in.B)
	} else if in.NO == 3 {
		in.SECT = gvf.TRAP
		readInput("BOTTOM WIDTH (mm)", file, &in.B)
		//Convert to m
		in.B = mmToM(in.B)
		readInput("ANGLE OF SIDE TO HORL (deg)", file, &in.FI)
	} else {
		check(errors.New("Error entering channel data"))
	}
}
func readInput(question string, file *os.File, input interface{}) {
	print(question + "\n")
	switch input.(type) {
	default:
		check(errors.New("Unsupported type"))
	case *uint:
		fmt.Fscanf(file, "%u", input)
	case *uint64:
		fmt.Fscanf(file, "%u", input)
	case *gvf.CompDir:
		fmt.Fscanf(file, "%s", input)
	case *string:
		fmt.Fscanf(file, "%s", input)
	case *int:
		fmt.Fscanf(file, "%d", input)
	case *float32:
		fmt.Fscanf(file, "%f", input)
	case *float64:
		fmt.Fscanf(file, "%f", input)
	}
}
func openInputFile() *os.File {
	file := os.Stdin
	args := os.Args[1:] //Args missing program name
	if len(args) > 0 {
		if !fileExists(args[0]) {
			fmt.Printf("File '%s' does not exist\n", args[0])
			os.Exit(-1)
		}
		f, err := os.Open(args[0])
		check(err)
		file = f
	}
	return file
}
func closeInputFile(file *os.File) {
	args := os.Args[1:]
	if len(args) > 0 && file != nil {
		file.Close()
	}
}
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
func mmToM(in float64) float64 {
	const THOUSAND float64 = 1000.0
	return (in / THOUSAND)
}
func check(e error) {
	if e != nil {
		panic(e)
	}
}
