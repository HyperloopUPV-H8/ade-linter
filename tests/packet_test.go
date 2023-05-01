package tests

import (
	"testing"

	ade_linter "github.com/HyperloopUPV-H8/ade-linter"
)

func TestPacket(t *testing.T) {
	t.Run("ids are uint16", func(t *testing.T) {
		mockIds := []MockTest[string, bool]{
			{"0", true},
			{"1", true},
			{"65535", true},
			{"65536", false},
			{"-1", false},
			{"-0", false},
		}

		AssertMocks(mockIds, ade_linter.CheckId, func(mock MockTest[string, bool], got bool) {
			t.Errorf("%s was incorrectly marked as %t", mock.Prompt, got)
		})
	})

	t.Run("type are corrects", func(t *testing.T) {
		mockTypes := []MockTest[string, bool]{
			{"data", true},
			{"order", true},
			{"fault", true},
			{"warning", true},
			{"", false},
			{"asdf", false},
		}

		AssertMocks(mockTypes, ade_linter.CheckPacketType, func(mock MockTest[string, bool], got bool) {
			t.Errorf("%s was incorrectly marked as %t", mock.Prompt, got)
		})
	})
}
