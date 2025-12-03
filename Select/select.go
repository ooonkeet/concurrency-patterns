// what a select statement does is it lets a goroutine wait on multiple communication operations


package main

import "fmt"

func main(){
	myChannel:=make(chan string)
	anotherChannel:=make(chan string)
	go func(){
		myChannel<-"data"
	}()
	go func(){
		anotherChannel<-"cow"
	}()
	select{
	case msgFromMyChannel := <-myChannel:
		fmt.Println(msgFromMyChannel)
	case msgFromAnotherChannel:= <-anotherChannel:
		fmt.Println(msgFromAnotherChannel)
		// select stmt is going to block until one of its cases can run
		// if select is able to receive messages from multiple channels at the same time it will choose one at random
	}
	
}