package snmp

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/achyar10/snmp-olt-zte/config"
	"github.com/achyar10/snmp-olt-zte/internal/utils"
	"github.com/gosnmp/gosnmp"
)

var (
	snmpHost      string // SNMP host
	snmpPort      uint16 // SNMP port
	snmpCommunity string // SNMP community
	//logSnmp       gosnmp.Logger // Logger for SNMP
)

// SetupSnmpConnection is a function to set up snmp connection
func SetupSnmpConnection(config *config.Config) (*gosnmp.GoSNMP, error) {
	var logSnmp gosnmp.Logger

	// Check if the application is running in development or production environment
	if os.Getenv("APP_ENV") == "development" || os.Getenv("APP_ENV") == "production" {
		snmpHost = os.Getenv("SNMP_HOST")
		snmpPort = utils.ConvertStringToUint16(os.Getenv("SNMP_PORT"))
		snmpCommunity = os.Getenv("SNMP_COMMUNITY")
		logSnmp = gosnmp.Logger{}
	} else {
		snmpHost = config.SnmpCfg.Ip
		snmpPort = config.SnmpCfg.Port
		snmpCommunity = config.SnmpCfg.Community
		logSnmp = gosnmp.NewLogger(log.New(os.Stdout, "", 0))
	}

	// Check if SNMP configuration is valid
	if snmpHost == "" || snmpPort == 0 || snmpCommunity == "" {
		return nil, fmt.Errorf("konfigurasi SNMP tidak valid")
	}

	// Create a new SNMP target instance
	target := &gosnmp.GoSNMP{
		Target:    snmpHost,
		Port:      snmpPort,
		Community: snmpCommunity,
		Version:   gosnmp.Version2c,
		Timeout:   time.Duration(30) * time.Second,
		Retries:   3,
		Logger:    logSnmp,
	}

	// Connect to the SNMP target
	err := target.Connect()
	if err != nil {
		return nil, fmt.Errorf("gagal terhubung ke SNMP: %w", err)
	}

	return target, nil
}
