/* dining savages problem in Go
*  from https://github.com/fsouza/lbos/blob/master/009-dining-savages.go
*/

package main

import (
	"fmt"
	"time"
	//"math/rand"
	"os"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

const M = 50

var cooksum = 0
var eatsum = 0

type Serving int

func SavagesThread(savages int, servings chan Serving) {
	for i := 0; i < savages; i++ {
		go func(){
			for {
				<-servings
				//fmt.Println("Eating : yum yum yum!")
				//time.Sleep(time.Duration(rand.Float64()) * time.Second)
				eatsum += 1
			}
		}()
	}
}

func CookThread(cooks int, servings chan Serving) {
	for i := 0; i < cooks; i++ {
		go func(){
			for {
				//fmt.Println("Cooking : the OTHER other white meat!")
				//time.Sleep(time.Duration(rand.Float64()) * time.Second)
				servings <- M
				cooksum += 1
			}
		}()
	}
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
    fmt.Printf("Total cooks: %d\tTotal eats: %d\n", cooksum, eatsum)
    os.Exit(0)
}

func main() {
	go monitor()
	time.Sleep(time.Second)

	finish := make(chan int)
	servings := make(chan Serving, M)
	CookThread(1, servings)
	SavagesThread(30, servings)

	<-finish
}
