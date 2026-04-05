// Package unit defines the SI unit constants used in DSMR OBIS references.
package unit

// Unit represents a physical measurement unit used in DSMR OBIS references.
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
