package main

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
	"math/rand"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"os"
)

func producer(np int) {
	event := WaitForEvent()
	log.Println("producer", np, "produces event", event)
	queue <- event
	time.Sleep(time.Duration(rand.Float64()) * time.Second)
}

func consumer(nc int) {
	for {
		event := <-queue
		log.Println("  consumer", nc, "gets event", event)
		event.process()
	}
}

var (
	last  int64                   // last event number
	queue = make(chan event64, 10) // 3 is arbitrary finite buffer size
	wg    sync.WaitGroup
)

type event64 int64

func (e event64) process() {
	log.Println("    processed: event", e)
	time.Sleep(time.Duration(rand.Float64()) * time.Second)
	wg.Done()
}

func WaitForEvent() event64 {
	return event64(atomic.AddInt64(&last, 1))
}

func monitor () {
    for i := 0; i < 10; i++ {
      c, err := cpu.Percent(time.Duration(200) * time.Millisecond, false)
      m, err := mem.VirtualMemory()
      if err == nil {
        log.Printf("CPU usage: %f%%\n", c)
        log.Printf("%v\n", m)
      }
      time.Sleep(5 * time.Second)
    }
    os.Exit(0)
}

const nEvents = 6

func main() {
	go monitor()
	time.Sleep(time.Second)

	go consumer(1)
	go consumer(2)
	wg.Add(nEvents)
	for i := 1; i <= nEvents; i++ {
		go producer(i)
	}
	wg.Wait()
}
