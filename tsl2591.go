//Package tsl2591 interacts with TSL2591 lux sensors
package tsl2591

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"time"

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

	err = tsl.Disable()
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

// GetFullLuminosity returns both visible and IR channel luminosity
func (tsl *TSL2591) GetFullLuminosity() (uint16, uint16, error) {
	err := tsl.Enable()
	if err != nil {
		return 0, 0, err
	}

	// Delay for ADC to complete
	for d := byte(0); d < tsl.timing; d++ {
		time.Sleep(120 * time.Millisecond)
	}

	bytes := make([]byte, 4)

	write := []byte{TSL2591_COMMAND_BIT | TSL2591_REGISTER_CHAN0_LOW}
	if err := tsl.dev.Tx(write, bytes); err != nil {
		return 0, 0, err
	}

	channel0 := binary.LittleEndian.Uint16(bytes[0:])
	channel1 := binary.LittleEndian.Uint16(bytes[2:])

	err = tsl.Disable()
	if err != nil {
		return 0, 0, err
	}

	return channel0, channel1, nil
}

// CalculateLux calculates lux from the provided intensities
func (tsl *TSL2591) CalculateLux(ch0, ch1 uint16) float64 {
	var (
		atime float64
		again float64

		cpl float64
		lux float64
	)

	// Return +Inf for overflow
	if ch0 == 0xFFFF || ch1 == 0xFFFF {
		return math.Inf(1)
	}

	switch tsl.timing {
	case TSL2591_INTEGRATIONTIME_100MS:
		atime = 100.0
	case TSL2591_INTEGRATIONTIME_200MS:
		atime = 200.0
	case TSL2591_INTEGRATIONTIME_300MS:
		atime = 300.0
	case TSL2591_INTEGRATIONTIME_400MS:
		atime = 400.0
	case TSL2591_INTEGRATIONTIME_500MS:
		atime = 500.0
	case TSL2591_INTEGRATIONTIME_600MS:
		atime = 600.0
	default:
		atime = 100.0
	}

	switch tsl.gain {
	case TSL2591_GAIN_LOW:
		again = 1.0
	case TSL2591_GAIN_MED:
		again = 25.0
	case TSL2591_GAIN_HIGH:
		again = 428.0
	case TSL2591_GAIN_MAX:
		again = 9876.0
	default:
		again = 1.0
	}

	cpl = (atime * again) / TSL2591_LUX_DF
	lux = (float64(ch0) - float64(ch1)) * (1.0 - (float64(ch1) / float64(ch0))) / cpl

	return lux
}
