package ade_linter

import (
	"fmt"
	"strconv"
)

func checkPackets(table Table, logger Logger) bool {
	packets := getPackets(table[1:], logger)
	return CheckAll(packets)
}

func getPackets(table Table, logger Logger) []Test {
	return Map(table, func(item []string) Test {
		packetLogger := logger.AddEntry("Packet", item[1])
		return Test(Packet(Packet{
			cells:  item,
			logger: packetLogger,
		}))
	})
}

type Packet struct {
	cells  []string
	logger Logger
}

func (packet Packet) Run() bool {
	if !CheckId(packet.cells[0]) {
		packet.logger.Error(fmt.Errorf("packet id is not valid: %s", packet.cells[0]))
		return false
	}

	if !CheckPacketType(packet.cells[2]) {
		packet.logger.Error(fmt.Errorf("packet type invalid: %s", packet.cells[2]))
		return false
	}

	return true
}

func CheckId(id string) bool {
	_, err := strconv.ParseUint(id, 10, 16)
	return err == nil
}

func CheckPacketType(kind string) bool {
	return kind == "data" || kind == "order" || kind == "warning" || kind == "fault"
}
