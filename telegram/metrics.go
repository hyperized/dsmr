package telegram

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"log"
	"strconv"
)

type Metric struct {
	Name  string
	Gauge prometheus.Gauge
}

type Metrics []Metric

var metrics = Metrics{
	Metric{
		Name: "MeterReadingElectricityDeliveredToClientTariff1",
		Gauge: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "electricity_delivered_to_client_tariff1_kwh",
			Help: "Meter Reading electricity delivered to client (Tariff 1) in 0,001 kWh",
		}),
	},
	Metric{
		Name: "MeterReadingElectricityDeliveredToClientTariff2",
		Gauge: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "electricity_delivered_to_client_tariff2_kwh",
			Help: "Meter Reading electricity delivered to client (Tariff 2) in 0,001 kWh",
		}),
	},
	Metric{
		Name: "ActualElectricityPowerDelivered",
		Gauge: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "actual_electricity_power_delivered_kw",
			Help: "Actual electricity power delivered (+P) in 1 Watt resolution.",
		}),
	},
}

func telegramToPrometheus(t Telegram) {
	for _, metric := range metrics {
		m, err := strconv.ParseFloat(t.data[metric.Name].Value()[0], 64) // Only float64 for now
		if err != nil {
			log.Println(err) // Don't want to do much with this.
		}
		metric.Gauge.Set(m) // Only gauge for now
	}
}
