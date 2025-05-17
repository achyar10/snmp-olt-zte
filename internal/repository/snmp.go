package repository

import (
	"fmt"
	"time"

	"github.com/gosnmp/gosnmp"
)

// SnmpRepositoryInterface is an interface that represents the SNMP repository contract
type SnmpRepositoryInterface interface {
	Get(oids []string) (result *gosnmp.SnmpPacket, err error)       // Get SNMP data for the given OIDs
	Walk(oid string, walkFunc func(pdu gosnmp.SnmpPDU) error) error // Walk SNMP to get all OIDs under the given OID
}

// snmpRepository is a struct that implements SnmpRepositoryInterface
type snmpRepository struct {
	target    string // SNMP target IP address
	community string // SNMP community string
	port      uint16 // SNMP port number
}

// NewPonRepository is a constructor function to create a new instance of snmpRepository
func NewPonRepository(target string, community string, port uint16) SnmpRepositoryInterface {
	return &snmpRepository{
		target:    target,    // SNMP target IP address
		community: community, // SNMP community string
		port:      port,      // SNMP port number
	}
}

// buildSNMPInstance for creating a new SNMP instance
func (r *snmpRepository) buildSNMPInstance() (*gosnmp.GoSNMP, error) {
	params := &gosnmp.GoSNMP{
		Target:    r.target,                       // SNMP target IP address
		Port:      r.port,                         // SNMP port number
		Community: r.community,                    // SNMP community string
		Version:   gosnmp.Version2c,               // SNMP version
		Timeout:   time.Duration(3) * time.Second, // SNMP timeout
		Retries:   1,                              // Number of retries for SNMP requests
	}

	// Set logger to nil to disable logging
	if err := params.Connect(); err != nil {
		return nil, fmt.Errorf("SNMP Connect error: %w", err) // Error connecting to SNMP target
	}
	return params, nil // Return the SNMP instance
}

// Get to get SNMP data for the given OIDs
func (r *snmpRepository) Get(oids []string) (*gosnmp.SnmpPacket, error) {
	snmp, err := r.buildSNMPInstance() // Create a new SNMP instance
	if err != nil {
		return nil, err
	}
	defer snmp.Conn.Close()

	result, err := snmp.Get(oids)
	if err != nil {
		return nil, fmt.Errorf("SNMP Get failed: %w", err)
	}
	return result, nil
}

// Walk for SNMP Walk to get all OIDs under the given OID
func (r *snmpRepository) Walk(oid string, walkFunc func(pdu gosnmp.SnmpPDU) error) error {
	snmp, err := r.buildSNMPInstance()
	if err != nil {
		return err
	}
	defer snmp.Conn.Close()

	err = snmp.Walk(oid, walkFunc)
	if err != nil {
		return fmt.Errorf("SNMP Walk failed: %w", err)
	}
	return nil
}
