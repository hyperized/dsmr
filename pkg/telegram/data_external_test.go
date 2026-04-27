package telegram_test

import (
	"testing"

	"github.com/hyperized/dsmr/pkg/telegram"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDataErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty string", ""},
		{"no obis code", "(123456)"},
		{"unknown obis code", "9-9:99.99.9(value)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := telegram.NewData(tt.input)
			assert.Error(t, err)
		})
	}
}

func TestTypedValueFloatingPoint(t *testing.T) {
	d, err := telegram.NewData("1-0:1.8.1(123456.789*kWh)")
	require.NoError(t, err)
	assert.Len(t, d.TypedValues(), 1)
	assert.NotNil(t, d.TypedValues()[0])
}

func TestTypedValueTimestamp(t *testing.T) {
	d, err := telegram.NewData("0-0:1.0.0(101209113020W)")
	require.NoError(t, err)
	assert.Len(t, d.TypedValues(), 1)
	assert.NotNil(t, d.TypedValues()[0])
}

func TestMBusOtherChannels(t *testing.T) {
	for _, ch := range []string{"2", "3", "4"} {
		t.Run("channel "+ch, func(t *testing.T) {
			d, err := telegram.NewData("0-" + ch + ":24.2.1(101209110000W)(12785.123*m3)")
			require.NoError(t, err)
			vals := d.Values()
			assert.Len(t, vals, 2)
			assert.Equal(t, "101209110000W", vals[0])
			assert.Equal(t, "12785.123", vals[1])
		})
	}
}

func TestDataIdentifier(t *testing.T) {
	d, err := telegram.NewData("1-0:1.8.1(123456.789*kWh)")
	require.NoError(t, err)
	assert.Equal(t, "1-0:1.8.1", d.Identifier())
}

func TestDataMetricName(t *testing.T) {
	d, err := telegram.NewData("1-0:1.8.1(123456.789*kWh)")
	require.NoError(t, err)
	assert.Equal(t, "electricity_delivered_to_client_tariff1_kwh", d.MetricName())
}

func TestDataString(t *testing.T) {
	d, err := telegram.NewData("1-0:1.8.1(123456.789*kWh)")
	require.NoError(t, err)
	s := d.String()
	assert.Contains(t, s, "123456.789")
}
