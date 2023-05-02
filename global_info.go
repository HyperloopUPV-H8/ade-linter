package ade_linter

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

const (
	Addresses  = "addresses"
	Units      = "units"
	Ports      = "ports"
	BoardIds   = "board_ids"
	MessageIds = "message_ids"
	BackendKey = "Backend"
)

var (
	UnitExp = regexp.MustCompile(fmt.Sprintf(`^(?:[/+*\-]%s+)+$`, FloatExp))
	IpExp   = regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`)
	IdExp   = regexp.MustCompile(`^\d+$`)
)

func checkGlobalInfo(sheet Sheet) bool {
	globalLogger := Log.AddEntry("GlobalInfo", "")
	globalInfo, err := getGlobalInfo(sheet, globalLogger)

	if err != nil {
		globalLogger.Error(err)
		return false
	}

	err = checkAddressTable(globalInfo.Addresses)

	if err != nil {
		globalLogger.Error(err)
		return false
	}

	err = checkUnits(globalInfo.Units)

	if err != nil {
		globalLogger.Error(err)
		return false
	}

	err = checkPorts(globalInfo.Ports)

	if err != nil {
		globalLogger.Error(err)
		return false
	}

	err = CheckBoardIds(globalInfo.BoardIds, globalInfo.Addresses)

	if err != nil {
		globalLogger.Error(err)
		return false
	}

	err = checkMessageIds(globalInfo.MessageIds)

	if err != nil {
		globalLogger.Error(err)
		return false
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

func getGlobalInfo(sheet Sheet, logger Logger) (GlobalInfo, error) {
	tables, err := getTables(sheet)

	if err != nil {
		return GlobalInfo{}, err
	}

	addresses, ok := tables[Addresses]

	if !ok {
		return GlobalInfo{}, errors.New("address table not found")
	}

	addresses = removeHeaders(addresses)

	units, ok := tables[Units]

	if !ok {
		return GlobalInfo{}, errors.New("units table not found")
	}

	units = removeHeaders(units)

	ports, ok := tables[Ports]

	if !ok {
		return GlobalInfo{}, errors.New("ports table not found")
	}

	ports = removeHeaders(ports)

	boardIds, ok := tables[BoardIds]

	if !ok {
		return GlobalInfo{}, errors.New("boardIds table not found")
	}

	boardIds = removeHeaders(boardIds)

	messageIds, ok := tables[MessageIds]

	if !ok {
		return GlobalInfo{}, errors.New("messageIds table not found")
	}

	messageIds = removeHeaders(messageIds)

	return GlobalInfo{
		Addresses:  tableToMap(addresses),
		Units:      tableToMap(units),
		Ports:      tableToMap(ports),
		BoardIds:   tableToMap(boardIds),
		MessageIds: tableToMap(messageIds),
	}, nil
}

func checkAddressTable(addresses map[string]string) error {
	_, ok := addresses[BackendKey]

	if !ok {
		return errors.New("address table doesn't exists")
	}

	return checkIps(addresses)
}

func removeHeaders(table Table) Table {
	if len(table) == 0 {
		return table
	}

	return table[1:]
}

func checkIps(ips map[string]string) error {

	for _, ipStr := range ips {
		if !CheckIp(ipStr) {
			return fmt.Errorf("incorrect IP: %s", ipStr)
		}
	}

	return nil
}

func CheckIp(ipStr string) bool {
	return IpExp.MatchString(ipStr)
}

func checkUnits(units map[string]string) error {
	for _, unitStr := range units {
		if !CheckUnit(unitStr) {
			return fmt.Errorf("incorrect units: %s", unitStr)
		}
	}

	return nil
}

func CheckUnit(unitStr string) bool {
	if unitStr == "" {
		return true
	}

	return UnitExp.MatchString(unitStr)
}

func checkPorts(ports map[string]string) error {
	for _, portStr := range ports {
		if !CheckPort(portStr) {
			return fmt.Errorf("incorrect port: %s", portStr)
		}
	}

	return nil
}

func CheckPort(portStr string) bool {
	_, err := strconv.ParseUint(portStr, 10, 16)
	return err == nil
}

func CheckBoardIds(boardIds map[string]string, addresses map[string]string) error {
	for board, idStr := range boardIds {
		if err := checkBoardId(board, idStr, addresses); err != nil {
			return err
		}
	}

	return nil
}

func checkBoardId(name string, id string, addresses map[string]string) error {
	if !IdExp.MatchString(id) {
		return fmt.Errorf("incorrect board id: %s - %s", name, id)
	}

	_, ok := addresses[name]

	if !ok {
		return fmt.Errorf("%s IP missing in address table", name)
	}

	return nil
}

func checkMessageIds(messageIds map[string]string) error {
	for _, idStr := range messageIds {
		if !CheckMessageId(idStr) {
			return fmt.Errorf("incorrect message id: %s", idStr)
		}
	}

	return nil
}

func CheckMessageId(id string) bool {
	return IdExp.MatchString(id)
}

func tableToMap(table [][]Cell) map[string]string {
	m := make(map[string]string, len(table))

	for _, row := range table {
		m[row[0]] = row[1]
	}

	return m
}
