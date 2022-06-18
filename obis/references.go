package obis

import (
	"github.com/hyperized/dsmr/cosem"
	"github.com/shopspring/decimal"
	"time"
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
	// Generic
	"1-3:0.2.8": {
		// 1-3:0.2.8(50)
		Name:        "VersionInformationForP1",
		Identifier:  "1-3:0.2.8",
		Description: "Version information for P1 output",
		Format: cosem.Format{
			Type:   []byte{},
			Tag:    cosem.OctetString,
			Length: 2,
		},
	},
	"0-0:1.0.0": {
		// 0-0:1.0.0(101209113020W)
		Name:        "Timestamp",
		Identifier:  "0-0:1.0.0",
		Description: "Date-time stamp of the P1 message",
		Format: cosem.Format{
			Type: time.Time{},
		},
	},

	"0-0:96.1.1": {
		// 0-0:96.1.1(4B384547303034303436333935353037)
		Name:        "EquipmentIdentifierElectricity",
		Identifier:  "0-0:96.1.1",
		Description: "Equipment identifier (Electricity)",
		Format: cosem.Format{
			Type:   []byte{},
			Tag:    cosem.OctetString,
			Length: 96,
		},
	},

	// Electrical
	"1-0:1.8.1": {
		// 1-0:1.8.1(123456.789*kWh)
		Name:        "MeterReadingElectricityDeliveredToClientTariff1",
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
		Name:        "MeterReadingElectricityDeliveredToClientTariff2",
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
		Name:        "MeterReadingElectricityDeliveredByClientTariff1",
		Identifier:  "1-0:2.8.1",
		Description: "Meter Reading electricity delivered by client (low tariff) in 0,001 kWh",
		Unit:        "kWh",
		Format: cosem.Format{
			Type:            decimal.Decimal{},
			Tag:             cosem.DoubleLongUnsigned,
			Length:          9,
			MinimumDecimals: 3,
			MaximumDecimals: 3,
		},
	},
	"1-0:2.8.2": {
		// 1-0:2.8.2(123456.789*kWh)
		Name:        "MeterReadingElectricityDeliveredByClientTariff2",
		Identifier:  "1-0:2.8.2",
		Description: "Meter Reading electricity delivered by client (Tariff 2) in 0,001 kWh",
		Unit:        "kWh",
		Format: cosem.Format{
			Type:            decimal.Decimal{},
			Tag:             cosem.DoubleLongUnsigned,
			Length:          9,
			MinimumDecimals: 3,
			MaximumDecimals: 3,
		},
	},
	"0-0:96.14.0": {
		// 0-0:96.14.0(0002)
		Name:        "TariffIndicatorElectricity",
		Identifier:  "0-0:96.14.0",
		Description: "Tariff indicator electricity.",
		Format: cosem.Format{
			Type:   []byte{},
			Tag:    cosem.OctetString,
			Length: 4,
		},
	},
	"1-0:1.7.0": {
		// 1-0:1.7.0(01.193*kW)
		Name:        "ActualElectricityPowerDelivered",
		Identifier:  "1-0:1.7.0.255",
		Description: "Actual electricity power delivered (+P) in 1 Watt resolution",
		Unit:        "kW",
		Format: cosem.Format{
			Type:            decimal.Decimal{},
			Tag:             cosem.LongUnsigned,
			MinimumDecimals: 3,
			MaximumDecimals: 3,
			Length:          5,
		},
	},
	"1-0:2.7.0": {
		// 1-0:2.7.0(00.000*kW)
		Name:        "ActualElectricityPowerReceived",
		Identifier:  "1-0:2.7.0.255",
		Description: "Actual electricity power received (-P) in 1 Watt resolution",
		Unit:        "kW",
		Format: cosem.Format{
			Type:            decimal.Decimal{},
			Tag:             cosem.LongUnsigned,
			MinimumDecimals: 3,
			MaximumDecimals: 3,
			Length:          5,
		},
	},

	"0-0:96.7.21": {
		// 0-0:96.7.21(00004)
		Name:        "NumbersOfPowerFailuresInAnyPhase",
		Identifier:  "0-0:96.7.21",
		Description: "Number of power failures in any phase",
		Format: cosem.Format{
			Type:            decimal.Decimal{},
			Tag:             cosem.LongUnsigned,
			MinimumDecimals: 0,
			MaximumDecimals: 0,
			Length:          5,
		},
	},

	"0-0:96.7.9": {
		// 0-0:96.7.9(00002)
		Name:        "NumberOfLongPowerFailuresInAnyPhase",
		Identifier:  "0-0:96.7.9",
		Description: "Number of long power failures in any phase",
		Format: cosem.Format{
			Type:            decimal.Decimal{},
			Tag:             cosem.LongUnsigned,
			MinimumDecimals: 0,
			MaximumDecimals: 0,
			Length:          5,
		},
	},
	// Skipping 1-0:99.97.0.255 power failure event log for now.
	"1-0:32.32.0": {
		// 1-0:32.32.0(00002)
		Name:        "NumberOfVoltageSagsInPhaseL1",
		Identifier:  "1-0:32.32.0",
		Description: "Number of voltage sags in Phase L1",
		Format: cosem.Format{
			Type:            decimal.Decimal{},
			Tag:             cosem.LongUnsigned,
			MinimumDecimals: 0,
			MaximumDecimals: 0,
			Length:          5,
		},
	},
	"1-0:52.32.0": {
		// 1-0:52.32.0(00001)
		Name:        "NumberOfVoltageSagsInPhaseL2",
		Identifier:  "1-0:52.32.0",
		Description: "Number of voltage sags in Phase L2",
		Format: cosem.Format{
			Type:            decimal.Decimal{},
			Tag:             cosem.LongUnsigned,
			MinimumDecimals: 0,
			MaximumDecimals: 0,
			Length:          5,
		},
	},
	"1-0:72.32.0": {
		// 1-0:72.32.0(00000)
		Name:        "NumberOfVoltageSagsInPhaseL3",
		Identifier:  "1-0:72.32.0",
		Description: "Number of voltage sags in Phase L3",
		Format: cosem.Format{
			Type:            decimal.Decimal{},
			Tag:             cosem.LongUnsigned,
			MinimumDecimals: 0,
			MaximumDecimals: 0,
			Length:          5,
		},
	},
	"1-0:32.36.0": {
		// 1-0:32.36.0(00000)
		Name:        "NumberOfVoltageSwellsInPhaseL1",
		Identifier:  "1-0:32.36.0",
		Description: "Number of voltage swells in phase L1",
		Format: cosem.Format{
			Type:            decimal.Decimal{},
			Tag:             cosem.LongUnsigned,
			MinimumDecimals: 0,
			MaximumDecimals: 0,
			Length:          5,
		},
	},
	"1-0:52.36.0": {
		// 1-0:52.36.0(00003)
		Name:        "NumberOfVoltageSwellsInPhaseL2",
		Identifier:  "1-0:52.36.0",
		Description: "Number of voltage swells in phase L2",
		Format: cosem.Format{
			Type:            decimal.Decimal{},
			Tag:             cosem.LongUnsigned,
			MinimumDecimals: 0,
			MaximumDecimals: 0,
			Length:          5,
		},
	},
	"1-0:72.36.0": {
		// 1-0:72.36.0(00000)
		Name:        "NumberOfVoltageSwellsInPhaseL3",
		Identifier:  "1-0:72.36.0",
		Description: "Number of voltage swells in phase L3",
		Format: cosem.Format{
			Type:            decimal.Decimal{},
			Tag:             cosem.LongUnsigned,
			MinimumDecimals: 0,
			MaximumDecimals: 0,
			Length:          5,
		},
	},
}
