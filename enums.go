package solpos

// SPFunctions defines the Solpos functionality as bitmask
type SPFunctions uint32

// HasFlag returns a boolean if the flag is already included
func (f SPFunctions) HasFlag(flag SPFunctions) bool { return f&flag != 0 }

// AddFlag adds a flag
func (f *SPFunctions) AddFlag(flag SPFunctions) { *f |= flag }

// ClearFlag removes a flag
func (f *SPFunctions) ClearFlag(flag SPFunctions) { *f &= ^flag }

// ToggleFlag toggles a flag.
func (f *SPFunctions) ToggleFlag(flag SPFunctions) { *f ^= flag }

//go:generate stringer -type=SPFunctions
const (
	NonFunction SPFunctions = 1 << 0  //
	LDoy        SPFunctions = 1 << 1  // L_DOY = 0x0001;
	LGeom       SPFunctions = 1 << 2  //L_GEOM = 0x0002;
	LZenetr     SPFunctions = 1 << 3  //L_ZENETR = 0x0004;
	LSsha       SPFunctions = 1 << 4  //L_SSHA = 0x0008;
	LSbcf       SPFunctions = 1 << 5  //L_SBCF = 0x0010;
	LTst        SPFunctions = 1 << 6  //L_TST = 0x0020;
	LSrss       SPFunctions = 1 << 7  // L_SRSS = 0x0040;
	LSolazm     SPFunctions = 1 << 8  //L_SOLAZM = 0x0080;
	LRefrac     SPFunctions = 1 << 9  //L_REFRAC = 0x0100;
	LAmass      SPFunctions = 1 << 10 //L_AMASS = 0x0200;
	LPrime      SPFunctions = 1 << 11 //L_PRIME = 0x0400;
	LTilt       SPFunctions = 1 << 12 //L_TILT = 0x0800;
	LEtr        SPFunctions = 1 << 13 //L_ETR = 0x1000;

	LDefault = LGeom | LZenetr | LSsha | LSbcf | LTst | LSrss | LSolazm | LRefrac | LAmass | LPrime | LTilt | LEtr
	LAll     = LDoy | LGeom | LZenetr | LSsha | LSbcf | LTst | LSrss | LSolazm | LRefrac | LAmass | LPrime | LTilt | LEtr // L_ALL = 0xFFFF;
	SDoy     = LDoy
	SGeom    = LGeom | SDoy
	SZenetr  = LZenetr | SGeom
	SSsha    = LSsha | SGeom
	SSbcf    = LSbcf | SSsha
	STst     = LTst | SGeom
	SSrss    = LSrss | SSsha | STst
	SSolazm  = LSolazm | SZenetr
	SRefrac  = LRefrac | SZenetr
	SAmass   = LAmass | SRefrac

	STilt = LTilt | SSolazm | SRefrac
	SEtr  = LEtr | SRefrac
	SAll  = LAll

	SdMask = LZenetr | LSsha | SSbcf | SSolazm
	SlMask = LZenetr | LSsha | SSbcf | SSolazm
	ClMask = LZenetr | LSsha | SSbcf | SSolazm
	CdMask = LZenetr | LSsha | SSbcf
	ChMask = LZenetr
)
