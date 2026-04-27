// Package prom is a telegram.Sink that translates parsed DSMR telegrams
// into Prometheus gauge updates.
package prom

import (
	"strconv"

	"github.com/hyperized/dsmr/pkg/cosem"
	"github.com/hyperized/dsmr/pkg/obis"
	"github.com/hyperized/dsmr/pkg/telegram"
)

const deviceTypeUnknown = "unknown"

// Sink updates a registered set of Prometheus metrics from incoming telegrams.
// Method calls are safe for concurrent use because the underlying Prometheus
// gauges are themselves thread safe.
type Sink struct {
	metrics *obis.Metrics
}

// New returns a Sink that writes to the given pre-registered Metrics.
func New(m *obis.Metrics) *Sink {
	return &Sink{metrics: m}
}

// Write applies a telegram's data lines to the underlying Prometheus metrics.
// Always returns nil; missing OBIS codes and unparsable values are skipped
// silently because they are non-fatal for downstream observability.
func (s *Sink) Write(t *telegram.Telegram) error {
	for id, d := range t.All() {
		tv := d.TypedValues()
		if len(tv) == 0 {
			continue
		}

		switch id {
		case "1-3:0.2.8":
			if v, ok := tv[0].(string); ok {
				s.metrics.DSMRInfo().WithLabelValues(v).Set(1)
			}
		case "0-0:96.1.1":
			if v, ok := tv[0].(string); ok {
				s.metrics.EquipInfo().WithLabelValues(v).Set(1)
			}
		case "0-1:24.2.1", "0-2:24.2.1", "0-3:24.2.1", "0-4:24.2.1":
			if len(tv) >= 2 {
				if vs, ok := tv[1].(string); ok {
					if f, err := strconv.ParseFloat(vs, 64); err == nil {
						ch := id[2:3]
						s.metrics.MBus().WithLabelValues(ch, mBusDeviceType(t, ch)).Set(f)
					}
				}
			}
		case "0-1:96.1.0", "0-2:96.1.0", "0-3:96.1.0", "0-4:96.1.0":
			if v, ok := tv[0].(string); ok {
				s.metrics.MBusEquipInfo().WithLabelValues(id[2:3], v).Set(1)
			}
		case "0-1:24.4.0", "0-2:24.4.0", "0-3:24.4.0", "0-4:24.4.0":
			if f, ok := toFloat64(tv[0]); ok {
				s.metrics.MBusValve().WithLabelValues(id[2:3]).Set(f)
			}
		default:
			metricName := d.MetricName()
			if metricName == "" {
				continue
			}
			if gauge, ok := s.metrics.Gauges()[metricName]; ok {
				if f, ok := toFloat64(tv[0]); ok {
					gauge.Set(f)
				}
			}
		}
	}
	return nil
}

// mBusDeviceType returns the raw device-type string from the telegram's
// 0-n:24.1.0 line for channel ch (e.g. "003"), or deviceTypeUnknown if the
// device-type line is absent.
func mBusDeviceType(t *telegram.Telegram, ch string) string {
	d := t.Get("0-" + ch + ":24.1.0")
	if d == nil {
		return deviceTypeUnknown
	}
	return d.Values()[0]
}

func toFloat64(v any) (float64, bool) {
	switch x := v.(type) {
	case *cosem.FloatingPoint:
		return x.Value(), true
	case *cosem.Integer:
		return float64(x.Value()), true
	case string:
		if f, err := strconv.ParseFloat(x, 64); err == nil {
			return f, true
		}
	}
	return 0, false
}
