package ade_linter

import (
	"errors"
	"fmt"
	"regexp"
)

const (
	Addresses  = "addresses"
	Units      = "units"
	Ports      = "ports"
	BoardIds   = "board_ids"
	MessageIds = "message_ids"
	BackendKey = "Backend"
)

func checkGlobalInfo(sheet Sheet) bool {
	globalLogger := Log.AddEntry("GlobalInfo", "")
	globalInfo, err := getGlobalInfo(sheet)

	if err != nil {
		globalLogger.Error(err)
	}

	err = checkAddressTable(globalInfo.Addresses)

	if err != nil {
		globalLogger.Error(err)
	}

	err = checkUnits(globalInfo.Units)

	if err != nil {
		globalLogger.Error(err)
	}

	err = checkPorts(globalInfo.Ports)

	if err != nil {
		globalLogger.Error(err)
	}

	err = checkBoardIds(globalInfo.BoardIds, globalInfo.Addresses)

	if err != nil {
		globalLogger.Error(err)
	}

	err = checkMessageIds(globalInfo.MessageIds)

	if err != nil {
		globalLogger.Error(err)
	}

	return true
}

type GlobalInfo struct {
	Addresses  map[string]string
	Units      map[string]string
	Ports      map[string]string
	BoardIds   map[string]string
	MessageIds map[string]string
}

func getGlobalInfo(sheet Sheet) (GlobalInfo, error) {
	tables := getTables(sheet)
	addresses, ok := tables[Addresses]

	if !ok {
		return GlobalInfo{}, errors.New("address table not found")
	}

	units, ok := tables[Units]

	if !ok {
		return GlobalInfo{}, errors.New("units table not found")
	}

	ports, ok := tables[Ports]

	if !ok {
		return GlobalInfo{}, errors.New("ports table not found")
	}

	boardIds, ok := tables[BoardIds]

	if !ok {
		return GlobalInfo{}, errors.New("boardIds table not found")
	}

	messageIds, ok := tables[MessageIds]

	if !ok {
		return GlobalInfo{}, errors.New("messageIds table not found")
	}

	return GlobalInfo{
		Addresses:  tableToMap(addresses[1:]),
		Units:      tableToMap(units[1:]),
		Ports:      tableToMap(ports[1:]),
		BoardIds:   tableToMap(boardIds[1:]),
		MessageIds: tableToMap(messageIds[1:]),
	}, nil
}

func checkAddressTable(addresses map[string]string) error {
	_, ok := addresses[BackendKey]

	if !ok {
		return errors.New("address table doesn't exists")
	}

	return checkIps(addresses)
}

func checkIps(ips map[string]string) error {
	ipExp := regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`)

	for _, ipStr := range ips {
		if !ipExp.MatchString(ipStr) {
			return fmt.Errorf("incorrect IP: %s", ipStr)
		}
	}

	return nil
}

func checkUnits(units map[string]string) error {
	unitExp := regexp.MustCompile(`^(?:[/+*-]\d+)$`)

	for _, unitStr := range units {
		if unitStr == "" {
			return nil
		}

		if !unitExp.MatchString(unitStr) {
			return fmt.Errorf("incorrect units: %s", unitStr)
		}
	}

	return nil
}

func checkPorts(ports map[string]string) error {
	portExp := regexp.MustCompile(`^\d{1,5}$`)

	for _, portStr := range ports {
		if !portExp.MatchString(portStr) {
			return fmt.Errorf("incorrect port: %s", portStr)
		}
	}

	return nil
}

func checkBoardIds(boardIds map[string]string, addresses map[string]string) error {
	idExp := regexp.MustCompile(`^\d+$`)

	for board, idStr := range boardIds {
		_, ok := addresses[board]

		if !ok {
			return fmt.Errorf("%s IP missing in address table", board)
		}

		if !idExp.MatchString(idStr) {
			return fmt.Errorf("incorrect board id: %s - %s", board, idStr)
		}
	}

	return nil
}

func checkMessageIds(messageIds map[string]string) error {
	idExp := regexp.MustCompile(`^\d+$`)

	for _, idStr := range messageIds {
		if !idExp.MatchString(idStr) {
			return fmt.Errorf("incorrect message id: %s", idStr)
		}
	}

	return nil
}

func tableToMap(table [][]Cell) map[string]string {
	m := make(map[string]string, len(table))

	for _, row := range table {
		m[row[0]] = row[1]
	}

	return m
}
