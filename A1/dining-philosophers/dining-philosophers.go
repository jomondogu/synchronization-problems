package main
/* dining philosophers problem in Go
* from https://rosettacode.org/wiki/Dining_philosophers#Go
*/

import (
    "hash/fnv"
    "log"
    "os"
    "time"
    "github.com/shirou/gopsutil/cpu"
    "github.com/shirou/gopsutil/mem"
)

// Number of philosophers is simply the length of this list.
// It is not otherwise fixed in the program.
var ph = []string{"Aristotle", "Kant", "Spinoza", "Marx", "Russell"}

const hunger = 3                // number of times each philosopher eats

var fmt = log.New(os.Stdout, "", 0) // for thread-safe output

var done = make(chan bool)

var thinksum = 0
var dinesum = 0

// This solution uses channels to implement synchronization.
// Sent over channels are "forks."
type fork byte

// A fork object in the program models a physical fork in the simulation.
// A separate channel represents each fork place.  Two philosophers
// have access to each fork.  The channels are buffered with capacity = 1,
// representing a place for a single fork.

// Goroutine for philosopher actions.  An instance is run for each
// philosopher.  Instances run concurrently.
func philosopher(phName string,
    dominantHand, otherHand chan fork, done chan bool) {

    for {
        <-dominantHand // pick up forks
        <-otherHand

        dinesum += 1

        dominantHand <- 'f' // put down forks
        otherHand <- 'f'

        thinksum += 1
    }
    done <- true
}

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
    fmt.Printf("Total dines: %d\tTotal thinks: %d\n", dinesum, thinksum)
    os.Exit(0)
}

func main() {
    go monitor()
    time.Sleep(time.Second)

    // Create fork channels and start philosopher goroutines,
    // supplying each goroutine with the appropriate channels
    place0 := make(chan fork, 1)
    place0 <- 'f' // byte in channel represents a fork on the table.
    placeLeft := place0
    for i := 1; i < len(ph); i++ {
        placeRight := make(chan fork, 1)
        placeRight <- 'f'
        go philosopher(ph[i], placeLeft, placeRight, done)
        placeLeft = placeRight
    }
    // Make one philosopher left handed by reversing fork place
    // supplied to philosopher's dominant hand.
    // This makes precedence acyclic, preventing deadlock.
    go philosopher(ph[0], place0, placeLeft, done)
    // they are all now busy eating
    for range ph {
        <-done // wait for philosphers to finish
    }

}
