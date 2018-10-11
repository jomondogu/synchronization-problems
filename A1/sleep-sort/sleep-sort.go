/*  sleep sort in Go
*   from https://github.com/arpitbbhayani/go-sleep-sort/blob/master/sleepsort.go
*/

package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
	//"github.com/shirou/gopsutil/cpu"
	//"github.com/shirou/gopsutil/mem"
)

var max = 10.0
var step = 1.0

// prints a number of sleeping for n seconds
func sleepAndPrint(x float64, wg *sync.WaitGroup) {
	defer wg.Done()

	// Sleeping for time proportional to value
	time.Sleep(time.Duration(x) * 1000.0 * time.Millisecond)

	// Printing the value
	fmt.Printf("%f\n", x)
}

// Sorts given integer slice using sleep sort
func Sort(numbers []float64) {
	var wg sync.WaitGroup

	// Creating wait group that waits of len(numbers) of go routines to finish
	wg.Add(len(numbers))

	for _, x := range numbers {
		// Spinning a Go routine
		go sleepAndPrint(x, &wg)
	}

	// Waiting for all go routines to finish
	wg.Wait()
}

/*
func monitor () {
    iterations := 10
    CPUusage := 0.0
    VMusage := 0.0
    for i := 0; i < iterations; i++ {
      c, err := cpu.Percent(time.Duration(200) * time.Millisecond, false)
      m, err := mem.VirtualMemory()
      if err == nil {
        cpu := c[0]
        vm := m.UsedPercent
        CPUusage += cpu
        VMusage += vm
        fmt.Printf("CPU usage: %f\n", cpu)
        fmt.Printf("%v\n", m)
      }
      time.Sleep(time.Second)
    }
    CPUavg := CPUusage / float64(iterations)
    fmt.Printf("Average CPU usage: %f \n", CPUavg)
    VMavg := VMusage / float64(iterations)
    fmt.Printf("Average VM usage: %f \n", VMavg)
		os.Exit(0);
}*/

func main() {
	//go monitor()
	//time.Sleep(time.Second)

	args := os.Args[1:]

	max, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		panic(err)
	}

	step, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		panic(err)
	}

	var numbers = []float64{}
	for xi := 0.0; xi < max; xi += step {
		numbers = append(numbers, xi)
	}

	Sort(numbers)
}
