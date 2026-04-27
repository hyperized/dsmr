package prom_test

import (
	"strings"
	"testing"

	"github.com/hyperized/dsmr/pkg/obis"
	"github.com/hyperized/dsmr/pkg/sink/prom"
	"github.com/hyperized/dsmr/pkg/telegram"
	"github.com/prometheus/client_golang/prometheus"
)

func newTestSink() *prom.Sink {
	return prom.New(obis.Register(prometheus.NewRegistry()))
}

// TestParseStreamWithSink exercises ParseStream → prom.Sink.Write using a
// telegram that covers every branch of Sink.Write:
//   - "1-3:0.2.8"          DSMRInfo GaugeVec branch
//   - "0-0:96.1.1"         EquipInfo GaugeVec branch
//   - "0-1:24.1.0"         device-type lookup for MBus label
//   - "0-1:96.1.0"         MBusEquipInfo GaugeVec branch
//   - "0-1:24.2.1"         MBus GaugeVec branch (channel + device_type labels)
//   - "0-1:24.4.0"         MBusValve GaugeVec branch
//   - "1-0:1.8.1"          default branch, FloatingPoint gauge (toFloat64 float path)
//   - "0-0:96.14.0"        default branch, string gauge   (toFloat64 string path)
//   - "0-0:1.0.0"          default branch, no metric name → continue
//   - "1-0:99.97.0"        GenericProfile → TypedValues nil → continue
func TestParseStreamWithSink(_ *testing.T) {
	const le = "\r\n"
	content := "/ISk5\\2MT382-1000" + le +
		le +
		"1-3:0.2.8(50)" + le +
		"0-0:96.1.1(4B384547303034303436333935353037)" + le +
		"0-1:24.1.0(003)" + le +
		"0-1:96.1.0(3232323241424344313233343536373839)" + le +
		"0-1:24.2.1(101209112500W)(12785.123*m3)" + le +
		"0-1:24.4.0(1)" + le +
		"1-0:1.8.1(123456.789*kWh)" + le +
		"0-0:96.14.0(0002)" + le +
		"0-0:1.0.0(101209113020W)" + le +
		"1-0:99.97.0(2)(0-0:96.7.19)" + le +
		"!0000\n"

	p := telegram.NewParser(strings.NewReader(content),
		telegram.WithLineEnding(le),
		telegram.WithCRCValidation(false),
		telegram.WithSink(newTestSink()),
	)
	p.ParseStream()
}
