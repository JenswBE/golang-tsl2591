// tsl2591 - A command for interacting with TSL2591 lux sensors.

package main

import (
	"flag"
	"log"
	"time"

	tsl2591 "github.com/JenswBE/golang-tsl2591"
)

const Interval = 1 * time.Second

func main() {
	bus := flag.String("bus", "", "Name of the bus")
	flag.Parse()

	opts := tsl2591.DefaultOptions()
	opts.Bus = *bus
	tsl, err := tsl2591.NewTSL2591(opts)
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		if disableErr := tsl.Disable(); disableErr != nil {
			log.Panic(err)
		}
	}()

	ticker := time.NewTicker(Interval)

	for {
		lux, err := tsl.Lux()
		if err != nil {
			log.Panic(err)
		}
		log.Printf("Total Light: %f lux\n", lux)

		ir, err := tsl.Infrared()
		if err != nil {
			log.Panic(err)
		}
		log.Printf("Infrared light: %d\n", ir)

		visible, err := tsl.Visible()
		if err != nil {
			log.Panic(err)
		}
		log.Printf("Visible light: %d\n", visible)

		full, err := tsl.FullSpectrum()
		if err != nil {
			log.Panic(err)
		}
		log.Printf("Full spectrum (IR + visible) light: %d\n", full)

		chan0, chan1, err := tsl.RawLuminosity()
		if err != nil {
			log.Panic(err)
		}
		log.Printf("Raw luminosity: %b (chan0), %b (chan1)\n\n", chan0, chan1)

		<-ticker.C
	}
}
