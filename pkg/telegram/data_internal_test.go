package telegram

import (
	"testing"

	"github.com/hyperized/dsmr/pkg/obis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDataFromLine(t *testing.T) {
	tests := []struct {
		input     string
		wantValue string
	}{
		{"1-0:32.32.0(00002)", "00002"},
		{"1-3:0.2.8(50)", "50"},
		{"0-0:1.0.0(101209113020W)", "101209113020W"},
		{"0-0:96.1.1(4B384547303034303436333935353037)", "4B384547303034303436333935353037"},
		{"1-0:1.8.1(123456.789*kWh)", "123456.789"},
		{"0-0:96.13.0(303132333435363738393A3B3C3D3E3F)", "303132333435363738393A3B3C3D3E3F"},
		{"0-1:24.2.1(101209110000W)(12785.123*m3)", "101209110000W"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			d, err := NewData(tt.input)
			require.NoError(t, err)
			assert.NotEmpty(t, d.raw)
			assert.Equal(t, tt.wantValue, d.value[0])
		})
	}
}

func TestMBusMultiValue(t *testing.T) {
	d, err := NewData("0-1:24.2.1(101209110000W)(12785.123*m3)")
	require.NoError(t, err)
	assert.Len(t, d.value, 2)
	assert.Equal(t, "101209110000W", d.value[0])
	assert.Equal(t, "12785.123", d.value[1])
	assert.Len(t, d.TypedValues(), 2)
}

func TestGenericProfileSkipsTypedValues(t *testing.T) {
	// 1-0:99.97.0 uses GenericProfile class — typed values are deliberately
	// not produced; raw values are still available via Values().
	d, err := NewData("1-0:99.97.0(2)(0-0:96.7.19)(101208152415W)(0000000240*s)")
	require.NoError(t, err)
	assert.Nil(t, d.TypedValues())
	assert.NotEmpty(t, d.Values())
}

func TestParseTypedIntegerFormat(t *testing.T) {
	// 0-1:24.1.0 uses I3 format.
	d, err := NewData("0-1:24.1.0(003)")
	require.NoError(t, err)
	assert.Len(t, d.TypedValues(), 1)
}

func TestParseTypedIntegerFormatError(t *testing.T) {
	// Wrong length for I3 → typed values empty.
	d, err := NewData("0-1:24.1.0(INVALID)")
	require.NoError(t, err)
	assert.Empty(t, d.TypedValues())
}

func TestParseTypedTSTExtendedRegisterError(t *testing.T) {
	// Bad timestamp for an M-Bus extended register.
	d, err := NewData("0-1:24.2.1(BADTS)(12785.123*m3)")
	require.NoError(t, err)
	tv := d.TypedValues()
	// Timestamp parse failed; only the metered value string is stored.
	assert.Len(t, tv, 1)
	assert.Equal(t, "12785.123", tv[0])
}

func TestParseTypedTSTError(t *testing.T) {
	// Bad timestamp for a plain TST (non-M-Bus) line.
	d, err := NewData("0-0:1.0.0(BADTS)")
	require.NoError(t, err)
	assert.Empty(t, d.TypedValues())
}

func TestNewDataEmptyValue(t *testing.T) {
	// Empty value inside parens → "no value could be parsed" error.
	d, err := NewData("0-0:96.13.0()")
	require.Error(t, err)
	assert.Nil(t, d)
}

func TestParseTypedFloatingPointError(t *testing.T) {
	// "BADVALUE" fails validation inside cosem.NewFloatingPoint → slog.Debug branch, empty TypedValues.
	d, err := NewData("1-0:1.8.1(BADVALUE)")
	require.NoError(t, err)
	assert.Empty(t, d.TypedValues())
}

func TestParseTypedEmptyFormatString(t *testing.T) {
	// Inject a zero-value Reference (empty formatString) to cover the fs==""
	// guard in parseTyped. This branch is unreachable via normal references.
	obis.References["9-9:99.00.0"] = &obis.Reference{}
	t.Cleanup(func() { delete(obis.References, "9-9:99.00.0") })

	d, err := NewData("9-9:99.00.0(somevalue)")
	require.NoError(t, err)
	assert.Empty(t, d.TypedValues())
}
