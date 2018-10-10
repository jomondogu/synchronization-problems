/* modus hall problem in Go
*  from https://blog.ksub.org/bytes/post/modus-hall/modus-hall.go
*/

package main

import (
	"fmt"
	"math/rand"
	"time"
	"os"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

var mods, reshalls int

type Mod struct{}

type ResHall struct{}

type Path struct {
	in       chan interface{}
	out      chan interface{}
	mods     []*Mod
	reshalls []*ResHall
}

func NewPath() *Path {
	return &Path{
		in:  make(chan interface{}),
		out: make(chan interface{}),
	}
}

func (p *Path) Run() {
	var (
		mods     chan interface{}
		reshalls chan interface{}
		mod      *Mod
		reshall  *ResHall
	)

	//timeout := time.NewTimer(2 * time.Second)

	for {
		select {
		case s := <-p.in:
			switch s.(type) {
			case *Mod:
				s := s.(*Mod)
				p.mods = append(p.mods, s)
				if len(p.mods) > len(p.reshalls) {
					mod = p.mods[0]
					mods, reshalls = p.out, nil
				}
			case *ResHall:
				s := s.(*ResHall)
				p.reshalls = append(p.reshalls, s)
				if len(p.reshalls) > len(p.mods) {
					reshall = p.reshalls[0]
					mods, reshalls = nil, p.out
				}
			}

		case mods <- mod:
			p.mods = p.mods[1:]
			if len(p.mods) <= len(p.reshalls) {
				mods = nil
			}

		case reshalls <- reshall:
			p.reshalls = p.reshalls[1:]
			if len(p.reshalls) <= len(p.mods) {
				reshalls = nil
			}
/*
		case <-timeout.C:
			return
			*/
		}
	}
}

func randMillisecond(low, high int) time.Duration {
	return time.Duration(low+rand.Intn(high-low)) * time.Millisecond

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
    fmt.Printf("Total heathens: %d\t Total prudes: %d\n", mods, reshalls)
    os.Exit(0)
}

func main() {
	go monitor()
	time.Sleep(time.Second)

	path := NewPath()

	for i := 5; i > 0; i-- {
		// Produce Mods
		go func(p *Path) {
			for {
				path.in <- new(Mod)
				//time.Sleep(randMillisecond(100, 500))
			}
		}(path)

		// Produce ResHalls
		go func(p *Path) {
			for {
				path.in <- new(ResHall)
				//time.Sleep(randMillisecond(100, 500))
			}
		}(path)
	}

	// Receive students on the other end
	go func(p *Path) {
		for {
			s := <-path.out
			switch s.(type) {
			case *Mod:
				mods++
			case *ResHall:
				reshalls++
			}
			//fmt.Println("Mods", mods, "ResHalls", reshalls)
		}
	}(path)

	path.Run()
}
