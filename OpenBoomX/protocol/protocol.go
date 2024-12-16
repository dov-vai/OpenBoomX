package protocol

const UBoomXName = "EarFun UBOOM X"
const UBoomXName2 = "EarFun Audio"
const UBoomXOUI = "F8:AB:E5"
const UBoomXOUI2 = "C7:AB:E5"

// Oluv's EQ Modes
var EQModes = map[string]string{
	"studio":   "efb046010102fe",
	"indoor":   "efb046010203fe",
	"indoor+":  "efb046010304fe",
	"outdoor":  "efb046010405fe",
	"outdoor+": "efb046010506fe",
	"boom":     "efb046010607fe",
	"ground":   "efb046010708fe",
}

// Light Actions
var LightActions = map[string]string{
	"default": "efb095040000000000fe",
	"off":     "efb095040100000000fe",
}

// Shutdown Timeout Modes
var ShutdownTimeouts = map[string]string{
	"5m":   "efb075010102fe",
	"10m":  "efb075010203fe",
	"30m":  "efb075010304fe",
	"60m":  "efb075010405fe",
	"90m":  "efb075010506fe",
	"120m": "efb075010607fe",
	"no":   "efb07501ff00fe",
}

const SpeakerPowerOff = "efb025010102fe"

var BluetoothPairing = map[string]string{
	"on":  "efb035010102fe",
	"off": "efb035010001fe",
}

// Beep Volume Levels
var BeepVolumes = map[int]string{
	0:   "efb065010102fe",
	25:  "efb065010203fe",
	50:  "efb065010304fe",
	75:  "efb065010405fe",
	100: "efb065010506fe",
}

// EQ Band Values
const (
	MaxBandValue = 120 // +10 dB is the max
	MinBandValue = 0   // -10 dB is the min
)

const RfcommChannel = 2

const BatteryLevelRequest = "efa0140000fe"
