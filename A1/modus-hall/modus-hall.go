/* modus hall problem in Go
*  from https://blog.ksub.org/bytes/post/modus-hall/modus-hall.go
*/

package main

import (
	"fmt"
	"math/rand"
	"time"
)

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

	timeout := time.NewTimer(2 * time.Second)

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

		case <-timeout.C:
			return
		}
	}
}

func randMillisecond(low, high int) time.Duration {
	return time.Duration(low+rand.Intn(high-low)) * time.Millisecond

}

func main() {
	path := NewPath()

	for i := 5; i > 0; i-- {
		// Produce Mods
		go func(p *Path) {
			for {
				path.in <- new(Mod)
				time.Sleep(randMillisecond(100, 500))
			}
		}(path)

		// Produce ResHalls
		go func(p *Path) {
			for {
				path.in <- new(ResHall)
				time.Sleep(randMillisecond(100, 500))
			}
		}(path)
	}

	// Receive students on the other end
	go func(p *Path) {
		var mods, reshalls int
		for {
			s := <-path.out
			switch s.(type) {
			case *Mod:
				mods++
			case *ResHall:
				reshalls++
			}
			fmt.Println("Mods", mods, "ResHalls", reshalls)
		}
	}(path)

	path.Run()
}
