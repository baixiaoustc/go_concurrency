package go_concurrency

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"
)

var (
	MAX        int = 100
	BUFFERSIZE int = 1000
)

func testBoring(msg string) {
	for i := 0; ; i++ {
		log.Printf("%s %d", msg, i)
		time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
	}
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

func testBoringWithChannelQuit(msg string, quit chan bool) <-chan string {
	c := make(chan string)
	go func() {
		for i := 0; ; i++ {
			select {
			case c <- fmt.Sprintf("%s %d", msg, i):
				time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
			case <-quit:
				return
			}

		}
	}()
	return c
}

func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func benchmarkList() []float64 {
	list := make([]float64, MAX)
	for n := 0; n < MAX; n++ {
		list[n] = float64(n)
	}
	return list
}

func cpubound(n float64) float64 {
	return math.Pow(3.1415926, n)
}

func genChan(nums []float64) <-chan float64 {
	out := make(chan float64)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func genChanBuffer(nums []float64) <-chan float64 {
	out := make(chan float64, BUFFERSIZE)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func cpuChan(in <-chan float64) <-chan float64 {
	out := make(chan float64)
	go func() {
		for n := range in {
			out <- cpubound(n)
		}
		close(out)
	}()
	return out
}

func cpuChanBuffer(in <-chan float64) <-chan float64 {
	out := make(chan float64, BUFFERSIZE)
	go func() {
		for n := range in {
			out <- cpubound(n)
		}
		close(out)
	}()
	return out
}

func mergeChan(cs ...<-chan float64) <-chan float64 {
	var wg sync.WaitGroup
	out := make(chan float64)

	output := func(c <-chan float64) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func mergeChanBuffer(cs ...<-chan float64) <-chan float64 {
	var wg sync.WaitGroup
	out := make(chan float64, BUFFERSIZE)

	output := func(c <-chan float64) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func iobound(n float64) float64 {
	time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
	return n
}

func ioChan(in <-chan float64) <-chan float64 {
	out := make(chan float64)
	go func() {
		for n := range in {
			out <- iobound(n)
		}
		close(out)
	}()
	return out
}

func ioChanBuffer(in <-chan float64) <-chan float64 {
	out := make(chan float64, BUFFERSIZE)
	go func() {
		for n := range in {
			out <- iobound(n)
		}
		close(out)
	}()
	return out
}
