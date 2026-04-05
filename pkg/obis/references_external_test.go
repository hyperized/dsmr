package obis_test

import (
	"testing"

	"github.com/hyperized/dsmr/pkg/cosem"
	"github.com/hyperized/dsmr/pkg/cosem/unit"
	"github.com/hyperized/dsmr/pkg/obis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMissing(t *testing.T) {
	assert.Nil(t, obis.New("0-0:00.00.0"))
}

func TestPolyPhaseOnlyL1NotMarked(t *testing.T) {
	codes := []string{
		"1-0:52.32.0", "1-0:72.32.0",
		"1-0:52.36.0", "1-0:72.36.0",
		"1-0:52.7.0", "1-0:72.7.0",
		"1-0:51.7.0", "1-0:71.7.0",
		"1-0:41.7.0", "1-0:61.7.0",
		"1-0:42.7.0", "1-0:62.7.0",
	}
	for _, code := range codes {
		t.Run(code, func(t *testing.T) {
			r := obis.New(code)
			require.NotNil(t, r)
			assert.True(t, r.PolyPhaseOnly(), "expected polyPhaseOnly=true for %s", code)
		})
	}
}

func TestMBusChannelsAllRegistered(t *testing.T) {
	for ch := 1; ch <= 4; ch++ {
		codes := []string{
			"0-" + string(rune('0'+ch)) + ":24.1.0",
			"0-" + string(rune('0'+ch)) + ":96.1.0",
			"0-" + string(rune('0'+ch)) + ":24.2.1",
			"0-" + string(rune('0'+ch)) + ":24.4.0",
		}
		for _, code := range codes {
			t.Run(code, func(t *testing.T) {
				assert.NotNil(t, obis.New(code))
			})
		}
	}
}

func TestMissingCodesRegistered(t *testing.T) {
	for _, code := range []string{"1-0:99.97.0", "0-0:96.13.0"} {
		t.Run(code, func(t *testing.T) {
			assert.NotNil(t, obis.New(code))
		})
	}
}

func TestReferenceAccessors(t *testing.T) {
	ref := obis.New("1-0:1.8.1")
	require.NotNil(t, ref)

	assert.Equal(t, "MeterReadingElectricityDeliveredToClientTariff1", ref.Name())
	assert.Equal(t, "1-0:1.8.1", ref.Identifier())
	assert.Equal(t, unit.KiloWattHour, ref.Unit())
	assert.NotEmpty(t, ref.Description())
	assert.Equal(t, "electricity_delivered_to_client_tariff1_kwh", ref.MetricName())
}

func TestFormatAccessors(t *testing.T) {
	ref := obis.New("1-0:1.8.1")
	require.NotNil(t, ref)

	f := ref.Format()
	// Explicit type conversions are required because many cosem constants are
	// untyped (only the first constant in each typed group has an explicit type).
	assert.Equal(t, cosem.Cosem(cosem.DoubleLongUnsigned), f.Tag())
	assert.Equal(t, cosem.Class(cosem.Register), f.Class())
	assert.Equal(t, cosem.Attribute(cosem.Value), f.Attribute())
	assert.Equal(t, "F9(3,3)", f.FormatString())
	assert.Equal(t, 9, f.Length())
	assert.Equal(t, 3, f.MinimumDecimals())
	assert.Equal(t, 3, f.MaximumDecimals())
}
