package main
/* producer-consumer problem in Go
*  from http://www.golangpatterns.info/concurrency/producer-consumer
*/

import ("fmt")
import ("time")
import ("math/rand")
import ("github.com/shirou/gopsutil/cpu")
import ("github.com/shirou/gopsutil/mem")
import ("os")

const MAX_NUM = 10

var done = make(chan bool)
var msgs = make(chan int, MAX_NUM)
var nums = 5

func produce () {
    for {
      num := rand.Intn(nums)
      msgs <- num
      //fmt.Printf("Produced token %v : it's a living!\n", num)
      time.Sleep(time.Duration(rand.Float64()) * time.Second)
    }
}

func consume () {
    for {
      <-msgs
      //msg := <-msgs
      //fmt.Printf("Consumed token %v : yum yum yum!\n", msg)
      time.Sleep(time.Duration(rand.Float64()) * time.Second)
   }
}

func monitor () {
    for i := 0; i < 10; i++ {
      c, err := cpu.Percent(time.Duration(200) * time.Millisecond, false)
      m, err := mem.VirtualMemory()
      if err == nil {
        fmt.Printf("CPU usage: %f%%\n", c)
        fmt.Printf("%v\n", m)
      }
      time.Sleep(5 * time.Second)
    }
    os.Exit(0)
}

func main () {

    go monitor()
    time.Sleep(time.Second)

    go produce()
    go consume()
    <- done
}
