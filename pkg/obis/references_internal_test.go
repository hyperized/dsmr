package obis

import (
	"testing"

	"github.com/hyperized/dsmr/pkg/cosem"
	"github.com/hyperized/dsmr/pkg/cosem/unit"
	"github.com/stretchr/testify/assert"
)

func TestNewInternalStructure(t *testing.T) {
	assert.Equal(t,
		&Reference{
			name:        "DSMRVersion",
			identifier:  "1-3:0.2.8",
			description: "DSMR version information for P1 output",
			unit:        unit.None,
			format: Format{
				tag:          cosem.OctetString,
				class:        cosem.Data,
				attribute:    cosem.Value,
				length:       2,
				formatString: "S2",
			},
		},
		New("1-3:0.2.8"),
	)
}
