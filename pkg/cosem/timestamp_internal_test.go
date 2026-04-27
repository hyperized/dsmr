package cosem

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestLoadLocationOrUTC covers both branches of the helper: a valid IANA name
// returns the requested location; an unknown name falls back to UTC.
func TestLoadLocationOrUTC(t *testing.T) {
	loc := loadLocationOrUTC("Europe/Amsterdam")
	assert.NotNil(t, loc)
	assert.Equal(t, "Europe/Amsterdam", loc.String())

	assert.Equal(t, time.UTC, loadLocationOrUTC("Not/A/Real/Zone"))
}
