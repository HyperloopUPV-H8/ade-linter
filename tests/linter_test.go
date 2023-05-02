package tests

import (
	"log"
	"path/filepath"
	"testing"

	ade_linter "github.com/HyperloopUPV-H8/ade-linter"
	"github.com/xuri/excelize/v2"
)

func TestLinter(t *testing.T) {
	t.Run("general linter test", func(t *testing.T) {
		file, err := excelize.OpenFile(filepath.Join("./", "ade.xlsx"))

		if err != nil {
			log.Fatal(err)
		}

		ade_linter.Lint(file)
	})
}
