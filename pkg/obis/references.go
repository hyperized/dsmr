// Package obis provides the OBIS reference registry and Prometheus metric
// registration for DSMR P1 telegrams.
//
// The References map is the single source of truth for every supported OBIS
// code: it records the human-readable name, COSEM format specification, unit,
// and (where applicable) the Prometheus metric name.  Call New to look up a
// reference by OBIS identifier string.
package obis

import (
	"github.com/hyperized/dsmr/pkg/cosem"
	"github.com/hyperized/dsmr/pkg/cosem/unit"
)

// Metric holds the Prometheus metric name for an OBIS reference.
// An empty name means no gauge is registered for that OBIS code.
type Metric struct {
	name string
}

// Format describes the COSEM data type and format constraints for an OBIS value.
type Format struct {
	tag             cosem.Cosem
	class           cosem.Class
	attribute       cosem.Attribute
	length          int
	minimumDecimals int
	maximumDecimals int
	formatString    string
}

// Reference is a fully-described OBIS code entry: name, identifier, unit,
// COSEM format, and optional Prometheus metric binding.
type Reference struct {
	name          string
	metric        Metric
	identifier    string
	description   string
	unit          unit.Unit
	format        Format
	polyPhaseOnly bool
}

// New looks up the OBIS reference for identifier (e.g. "1-0:1.8.1").
// Returns nil when the identifier is not in the registry.
func New(reference string) *Reference {
	r, ok := References[reference]
	if !ok {
		return nil
	}
	return &r
}

// Name returns the human-readable name for this OBIS code.
func (r *Reference) Name() string {
	return r.name
}

// Identifier returns the OBIS identifier string (e.g. "1-0:1.8.1").
func (r *Reference) Identifier() string {
	return r.identifier
}

// Unit returns the physical unit for this OBIS value (e.g. "kWh").
func (r *Reference) Unit() unit.Unit {
	return r.unit
}

// Description returns a human-readable description of this OBIS code.
func (r *Reference) Description() string {
	return r.description
}

// PolyPhaseOnly reports whether this OBIS code only appears in poly-phase
// (three-phase) meter configurations.
func (r *Reference) PolyPhaseOnly() bool {
	return r.polyPhaseOnly
}

// MetricName returns the Prometheus metric name, or "" if none is registered.
func (r *Reference) MetricName() string {
	return r.metric.name
}

// Format returns the COSEM format descriptor for this reference.
func (r *Reference) Format() Format {
	return r.format
}

// Tag returns the COSEM data type tag for this format.
func (f Format) Tag() cosem.Cosem {
	return f.tag
}

// Class returns the DLMS/COSEM interface class for this format.
func (f Format) Class() cosem.Class {
	return f.class
}

// Attribute returns the DLMS/COSEM attribute number for this format.
func (f Format) Attribute() cosem.Attribute {
	return f.attribute
}

// FormatString returns the format string (e.g. "F9(3,3)", "TST", "I3").
func (f Format) FormatString() string {
	return f.formatString
}

// Length returns the total field width from the format string.
func (f Format) Length() int {
	return f.length
}

// MinimumDecimals returns the minimum number of digits after the decimal point.
func (f Format) MinimumDecimals() int {
	return f.minimumDecimals
}

// MaximumDecimals returns the maximum number of digits after the decimal point.
func (f Format) MaximumDecimals() int {
	return f.maximumDecimals
}

// References is the registry of all supported OBIS codes, keyed by identifier.
var References = map[string]Reference{
	// Generic
	"1-3:0.2.8": {
		// 1-3:0.2.8(50)
		name:        "DSMRVersion",
		identifier:  "1-3:0.2.8",
		description: "DSMR version information for P1 output",
		format: Format{
			tag:          cosem.OctetString,
			class:        cosem.Data,
			attribute:    cosem.Value,
			length:       2,
			formatString: "S2",
		},
	},
	"0-0:1.0.0": {
		// 0-0:1.0.0(101209113020W)
		name:        "Timestamp",
		identifier:  "0-0:1.0.0",
		description: "Date-time stamp of the P1 message",
		format: Format{
			tag:          cosem.OctetString,
			class:        cosem.Data,
			attribute:    cosem.Value,
			formatString: "TST",
		},
	},
	"0-0:96.1.1": {
		// 0-0:96.1.1(4B384547303034303436333935353037)
		name:        "EquipmentIdentifierElectricity",
		identifier:  "0-0:96.1.1",
		description: "Equipment identifier (Electricity)",
		format: Format{
			tag:          cosem.OctetString,
			class:        cosem.Data,
			attribute:    cosem.Value,
			length:       96,
			formatString: "S0..96",
		},
	},

	// Electrical — energy
	"1-0:1.8.1": {
		// 1-0:1.8.1(123456.789*kWh)
		name: "MeterReadingElectricityDeliveredToClientTariff1",
		metric: Metric{
			name: "electricity_delivered_to_client_tariff1_kwh",
		},
		identifier:  "1-0:1.8.1",
		description: "Meter Reading electricity delivered to client (Tariff 1) in 0,001 kWh",
		unit:        unit.KiloWattHour,
		format: Format{
			tag:             cosem.DoubleLongUnsigned,
			class:           cosem.Register,
			attribute:       cosem.Value,
			length:          9,
			minimumDecimals: 3,
			maximumDecimals: 3,
			formatString:    "F9(3,3)",
		},
	},
	"1-0:1.8.2": {
		// 1-0:1.8.2(123456.789*kWh)
		name: "MeterReadingElectricityDeliveredToClientTariff2",
		metric: Metric{
			name: "electricity_delivered_to_client_tariff2_kwh",
		},
		identifier:  "1-0:1.8.2",
		description: "Meter Reading electricity delivered to client (Tariff 2) in 0,001 kWh",
		unit:        unit.KiloWattHour,
		format: Format{
			tag:             cosem.DoubleLongUnsigned,
			class:           cosem.Register,
			attribute:       cosem.Value,
			length:          9,
			minimumDecimals: 3,
			maximumDecimals: 3,
			formatString:    "F9(3,3)",
		},
	},
	"1-0:2.8.1": {
		// 1-0:2.8.1(123456.789*kWh)
		name: "MeterReadingElectricityDeliveredByClientTariff1",
		metric: Metric{
			name: "electricity_delivered_by_client_tariff1_kwh",
		},
		identifier:  "1-0:2.8.1",
		description: "Meter Reading electricity delivered by client (Tariff 1) in 0,001 kWh",
		unit:        unit.KiloWattHour,
		format: Format{
			tag:             cosem.DoubleLongUnsigned,
			class:           cosem.Register,
			attribute:       cosem.Value,
			length:          9,
			minimumDecimals: 3,
			maximumDecimals: 3,
			formatString:    "F9(3,3)",
		},
	},
	"1-0:2.8.2": {
		// 1-0:2.8.2(123456.789*kWh)
		name: "MeterReadingElectricityDeliveredByClientTariff2",
		metric: Metric{
			name: "electricity_delivered_by_client_tariff2_kwh",
		},
		identifier:  "1-0:2.8.2",
		description: "Meter Reading electricity delivered by client (Tariff 2) in 0,001 kWh",
		unit:        unit.KiloWattHour,
		format: Format{
			tag:             cosem.DoubleLongUnsigned,
			class:           cosem.Register,
			attribute:       cosem.Value,
			length:          9,
			minimumDecimals: 3,
			maximumDecimals: 3,
			formatString:    "F9(3,3)",
		},
	},

	// Electrical — tariff / actual power
	"0-0:96.14.0": {
		// 0-0:96.14.0(0002)
		name: "TariffIndicatorElectricity",
		metric: Metric{
			name: "electricity_tariff",
		},
		identifier:  "0-0:96.14.0",
		description: "Tariff indicator electricity",
		format: Format{
			tag:          cosem.OctetString,
			class:        cosem.Data,
			attribute:    cosem.Value,
			length:       4,
			formatString: "S4",
		},
	},
	"1-0:1.7.0": {
		// 1-0:1.7.0(01.193*kW)
		name: "ActualElectricityPowerDelivered",
		metric: Metric{
			name: "actual_electricity_power_delivered_kw",
		},
		identifier:  "1-0:1.7.0",
		description: "Actual electricity power delivered (+P) in 1 Watt resolution",
		unit:        unit.KiloWatt,
		format: Format{
			tag:             cosem.LongUnsigned,
			class:           cosem.Register,
			attribute:       cosem.Value,
			length:          5,
			minimumDecimals: 3,
			maximumDecimals: 3,
			formatString:    "F5(3,3)",
		},
	},
	"1-0:2.7.0": {
		// 1-0:2.7.0(00.000*kW)
		name: "ActualElectricityPowerReceived",
		metric: Metric{
			name: "actual_electricity_power_received_kw",
		},
		identifier:  "1-0:2.7.0",
		description: "Actual electricity power received (-P) in 1 Watt resolution",
		unit:        unit.KiloWatt,
		format: Format{
			tag:             cosem.LongUnsigned,
			class:           cosem.Register,
			attribute:       cosem.Value,
			length:          5,
			minimumDecimals: 3,
			maximumDecimals: 3,
			formatString:    "F5(3,3)",
		},
	},

	// Electrical — power failures
	"0-0:96.7.21": {
		// 0-0:96.7.21(00004)
		name: "NumberOfPowerFailuresInAnyPhase",
		metric: Metric{
			name: "power_failures_in_any_phase",
		},
		identifier:  "0-0:96.7.21",
		description: "Number of power failures in any phase",
		format: Format{
			tag:          cosem.LongUnsigned,
			class:        cosem.Data,
			attribute:    cosem.Value,
			length:       5,
			formatString: "F5(0,0)",
		},
	},
	"0-0:96.7.9": {
		// 0-0:96.7.9(00002)
		name: "NumberOfLongPowerFailuresInAnyPhase",
		metric: Metric{
			name: "long_power_failures_in_any_phase",
		},
		identifier:  "0-0:96.7.9",
		description: "Number of long power failures in any phase",
		format: Format{
			tag:          cosem.LongUnsigned,
			class:        cosem.Data,
			attribute:    cosem.Value,
			length:       5,
			formatString: "F5(0,0)",
		},
	},
	"1-0:99.97.0": {
		// 1-0:99.97.0(002)(0-0:96.7.19)(210101000001W)(0000000700*s)...
		name:        "PowerFailureEventLog",
		identifier:  "1-0:99.97.0",
		description: "Power failure event log",
		format: Format{
			tag:          cosem.OctetString,
			class:        cosem.GenericProfile,
			attribute:    cosem.Buffer,
			formatString: "TST",
		},
	},

	// Electrical — voltage sags
	"1-0:32.32.0": {
		// 1-0:32.32.0(00002)
		name: "NumberOfVoltageSagsInPhaseL1",
		metric: Metric{
			name: "voltage_sags_phase_l1",
		},
		identifier:  "1-0:32.32.0",
		description: "Number of voltage sags in Phase L1",
		format: Format{
			tag:          cosem.LongUnsigned,
			class:        cosem.Data,
			attribute:    cosem.Value,
			length:       5,
			formatString: "F5(0,0)",
		},
	},
	"1-0:52.32.0": {
		// 1-0:52.32.0(00001)
		name: "NumberOfVoltageSagsInPhaseL2",
		metric: Metric{
			name: "voltage_sags_phase_l2",
		},
		identifier:    "1-0:52.32.0",
		description:   "Number of voltage sags in Phase L2",
		polyPhaseOnly: true,
		format: Format{
			tag:          cosem.LongUnsigned,
			class:        cosem.Data,
			attribute:    cosem.Value,
			length:       5,
			formatString: "F5(0,0)",
		},
	},
	"1-0:72.32.0": {
		// 1-0:72.32.0(00000)
		name: "NumberOfVoltageSagsInPhaseL3",
		metric: Metric{
			name: "voltage_sags_phase_l3",
		},
		identifier:    "1-0:72.32.0",
		description:   "Number of voltage sags in Phase L3",
		polyPhaseOnly: true,
		format: Format{
			tag:          cosem.LongUnsigned,
			class:        cosem.Data,
			attribute:    cosem.Value,
			length:       5,
			formatString: "F5(0,0)",
		},
	},

	// Electrical — voltage swells
	"1-0:32.36.0": {
		// 1-0:32.36.0(00000)
		name: "NumberOfVoltageSwellsInPhaseL1",
		metric: Metric{
			name: "voltage_swells_phase_l1",
		},
		identifier:  "1-0:32.36.0",
		description: "Number of voltage swells in Phase L1",
		format: Format{
			tag:          cosem.LongUnsigned,
			class:        cosem.Data,
			attribute:    cosem.Value,
			length:       5,
			formatString: "F5(0,0)",
		},
	},
	"1-0:52.36.0": {
		// 1-0:52.36.0(00003)
		name: "NumberOfVoltageSwellsInPhaseL2",
		metric: Metric{
			name: "voltage_swells_phase_l2",
		},
		identifier:    "1-0:52.36.0",
		description:   "Number of voltage swells in Phase L2",
		polyPhaseOnly: true,
		format: Format{
			tag:          cosem.LongUnsigned,
			class:        cosem.Data,
			attribute:    cosem.Value,
			length:       5,
			formatString: "F5(0,0)",
		},
	},
	"1-0:72.36.0": {
		// 1-0:72.36.0(00000)
		name: "NumberOfVoltageSwellsInPhaseL3",
		metric: Metric{
			name: "voltage_swells_phase_l3",
		},
		identifier:    "1-0:72.36.0",
		description:   "Number of voltage swells in Phase L3",
		polyPhaseOnly: true,
		format: Format{
			tag:          cosem.LongUnsigned,
			class:        cosem.Data,
			attribute:    cosem.Value,
			length:       5,
			formatString: "F5(0,0)",
		},
	},

	// Text message
	"0-0:96.13.0": {
		// 0-0:96.13.0(303132333435363738393A3B3C3D3E3F)
		name:        "TextMessage",
		identifier:  "0-0:96.13.0",
		description: "Text message max 2048 characters",
		format: Format{
			tag:          cosem.OctetString,
			class:        cosem.Data,
			attribute:    cosem.Value,
			length:       2048,
			formatString: "S0..2048",
		},
	},

	// Electrical — instantaneous voltage
	"1-0:32.7.0": {
		// 1-0:32.7.0(220.1*V)
		name: "InstantVoltageL1",
		metric: Metric{
			name: "instant_voltage_l1",
		},
		identifier:  "1-0:32.7.0",
		description: "Instantaneous voltage L1 in V resolution",
		unit:        unit.Volt,
		format: Format{
			tag:             cosem.LongUnsigned,
			class:           cosem.Register,
			attribute:       cosem.Value,
			length:          4,
			minimumDecimals: 1,
			maximumDecimals: 1,
			formatString:    "F4(1,1)",
		},
	},
	"1-0:52.7.0": {
		// 1-0:52.7.0(220.2*V)
		name: "InstantVoltageL2",
		metric: Metric{
			name: "instant_voltage_l2",
		},
		identifier:    "1-0:52.7.0",
		description:   "Instantaneous voltage L2 in V resolution",
		unit:          unit.Volt,
		polyPhaseOnly: true,
		format: Format{
			tag:             cosem.LongUnsigned,
			class:           cosem.Register,
			attribute:       cosem.Value,
			length:          4,
			minimumDecimals: 1,
			maximumDecimals: 1,
			formatString:    "F4(1,1)",
		},
	},
	"1-0:72.7.0": {
		// 1-0:72.7.0(220.3*V)
		name: "InstantVoltageL3",
		metric: Metric{
			name: "instant_voltage_l3",
		},
		identifier:    "1-0:72.7.0",
		description:   "Instantaneous voltage L3 in V resolution",
		unit:          unit.Volt,
		polyPhaseOnly: true,
		format: Format{
			tag:             cosem.LongUnsigned,
			class:           cosem.Register,
			attribute:       cosem.Value,
			length:          4,
			minimumDecimals: 1,
			maximumDecimals: 1,
			formatString:    "F4(1,1)",
		},
	},

	// Electrical — instantaneous current
	"1-0:31.7.0": {
		// 1-0:31.7.0(001*A)
		name: "InstantCurrentL1",
		metric: Metric{
			name: "instant_current_l1",
		},
		identifier:  "1-0:31.7.0",
		description: "Instantaneous current L1 in A resolution",
		unit:        unit.Ampere,
		format: Format{
			tag:          cosem.LongUnsigned,
			class:        cosem.Register,
			attribute:    cosem.Value,
			length:       3,
			formatString: "F3(0,0)",
		},
	},
	"1-0:51.7.0": {
		// 1-0:51.7.0(002*A)
		name: "InstantCurrentL2",
		metric: Metric{
			name: "instant_current_l2",
		},
		identifier:    "1-0:51.7.0",
		description:   "Instantaneous current L2 in A resolution",
		unit:          unit.Ampere,
		polyPhaseOnly: true,
		format: Format{
			tag:          cosem.LongUnsigned,
			class:        cosem.Register,
			attribute:    cosem.Value,
			length:       3,
			formatString: "F3(0,0)",
		},
	},
	"1-0:71.7.0": {
		// 1-0:71.7.0(003*A)
		name: "InstantCurrentL3",
		metric: Metric{
			name: "instant_current_l3",
		},
		identifier:    "1-0:71.7.0",
		description:   "Instantaneous current L3 in A resolution",
		unit:          unit.Ampere,
		polyPhaseOnly: true,
		format: Format{
			tag:          cosem.LongUnsigned,
			class:        cosem.Register,
			attribute:    cosem.Value,
			length:       3,
			formatString: "F3(0,0)",
		},
	},

	// Electrical — instantaneous active power delivered (+P)
	"1-0:21.7.0": {
		// 1-0:21.7.0(01.111*kW)
		name: "InstantActivePowerDeliveredL1",
		metric: Metric{
			name: "instant_active_power_delivered_l1",
		},
		identifier:  "1-0:21.7.0",
		description: "Instantaneous active power L1 (+P) in W resolution",
		unit:        unit.KiloWatt,
		format: Format{
			tag:             cosem.LongUnsigned,
			class:           cosem.Register,
			attribute:       cosem.Value,
			length:          5,
			minimumDecimals: 3,
			maximumDecimals: 3,
			formatString:    "F5(3,3)",
		},
	},
	"1-0:41.7.0": {
		// 1-0:41.7.0(02.222*kW)
		name: "InstantActivePowerDeliveredL2",
		metric: Metric{
			name: "instant_active_power_delivered_l2",
		},
		identifier:    "1-0:41.7.0",
		description:   "Instantaneous active power L2 (+P) in W resolution",
		unit:          unit.KiloWatt,
		polyPhaseOnly: true,
		format: Format{
			tag:             cosem.LongUnsigned,
			class:           cosem.Register,
			attribute:       cosem.Value,
			length:          5,
			minimumDecimals: 3,
			maximumDecimals: 3,
			formatString:    "F5(3,3)",
		},
	},
	"1-0:61.7.0": {
		// 1-0:61.7.0(03.333*kW)
		name: "InstantActivePowerDeliveredL3",
		metric: Metric{
			name: "instant_active_power_delivered_l3",
		},
		identifier:    "1-0:61.7.0",
		description:   "Instantaneous active power L3 (+P) in W resolution",
		unit:          unit.KiloWatt,
		polyPhaseOnly: true,
		format: Format{
			tag:             cosem.LongUnsigned,
			class:           cosem.Register,
			attribute:       cosem.Value,
			length:          5,
			minimumDecimals: 3,
			maximumDecimals: 3,
			formatString:    "F5(3,3)",
		},
	},

	// Electrical — instantaneous active power received (-P)
	"1-0:22.7.0": {
		// 1-0:22.7.0(04.444*kW)
		name: "InstantActivePowerReceivedL1",
		metric: Metric{
			name: "instant_active_power_received_l1",
		},
		identifier:  "1-0:22.7.0",
		description: "Instantaneous active power L1 (-P) in W resolution",
		unit:        unit.KiloWatt,
		format: Format{
			tag:             cosem.LongUnsigned,
			class:           cosem.Register,
			attribute:       cosem.Value,
			length:          5,
			minimumDecimals: 3,
			maximumDecimals: 3,
			formatString:    "F5(3,3)",
		},
	},
	"1-0:42.7.0": {
		// 1-0:42.7.0(05.555*kW)
		name: "InstantActivePowerReceivedL2",
		metric: Metric{
			name: "instant_active_power_received_l2",
		},
		identifier:    "1-0:42.7.0",
		description:   "Instantaneous active power L2 (-P) in W resolution",
		unit:          unit.KiloWatt,
		polyPhaseOnly: true,
		format: Format{
			tag:             cosem.LongUnsigned,
			class:           cosem.Register,
			attribute:       cosem.Value,
			length:          5,
			minimumDecimals: 3,
			maximumDecimals: 3,
			formatString:    "F5(3,3)",
		},
	},
	"1-0:62.7.0": {
		// 1-0:62.7.0(06.666*kW)
		name: "InstantActivePowerReceivedL3",
		metric: Metric{
			name: "instant_active_power_received_l3",
		},
		identifier:    "1-0:62.7.0",
		description:   "Instantaneous active power L3 (-P) in W resolution",
		unit:          unit.KiloWatt,
		polyPhaseOnly: true,
		format: Format{
			tag:             cosem.LongUnsigned,
			class:           cosem.Register,
			attribute:       cosem.Value,
			length:          5,
			minimumDecimals: 3,
			maximumDecimals: 3,
			formatString:    "F5(3,3)",
		},
	},

	// M-Bus — channel 1
	"0-1:24.1.0": {
		// 0-1:24.1.0(003)
		name:        "MBusDeviceTypeChannel1",
		identifier:  "0-1:24.1.0",
		description: "M-Bus Device-Type channel 1",
		format: Format{
			tag:          cosem.LongUnsigned,
			class:        cosem.Data,
			attribute:    cosem.Value,
			length:       3,
			formatString: "I3",
		},
	},
	"0-1:96.1.0": {
		// 0-1:96.1.0(3232323241424344313233343536373839)
		name:        "MBusEquipmentIdentifierChannel1",
		identifier:  "0-1:96.1.0",
		description: "M-Bus equipment identifier channel 1",
		format: Format{
			tag:          cosem.OctetString,
			class:        cosem.Data,
			attribute:    cosem.Value,
			length:       96,
			formatString: "S0..96",
		},
	},
	"0-1:24.2.1": {
		// 0-1:24.2.1(150117180000W)(00423.422*m3)
		name:        "MBusLastValueChannel1",
		identifier:  "0-1:24.2.1",
		description: "M-Bus last 5-minute value channel 1",
		format: Format{
			class:        cosem.ExtendedRegister,
			attribute:    cosem.Value,
			formatString: "TST",
		},
	},

	// M-Bus — channel 2
	"0-2:24.1.0": {
		name:        "MBusDeviceTypeChannel2",
		identifier:  "0-2:24.1.0",
		description: "M-Bus Device-Type channel 2",
		format: Format{
			tag:          cosem.LongUnsigned,
			class:        cosem.Data,
			attribute:    cosem.Value,
			length:       3,
			formatString: "I3",
		},
	},
	"0-2:96.1.0": {
		name:        "MBusEquipmentIdentifierChannel2",
		identifier:  "0-2:96.1.0",
		description: "M-Bus equipment identifier channel 2",
		format: Format{
			tag:          cosem.OctetString,
			class:        cosem.Data,
			attribute:    cosem.Value,
			length:       96,
			formatString: "S0..96",
		},
	},
	"0-2:24.2.1": {
		name:        "MBusLastValueChannel2",
		identifier:  "0-2:24.2.1",
		description: "M-Bus last 5-minute value channel 2",
		format: Format{
			class:        cosem.ExtendedRegister,
			attribute:    cosem.Value,
			formatString: "TST",
		},
	},

	// M-Bus — channel 3
	"0-3:24.1.0": {
		name:        "MBusDeviceTypeChannel3",
		identifier:  "0-3:24.1.0",
		description: "M-Bus Device-Type channel 3",
		format: Format{
			tag:          cosem.LongUnsigned,
			class:        cosem.Data,
			attribute:    cosem.Value,
			length:       3,
			formatString: "I3",
		},
	},
	"0-3:96.1.0": {
		name:        "MBusEquipmentIdentifierChannel3",
		identifier:  "0-3:96.1.0",
		description: "M-Bus equipment identifier channel 3",
		format: Format{
			tag:          cosem.OctetString,
			class:        cosem.Data,
			attribute:    cosem.Value,
			length:       96,
			formatString: "S0..96",
		},
	},
	"0-3:24.2.1": {
		name:        "MBusLastValueChannel3",
		identifier:  "0-3:24.2.1",
		description: "M-Bus last 5-minute value channel 3",
		format: Format{
			class:        cosem.ExtendedRegister,
			attribute:    cosem.Value,
			formatString: "TST",
		},
	},

	// M-Bus — channel 4
	"0-4:24.1.0": {
		name:        "MBusDeviceTypeChannel4",
		identifier:  "0-4:24.1.0",
		description: "M-Bus Device-Type channel 4",
		format: Format{
			tag:          cosem.LongUnsigned,
			class:        cosem.Data,
			attribute:    cosem.Value,
			length:       3,
			formatString: "I3",
		},
	},
	"0-4:96.1.0": {
		name:        "MBusEquipmentIdentifierChannel4",
		identifier:  "0-4:96.1.0",
		description: "M-Bus equipment identifier channel 4",
		format: Format{
			tag:          cosem.OctetString,
			class:        cosem.Data,
			attribute:    cosem.Value,
			length:       96,
			formatString: "S0..96",
		},
	},
	"0-4:24.2.1": {
		name:        "MBusLastValueChannel4",
		identifier:  "0-4:24.2.1",
		description: "M-Bus last 5-minute value channel 4",
		format: Format{
			class:        cosem.ExtendedRegister,
			attribute:    cosem.Value,
			formatString: "TST",
		},
	},
}
