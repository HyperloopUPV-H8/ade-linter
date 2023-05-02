package ade_linter

import (
	"errors"
	"fmt"

	"golang.org/x/exp/slices"
)

func checkStructures(packets Table, measurements Table, structuresColumns Table, logger Logger) bool {
	structures := getStructures(structuresColumns, logger)
	packetNames := getPacketNames(packets)
	measurementIds := getMeasurementIds(measurements)

	for _, structure := range structures {
		if !checkStructure(structure, packetNames, measurementIds) {
			return false
		}
	}

	return true
}

type Structure struct {
	column []string
	logger Logger
}

func getStructures(tableColumns Table, logger Logger) []Structure {
	structures := make([]Structure, 0)

	for _, column := range tableColumns {
		structure, err := getStructure(column, logger)

		if err != nil {
			logger.Error(err)
			continue
		}

		structures = append(structures, structure)
	}

	return structures
}

func getStructure(column []string, logger Logger) (Structure, error) {
	if len(column) == 0 {
		return Structure{}, errors.New("structure column is empty")
	}

	return Structure{
		column: column,
		logger: logger.AddEntry("Structure", column[0]),
	}, nil
}

func checkStructure(structure Structure, packets []string, measurements []string) bool {
	if !isPacketDefined(structure.column[0], packets) {
		structure.logger.Error(fmt.Errorf("packet %s is not defined in the packet table", structure.column[0]))
		return false
	}

	err := areMeasurementsDefined(structure.column[1:], measurements)

	if err != nil {
		structure.logger.Error(err)
	}

	return true
}

func isPacketDefined(packetName string, packets []string) bool {
	return slices.Contains(packets, packetName)
}

func areMeasurementsDefined(structureMeasurements []string, mNames []string) error {
	for _, name := range structureMeasurements {
		if !slices.Contains(mNames, name) {
			return fmt.Errorf("measurement %s is not defined in the measurement table", name)
		}
	}

	return nil
}
