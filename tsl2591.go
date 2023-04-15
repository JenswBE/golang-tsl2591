// Package tsl2591 interacts with TSL2591 lux sensors
//
// Heavily inspired by https://github.com/mstahl/tsl2591
// and https://github.com/adafruit/Adafruit_TSL2591_Library/
// as well as https://github.com/adafruit/Adafruit_TSL2591_Library/blob/master/Adafruit_TSL2591.cpp
package tsl2591

import (
	"fmt"
	"math"

	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
)

// Opts holds various configuration options for the sensor
type Opts struct {
	// Bus name, alias or its number.
	// See https://pkg.go.dev/periph.io/x/conn/v3/i2c/i2creg#Open for more info.
	Bus    string
	Gain   Gain
	Timing IntegrationTime
}

func DefaultOptions() *Opts {
	return &Opts{
		Bus:    "",
		Gain:   GainMed,
		Timing: IntegrationTime100MS,
	}
}

// TSL2591 holds board setup detail
type TSL2591 struct {
	dev    i2c.Dev
	gain   Gain
	timing IntegrationTime
}

// NewTSL2591 sets up a TSL2591 chip via the I2C protocol, sets its gain and timing
// attributes, and returns an error if any occurred in that process or if the
// TSL2591 was not found
func NewTSL2591(opts *Opts) (*TSL2591, error) {
	// Use default opts if not set
	if opts == nil {
		opts = DefaultOptions()
	}

	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		return nil, fmt.Errorf("unable to init host: %w", err)
	}

	// Open the first available I2C bus:
	bus, err := i2creg.Open(opts.Bus)
	if err != nil {
		return nil, fmt.Errorf("unable to open I2C bus: %w", err)
	}

	// Address the device with address TSL2591_ADDR on the I2C bus:
	dev := i2c.Dev{Addr: Addr, Bus: bus}
	tsl := &TSL2591{dev: dev}

	// Read the device ID from the TSL2591. It should be 0x50.
	deviceID, err := tsl.readU8(RegisterDeviceID)
	if err != nil {
		return nil, fmt.Errorf("unable to read device ID from I2C bus: %w", err)
	}
	if deviceID != DeviceID {
		return nil, UnexpectedDeviceIDError{Actual: deviceID, Expected: DeviceID}
	}

	if err = tsl.SetGain(opts.Gain); err != nil {
		return nil, fmt.Errorf("unable to set gain: %w", err)
	}

	if err = tsl.SetTiming(opts.Timing); err != nil {
		return nil, fmt.Errorf("unable to set timing: %w", err)
	}

	if err = tsl.Enable(); err != nil {
		return nil, fmt.Errorf("unable to enable sensor: %w", err)
	}

	return tsl, nil
}

// Enable enables the TSL2591 chip
func (tsl *TSL2591) Enable() error {
	err := tsl.writeU8(RegisterEnable, EnablePowerOn|EnableAEN|EnableAIEN|EnableNPIEN)
	if err != nil {
		return fmt.Errorf("failed to enable sensor: %w", err)
	}
	return nil
}

// Disable disables the TSL2591 chip
func (tsl *TSL2591) Disable() error {
	err := tsl.writeU8(RegisterEnable, EnablePowerOff)
	if err != nil {
		return fmt.Errorf("failed to disable sensor: %w", err)
	}
	return nil
}

// SetGain sets TSL2591 gain
func (tsl *TSL2591) SetGain(gain Gain) error {
	// Get control
	control, err := tsl.readU8(RegisterControl)
	if err != nil {
		return fmt.Errorf("failed to read current sensor control: %w", err)
	}

	// Update control
	control &= 0b11001111
	control |= byte(gain)

	// Write control
	if err = tsl.writeU8(RegisterControl, control); err != nil {
		return fmt.Errorf("failed to write sensor control: %w", err)
	}
	tsl.gain = gain
	return nil
}

// SetTiming sets TSL2591 timing. Chip is enabled, timing set, then disabled
func (tsl *TSL2591) SetTiming(timing IntegrationTime) error {
	// Get control
	control, err := tsl.readU8(RegisterControl)
	if err != nil {
		return fmt.Errorf("failed to read current sensor control: %w", err)
	}

	// Update control
	control &= 0b11111000
	control |= byte(timing)

	// Write control
	if err = tsl.writeU8(RegisterControl, control); err != nil {
		return fmt.Errorf("failed to write sensor control: %w", err)
	}
	tsl.timing = timing
	return nil
}

// RawLuminosity reads from the sensor
func (tsl *TSL2591) RawLuminosity() (uint16, uint16, error) {
	// The first value is IR + visible luminosity (channel 0)
	// and the second is the IR only (channel 1). Both values
	// are 16-bit unsigned numbers (0-65535)
	c0, err := tsl.readU16(RegisterChan0Low)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read channel 0 of raw luminosity: %w", err)
	}

	c1, err := tsl.readU16(RegisterChan1Low)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read channel 1 of raw luminosity: %w", err)
	}

	return c0, c1, nil
}

// FullSpectrum returns the full spectrum value
func (tsl *TSL2591) FullSpectrum() (uint32, error) {
	// Full spectrum (IR + visible) light and return its value
	// as a 32-bit unsigned number
	c0, c1, err := tsl.RawLuminosity()
	if err != nil {
		return 0, err
	}
	return uint32(c1)<<16 | uint32(c0), nil
}

// Infrared returns infrared value
func (tsl *TSL2591) Infrared() (uint16, error) {
	_, c1, err := tsl.RawLuminosity()
	if err != nil {
		return 0, err
	}
	return c1, nil
}

// Visible returns visible value
func (tsl *TSL2591) Visible() (uint32, error) {
	c0, c1, err := tsl.RawLuminosity()
	if err != nil {
		return 0, err
	}
	full := uint32(c1)<<16 | uint32(c0)
	return full - uint32(c1), nil
}

// Lux calculates a lux value from both the infrared and visible channels
func (tsl *TSL2591) Lux() (float64, error) {
	c0, c1, err := tsl.RawLuminosity()
	if err != nil {
		return 0, err
	}

	// Compute the atime in milliseconds
	atime := 100*uint16(tsl.timing) + 100

	// Set the maximum sensor counts based on the integration time (atime) setting
	var maxCounts uint16
	if tsl.timing == IntegrationTime100MS {
		maxCounts = MaxCount100ms
	} else {
		maxCounts = MaxCount
	}

	// Handle overflow.
	if c0 >= maxCounts || c1 >= maxCounts {
		return 0, ErrOverflow
	}

	// Calculate lux
	var again uint16
	switch tsl.gain {
	case GainLow:
		again = 1
	case GainMed:
		again = 25
	case GainHigh:
		again = 428
	case GainMax:
		again = 9876
	}

	cpl := float64((atime * again)) / LuxDF
	lux1 := (float64(c0) - (LuxCoefB * float64(c1))) / cpl
	lux2 := ((LuxCoefC * float64(c0)) - (LuxCoefD * float64(c1))) / cpl

	return math.Max(lux1, lux2), nil
}
