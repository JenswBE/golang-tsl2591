package tsl2591

import (
	"encoding/binary"
	"fmt"
)

// readU8 reads an 8-bit unsigned value from the specified 8-bit address.
func (tsl *TSL2591) readU8(address byte) (uint8, error) {
	readBuffer := make([]byte, 1)
	cmd := []byte{CommandBit | address}
	if err := tsl.dev.Tx(cmd, readBuffer); err != nil {
		return 0, fmt.Errorf("failed to read uint8: %w", err)
	}
	return readBuffer[0], nil
}

// writeU8 writes an 8-bit unsigned value to the specified 8-bit address.
func (tsl *TSL2591) writeU8(address, value byte) error {
	data := []byte{
		CommandBit | address,
		value,
	}
	if _, err := tsl.dev.Write(data); err != nil {
		return fmt.Errorf("failed to write uint8 %x to address %x: %w", value, address, err)
	}
	return nil
}

// readU16 reads a 16-bit little-endian unsigned value from the specified 8-bit address
func (tsl *TSL2591) readU16(address byte) (uint16, error) {
	readBuffer := make([]byte, 2)
	cmd := []byte{CommandBit | address}
	if err := tsl.dev.Tx(cmd, readBuffer); err != nil {
		return 0, fmt.Errorf("failed to read uint16: %w", err)
	}
	return binary.LittleEndian.Uint16(readBuffer), nil
}
