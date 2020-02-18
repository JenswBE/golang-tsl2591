/**
 * tsl2591 - A command for interacting with TSL2591 lux sensors.
 */

package main

import (
	"fmt"
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

	fmt.Println("finish")

	ticker := time.NewTicker(Interval)

	for {
		channel0, channel1, err := tsl.GetFullLuminosity()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Total Light: %f lux\n", tsl.CalculateLux(channel0, channel1))
		fmt.Printf("Visible: %d\n", channel0)
		fmt.Printf("Infrared: %d\n", channel1)
		<-ticker.C
	}

}
