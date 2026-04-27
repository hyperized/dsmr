package obis_test

import (
	"testing"

	"github.com/hyperized/dsmr/pkg/obis"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	m := obis.Register(prometheus.NewRegistry())
	assert.NotNil(t, m)
	assert.NotNil(t, m.Gauges())
	assert.NotEmpty(t, m.Gauges())
	assert.NotNil(t, m.MBus())
	assert.NotNil(t, m.MBusValve())
	assert.NotNil(t, m.MBusEquipInfo())
	assert.NotNil(t, m.DSMRInfo())
	assert.NotNil(t, m.EquipInfo())
}

func TestRegisterTwiceDoesNotPanic(t *testing.T) {
	// Two separate registries must not collide — this used to panic when
	// Register() shared the global DefaultRegisterer via promauto.
	assert.NotNil(t, obis.Register(prometheus.NewRegistry()))
	assert.NotNil(t, obis.Register(prometheus.NewRegistry()))
}
