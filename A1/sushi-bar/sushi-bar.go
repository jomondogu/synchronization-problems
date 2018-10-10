package main
/*  sushi bar problem in go
*   from https://blog.ksub.org/bytes/post/sushi-bar/sushi-bar.go
*/

import (
	"fmt"
	//"log"
	"math/rand"
	"time"
	"os"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

var waitsum = 0
var eatsum = 0

type Customer struct {
	Name string
}

func (c *Customer) Sit(r *Restaurant) {
	//log.Printf("Customer %s sat down\n", c.Name)
	eatsum++
	go func() {
		//time.Sleep(randMillisecond(500, 600))
		r.Leave(c)
		//log.Printf("Customer %s left the restaurant\n", c.Name)
	}()
}

type Restaurant struct {
	in  chan *Customer
	out chan *Customer
}

func NewRestaurant() *Restaurant {
	return &Restaurant{
		in:  make(chan *Customer),
		out: make(chan *Customer),
	}
}

func (r *Restaurant) Run() {
	n := 0
	in := r.in
	//timeout := time.After(5 * time.Second)
	for {
		select {
		case c := <-in:
			n++
			if n == 5 {
				in = nil
			}
			c.Sit(r)
		case <-r.out:
			n--
			if n == 0 {
				in = r.in
			}
		/*
		case <-timeout:
			return
			*/
		}
	}
}

func (r *Restaurant) Enter(c *Customer) {
	r.in <- c
}

func (r *Restaurant) Leave(c *Customer) {
	r.out <- c
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
    fmt.Printf("Total wait loops: %d\t Total eats: %d\n", waitsum, eatsum)
    os.Exit(0)
}

func main() {
	go monitor()
	time.Sleep(time.Second)

	r := NewRestaurant()

	go func(r *Restaurant) {
		for i := 0; true; i++ {
			//time.Sleep(randMillisecond(0, 250))
			c := &Customer{Name: fmt.Sprintf("%d", i)}
			//log.Printf("Customer %s waiting\n", c.Name)
			waitsum++
			r.Enter(c)
		}
	}(r)

	r.Run()
}

func randMillisecond(low, high int) time.Duration {
	return time.Duration(low+rand.Intn(high-low)) * time.Millisecond
}
