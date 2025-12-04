/*
		start
		  |
		  | data
		  |
	---------------
	|   stage 1   |  operation 1
	---------------
	      |
		  | data
		  |
	---------------
	|   stage 2   |  operation 2
	---------------
	      |
		  | data
		  |
		 end

All of the stages working together form the pipeline.

Using pipeline we can seperate the concerns of each stage.
*/

package main

import "fmt"

func sliceToChannel(nums[] int) <-chan int{
	// return a read only channel
	out:=make(chan int) //unbuff channel
	go func(){
		for _,n:=range nums{
			out<-n //when it puts one number from slice to the out channel, it actually blocks until that number is read from the out channel
		}
		close(out)
	}()
	// eash item in that slice we'll add it to the out channel -> remember this is an asynchronous goRoutine.
		return out //we return this channel while the goRoutine is still running, this out is passed to square function
}

func sq(in<-chan int) <-chan int{
	// take and return a read only channel
	out:=make(chan int) //unbuffered channel, reading values from out channel passed from sliceToChannel function
	go func ()  {
		for n:=range in{
			out<-n*n
		}
		close(out)
	}()
	//this goRoutine range over in channel that is returned from the sliceToChannel function
	// this go func would not block rest of this function
	// this goRoutine could still be sending values to the out channel even at the time we return it
	return out //the function of SliceToChannel 's GoRoutine is blocked until its value is read by this function's go routine, and this function's go routine remains blocked until another value is passed from 'sliceToChannel'.
	// therefore closing a channel is important as it won't wait and form a deadlock situation.
}

func main(){
	// start and end of our pipeline is going to be orchestrated via our main function.

	nums:=[]int{2,3,4,7,1} //input
	// stage1 -> converts a slice to channel
	dataChannel:=sliceToChannel(nums)
	// stage 2 -> squares the data
	finalChannel:=sq(dataChannel)
	// stage 3 -> output the result of entire pipeline
	for n:=range finalChannel{
		// ranging over an unbuffered channel just what sq function did -> range over the unbuffered channel passed by sliceToChannel
		fmt.Println(n) //we can see the squares of the int slice is in order as the communication between two goRoutine is synchronous
	}
	// remember stage 1 and stage 2 are communicating synchronously, unbuffered channel default capacity -> 1
}