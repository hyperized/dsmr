package obis_test

import (
	"sync"
	"testing"

	"github.com/hyperized/dsmr/pkg/obis"
	"github.com/stretchr/testify/assert"
)

var (
	metricsOnce sync.Once
	sharedM     *obis.Metrics
)

func getTestMetrics() *obis.Metrics {
	metricsOnce.Do(func() { sharedM = obis.Register() })
	return sharedM
}

func TestRegister(t *testing.T) {
	m := getTestMetrics()
	assert.NotNil(t, m)
	assert.NotNil(t, m.Gauges())
	assert.NotEmpty(t, m.Gauges())
	assert.NotNil(t, m.MBus())
	assert.NotNil(t, m.DSMRInfo())
	assert.NotNil(t, m.EquipInfo())
}
