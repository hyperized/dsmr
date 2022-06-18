package cosem

const (
	NullData           int = 0
	Boolean                = 3
	BitString              = 4
	DoubleLong             = 5
	DoubleLongUnsigned     = 6
	FloatingPo             = 7
	OctetString            = 9
	VisibleString          = 10
	Bcd                    = 13
	Integer                = 15
	Long                   = 16
	Unsigned               = 17
	LongUnsigned           = 18
	Long64                 = 20
	Long64Unsigned         = 21
	Enum                   = 22
	Float32                = 23
	Float64                = 24
)

type Format struct {
	Type            interface{}
	Tag             int
	MinimumDecimals int
	MaximumDecimals int
	Length          int
}
