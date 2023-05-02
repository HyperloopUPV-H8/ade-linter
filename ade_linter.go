package ade_linter

import (
	"fmt"
	"log"
	"strings"

	"github.com/xuri/excelize/v2"
)

var Log = NewLogger()

func Lint(file *excelize.File) bool {
	fmt.Println("🧼 Cleaning ADE...")
	doc := getDocument(file)

	boardSheets := getBoardSheets(doc)

	if !areTitlesCorrect(doc, boardSheets) {
		return false
	}

	if !checkGlobalInfo(doc.Sheets[GLOBAL_INFO]) {
		return false
	}

	if !checkBoards(boardSheets) {
		return false
	}

	Log.Pass("ADE has been validated")
	return true
}

func getBoardSheets(doc Document) map[string]Sheet {
	boards := make(map[string]Sheet)

	for name, sheet := range doc.Sheets {
		if !strings.HasPrefix(name, BOARD_PREFIX) {
			if name != GLOBAL_INFO {
				Log.Warn(fmt.Errorf("sheet name %s doesn't have board prefix %s", name, BOARD_PREFIX))
			}
			continue
		}

		boards[name] = sheet
	}

	return boards
}

func getDocument(file *excelize.File) Document {
	document := NewDocument()

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
