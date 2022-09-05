package cosem

type Cosem uint8
type Class uint8
type Attribute uint8
type Unit string

type Type interface {
	WithString(value string) error
}

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
