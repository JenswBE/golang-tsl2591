/**
 * tsl2591 - A command for interacting with TSL2591 lux sensors.
 */

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jimnelson2/tsl2591"
)

const Interval = 1 * time.Second

func main() {

	fmt.Println("start")
	tsl, err := tsl2591.NewTSL2591(&tsl2591.Opts{
		Gain:   tsl2591.TSL2591_GAIN_LOW,
		Timing: tsl2591.TSL2591_INTEGRATIONTIME_600MS,
	})
	if err != nil {
		panic(err)
	}
	defer tsl.Disable()

	fmt.Println("finish")

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

		<-ticker.C
	}

}
