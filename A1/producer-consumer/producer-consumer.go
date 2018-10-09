package main
/* producer-consumer problem in Go
*  from http://www.golangpatterns.info/concurrency/producer-consumer
*/

import ("fmt")
import ("time")
import ("math/rand")

var done = make(chan bool)
var msgs = make(chan int)
var nums = 5

func produce () {

    for {
      num := rand.Intn(nums)
      msgs <- num
      fmt.Printf("Produced token %v : it's a living!\n", num)
      time.Sleep(time.Duration(rand.Float64()) * time.Second)
    }
    /*
    for i := 0; i < 10; i++ {
        msgs <- i
        fmt.Printf("Produced token %v, it's a living!\n", i)
    }
    done <- true
    */
}

func consume () {
    for {
      msg := <-msgs
      fmt.Printf("Consumed token %v : yum yum yum!\n", msg)
      time.Sleep(time.Duration(rand.Float64()) * time.Second)
   }
}

func main () {
   go produce()
   go consume()
   <- done
}
