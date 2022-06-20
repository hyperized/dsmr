package telegram

import (
	"github.com/hyperized/dsmr/obis"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"log"
	"strconv"
)

var gauges map[string]prometheus.Gauge

func RegisterGauges() {
	gauges = make(map[string]prometheus.Gauge)
	for _, reference := range obis.References {
		if reference.Metric.Name != "" {
			gauges[reference.Metric.Name] = promauto.NewGauge(prometheus.GaugeOpts{
				Name: reference.Metric.Name,
				Help: reference.Description,
			})
		}
	}

	log.Printf("Gauges: %v\n", gauges)
}

func telegramToPrometheus(t Telegram) {
	for _, data := range t.data {
		if data.Metric.Name != "" {
			m, err := strconv.ParseFloat(data.Value[0], 64) // Only float64 for now

			if err != nil {
				log.Printf("%#v\n", t)
				log.Println(err) // Don't want to do much with this.
			}

			gauges[data.Metric.Name].Set(m)
		}
	}
}
