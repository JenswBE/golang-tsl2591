package tsl2591

// General purpose constants
const (
	// FullSpectrum is channel 0
	FullSpectrum byte = 0

	// Infrared is channel 1
	Infrared byte = 1

	// Visible is FullSpectrum minus Infrared, i.e. channel 0 - channel 1
	Visible byte = 2

	// Addr is the default I2C address for the TSL2591
	Addr uint16 = 0x29

	// CommandBits is 1010 0000 - sets bits 7 and 5 to indicate 'command normal'
	CommandBit byte = 0xa0

	// ClearInt command for 'Clear ALS and no persist ALS interrupt'
	ClearInt byte = 0xe7

	// TestInt command for 'Interrupt set - forces an interrupt'
	TestInt byte = 0xe4

	// WordBit to read/write word rather than byte
	WordBit byte = 0x20

	// BlockBit to block read/write
	BlockBit byte = 0x10

	// EnablePowerOff to set 'enable' register to disabled
	EnablePowerOff byte = 0x00

	// EnablePowerOn to set 'enable' register to enabled
	EnablePowerOn byte = 0x01

	// EnableAEN commands the ALS function. 1 enables, 0 disables
	EnableAEN byte = 0x02

	// EnableAIEN permits ALS interrupts to be generated, subject to the persist filter
	EnableAIEN byte = 0x10

	// EnableNPIEN commands that NP Threshold conditions will generate an interrupt, bypassing the persist filter
	EnableNPIEN byte = 0x80

	// LuxDF is the Lux cooefficient
	LuxDF float64 = 408.0

	// LuxCoefB is the channel0 coefficient
	LuxCoefB float64 = 1.64

	// LuxCoefC is channel1 coefficient A
	LuxCoefC float64 = 0.59

	// LuxCoefD is channel2 coefficient B
	LuxCoefD      float64 = 0.86
	MaxCount100ms uint16  = 0x8fff
	MaxCount      uint16  = 0xffff
)

// Register maps
const (
	// RegisterEnable is the enable register
	RegisterEnable byte = 0x00

	// RegisterControl is the control register
	RegisterControl byte = 0x01

	// RegisterThresholdAILTL is the ALS low threshold lower byte
	RegisterThresholdAILTL byte = 0x04

	// RegisterThresholdAILTH is the ALS low threshold upper byte
	RegisterThresholdAILTH byte = 0x05

	// RegisterThresholdAIHTL is the ALS high threshold lower byte
	RegisterThresholdAIHTL byte = 0x06

	// RegisterThresholdAIHTH is the ALS high threshold upper byte
	RegisterThresholdAIHTH byte = 0x07

	// RegisterThresholdNPAILTL is the no-persist ALS low threshold lower byte
	RegisterThresholdNPAILTL byte = 0x08

	// RegisterThresholdNPAILTH is the no-persist ALS low threshold higher byte
	RegisterThresholdNPAILTH byte = 0x09

	// RegisterThresholdNPAIHTL is the no-persist ALS high threshold lower byte
	RegisterThresholdNPAIHTL byte = 0x0a

	// RegisterThresholdNPAIHTH is the no-persist ALS high threshold higher byte
	RegisterThresholdNPAIHTH byte = 0x0b

	// RegisterPersistFilter is the interrupt persistence filter
	RegisterPersistFilter byte = 0x0c

	// RegisterPackagePID is for package identification
	RegisterPackagePID byte = 0x11

	// RegisterDeviceID is for device identification
	RegisterDeviceID byte = 0x12

	// RegisterDeviceStatus is for internal status
	RegisterDeviceStatus byte = 0x13

	// RegisterChan0Low is channel 0 data, low byte
	RegisterChan0Low byte = 0x14

	// RegisterChan0High is channel 0 data, high byte
	RegisterChan0High byte = 0x15

	// RegisterChan1Low is channel 1 data, low byte
	RegisterChan1Low byte = 0x16

	// RegisterChan1High is channel 1 data, high byte
	RegisterChan1High byte = 0x17
)

// Constants for sensor integration timing
const (
	// Integrationtime100MS is 100 millis
	Integrationtime100MS byte = 0x00

	// Integrationtime200MS is 200 millis
	Integrationtime200MS byte = 0x01

	// Integrationtime300MS is 300 millis
	Integrationtime300MS byte = 0x02

	// Integrationtime400MS is 400 millis
	Integrationtime400MS byte = 0x03

	// Integrationtime500MS is 500 millis
	Integrationtime500MS byte = 0x04

	// Integrationtime600MS is 600 millis
	Integrationtime600MS byte = 0x05
)

// Constants for adjusting the persistance filter
const (
	// PersistEvery is every ALS cycle generates an interrupt
	PersistEvery byte = 0x00

	// PersistAny for any value outside of threshold range
	PersistAny byte = 0x01

	// Persist2 for 2 consecutive values out of range
	Persist2 byte = 0x02

	// Persist3 for 3 consecutive values out of range
	Persist3 byte = 0x03

	// Persist5 for 5 consecutive values out of range
	Persist5 byte = 0x04

	// Persist10 for 10 consecutive values out of range
	Persist10 byte = 0x05

	// Persist15 for 15 consecutive values out of range
	Persist15 byte = 0x06

	// Persist20 for 20 consecutive values out of range
	Persist20 byte = 0x07

	// Persist25 for 25 consecutive values out of range
	Persist25 byte = 0x08

	// Persist30 for 30 consecutive values out of range
	Persist30 byte = 0x09

	// Persist35 for 35 consecutive values out of range
	Persist35 byte = 0x0a

	// Persist40 for 40 consecutive values out of range
	Persist40 byte = 0x0b

	// Persist45 for 45 consecutive values out of range
	Persist45 byte = 0x0c

	// Persist50 for 50 consecutive values out of range
	Persist50 byte = 0x0d

	// Persist55 for 55 consecutive values out of range
	Persist55 byte = 0x0e

	// Persist60 for 60 consecutive values out of range
	Persist60 byte = 0x0f
)

// Constants for adjusting the sensor gain
const (
	// GainLow is low gain (1x)
	GainLow byte = 0x00

	// GainMed is medium gain (25x)
	GainMed byte = 0x10

	// GainHigh is high gain (428x)
	GainHigh byte = 0x20

	// GainMax is max gain (9876x)
	GainMax byte = 0x30
)
