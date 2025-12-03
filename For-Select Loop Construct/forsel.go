// Example of for-select pattern being used to have an infinitely running go routine

package main

import (
	"fmt"
	"time"
)

func main() {
	go func() {
		for {
			select {
			default:
				fmt.Println("Doing Work")
			}
		}
	}()
	time.Sleep(time.Second*10)
}

// infinite iteration of loop until 10 sec is up. because main func would cut off after 10 seconds.