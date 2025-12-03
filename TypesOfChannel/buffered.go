package main

import "fmt"

func main() {
	charChannel := make(chan string, 3) //buffered channel with limited capacity of 3. We don't need any corresponsing receiver for these 3 values.
	// default channel is unbuffered having no capacity - perform synchronous communication b/w goRoutines

	// unbuff channel provides a guarentee that an exchange between 2 goRoutines is performed at the instant the send and receive takes place
	// buffered channel we're actually using the cue-like functionality where we can basically send data onto the channel and just forget about it upto the alloted capacity.

	/*
			                           -----------------------------
								data   |   Sending                 |
							   --------|          Go               |
							   |	   |            Routine        |
							   |	   -----------------------------
		                       |
							  \ /
						-----------------------
						| Unbuffered Channel  |
						-----------------------
				                        |
				---------------------   |
				| Receiving         |/  |
				|          Go       | ---
				|            Routine|\
				---------------------

		------------------------------------------------------------------------------------------------------------

			                           -----------------------------
								data   |   Sending                 |
							   --------|          Go               |
							   |	   |            Routine        |
							   |	   -----------------------------
		                       |
							  \ /                      Limit stored inform of queue
						-----------------------  data  ----------------------------
						| buffered Channel    | -----> | data -> | data -> | data | ------
						-----------------------        ----------------------------      |
		                                                                                 |
				                                                                         |
				---------------------                                                    |
				| Receiving         |/                                                   |
				|          Go       | ----------------------------------------------------
				|            Routine|\
				---------------------

				1. when sending go routine sends it data to unbuffered channel, it would go to a waiting state. That means sending go routine would be blocked and it would stay blocked until there's a receiver receiving the data from the unbuffered channel. - Synchronous communication

				2. In case of unbuffered channel the sending go routine would not be blocked until the queue is filled with data, the buffered channel sends its data to queue whose capacity is set in the declaration of channel. The receiving go routine accepts the data from the queue. - Asynchronous communication
	*/

	chars := []string{"a", "b", "c"}

	for _, s := range chars {
		select {
		case charChannel <- s:
		}
	}
	close(charChannel)
	for result := range charChannel {
		fmt.Println(result)
	}
	// this shows we are actually able to loop over a closed channel ans still receive a residul data put onto a channel. This loop will end because internally it knows that the channel was closed at a certain point.
}