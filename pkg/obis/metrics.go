package obis

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// GaugesMap maps a Prometheus metric name to its Gauge.
type GaugesMap map[string]prometheus.Gauge

// Metrics holds all registered Prometheus metrics for DSMR data.
type Metrics struct {
	gauges    GaugesMap
	mbus      *prometheus.GaugeVec
	mbusValve *prometheus.GaugeVec
	dsmrInfo  *prometheus.GaugeVec
	equipInfo *prometheus.GaugeVec
}

// Gauges returns the per-OBIS scalar gauge map.
func (m *Metrics) Gauges() GaugesMap { return m.gauges }

// MBus returns the GaugeVec for M-Bus channel metered values (labels: channel, device_type).
func (m *Metrics) MBus() *prometheus.GaugeVec { return m.mbus }

// MBusValve returns the GaugeVec for M-Bus valve/switch state (label: channel).
func (m *Metrics) MBusValve() *prometheus.GaugeVec { return m.mbusValve }

// DSMRInfo returns the GaugeVec for DSMR version information labels.
func (m *Metrics) DSMRInfo() *prometheus.GaugeVec { return m.dsmrInfo }

// EquipInfo returns the GaugeVec for electricity equipment identifier labels.
func (m *Metrics) EquipInfo() *prometheus.GaugeVec { return m.equipInfo }

// Register creates and registers all DSMR Prometheus metrics.
func Register() *Metrics {
	gauges := make(GaugesMap)
	for _, ref := range References {
		if ref.metric.name != "" {
			gauges[ref.metric.name] = promauto.NewGauge(prometheus.GaugeOpts{
				Name: ref.metric.name,
				Help: ref.description,
			})
		}
	}

	mbus := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "mbus_last_value",
		Help: "M-Bus last 5-minute metered value",
	}, []string{"channel", "device_type"})

	mbusValve := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "mbus_valve_state",
		Help: "M-Bus valve/switch position (0=disconnected, 1=connected, 2=ready)",
	}, []string{"channel"})

	dsmrInfo := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "dsmr_info",
		Help: "DSMR version information (value always 1; version in label)",
	}, []string{"version"})

	equipInfo := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "electricity_equipment_info",
		Help: "Electricity meter equipment identifier (value always 1; identifier in label)",
	}, []string{"identifier"})

	return &Metrics{
		gauges:    gauges,
		mbus:      mbus,
		mbusValve: mbusValve,
		dsmrInfo:  dsmrInfo,
		equipInfo: equipInfo,
	}
}
