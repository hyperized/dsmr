package obis

import (
	"github.com/hyperized/dsmr/cosem"
	"github.com/shopspring/decimal"
	"time"
)

type Reference struct {
	Name        string
	Metric      Metric
	Identifier  string
	Description string
	Unit        string
	Format      cosem.Format
	Parser      interface{}
}

type Metric struct {
	Name string
}

var References = map[string]Reference{
	// Generic
	"1-3:0.2.8": {
		// 1-3:0.2.8(50)
		Name: "DSMRVersion",
		Metric: Metric{
			Name: "dsmr_version",
		},
		Identifier:  "1-3:0.2.8",
		Description: "DSMR version information for P1 output",
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
		Name: "EquipmentIdentifierElectricity",
		Metric: Metric{
			Name: "electricity_equipment_identifier",
		},
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
		Name: "MeterReadingElectricityDeliveredToClientTariff1",
		Metric: Metric{
			Name: "electricity_delivered_to_client_tariff1_kwh",
		},
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
		Name: "MeterReadingElectricityDeliveredToClientTariff2",
		Metric: Metric{
			Name: "electricity_delivered_to_client_tariff2_kwh",
		},
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
		Name: "MeterReadingElectricityDeliveredByClientTariff1",
		Metric: Metric{
			Name: "electricity_delivered_by_client_tariff1_kwh",
		},
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
		Name: "MeterReadingElectricityDeliveredByClientTariff2",
		Metric: Metric{
			Name: "electricity_delivered_by_client_tariff2_kwh",
		},
		Identifier:  "1-0:2.8.2",
		Description: "Meter Reading electricity delivered by client (high 2) in 0,001 kWh",
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
		Name: "TariffIndicatorElectricity",
		Metric: Metric{
			Name: "electricity_tariff",
		},
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
		Name: "ActualElectricityPowerDelivered",
		Metric: Metric{
			Name: "actual_electricity_power_delivered_kw",
		},
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
		Name: "ActualElectricityPowerReceived",
		Metric: Metric{
			Name: "actual_electricity_power_received_kw",
		},
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
		Name: "NumberOfPowerFailuresInAnyPhase",
		Metric: Metric{
			Name: "power_failures_in_any_phase",
		},
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
		Name: "NumberOfLongPowerFailuresInAnyPhase",
		Metric: Metric{
			Name: "long_power_failures_in_any_phase",
		},
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
	// "1-0:99.97.0.255": {}, Skipping power failure event log for now.
	"1-0:32.32.0": {
		// 1-0:32.32.0(00002)
		Name: "NumberOfVoltageSagsInPhaseL1",
		Metric: Metric{
			Name: "voltage_sags_phase_l1",
		},
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
		Name: "NumberOfVoltageSagsInPhaseL2",
		Metric: Metric{
			Name: "voltage_sags_phase_l2",
		},
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
		Name: "NumberOfVoltageSagsInPhaseL3",
		Metric: Metric{
			Name: "voltage_sags_phase_l3",
		},
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
		Name: "NumberOfVoltageSwellsInPhaseL1",
		Metric: Metric{
			Name: "voltage_swells_phase_l1",
		},
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
		Name: "NumberOfVoltageSwellsInPhaseL2",
		Metric: Metric{
			Name: "voltage_swells_phase_l2",
		},
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
		Name: "NumberOfVoltageSwellsInPhaseL3",
		Metric: Metric{
			Name: "voltage_swells_phase_l3",
		},
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
	// "0-0:96.13.0": {}, Skip for now
	"1-0:32.7.0": {
		// 1-0:32.7.0(220.1*V)
		Name: "InstantVoltageL1",
		Metric: Metric{
			Name: "instant_voltage_l1",
		},
		Identifier:  "1-0:32.7.0",
		Description: "Instantaneous voltage L1 in V resolution",
		Format: cosem.Format{
			Type:            decimal.Decimal{},
			Tag:             cosem.LongUnsigned,
			MinimumDecimals: 1,
			MaximumDecimals: 1,
			Length:          4,
		},
		Unit: "V",
	},
	"1-0:52.7.0": {
		// 1-0:52.7.0(220.2*V)
		Name: "InstantVoltageL2",
		Metric: Metric{
			Name: "instant_voltage_l2",
		},
		Identifier:  "1-0:52.7.0",
		Description: "Instantaneous voltage L2 in V resolution",
		Format: cosem.Format{
			Type:            decimal.Decimal{},
			Tag:             cosem.LongUnsigned,
			MinimumDecimals: 1,
			MaximumDecimals: 1,
			Length:          4,
		},
		Unit: "V",
	},
	"1-0:72.7.0": {
		// 1-0:72.7.0(220.3*V)
		Name: "InstantVoltageL3",
		Metric: Metric{
			Name: "instant_voltage_l3",
		},
		Identifier:  "1-0:72.7.0",
		Description: "Instantaneous voltage L3 in V resolution",
		Format: cosem.Format{
			Type:            decimal.Decimal{},
			Tag:             cosem.LongUnsigned,
			MinimumDecimals: 1,
			MaximumDecimals: 1,
			Length:          4,
		},
		Unit: "V",
	},
	"1-0:31.7.0": {
		// 1-0:31.7.0(001*A)
		Name: "InstantCurrentL1",
		Metric: Metric{
			Name: "instant_current_l1",
		},
		Identifier:  "1-0:31.7.0",
		Description: "Instantaneous current L1 in A resolution",
		Format: cosem.Format{
			Type:            decimal.Decimal{},
			Tag:             cosem.LongUnsigned,
			MinimumDecimals: 0,
			MaximumDecimals: 0,
			Length:          3,
		},
		Unit: "A",
	},
	"1-0:51.7.0": {
		// 1-0:51.7.0(002*A)
		Name: "InstantCurrentL2",
		Metric: Metric{
			Name: "instant_current_l2",
		},
		Identifier:  "1-0:51.7.0",
		Description: "Instantaneous current L2 in A resolution",
		Format: cosem.Format{
			Type:            decimal.Decimal{},
			Tag:             cosem.LongUnsigned,
			MinimumDecimals: 0,
			MaximumDecimals: 0,
			Length:          3,
		},
		Unit: "A",
	},
	"1-0:71.7.0": {
		// 1-0:71.7.0(003*A)
		Name: "InstantCurrentL3",
		Metric: Metric{
			Name: "instant_current_l3",
		},
		Identifier:  "1-0:71.7.0",
		Description: "Instantaneous current L3 in A resolution",
		Format: cosem.Format{
			Type:            decimal.Decimal{},
			Tag:             cosem.LongUnsigned,
			MinimumDecimals: 0,
			MaximumDecimals: 0,
			Length:          3,
		},
		Unit: "A",
	},
	"1-0:21.7.0": {
		// 1-0:21.7.0(01.111*kW)
		Name: "InstantActivePowerDeliveredL1",
		Metric: Metric{
			Name: "instant_active_power_delivered_l1",
		},
		Identifier:  "1-0:21.7.0",
		Description: "Instantaneous active power L1 (+P) in W resolution",
		Format: cosem.Format{
			Type:            decimal.Decimal{},
			Tag:             cosem.LongUnsigned,
			MinimumDecimals: 3,
			MaximumDecimals: 3,
			Length:          5,
		},
		Unit: "kW",
	},
	"1-0:41.7.0": {
		// 1-0:41.7.0(02.222*kW)
		Name: "InstantActivePowerDeliveredL2",
		Metric: Metric{
			Name: "instant_active_power_delivered_l2",
		},
		Identifier:  "1-0:41.7.0",
		Description: "Instantaneous active power L2 (+P) in W resolution",
		Format: cosem.Format{
			Type:            decimal.Decimal{},
			Tag:             cosem.LongUnsigned,
			MinimumDecimals: 3,
			MaximumDecimals: 3,
			Length:          5,
		},
		Unit: "kW",
	},
	"1-0:61.7.0": {
		// 1-0:61.7.0(03.333*kW)
		Name: "InstantActivePowerDeliveredL3",
		Metric: Metric{
			Name: "instant_active_power_delivered_l3",
		},
		Identifier:  "1-0:61.7.0",
		Description: "Instantaneous active power L3 (+P) in W resolution",
		Format: cosem.Format{
			Type:            decimal.Decimal{},
			Tag:             cosem.LongUnsigned,
			MinimumDecimals: 3,
			MaximumDecimals: 3,
			Length:          5,
		},
		Unit: "kW",
	},
	"1-0:22.7.0": {
		// 1-0:22.7.0(04.444*kW)
		Name: "InstantCurrentReceivedL1",
		Metric: Metric{
			Name: "instant_active_power_received_l1",
		},
		Identifier:  "1-0:22.7.0",
		Description: "Instantaneous active power L1 (-P) in W resolution",
		Format: cosem.Format{
			Type:            decimal.Decimal{},
			Tag:             cosem.LongUnsigned,
			MinimumDecimals: 0,
			MaximumDecimals: 0,
			Length:          3,
		},
		Unit: "A",
	},
	"1-0:42.7.0": {
		// 1-0:42.7.0(05.555*kW)
		Name: "InstantCurrentReceivedL2",
		Metric: Metric{
			Name: "instant_active_power_received_l2",
		},
		Identifier:  "1-0:42.7.0",
		Description: "Instantaneous active power L2 (-P) in W resolution",
		Format: cosem.Format{
			Type:            decimal.Decimal{},
			Tag:             cosem.LongUnsigned,
			MinimumDecimals: 0,
			MaximumDecimals: 0,
			Length:          3,
		},
		Unit: "A",
	},
	"1-0:62.7.0": {
		// 1-0:62.7.0(06.666*kW)
		Name: "InstantCurrentReceivedL3",
		Metric: Metric{
			Name: "instant_active_power_received_l3",
		},
		Identifier:  "1-0:62.7.0",
		Description: "Instantaneous active power L3 (-P) in W resolution",
		Format: cosem.Format{
			Type:            decimal.Decimal{},
			Tag:             cosem.LongUnsigned,
			MinimumDecimals: 0,
			MaximumDecimals: 0,
			Length:          3,
		},
		Unit: "A",
	},

	// Text messages

	// Gas data

	// Thermal data

	// Water data

	// M-Bus data
}
