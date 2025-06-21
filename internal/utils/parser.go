package utils

import (
	"strings"

	"github.com/achyar10/snmp-olt-zte/internal/model"
)

func ParseONULineOutput(output string) []model.ONUItem {
	// Tangani jika output error dari ZTE OLT
	if strings.Contains(output, "No related information to show") || strings.HasPrefix(strings.TrimSpace(output), "%Code") {
		return []model.ONUItem{} // kosong, tidak ada ONU terdeteksi
	}

	lines := strings.Split(output, "\n")
	var results []model.ONUItem

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.Contains(line, "---") || strings.HasPrefix(line, "OltIndex") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 3 {
			results = append(results, model.ONUItem{
				OltIndex:     parts[0],
				Model:        parts[1],
				SerialNumber: parts[2],
				Status:       "unactivated",
			})
		}
	}
	return results
}
