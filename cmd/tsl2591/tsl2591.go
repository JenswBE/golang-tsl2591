/**
 * tsl2591 - A command for interacting with TSL2591 lux sensors.
 */

package main

import (
	"fmt"
	"log"
	"time"

	tsl2591 "github.com/JenswBE/golang-tsl2591"
)

const Interval = 1 * time.Second

func main() {

	tsl, err := tsl2591.NewTSL2591(nil)
	if err != nil {
		panic(err)
	}
	defer tsl.Disable()

	ticker := time.NewTicker(Interval)

	for {
		lux, err := tsl.Lux()
		if err != nil {
			log.Panic(err)
		}
		fmt.Printf("Total Light: %f lux\n", lux)

		ir, err := tsl.Infrared()
		if err != nil {
			log.Panic(err)
		}
		fmt.Printf("Infrared light: %d\n", ir)

		visible, err := tsl.Visible()
		if err != nil {
			log.Panic(err)
		}
		fmt.Printf("Visible light: %d\n", visible)

		full, err := tsl.FullSpectrum()
		if err != nil {
			log.Panic(err)
		}
		fmt.Printf("Full spectrum (IR + visible) light: %d\n", full)

		chan0, chan1, err := tsl.RawLuminosity()
		if err != nil {
			log.Panic(err)
		}
		fmt.Printf("Raw luminosity: %b (chan0), %b (chan1)\n", chan0, chan1)

		fmt.Println()

		<-ticker.C
	}

}
