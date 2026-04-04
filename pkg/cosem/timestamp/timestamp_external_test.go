package timestamp_test

import (
	"testing"
	"time"

	"github.com/hyperized/dsmr/pkg/cosem"
	"github.com/hyperized/dsmr/pkg/cosem/timestamp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimestampNew(t *testing.T) {
	loc, _ := time.LoadLocation("Europe/Amsterdam")

	tests := []struct {
		name     string
		input    string
		fails    bool
		wantDST  timestamp.DSTFlag
		wantTime time.Time
		wantSumm bool
	}{
		{
			name:     "winter timestamp",
			input:    "210101120000W",
			wantDST:  timestamp.Winter,
			wantTime: time.Date(2021, 1, 1, 12, 0, 0, 0, loc),
			wantSumm: false,
		},
		{
			name:     "summer timestamp",
			input:    "210701140530S",
			wantDST:  timestamp.Summer,
			wantTime: time.Date(2021, 7, 1, 14, 5, 30, 0, loc),
			wantSumm: true,
		},
		{
			name:  "too short",
			input: "2101011200W",
			fails: true,
		},
		{
			name:  "invalid DST flag",
			input: "210101120000X",
			fails: true,
		},
		{
			name:  "invalid datetime",
			input: "999999999999W",
			fails: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts, err := timestamp.New(tt.input)
			if tt.fails {
				require.Error(t, err)
				assert.Nil(t, ts)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantDST, ts.DST())
			assert.Equal(t, tt.wantSumm, ts.IsSummer())
			assert.True(t, tt.wantTime.Equal(ts.Value()),
				"expected %v, got %v", tt.wantTime, ts.Value())
		})
	}
}

func TestTimestampWithTag(t *testing.T) {
	ts, err := timestamp.New("210101120000W", timestamp.WithTag(cosem.OctetString))
	require.NoError(t, err)
	assert.NotNil(t, ts)
	assert.Equal(t, timestamp.Winter, ts.DST())
}
