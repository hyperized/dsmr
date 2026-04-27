// Package cosem defines the COSEM (Companion Specification for Energy Metering)
// type system used by DSMR P1 telegrams: data tag values, DLMS/COSEM interface
// classes, attribute identifiers, SI units, and the concrete value types
// (Timestamp, Integer, FloatingPoint, String, OctetString, Enum) that the
// parser produces.
package cosem

// Cosem is a DLMS/COSEM data type tag (IEC 62056-62 §4).
type Cosem uint8

// Class is a DLMS/COSEM interface class identifier.
type Class uint8

// Attribute is a DLMS/COSEM attribute number within a class.
type Attribute uint8

// Unit is an SI unit label used in OBIS references.
type Unit string

// SI unit constants used in DSMR OBIS references.
const (
	KiloWattHour Unit = "kWh"
	KiloWatt     Unit = "kW"
	Volt         Unit = "V"
	Ampere       Unit = "A"
	CubicMeter   Unit = "m3"
	GigaJoule    Unit = "GJ"
	Second       Unit = "s"
	None         Unit = ""
)

// DLMS/COSEM data type tag values (IEC 62056-62). All constants share the
// "Tag" prefix to avoid colliding with the COSEM value types defined in this
// package (Integer, Enum, OctetString, FloatingPoint).
const (
	TagNullData           Cosem = 0
	TagBoolean            Cosem = 3
	TagBitString          Cosem = 4
	TagDoubleLong         Cosem = 5
	TagDoubleLongUnsigned Cosem = 6
	TagFloatingPoint      Cosem = 7
	TagOctetString        Cosem = 9
	TagVisibleString      Cosem = 10
	TagBcd                Cosem = 13
	TagInteger            Cosem = 15
	TagLong               Cosem = 16
	TagUnsigned           Cosem = 17
	TagLongUnsigned       Cosem = 18
	TagLong64             Cosem = 20
	TagLong64Unsigned     Cosem = 21
	TagEnum               Cosem = 22
	TagFloat32            Cosem = 23
	TagFloat64            Cosem = 24
)

// DLMS/COSEM attribute numbers used by DSMR references.
const (
	Value  Attribute = 1
	Buffer Attribute = 2
)

// DLMS/COSEM interface class identifiers used by DSMR references.
const (
	Data             Class = 1
	Register         Class = 3
	ExtendedRegister Class = 4
	GenericProfile   Class = 7
	MBusClient       Class = 72
)

// DateTimeFormat is the Go time-parse layout for the YYMMDDhhmmss portion of
// a DSMR timestamp (without the trailing DST flag character).
const DateTimeFormat = "060102150405"
