package main
/*  sushi bar problem in go
*   from https://blog.ksub.org/bytes/post/sushi-bar/sushi-bar.go
*/

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

type Customer struct {
	Name string
}

func (c *Customer) Sit(r *Restaurant) {
	log.Printf("Customer %s sat down\n", c.Name)
	go func() {
		time.Sleep(randMillisecond(500, 600))
		r.Leave(c)
		log.Printf("Customer %s left the restaurant\n", c.Name)
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
	timeout := time.After(5 * time.Second)
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
		case <-timeout:
			return
		}
	}
}

func (r *Restaurant) Enter(c *Customer) {
	r.in <- c
}

func (r *Restaurant) Leave(c *Customer) {
	r.out <- c
}

func main() {
	r := NewRestaurant()

	go func(r *Restaurant) {
		for i := 0; true; i++ {
			time.Sleep(randMillisecond(0, 250))
			c := &Customer{Name: fmt.Sprintf("%d", i)}
			log.Printf("Customer %s waiting\n", c.Name)
			r.Enter(c)
		}
	}(r)

	r.Run()
}

func randMillisecond(low, high int) time.Duration {
	return time.Duration(low+rand.Intn(high-low)) * time.Millisecond
}
