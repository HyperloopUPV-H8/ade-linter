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

		document := getDocument(file)
		ade_linter.Lint(document)
	})
}

func getDocument(file *excelize.File) ade_linter.Document {
	document := ade_linter.NewDocument()

	sheetNames := file.GetSheetMap()
	for _, name := range sheetNames {
		rows, err := file.GetRows(name)

		if err != nil {
			log.Fatal("sheet not found")
		}

		document.Sheets[name] = makeRowsSameLength(rows)
	}

	return document
}

func makeRowsSameLength(rows [][]string) [][]string {
	maxLength := 0

	for _, row := range rows {
		if len(row) > maxLength {
			maxLength = len(row)
		}
	}

	fullRows := make([][]string, 0)
	for _, row := range rows {
		fullRow := make([]string, maxLength)
		copy(fullRow, row)
		fullRows = append(fullRows, fullRow)
	}

	return fullRows
}
