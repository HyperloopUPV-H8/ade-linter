package ade_linter

import (
	"fmt"
	"strings"
)

const TablePrefix = "[TABLE]"

func getTables(sheet Sheet) (Tables, error) {
	tables := make(map[string]Table)

	for i, row := range sheet {
		for j, cell := range row {
			if strings.HasPrefix(cell, TablePrefix) {
				name := strings.TrimPrefix(cell, fmt.Sprintf("%s ", TablePrefix))
				table, err := getTableAutomaticWidth(sheet, i, j)

				if err != nil {
					return map[string]Table{}, err
				}

				tables[name] = table
			}
		}
	}

	return tables, nil
}

func findTableWithWidth(sheet Sheet, name string, width int) (Table, bool) {
	row, col := findTableHeader(sheet, name)

	if row == -1 || col == -1 {
		return [][]string{}, false
	}

	if row == len(sheet)-1 {
		return make([][]string, 0), true
	}

	return getTableWithWidth(sheet, row+1, col, width), true
}

func findTableAutoWidth(sheet Sheet, name string) (Table, error) {
	row, col := findTableHeader(sheet, name)

	if row == -1 || col == -1 {
		return [][]string{}, fmt.Errorf("table %s not found", name)
	}

	table, err := getTableAutomaticWidth(sheet, row, col)

	if err != nil {
		return [][]string{}, err
	}

	return table, nil
}

func findTableHeader(sheet Sheet, name string) (int, int) { // returns row, col
	for i, row := range sheet {
		for k, cell := range row {
			if cell == fmt.Sprint(TablePrefix, " ", name) {
				return i, k
			}
		}
	}

	return -1, -1
}

func getTableWithWidth(sheet Sheet, row int, col int, width int) Table {
	submatrix := make([][]string, 0)

	for i := row; i < len(sheet) && !isRowEmpty(sheet[i][col:col+width]); i++ {
		submatrix = append(submatrix, sheet[i][col:col+width])
	}

	return submatrix
}

func getTableAutomaticWidth(sheet Sheet, row int, column int) ([][]Cell, error) {
	if row == len(sheet)-1 {
		return [][]Cell{}, fmt.Errorf("table header is in the last row")
	}

	rowLength := getRowLength(sheet[row+1], column)
	columnLength := getColumnLength(sheet, row+1, column)

	return getSubMatrix(sheet, row+1, rowLength, column, columnLength), nil
}

func getRowLength(row []Cell, start int) int {
	for i := start; i < len(row); i++ {
		if isRowEmpty(row) {
			return i - start
		}
	}

	return len(row) - start
}

func isRowEmpty(row []string) bool {
	empty := true

	for _, cell := range row {
		if cell != "" {
			empty = false
		}
	}

	return empty
}

func getColumnLength(sheet Sheet, row int, col int) int {
	for i := row; i < len(sheet); i++ {
		if sheet[i][col] == "" {
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
