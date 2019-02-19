package go_concurrency

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

var quit chan int

func testBoring(msg string) {
	for i := 0; ; i++ {
		log.Printf("%s %d", msg, i)
		time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
	}
}

func loop() {
	for i := 0; i < 10; i++ {
		log.Printf("%d ", i)
	}
	quit <- 0
}

func testBoringWithChannel(msg string, c chan string) {
	for i := 0; ; i++ {
		c <- fmt.Sprintf("%s %d", msg, i)
		time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
	}
}

func testBoringWithChannelClose(msg string, c chan string) {
	for i := 0; i < 5; i++ {
		c <- fmt.Sprintf("%s %d", msg, i)
		time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
	}
	close(c)
}

func testBoringWithChannelGenerate(msg string) <-chan string {
	c := make(chan string)
	go func() {
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}()
	return c
}
