# tsl2591

This is a Go module for the TSL2591 lux sensor, available from http://www.adafruit.com/products/1980 among other places.

## Why this exists

We already have http://github.com/mstahl/tsl2591, but I wanted to make a few changes

* Use modules
* Switch from the use of the deprecated `golang.org/x/exp/io/i2c` import to the recommended `periph.io/x/periph`
* I had a small bit of trouble working with the original, so I figured I'd have better success re-writing from scratch as that would force me to learn a bit more. That being said, a *significant* amount of code in this module has been copied from other sources as noted in the package header, i.e. Adafruit's original cpp and python implementations, as well as mstahl's work.

Note this module does NOT provide FULL control of the tsl2581, i.e. interrupts and alerts have not been exposed. That being said, it does all I need.

* enable/disable
* set gain
* set timing
* read visible, IR, and full spectrum
* calculate lux

If anyone would like to implement missing functionality, or discovers problems with what's here - I sure there are problems somewhere - please submit and issue or even better, a PR.

## Example Usage

Import this module and retrieve a lux value, e.g.

```go
import ("github.com/jimnelson2/tsl2591")

    // connect the the tsl2591
	tsl, err := tsl2591.NewTSL2591(&tsl2591.Opts{
		Gain:   tsl2591.GainMed,
		Timing: tsl2591.Integrationtime600MS,
	})
	if err != nil {
		panic(err)
	}

	// read lux
	lux, _ := tsl.Lux()
```

## Sample code

Sample code is [here](cmd/tsl2591/tsl2591.go), intended for use on a Raspberry Pi Zero.


Compile the code - in this case for Raspberry Pi Zero

```sh
env GOOS=linux GOARCH=arm GOARM=5 go build -o tsltest cmd/tsl2591/tsl2591.go
```

And when executed on a RPi with a connected TSL2591, output like the following will be printed every second.

```
Total Light: 12.451616 lux
Infrared light: 1921
Visible light: 125894656
Full spectrum (IR + visible) light: 125898232
Total Light: 12.426864 lux
Infrared light: 1920
Visible light: 125829120
Full spectrum (IR + visible) light: 125832693
```

## Acknowledgements

As noted above and in the code - substantial work from the following is included from 

* https://github.com/mstahl/tsl2591
* https://github.com/adafruit/Adafruit_TSL2591_Library/
* https://github.com/adafruit/Adafruit_TSL2591_Library/blob/master/Adafruit_TSL2591.cpp

