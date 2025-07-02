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

func BuildZTERegisterCommand(slot, port int, region, serialNumber, code string, onu int) string {
	return fmt.Sprintf(
		`con t
interface gpon-olt_1/%d/%d
onu %d type ZTE sn %s
exit
interface gpon-onu_1/%d/%d:%d
name %s
description zone %s
tcont 3 profile 1G
gemport 1 tcont 3
service-port 1 vport 1 user-vlan 800 vlan 800
service-port 2 vport 1 user-vlan 100 vlan 100
exit

pon-onu-mng gpon-onu_1/%d/%d:%d
service 1 gemport 1 vlan 800
service TR069 gemport 1 vlan 100
wan-ip 2 mode dhcp vlan-profile 100 host 2
tr069-mgmt 1 acs http://10.0.0.3:7547
security-mgmt 212 state enable mode forward protocol web
wan-ip 1 mode pppoe username %s password 101094 vlan-profile 800 host 1
wan 1 ethuni 1 ssid 1 service internet host 1
end
wr`,
		slot, port,
		onu, serialNumber,
		slot, port, onu,
		code, region,
		slot, port, onu,
		code,
	)
}
