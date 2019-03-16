package go_concurrency

import (
	"log"
	"testing"
	"time"
)

func TestBasicGo(t *testing.T) {
	go testBoring("boring!")
}

func TestGoAndWait(t *testing.T) {
	go testBoring("boring!")

	log.Printf("I'm listening.")
	time.Sleep(time.Second)
	log.Printf("I'm leaving.")
}

func TestGoFuncParam(t *testing.T) {
	for i := 0; i < 10; i++ {
		go func() {
			log.Println(i)
		}()
	}
	time.Sleep(100 * time.Millisecond)
	log.Print("xxxxxxxxxxxxxx")

	for i := 0; i < 10; i++ {
		go func(n *int) {
			log.Println(*n)
		}(&i)
	}
	time.Sleep(100 * time.Millisecond)
}

func TestGoWithChannel(t *testing.T) {
	ch := make(chan string)
	go testBoringWithChannel("boring!", ch)
	for c := range ch {
		log.Printf("You say: %s", c)
	}
}

func TestGoWithChannelClose(t *testing.T) {
	ch := make(chan string)
	go testBoringWithChannelClose("boring!", ch)
	for c := range ch {
		log.Printf("You say: %s", c)
	}
}

func TestGoWithChannelCloseThenWrite(t *testing.T) {
	ch := make(chan string)
	go testBoringWithChannelClose("boring!", ch)
	for c := range ch {
		log.Printf("You say: %s", c)
	}
	ch <- "after close"
}

//func TestReceiveChannelWrong(t *testing.T) {
//	ch := make(<-chan int)
//	go func() {
//		ch <- 1
//	}()
//	a := <-ch
//	log.Println(a)
//}

func TestReceiveChannelRight(t *testing.T) {
	ch := func() <-chan int {
		ch := make(chan int)
		go func() {
			ch <- 2
		}()
		return ch
	}()
	a := <-ch
	log.Println(a)
}

func TestMultiChannelsSerial(t *testing.T) {
	ch1 := make(chan string)
	go testBoringWithChannelClose("boring!", ch1)
	ch2 := make(chan string)
	go testBoringWithChannelClose("funning!", ch2)

	for b := range ch1 {
		log.Printf("You say: %s", b)
	}
	for b := range ch2 {
		log.Printf("You say: %s", b)
	}
}

func TestMultiChannelsConcurrently(t *testing.T) {
	ch1 := make(chan string)
	go testBoringWithChannelClose("boring!", ch1)
	ch2 := make(chan string)
	go testBoringWithChannelClose("funning!", ch2)

	for {
		select {
		case c1, ok := <-ch1:
			if !ok {
				ch1 = nil
				log.Print("close ch1")
				continue
			}
			log.Printf("You say: %s", c1)
		case c2, ok := <-ch2:
			if !ok {
				ch2 = nil
				log.Print("close ch2")
				continue
			}
			log.Printf("You say: %s", c2)
		default:
			break
		}

		if ch1 == nil && ch2 == nil {
			break
		}
	}

	log.Print("go on")
}

func TestTimeOutEach(t *testing.T) {
	ch := make(chan string)
	go testBoringWithChannel("boring!", ch)
	for i := 0; i < 5; i++ {
		select {
		case c := <-ch:
			log.Printf("You say: %s", c)
		case <-time.After(500 * time.Millisecond):
			log.Println("You talk too slow.")
		}
	}
}

func TestTimeOutWhole(t *testing.T) {
	ch := make(chan string)
	timeout := time.After(1 * time.Second)
	go testBoringWithChannel("boring!", ch)
	for i := 0; i < 5; i++ {
		select {
		case c := <-ch:
			log.Printf("You say: %s", c)
		case <-timeout:
			log.Println("You talk too much.")
			return
		}
	}
}

func TestChannelGenerator(t *testing.T) {
	ch := testBoringWithChannelGenerate("boring!")
	for i := 0; i < 5; i++ {
		log.Printf("You say: %q\n", <-ch)
	}
	log.Println("I'm leaving.")
}

func TestChannelQuit(t *testing.T) {
	quit := make(chan bool)
	ch := testBoringWithChannelQuit("boring!", quit)
	for i := 0; i < 5; i++ {
		log.Printf("You say: %q\n", <-ch)
	}
	log.Println("I'm leaving.")
	quit <- true
}
