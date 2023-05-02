package ade_linter

import (
	"fmt"
	"strings"
)

const (
	PacketTable      = "Packets"
	MeasurementTable = "Measurements"
	Structures       = "Structures"
)

func checkBoards(sheets map[string]Sheet) bool {
	return Every(getBoards(sheets))
}

func getBoards(sheets map[string]Sheet) []Test {
	boards := make([]Test, len(sheets))

	i := 0
	for name, sheet := range sheets {
		boardName := strings.TrimPrefix(name, BOARD_PREFIX)
		boardLogger := Log.AddEntry("Board", boardName)
		board, err := getBoard(sheet, boardLogger)

		if err != nil {
			//TODO: error
			return []Test{
				FalseTest{},
			}
		}

		boards[i] = board
		i++
	}

	return boards
}

type Board struct {
	Packets      Table
	Measurements Table
	Structures   Table
	logger       Logger
}

func (board Board) Run() bool {
	if !checkPackets(board.Packets, board.logger) {
		return false
	}

	if !checkMeasurements(board.Measurements, board.logger) {
		return false
	}

	return checkStructures(board.Packets, board.Measurements, board.Structures, board.logger)
}

func getBoard(sheet Sheet, logger Logger) (Board, error) {
	packets, err := getTable(PacketTable, sheet, []string{
		"ID",
		"Name",
		"Type",
	})

	if err != nil {
		return Board{}, err
	}

	measurements, err := getTable(MeasurementTable, sheet, []string{"ID", "Name", "Type", "PodUnits", "DisplayUnits", "SafeRange", "WarningRange"})

	if err != nil {
		return Board{}, err
	}

	structures, err := findTableAutoWidth(sheet, Structures)

	if err != nil {
		logger.Error(err)
		return Board{}, err
	}

	structures = getStructureColumns(structures)

	return Board{
		Packets:      packets,
		Measurements: measurements,
		Structures:   structures,
		logger:       logger,
	}, nil
}

func getTable(name string, sheet Sheet, headers []string) (Table, error) {
	table, ok := findTableWithWidth(sheet, name, len(headers))

	if !ok {
		return Table{}, fmt.Errorf("table %s not found", name)
	}

	if len(table) == 0 {
		return Table{}, fmt.Errorf("table %s is empty (not even headers)", name)
	}

	if !areHeadersCorrect(table[0], headers) {
		return Table{}, fmt.Errorf("incorrect headers: %v", table[0])
	}

	if len(table) == 1 {
		return Table{}, fmt.Errorf("table %s has no fields", name)
	}

	return table[1:], nil
}

func areHeadersCorrect(headerRow []string, headers []string) bool {
	if len(headerRow) != len(headers) {
		return false
	}

	for index, header := range headers {
		if !(header == headerRow[index]) {
			return false
		}
	}

	return true
}

func getPacketNames(packets Table) []string {
	return Map(packets, func(item []string) string {
		return item[1]
	})
}

func getMeasurementIds(measurements Table) []string {
	return Map(measurements, func(item []string) string {
		return item[0]
	})
}

func getStructureColumns(rows [][]string) [][]string {
	structures := rows[1:]
	structures = toColumns(structures)

	croppedStructures := make([][]string, 0)
outer:
	for _, column := range structures {
		for j, cell := range column {
			if cell == "" {
				croppedStructures = append(croppedStructures, column[:j])
				continue outer
			}
		}
		croppedStructures = append(croppedStructures, column)
	}

	return croppedStructures
}

func toColumns(rows [][]string) [][]string {
	if len(rows) == 0 {
		return rows
	}

	columns := make([][]string, 0)
	for i := 0; i < len(rows[0]); i++ {
		columns = append(columns, toColumn(rows, i))
	}

	return columns
}

func toColumn(rows [][]string, col int) []string {
	if len(rows) == 0 {
		return make([]string, 0)
	}

	column := make([]string, 0)

	for i := 0; i < len(rows); i++ {
		column = append(column, rows[i][col])
	}

	return column
}
