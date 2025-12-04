package main

import (
	"fmt"
	"time"
)


func doWork(done <- chan bool) {  //implementing done as read only channel
		for {
			select {
			case <-done:
				return
			
			// this means by default this goRoutine will continue to do whatever work that it needs to do and the parent go Routine will have the power to stop this goroutine from doing work when it deems necessary
			
			default:
				fmt.Println("Doing Work")
			}
		}
	/* here the reasoning of for select loop comes to play
	
	We want our Go Routine to do work by default if the parent cancels this goRoutine we want that to be a case as well
	*/
}

func main() {

	done:=make(chan bool)
	go doWork(done)
	time.Sleep(time.Second*3)
	close(done) //parent goRoutine closes the done channel which triggers the "case<-done" in child done channel and the function returns and stops.
}

/* Suppose we have a go routine running that we don't want to run for the lifetime of application and we don't actually realize that the goRoutine is running in the background, consuming memory and processing power and resources - Go Routine Leak.

To avoid it we should implement a mechanism that will allow the parent goRoutine to cancel its children goRoutine, here done channel comes to play.


*/