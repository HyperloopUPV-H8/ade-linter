package tests

import (
	"testing"

	ade_linter "github.com/HyperloopUPV-H8/ade-linter"
)

func TestGlobalInfo(t *testing.T) {
	t.Run("addresses are validated correctly", func(t *testing.T) {
		addresses := []MockTest[string, bool]{
			{"127.0.0.1", true},
			{"255.255.255.255", true},
			{"0.0.0.0", true},
			{"0.0.0.", false},
			{"0.0.0", false},
			{"0", false},
			{".0.", false},
			{".0.", false},
		}

		AssertMocks(addresses, ade_linter.CheckIp, func(mock MockTest[string, bool], got bool) {
			t.Errorf("%s was incorrectly marked as %t", mock.Prompt, got)
		})
	})

	t.Run("units are validated correctly", func(t *testing.T) {
		units := []MockTest[string, bool]{
			{"+123-123/123*1231", true},
			{"+123-123//123**1231", false},
			{"123", false},
		}

		AssertMocks(units, ade_linter.CheckUnit, func(mock MockTest[string, bool], got bool) {
			t.Errorf("%s was incorrectly marked as %t", mock.Prompt, got)
		})
	})

	t.Run("ports are validated correctly", func(t *testing.T) {
		ports := []MockTest[string, bool]{
			{"0", true},
			{"123", true},
			{"65535", true},
			{"65536", false},
			{"-1", false},
		}

		AssertMocks(ports, ade_linter.CheckPort, func(mock MockTest[string, bool], got bool) {
			t.Errorf("%s was incorrectly marked as %t", mock.Prompt, got)
		})
	})

	t.Run("board_ids are validated correctly", func(t *testing.T) {

		type Case struct {
			boardIds  map[string]string
			addresses map[string]string
		}

		cases := []MockTest[Case, bool]{
			{Case{boardIds: map[string]string{
				"LCU_MASTER": "1",
				"LCU_SLAVE":  "2",
			},
				addresses: map[string]string{
					"LCU_MASTER": "127.0.0.1",
					"LCU_SLAVE":  "127.0.0.1",
				},
			}, true},
		}

		for _, myCase := range cases {
			if err := ade_linter.CheckBoardIds(myCase.Prompt.boardIds, myCase.Prompt.addresses); err != nil {
				t.Errorf(err.Error())
			}
		}

	})

	t.Run("message_ids are validated correctly", func(t *testing.T) {
		ids := []MockTest[string, bool]{
			{"1", true},
			{"2", true},
			{"12341235", true},
			{"-1", false},
			{"-10", false},
		}

		AssertMocks(ids, ade_linter.CheckMessageId, func(mock MockTest[string, bool], got bool) {
			t.Errorf("%s was incorrectly marked as %t", mock.Prompt, got)
		})
	})
}
