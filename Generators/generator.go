// generators generate a stream of data on a channel
// CAUTION -> THIS PROGRAM CONTAINS AN INFINITE LOOP.
package main

import (
	"fmt"
	"math/rand"
)

/* 
	repeatFunc -> we want function to be generic (we dont want to limit the function that gets passed into this function, we ant this function to be able to repeat a call to any function that gets passed into it that has a return value - any non-void function can be passed irrespective of datatype)
*/
func repeatfunc[T any, K any](done <- chan K, fn func() T)<-chan T{
	//used generics and done channel
	stream := make(chan T)
	go func ()  {
		defer close(stream) //defer the closure of channel that we created above the stream
		for{
			select{
			case <- done:
				return //done channel ensures that go routines are shut down in the event that we call done from the main function
			case stream <- fn(): //go has first class functions so we can actually pass functions as parameters to other functions, fn is going to return type T which is the value put on the channel, it would be written onto the channel
			}
		}
	}()
	return stream //as long as the go routine is running this return stream is going to continue to get values put onto it even after this function completes
}
func main(){
	done:=make(chan int)
	defer close(done)
	// generator taked in a function that returns a value of type T
	randomNumFetcher := func() int {return rand.Intn(500000000)} //returns a random integer value till the limit we pass, here it is 500M, repeatfunc() is the generator we created.
	for rando:= range repeatfunc(done, randomNumFetcher){
		fmt.Println(rando)
	}//creates an infinte stream of numbers, we will funnel this infinite stream into our pipeline.
}