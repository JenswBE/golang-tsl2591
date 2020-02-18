// TODO JIM YOU ARE HERE
I think our starting code had the wrong size for visible light - should be 32-bit, same for full spectrum?

// Python
Total light: 0.762688lux
Infrared light: 14
Visible light: 917541
Full spectrum (IR + visible) light: 917555
Total light: 0.762688lux
Infrared light: 14
Visible light: 917541
Full spectrum (IR + visible) light: 917555
Total light: 0.762688lux
Infrared light: 14
Visible light: 917541
Full spectrum (IR + visible) light: 917555

// Go
Total Light: 18.253333 lux
Visible: 51
Infrared: 14
Total Light: 18.253333 lux
Visible: 51
Infrared: 14
Total Light: 18.253333 lux
Visible: 51
Infrared: 14
Total Light: 18.253333 lux
Visible: 51
Infrared: 14

This is a Golang driver for the TSL2591 lux sensor.

## Installation

    go get -u github.com/mstahl/tsl2591

## Usage

    import "github.com/mstahl/tsl2591"

For now, `tsl2591` only supports retrieving luminosity data, so no interrupts
or alerts yet.

```go
	tsl, err := tsl2591.NewTSL2591(&tsl2591.Opts{
		Gain:   tsl2591.TSL2591_GAIN_LOW,
		Timing: tsl2591.TSL2591_INTEGRATIONTIME_600MS,
	})
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(1*time.Second)

	for {
		channel0, channel1 := tsl.GetFullLuminosity()
		log.Printf("0x%04x 0x%04x\n", channel0, channel1)
		<-ticker.C
	}
```

## Acknowledgements

This library is basically a golang port of [Adafruit's TSL2591 library](https://github.com/adafruit/Adafruit_TSL2591_Library/)
