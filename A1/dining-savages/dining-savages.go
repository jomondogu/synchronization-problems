/* dining savages problem in Go
*  from https://github.com/fsouza/lbos/blob/master/009-dining-savages.go
*/

package main

import (
	"fmt"
	"time"
	"math/rand"
)

const M = 50

type Serving int

func SavagesThread(savages int, servings chan Serving) {
	for i := 0; i < savages; i++ {
		go func(){
			for {
				<-servings
				fmt.Println("Eating : yum yum yum!")
				time.Sleep(time.Duration(rand.Float64()) * time.Second)
			}
		}()
	}
}

func CookThread(cooks int, servings chan Serving) {
	for i := 0; i < cooks; i++ {
		go func(){
			for {
				fmt.Println("Cooking : the OTHER other white meat!")
				time.Sleep(time.Duration(rand.Float64()) * time.Second)
				servings <- M
			}
		}()
	}
}

func main() {
	finish := make(chan int)
	servings := make(chan Serving, M)
	CookThread(1, servings)
	SavagesThread(30, servings)

	<-finish
}
