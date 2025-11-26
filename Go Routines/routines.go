package main

import (
	"fmt"
	"time"
)

func someFunc(num string){
	fmt.Println(num)
}

func main(){
	someFunc("1")
	// this is an example of synchronous coding where sequence is followed.
	fmt.Println("hi")
	// next is an example of async coding
	go someFunc("3")
	go someFunc("7")
	go someFunc("9")
	time.Sleep(time.Second*2)
	// as it is async code with main, time.Sleep freezes the main func till then all the async tasks are executed.
	fmt.Println("hello")
}