package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/achyar10/snmp-olt-zte/internal/model"
)

func GetAvailableONUOnly(oltIndex string, max int) ([]model.ONUStatus, error) {

	slot, port, err := ParseOltIndex(oltIndex)
	if err != nil {
		return nil, err
	}

	cmd := fmt.Sprintf("show gpon onu state gpon-olt_1/%d/%d", slot, port)
	output, err := RunTelnetCommand(cmd)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(output, "\n")
	usedMap := make(map[int]bool)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "OnuIndex") || strings.Contains(line, "---") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 1 {
			continue
		}

		idxParts := strings.Split(parts[0], ":")
		if len(idxParts) != 2 {
			continue
		}

		onuID, err := strconv.Atoi(idxParts[1])
		if err != nil {
			continue
		}

		usedMap[onuID] = true
	}

	// Hanya ambil yang available
	var results []model.ONUStatus
	for i := 1; i <= max; i++ {
		if _, used := usedMap[i]; !used {
			results = append(results, model.ONUStatus{
				ID:     i,
				Status: "available",
			})
		}
	}

	return results, nil
}

func ParseOltIndex(index string) (int, int, error) {
	re := regexp.MustCompile(`gpon-olt_1/(\d+)/(\d+)`)
	match := re.FindStringSubmatch(index)
	if len(match) != 3 {
		return 0, 0, fmt.Errorf("invalid olt_index format: %s", index)
	}

	slot, err1 := strconv.Atoi(match[1])
	port, err2 := strconv.Atoi(match[2])
	if err1 != nil || err2 != nil {
		return 0, 0, fmt.Errorf("invalid numbers in olt_index: %s", index)
	}
	return slot, port, nil
}

func BuildZTERegisterCommand(slot, port int, region, serialNumber, code string, onu, vlanID int) string {
	return fmt.Sprintf(
		`con t
interface gpon-olt_1/%d/%d
onu %d type ALL sn %s
exit
interface gpon-onu_1/%d/%d:%d
name %s
description zone %s
tcont 3 profile 10m
gemport 1 tcont 1
gemport 1 traffic-limit upstream 100m downstream 100m
service-port 1 vport 1 user-vlan %d vlan %d
service-port 2 vport 1 user-vlan 100 vlan 100
exit

pon-onu-mng gpon-onu_1/%d/%d:%d
service 1 gemport 1 vlan %d
service 2 gemport 1 vlan 100
security-mgmt 212 state enable mode forward protocol web
wan-ip 1 mode pppoe username %s password aba vlan-profile netmedia143 host 1
wan 1 service internet host 1
end
wr`,
		slot, port,
		onu, serialNumber,
		slot, port, onu,
		code, region,
		vlanID, vlanID,
		slot, port, onu,
		vlanID,
		code,
	)
}

func BuildZTERebootCommand(slot, port, onu int) string {
	return fmt.Sprintf(
		`config terminal
pon-onu
pon-onu-mng gpon-onu_1/%d/%d:%d
reboot`,
		slot, port, onu,
	)
}

func BuildZTERemoveCommand(slot, port, onu int) string {
	return fmt.Sprintf(
		`conf t
interface gpon-olt_1/%d/%d
no onu %d
end
wr`,
		slot, port, onu,
	)
}
