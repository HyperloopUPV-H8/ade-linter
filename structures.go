package ade_linter

import (
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
		structures = append(structures,
			Structure{
				column: column,
				logger: logger.AddEntry("Structure", column[1]),
			},
		)
	}

	return structures
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
