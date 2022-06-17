package obis

import (
	"github.com/hyperized/dsmr/cosem"
	"github.com/shopspring/decimal"
)

type Reference struct {
	Name        string
	Identifier  string
	Description string
	Unit        string
	Format      cosem.Format
	Parser      interface{}
}

var References = map[string]Reference{
	"1-3:0.2.8": {
		// 1-3:0.2.8(50)
		Identifier:  "1-3:0.2.8",
		Description: "Version information for P1 output",
		Format: cosem.Format{
			Type:   []byte{},
			Tag:    cosem.OctetString,
			Length: 2,
		},
	},

	// Electrical
	"1-0:1.8.1": {
		// 1-0:1.8.1(123456.789*kWh)
		Identifier:  "1-0:1.8.1",
		Description: "Meter Reading electricity delivered to client (Tariff 1) in 0,001 kWh",
		Unit:        "kWh",
		Format: cosem.Format{
			Type:            decimal.Decimal{},
			Tag:             cosem.DoubleLongUnsigned,
			Length:          9,
			MinimumDecimals: 3,
			MaximumDecimals: 3,
		},
	},
	"1-0:1.8.2": {
		// 1-0:1.8.2(123456.789*kWh)
		Identifier:  "1-0:1.8.1",
		Description: "Meter Reading electricity delivered to client (Tariff 2) in 0,001 kWh",
		Unit:        "kWh",
		Format: cosem.Format{
			Type:            decimal.Decimal{},
			Tag:             cosem.DoubleLongUnsigned,
			Length:          9,
			MinimumDecimals: 3,
			MaximumDecimals: 3,
		},
	},
	"1-0:2.8.1": {
		// 1-0:2.8.1(123456.789*kWh)
		Identifier:  "1-0:1.8.1",
		Description: "Meter Reading electricity delivered to client (Tariff 2) in 0,001 kWh",
		Unit:        "kWh",
		Format: cosem.Format{
			Type:            decimal.Decimal{},
			Tag:             cosem.DoubleLongUnsigned,
			Length:          9,
			MinimumDecimals: 3,
			MaximumDecimals: 3,
		},
	},
	"1-0:32.32.0": {
		// 1-0:32.32.0(00002)
		Name:        "NumberOfVoltageSagsInPhaseL1",
		Identifier:  "1-0:32.32.0",
		Description: "Number of voltage sags in Phase L1",
		Format: cosem.Format{
			Type:   decimal.Decimal{},
			Tag:    cosem.LongUnsigned,
			Length: 5,
		},
	},
	"1-0:52.32.0": {
		// 1-0:52.32.0(00001)
		Name:        "NumberOfVoltageSagsInPhaseL2",
		Identifier:  "1-0:52.32.0",
		Description: "Number of voltage sags in Phase L2",
		Format: cosem.Format{
			Type:   decimal.Decimal{},
			Tag:    cosem.LongUnsigned,
			Length: 5,
		},
	},
	// 1-0:72.32.0(00000)
	"1-0:72.32.0": {
		Name:        "NumberOfVoltageSagsInPhaseL3",
		Identifier:  "1-0:72.32.0",
		Description: "Number of voltage sags in Phase L3",
		Format: cosem.Format{
			Type:   decimal.Decimal{},
			Tag:    cosem.LongUnsigned,
			Length: 5,
		},
	},
}
