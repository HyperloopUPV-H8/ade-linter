package ade_linter

import (
	"errors"
	"fmt"
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
		boardLogger := Log.AddEntry("Board", name)
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
	tables := getTables(sheet)

	packets, ok := tables[PacketTable]

	if !ok {
		logger.Error(fmt.Errorf("packet table not found"))
		return Board{}, errors.New("packet table not found")
	}

	packets = packets[1:]

	measurements, ok := tables[MeasurementTable]

	if !ok {
		logger.Error(fmt.Errorf("measurement table not found"))
		return Board{}, errors.New("measurements table not found")
	}

	measurements = measurements[1:]

	structures, ok := tables[Structures]

	if !ok {
		logger.Error(fmt.Errorf("structures table not found"))
		return Board{}, errors.New("structures table not found")
	}

	structures = getStructureColumns(structures)

	return Board{
		Packets:      packets,
		Measurements: measurements,
		Structures:   structures,
		logger:       logger,
	}, nil
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
