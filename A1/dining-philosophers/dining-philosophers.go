package main
/* dining philosophers problem in Go
* from https://rosettacode.org/wiki/Dining_philosophers#Go
*/

import (
    "hash/fnv"
    "log"
    //"math/rand"
    "os"
    "time"
    "github.com/shirou/gopsutil/cpu"
    "github.com/shirou/gopsutil/mem"
)

// Number of philosophers is simply the length of this list.
// It is not otherwise fixed in the program.
var ph = []string{"Aristotle", "Kant", "Spinoza", "Marx", "Russell"}

//const hunger = 3                // number of times each philosopher eats
//const think = time.Second // mean think time
//const eat = time.Second   // mean eat time

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
    //fmt.Println(phName, "seated")
    // each philosopher goroutine has a random number generator,
    // seeded with a hash of the philosopher's name.
    h := fnv.New64a()
    h.Write([]byte(phName))
    //rg := rand.New(rand.NewSource(int64(h.Sum64())))
    // utility function to sleep for a randomized nominal time
    for {
        //fmt.Println(phName, "hungry")
        <-dominantHand // pick up forks
        <-otherHand
        //fmt.Println(phName, "eating")
        dinesum += 1
        //time.Sleep(eat)
        dominantHand <- 'f' // put down forks
        otherHand <- 'f'
        //fmt.Println(phName, "thinking")
        thinksum += 1
        //time.Sleep(think)
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

    //fmt.Println("table empty")
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
    //fmt.Println("table empty")
}
