package tests

import (
	"testing"

	ade_linter "github.com/HyperloopUPV-H8/ade-linter"
)

func TestMeasurement(t *testing.T) {
	t.Run("enum is checked correctly", func(t *testing.T) {
		enumMocks := []MockTest[string, bool]{
			{"enum(A,B)", true},
			{"Enum(A,B,C,D,FADF)", false},
			{"ENUM(A,B,C,D,FADF)", false},
			{"enum(,,D,)", false},
			{"enum(A)", true},
			{"enum()", false},
		}

		AssertMocks(enumMocks, ade_linter.IsEnum, func(mock MockTest[string, bool], got bool) {
			t.Errorf("%s was incorrectly marked as %t", mock.Prompt, got)

		})
	})

	t.Run("units are checked correctly", func(t *testing.T) {
		unitMocks := []MockTest[string, bool]{
			{"A", true},
			{"A  ", false},
			{"V#-123", true},
			{"V#-123+123", true},
			{"V#/123-30*10-0", true},
			{"V##/1", false},
			{"V#132", false},
			{"", true},
			{"#+11", false},
		}

		AssertMocks(unitMocks, ade_linter.CheckMeasurementUnits, func(mock MockTest[string, bool], got bool) {
			t.Errorf("%s was incorrectly marked as %t", mock.Prompt, got)
		})
	})

	t.Run("range is checked correctly", func(t *testing.T) {
		rangeMocks := []MockTest[string, bool]{
			{"[10,20]", true},
			{"[10, 20]", true},
			{"[10,]", true},
			{"[10, ]", true},
			{"[10 , ]", true},
			{"[,10]", true},
			{"[, 10]", true},
			{"", true},
			{"[-10, 10]", true},
			{"[10, -10]", true},
			{"[-10, -10]", true},
			{"[,]", false},
			{"[]", false},
		}

		AssertMocks(rangeMocks, ade_linter.CheckRange, func(mock MockTest[string, bool], got bool) {
			t.Errorf("%s was incorrectly marked as %t", mock.Prompt, got)

		})
	})
}
