//Package tsl2591 interacts with TSL2591 lux sensors
package tsl2591

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"

	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"
)

// Opts holds various configuration options for the sensor
type Opts struct {
	Gain   byte
	Timing byte
}

// TSL2591 holds board setup detail
type TSL2591 struct {
	enabled bool
	timing  byte
	gain    byte
	dev     i2c.Dev
}

// NewTSL2591 sets up a TSL2591 chip via the I2C protocol, sets its gain and timing
// attributes, and returns an error if any occurred in that process or if the
// TSL2591 was not found.
func NewTSL2591(opts *Opts) (*TSL2591, error) {

	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		return nil, err
	}

	// Open the first available I²C bus:
	bus, err := i2creg.Open("")
	if err != nil {
		return nil, err
	}

	// Address the device with address TSL2591_ADDR on the I²C bus:
	dev := i2c.Dev{Addr: TSL2591_ADDR, Bus: bus}
	tsl := &TSL2591{
		dev: dev,
	}

	// Read the device ID from the TSL2591. It should be 0x50
	write := []byte{TSL2591_COMMAND_BIT | TSL2591_REGISTER_DEVICE_ID}
	read := make([]byte, 1)
	if err := tsl.dev.Tx(write, read); err != nil {
		return nil, err
	}
	if read[0] != 0x50 {
		fmt.Printf("%v\n", read)
		return nil, errors.New("can't find a TSL2591 on I2C bus /dev/i2c-1")
	}

	err = tsl.SetTiming(opts.Timing)
	if err != nil {
		return nil, err
	}

	err = tsl.SetGain(opts.Gain)
	if err != nil {
		return nil, err
	}

	err = tsl.Enable()
	if err != nil {
		return nil, err
	}

	return tsl, nil
}

// Enable enables the TSL2591 chip
func (tsl *TSL2591) Enable() error {

	write := []byte{TSL2591_COMMAND_BIT | TSL2591_REGISTER_ENABLE |
		TSL2591_ENABLE_POWERON | TSL2591_ENABLE_AEN |
		TSL2591_ENABLE_AIEN | TSL2591_ENABLE_NPIEN}
	if _, err := tsl.dev.Write(write); err != nil {
		return err
	}

	tsl.enabled = true
	return nil
}

// Disable disables the TSL2591 chip
func (tsl *TSL2591) Disable() error {

	write := []byte{TSL2591_COMMAND_BIT | TSL2591_REGISTER_ENABLE |
		TSL2591_ENABLE_POWEROFF}
	if _, err := tsl.dev.Write(write); err != nil {
		return err
	}

	tsl.enabled = false
	return nil

}

// SetGain sets TSL2591 gain. Chip is enabled, gain set, then disabled
func (tsl *TSL2591) SetGain(gain byte) error {

	err := tsl.Enable()
	if err != nil {
		return err
	}

	write := []byte{TSL2591_COMMAND_BIT | TSL2591_REGISTER_ENABLE |
		tsl.timing | gain}
	if _, err := tsl.dev.Write(write); err != nil {
		return err
	}

	err = tsl.Disable()
	if err != nil {
		return err
	}

	tsl.gain = gain

	return nil
}

// SetTiming sets TSL2591 timing. Chip is enabled, timing set, then disabled
func (tsl *TSL2591) SetTiming(timing byte) error {
	err := tsl.Enable()
	if err != nil {
		return err
	}

	write := []byte{TSL2591_COMMAND_BIT | TSL2591_REGISTER_ENABLE |
		timing | tsl.gain}
	if _, err := tsl.dev.Write(write); err != nil {
		return err
	}

	err = tsl.Disable()
	if err != nil {
		return err
	}
	tsl.timing = timing
	return nil
}

// Read a 16-bit little-endian unsigned value from the specified 8-bit address
func (tsl *TSL2591) readU16(address byte) (uint16, error) {
	readBuffer := make([]byte, 2)
	cmd := []byte{TSL2591_COMMAND_BIT | address}
	if err := tsl.dev.Tx(cmd, readBuffer); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(readBuffer), nil
}

// RawLuminosity reads from the sensor
func (tsl *TSL2591) RawLuminosity() (uint16, uint16, error) {
	// The first value is IR + visible luminosity (channel 0)
	// and the second is the IR only (channel 1). Both values
	// are 16-bit unsigned numbers (0-65535)
	c0, err := tsl.readU16(TSL2591_REGISTER_CHAN0_LOW)
	if err != nil {
		return 0, 0, err
	}

	c1, err := tsl.readU16(TSL2591_REGISTER_CHAN1_LOW)
	if err != nil {
		return 0, 0, err
	}

	return c0, c1, nil
}

// FullSpectrum returns the full spectrum value
func (tsl *TSL2591) FullSpectrum() (uint32, error) {
	// Full spectrum (IR + visible) light and return its value
	// as a 32-bit unsigned number

	c0, c1, err := tsl.RawLuminosity()
	if err != nil {
		return 0, nil
	}

	return uint32(c1)<<16 | uint32(c0), nil

}

// Infrared returns infrared value
func (tsl *TSL2591) Infrared() (uint16, error) {
	_, c1, err := tsl.RawLuminosity()
	if err != nil {
		return 0, nil
	}
	return c1, nil
}

// Visible returns visible value
func (tsl *TSL2591) Visible() (uint32, error) {
	_, c1, err := tsl.RawLuminosity()
	if err != nil {
		return 0, nil
	}
	full := uint32(c1)<<16 | uint32(c1)
	return full - uint32(c1), nil
}

// Lux calculates a lux value from both the infrared and visible channels
func (tsl *TSL2591) Lux() (float64, error) {

	c0, c1, err := tsl.RawLuminosity()
	if err != nil {
		return 0, nil
	}

	// Compute the atime in milliseconds
	atime := 100.0*tsl.timing + 100.0

	// Set the maximum sensor counts based on the integration time (atime) setting
	var maxCounts uint16
	if tsl.timing == TSL2591_INTEGRATIONTIME_100MS {
		maxCounts = TSL2591_MAX_COUNT_100MS
	} else {
		maxCounts = TSL2591_MAX_COUNT
	}

	// Handle overflow.
	if c0 >= maxCounts || c1 >= maxCounts {
		return 0, errors.New("overflow reading light channels")
	}

	// Calculate lux
	// https://github.com/adafruit/Adafruit_TSL2591_Library/blob/master/Adafruit_TSL2591.cpp
	again := uint16(tsl.gain)
	cpl := float64((uint16(atime) * again)) / TSL2591_LUX_DF
	lux1 := (float64(c0) - (TSL2591_LUX_COEFB * float64(c1))) / cpl
	lux2 := ((TSL2591_LUX_COEFC * float64(c0)) - (TSL2591_LUX_COEFD * float64(c1))) / cpl

	return math.Max(lux1, lux2), nil

}
