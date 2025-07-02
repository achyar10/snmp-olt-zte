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
		return "", fmt.Errorf("error loading config: %w", err)
	}

	address := formatAddress(cfg.TelnetCfg.Ip, cfg.TelnetCfg.Port)
	log.Println("🔌 Connecting to", address)

	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return "", fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()
	log.Println("✅ Connected to", address)

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// Baca banner awal
	readAndLogRaw(conn, "initial-banner", 5*time.Second)

	// Login
	if err := writeAndLog(writer, cfg.TelnetCfg.Username); err != nil {
		return "", err
	}
	time.Sleep(300 * time.Millisecond)

	if err := writeAndLog(writer, cfg.TelnetCfg.Password); err != nil {
		return "", err
	}
	time.Sleep(500 * time.Millisecond)

	if err := writeAndLog(writer, ""); err != nil {
		return "", err
	}

	// Tunggu prompt
	if err := readUntil(conn, "GPON-D1-JKT-PSR#", 8*time.Second, "wait-prompt"); err != nil {
		return "", fmt.Errorf("login failed: %w", err)
	}

	// Kirim perintah satu per satu
	log.Println("🚀 Executing multiline command:")
	var result strings.Builder

	for _, line := range strings.Split(command, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		log.Printf("📤 Sending: %s", line)
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return "", fmt.Errorf("failed to write command: %w", err)
		}
		if err := writer.Flush(); err != nil {
			return "", fmt.Errorf("flush failed: %w", err)
		}

		// Set timeout untuk setiap respon
		conn.SetReadDeadline(time.Now().Add(3 * time.Second))

		// Baca respons baris demi baris
		for {
			respLine, err := reader.ReadString('\n')
			if err != nil {
				break // lanjut ke command berikutnya
			}
			decoded := decodeGBK(respLine)
			result.WriteString(decoded)
			log.Printf("📥 Output: %s", strings.TrimSpace(decoded))
		}

		time.Sleep(150 * time.Millisecond)
	}

	return result.String(), nil
}

func formatAddress(ip string, port uint16) string {
	if strings.Contains(ip, ":") && !strings.HasPrefix(ip, "[") {
		return fmt.Sprintf("[%s]:%d", ip, port) // IPv6 safe
	}
	return fmt.Sprintf("%s:%d", ip, port) // IPv4
}

func writeAndLog(writer *bufio.Writer, cmd string) error {
	log.Printf("📤 Sending: %s", strings.TrimSpace(cmd))
	_, err := writer.WriteString(cmd + "\n")
	if err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}
	return writer.Flush()
}
