package utils

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/achyar10/snmp-olt-zte/config"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func decodeGBK(input string) string {
	reader := transform.NewReader(strings.NewReader(input), simplifiedchinese.GBK.NewDecoder())
	decoded, err := bufio.NewReader(reader).ReadString('\n')
	if err != nil {
		return input
	}
	return decoded
}

func readUntil(conn net.Conn, expect string, timeout time.Duration, label string) error {
	conn.SetReadDeadline(time.Now().Add(timeout))
	reader := bufio.NewReader(conn)

	start := time.Now()
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("[%s] read error: %v", label, err)
		}
		lineDecoded := decodeGBK(line)
		log.Printf("[%s] %s", label, strings.TrimSpace(lineDecoded))
		if strings.Contains(line, expect) {
			log.Printf("[%s] found: %s (%.2fs)", label, expect, time.Since(start).Seconds())
			break
		}
	}
	return nil
}

func readAndLogRaw(conn net.Conn, label string, maxTime time.Duration) {
	conn.SetReadDeadline(time.Now().Add(maxTime))
	reader := bufio.NewReader(conn)
	log.Println("📡 Raw output start:")
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		log.Printf("[%s] %s", label, decodeGBK(strings.TrimSpace(line)))
	}
	log.Println("📡 Raw output end")
}

func RunTelnetCommand(command string) (string, error) {
	configPath := GetConfigPath(os.Getenv("APP_ENV"))
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	address := fmt.Sprintf("%s:23", cfg.TelnetCfg.Ip)
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return "", fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()
	log.Println("✅ Connected to", address)

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	readAndLogRaw(conn, "initial-banner", 5*time.Second)

	writer.WriteString(cfg.TelnetCfg.Username + "\n")
	writer.Flush()
	time.Sleep(1 * time.Second)

	writer.WriteString(cfg.TelnetCfg.Password + "\n")
	writer.Flush()
	time.Sleep(1 * time.Second)

	writer.WriteString("\n")
	writer.Flush()

	if err := readUntil(conn, "GPON-D1-JKT-PSR#", 8*time.Second, "wait-prompt"); err != nil {
		return "", fmt.Errorf("login failed: %w", err)
	}

	writer.WriteString(command + "\n")
	writer.Flush()

	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	var result strings.Builder
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		result.WriteString(decodeGBK(line))
	}

	return result.String(), nil
}
