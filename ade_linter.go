package ade_linter

import (
	"fmt"
	"log"

	"github.com/xuri/excelize/v2"
)

var Log = NewLogger()

func Lint(file *excelize.File) bool {
	fmt.Println("🧼 Cleaning ADE...")

	doc := getDocument(file)
	if areTitlesCorrect(doc) && checkGlobalInfo(doc.Sheets[GLOBAL_INFO]) && checkBoards(getBoardSheets(doc)) {
		Log.Pass("ADE has been validated")
		return true
	}

	return false
}

func getBoardSheets(doc Document) map[string]Sheet {
	boards := NewMap(doc.Sheets)
	delete(boards, GLOBAL_INFO)
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
