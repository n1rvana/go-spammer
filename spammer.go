// Copyright 2016 Nirvana Project
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"time"
)

// The point of this is to discover that attempting to pull data off of a buffered
// channel blocks when there is no data.  So maybe a select is the way to go when you want
// to process multiple channels in a loop.  Otherwsie you block on the first, and never
// seem to hit the second.

var (
	littleSpam chan string
	rangeSpam  chan string
	forSpam    chan string
	selectA    chan string
	selectB    chan string
	collection []chan string
)

func main() {
	littleSpam = make(chan string, 100)
	rangeSpam = make(chan string, 100)
	forSpam = make(chan string, 100)
	selectA = make(chan string, 100)
	selectB = make(chan string, 100)
	for index := 0; index < 50; index++ {
		collection = append(collection, make(chan string, 100))
	}
	go spammer()
	go little()
	go ranger()
	go forloop()
	go selectLoop()
	go lenTest()
	time.Sleep(time.Second * 30)

}

func spammer() {
	count := 0
	for {
		message := fmt.Sprintf("Spam #%d", count)
		if 15 > count {
			littleSpam <- message
			for _, channel := range collection {
				channel <- message
			}
		}
		selectA <- message
		selectB <- message
		forSpam <- message
		rangeSpam <- message
		count++
		time.Sleep(time.Second * 1)

	}
}

func little() {
	for message := range littleSpam {
		fmt.Printf("littleSpam: %s\n", message)
	}
	fmt.Printf("littleSpam: %s\n", "Never gets here because when the channel is empty, range blocks!")

}

func ranger() {
	for message := range rangeSpam {
		fmt.Printf("ranger: %s\n", message)
	}
}

func forloop() {
	for {
		message := <-forSpam
		fmt.Printf("forSpam: %s\n", message)
	}
}

func selectLoop() {
	for {
		select {
		case messageA := <-selectA:
			fmt.Printf("selectLoop: messageA %s\n", messageA)
		case messageB := <-selectB:
			fmt.Printf("selectLoop: messageB %s\n", messageB)

		}

	}
}

// But maybe you have an unknown number of channels, say in an array, and you need to process them
// all regularly, so you can't use select.  In this case you use len(channel) to determine if there
// are any messages waiting

func lenTest() {
	time.Sleep(time.Second * 10) // let some messages accumulate
	for {                        // keep checking the channels
		for index, currChannel := range collection {
			for 0 < len(currChannel) { // effectively "While there are messages"
				message := <-currChannel
				fmt.Printf("lenTest: channel %d got message %s\n", index, message)
			}
		}
		fmt.Printf("lenTest: Outside the channel loop - should keep printing even after 15 messages (when we're getting no more messages)")
		time.Sleep(time.Second * 1) // slow the loop down
	}
}
