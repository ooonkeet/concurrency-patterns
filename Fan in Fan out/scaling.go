package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func fanIn[T any](done <-chan int, channels ...<-chan T) <-chan T{
	var wg sync.WaitGroup
	// as we are spinning up a go routine for each transfer we need to use wait group - we want to wait for each transfer to finish
	fanndedInStream := make(chan T) //single channel where all of the values will be fanned into [all results compiled in this single channel]
	transfer := func(c <-chan T){ //transfer data from one channel to another (read values from past in channel and put them on the fan in stream)
	// transfer all the channels fanned out data to one stream fannedInStream
		defer wg.Done()
		for i:=range c{
			select{
			case <-done:
				return 
			case fanndedInStream <-i:
			}
		}
	}
	for _,c := range channels{ //range through the channels in function parameter (we are getting all of the resulting channels from the prime find functions that are running concurrently and all of those channels is going to contain the primes)
		wg.Add(1)
		go transfer(c) // spin up go routine to undergo all the transfers, transfer is going to take in c
	}

	go func(){
		wg.Wait()
		// we want to wait for each transfer to finish
		close(fanndedInStream)
		// as all go routine is finish the fan in stream needs to be closed
	}()
	return fanndedInStream
}

/*
	repeatFunc -> we want function to be generic (we dont want to limit the function that gets passed into this function, we ant this function to be able to repeat a call to any function that gets passed into it that has a return value - any non-void function can be passed irrespective of datatype)
*/
func repeatfunc[T any, K any](done <- chan K, fn func() T)<-chan T{
	//used generics and done channel
	stream := make(chan T) //unbuffered channel so the next go routine is blocked until next go routine in take function accepts the value - unbuffered channel is also synchronized channel
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

// func take - use generics and take in done channel and take in the stream that we want to take values from and that's going to be stream produced by generator, and an int for storing number of items we want to take from the stream produced by the generator, and it would return a channel as well
func take[T any, K any](done<-chan K, stream <-chan T, n int) <-chan T{
	taken:=make(chan T)
	go func ()  {
		defer close(taken)
		for i:=0;i<n;i++{
			select{
			case <-done:
				return 
			case taken <- <-stream:
			// this syntax write to the taken stream, so this is the stream that contains value that we have taken from the stream produced by the generator
			// write a value from stream onto the taken stream - "<-stream" reads value out of stream and that value is written to taken "taken<-(<-stream)"
			}
		}
	}()
	return taken
	// this function thus allow us to control the amount of data we take from the stream produced by the generator
}

func primeFinder(done <-chan int,randIntStream <-chan int) <-chan int{ //taking done channel with a specific type to avoid generics, randomintstream channel of int (read only channel) return channel of int
	isPrime := func(randomInt int) bool{ //check every number in between random int and 1 and if that number divides evenly into the random int then it means the number isn't prime because a prime number is only divisible by 1 and itself
		for i:=randomInt-1;i>1;i--{
			if randomInt%i==0{
				return  false
			}
		}
		return true
		// this essentially created a slow pipeline stage artificially
	}
	primes:=make(chan int)
	go func(){
		defer close(primes)
		for{
			select{
			case <-done:
				return 
			case randomInt := <-randIntStream:
				if isPrime(randomInt){
					primes <- randomInt
				}

			}
		}
	}()
	return primes
}

func main(){
	start:=time.Now()
	done:=make(chan int)
	defer close(done)
	// generator taked in a function that returns a value of type T
	randomNumFetcher := func() int {return rand.Intn(500000000)} //returns a random integer value till the limit we pass, here it is 500M, repeatfunc() is the generator we created.
	randIntStream:=repeatfunc(done,randomNumFetcher)
	// fan out
	CPUcount := runtime.NumCPU() //cpu count of system
	primeFinderChannels := make([]<-chan int,CPUcount) //create a slice of resulting channels from all of instances of prime finder that we spin up - slice of read only channels and the length of the slice should be cpu count -> resulting channels proportionate to number of cpus
	// HERE WE WANT TO put the channels that come from the calls of prime finder
	for i:=0;i<CPUcount;i++{ //every cpu available we are going to spin up our prime finder function and result will be inserted into above slice
		primeFinderChannels[i]=primeFinder(done,randIntStream) //leave us with a slice of channels that comes from Prime Finder(each cpu spin up this prime finder)
		// prime finder channels will consist of results of prime finder calls
	}	
	// fan in
	fannedInStream := fanIn(done,primeFinderChannels...)
	for rando:=range take(done, fannedInStream, 10){
		fmt.Println(rando)
	}
	// substantial improve in performance
	fmt.Println(time.Since(start)) //tells us how much time it took
}