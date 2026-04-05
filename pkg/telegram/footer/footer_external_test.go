package footer_test

import (
	"testing"

	"github.com/hyperized/dsmr/pkg/telegram/footer"
	"github.com/stretchr/testify/assert"
)

func TestFooterNew(t *testing.T) {
	f := footer.New(footer.WithCRC("F46A"))
	assert.Equal(t, "Footer [CRC: F46A]", f.String())
}

func TestFooterCRC(t *testing.T) {
	f := footer.New(footer.WithCRC("ABCD"))
	assert.Equal(t, "ABCD", f.CRC())
}
