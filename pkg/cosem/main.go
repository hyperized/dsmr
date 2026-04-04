// Package cosem defines the COSEM (Companion Specification for Energy Metering)
// type system used by DSMR P1 telegrams. It provides typed constants for data
// tags, DLMS/COSEM classes, and attribute identifiers as specified in
// IEC 62056-62 and the DSMR P1 companion standard.
package cosem

// Cosem represents a DLMS/COSEM data type tag (IEC 62056-62 §4).
type Cosem uint8

// Class represents a DLMS/COSEM interface class identifier.
type Class uint8

// Attribute represents a DLMS/COSEM attribute number within a class.
type Attribute uint8

// Unit is a string-typed SI unit label used in OBIS references (e.g. "kWh").
type Unit string

// DLMS/COSEM class identifiers and attribute numbers used by DSMR references.
const (
	Value              Attribute = 1
	Buffer                       = 2
	Data               Class     = 1
	Register                     = 3
	ExtendedRegister             = 4
	GenericProfile               = 7
	MBusClient                   = 72
	NullData           Cosem     = 0
	Boolean                      = 3
	BitString                    = 4
	DoubleLong                   = 5
	DoubleLongUnsigned           = 6
	FloatingPo                   = 7
	OctetString                  = 9
	VisibleString                = 10
	Bcd                          = 13
	Integer                      = 15
	Long                         = 16
	Unsigned                     = 17
	LongUnsigned                 = 18
	Long64                       = 20
	Long64Unsigned               = 21
	Enum                         = 22
	Float32                      = 23
	Float64                      = 24
)
