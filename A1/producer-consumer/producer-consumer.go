package main
/* producer-consumer problem in Go
*  from http://www.golangpatterns.info/concurrency/producer-consumer
*/

import (
    "fmt"
    "time"
    "math/rand"
    "github.com/shirou/gopsutil/cpu"
    "github.com/shirou/gopsutil/mem"
    "os"
)

const MAX_NUM = 10

var done = make(chan bool)
var msgs = make(chan int, MAX_NUM)
var nums = 5
var psum = 0
var csum = 0

func produce () {
    for {
      num := rand.Intn(nums)
      msgs <- num
      psum += 1
      //fmt.Printf("Produced token %v : it's a living!\n", num)
    }
}

func consume () {
    for {
      <-msgs
      csum += 1
      //msg := <-msgs
      //fmt.Printf("Consumed token %v : yum yum yum!\n", msg)
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
    fmt.Printf("Total productions: %d\tTotal consumptions: %d\n", psum, csum)
    os.Exit(0)
}

func main () {

    go monitor()
    time.Sleep(time.Second)

    go produce()
    go consume()
    <- done
}
