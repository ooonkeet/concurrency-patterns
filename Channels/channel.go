package main

import "fmt"

func main() {
	// initialize a channel
	myChannel := make(chan string)

	// GoRoutine send data to channel
	go func() {
		myChannel <- "data"
	}()
	// read data from channel
	msg := <-myChannel //this line blocks the execution and main function will actually wait for either this chanel to be closed or for a message to be received from this channel
	// go routine is putting the data onto the channel via this arrow and main function is reading the data from the channel
	// so data is going from channel to msg variable
	// this is an example of join point implementation
	fmt.Println(msg)
}