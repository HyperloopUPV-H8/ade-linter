package ade_linter

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	RangeExp = regexp.MustCompile(fmt.Sprintf(`^\[(?:(?:%s,%s)|(?:,%s)|(?:%s,))\]$`, FloatExp, FloatExp, FloatExp, FloatExp))
)

func checkMeasurements(table Table, logger Logger) bool {
	measurements := getMeasurements(table, logger)
	return Every(measurements)
}

func getMeasurements(table Table, logger Logger) []Test {
	return Map(table, func(item []string) Test {
		return Test(Measurement(Measurement{
			cells:  item,
			logger: logger.AddEntry("Measurement", item[0]),
		}))
	})
}

type Measurement struct {
	cells  []string
	logger Logger
}

func (m Measurement) Run() bool {

	if m.cells[0] == "" {
		m.logger.Error(fmt.Errorf("empty id"))
		return false
	}

	if m.cells[1] == "" {
		m.logger.Error(fmt.Errorf("empty name"))
		return false
	}

	if !checkMeasurementType(m.cells[2]) {
		m.logger.Error(fmt.Errorf("invalid type: %s", m.cells[2]))
		return false
	}

	if !CheckMeasurementUnits(m.cells[3]) {
		m.logger.Error(fmt.Errorf("invalid units: %s", m.cells[3]))
		return false
	}

	if !CheckMeasurementUnits(m.cells[4]) {
		m.logger.Error(fmt.Errorf("invalid units: %s", m.cells[4]))
		return false
	}

	if !CheckRange(m.cells[5]) {
		m.logger.Error(fmt.Errorf("invalid range: %s", m.cells[5]))
		return false
	}

	if !CheckRange(m.cells[6]) {
		m.logger.Error(fmt.Errorf("invalid range: %s", m.cells[6]))
		return false
	}

	return true
}

func checkMeasurementType(kind string) bool {
	return isNumericType(kind) || isBoolean(kind) || IsEnum(kind)
}

func isNumericType(kind string) bool {
	return kind == "uint8" || kind == "uint16" || kind == "uint32" || kind == "uint64" || kind == "int8" || kind == "int16" || kind == "int32" || kind == "int64" || kind == "float32" || kind == "float64"
}

func isBoolean(kind string) bool {
	return kind == "bool"
}

func IsEnum(kind string) bool {
	if !strings.HasPrefix(kind, "enum") {
		return false
	}

	optionsStr := strings.TrimPrefix(kind, "enum(")
	optionsStr = strings.TrimSuffix(optionsStr, ")")
	optionsStr = strings.ReplaceAll(optionsStr, " ", "")
	optionsExp := regexp.MustCompile(`^(?:\w+,)*\w+$`)
	return optionsExp.MatchString(optionsStr)
}

func CheckMeasurementUnits(unitStr string) bool {
	unitExp := regexp.MustCompile(`^(?:[^# ]+(?:#(?:[+-/*]\d+)+)?)?$`)
	return unitExp.MatchString(unitStr)
}

func CheckRange(rangeStr string) bool {
	if rangeStr == "" {
		return true
	}

	return RangeExp.MatchString(strings.ReplaceAll(rangeStr, " ", ""))
}
