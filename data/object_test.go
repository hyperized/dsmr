package data

import (
	"testing"
)

func TestNumberOfVoltageSagsInPhaseL1(t *testing.T) {
	t.Run("NumberOfVoltageSagsInPhaseL1", func(t *testing.T) {
		object, err := NewFromLine("1-0:32.32.0(00002)")
		t.Logf("%#v \n", object)
		if err != nil {
			return
		}
	})
}

func TestNumberOfVoltageSagsInPhaseL2(t *testing.T) {
	t.Run("NumberOfVoltageSagsInPhaseL2", func(t *testing.T) {
		object, err := NewFromLine("1-0:52.32.0(00001)")
		t.Logf("%#v \n", object)
		if err != nil {
			return
		}
	})
}
