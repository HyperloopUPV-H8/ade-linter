package ade_linter

import (
	"fmt"
	"strings"
)

const TablePrefix = "[TABLE]"

func getTables(sheet Sheet) Tables {
	tables := make(map[string]Table)

	for i, row := range sheet {
		for j, cell := range row {
			if strings.HasPrefix(cell, TablePrefix) {
				name, cells := getTable(sheet, i, j)
				tables[name] = cells
			}
		}
	}

	return tables
}

func getTable(sheet Sheet, row int, column int) (string, [][]Cell) {
	name := strings.TrimPrefix(sheet[row][column], fmt.Sprintf("%s ", TablePrefix))

	rowLength := getRowLength(sheet[row+1], column)
	columnLength := getColumnLength(sheet, row+1, column)

	return name, getSubMatrix(sheet, row+1, rowLength, column, columnLength)
}

func getRowLength(row []Cell, start int) int {
	for i := start; i < len(row); i++ {
		if row[i] == "" {
			return i - start
		}
	}

	return len(row) - start
}

func getColumnLength(sheet Sheet, row int, col int) int {
	for i := row; i < len(sheet); i++ {
		if len(sheet[i])-1 < col || sheet[i][col] == "" {
			return i - row
		}
	}

	return len(sheet) - row
}

func getSubMatrix[T any](matrix [][]T, startRow int, rowLength int, startCol int, colLength int) [][]T {
	rows := matrix[startRow : startRow+colLength]

	submatrix := make([][]T, len(rows))

	for index, row := range rows {
		submatrix[index] = row[startCol : startCol+rowLength]
	}

	return submatrix
}
