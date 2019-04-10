package go_concurrency

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"testing"
	"time"
)

var WorkPool chan work

type work struct {
	sth string
}

func (w work) Working() {
	log.Printf("start %s", w.sth)

	time.Sleep(10 * time.Second)

	log.Printf("end %s", w.sth)
}

func (w work) Saving() {
	log.Printf("save %s", w.sth)
}

//fake, just a example
func (w work) SavingDB(ctx context.Context) {
	log.Printf("save %s", w.sth)

	stmt := "select name from db.table"
	db := sql.DB{}
	conn, _ := db.Conn(ctx)
	rows, err := conn.QueryContext(ctx, stmt)
	if err != nil {
		if err == context.DeadlineExceeded {
			// context canceled
		}
		return
	}

	var name string
	for rows.Next() {
		if err := rows.Scan(&name); err != nil {
			if err == context.DeadlineExceeded {
				log.Println("scan canceled")
			}
		}
	}
}

func RunWorkerSimple() {
	workers := 2
	for i := 0; i < workers; i++ {
		go HandleWorkerSimple()
	}
}

func HandleWorkerSimple() {
	for {
		select {
		case work := <-WorkPool:
			work.Working()
		}
	}

	return
}

func TestWorkerSimple(t *testing.T) {
	WorkPool = make(chan work)
	go RunWorkerSimple()

	go func() {
		list := benchmarkList()
		for _, l := range list {
			WorkPool <- work{fmt.Sprint(l)}
		}
	}()

	select {}
}

func RunWorkerHold(stop, stopped chan struct{}) {
	var wg sync.WaitGroup
	workers := 2
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go HandleWorkerHold(stop, &wg)
	}
	wg.Wait()
	stopped <- struct{}{}
}

func HandleWorkerHold(stop chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case work := <-WorkPool:
			work.Working()
		case <-stop:
			log.Println("worker: caller has told us to stop")
			return
		}
	}

	return
}

func TestWorkerHold(t *testing.T) {
	WorkPool = make(chan work)
	stopChan := make(chan struct{})
	stoppedChan := make(chan struct{})
	go RunWorkerHold(stopChan, stoppedChan)

	go func() {
		list := benchmarkList()
		for _, l := range list {
			WorkPool <- work{fmt.Sprint(l)}
		}
	}()

	// listen for C-c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println("main: received C-c - shutting down")

	// tell the goroutine to stop
	log.Println("main: telling workers to stop")
	close(stopChan)
	// and wait for them to reply back
	<-stoppedChan
	log.Println("main: workers has told us they've finished")

	for work := range WorkPool {
		work.Saving()
	}

	return
}

func RunWorkerContext(ctx context.Context, stopped chan struct{}) {
	WorkPool = make(chan work)
	var wg sync.WaitGroup
	workers := 2
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go HandleWorkerContext(ctx, &wg)
	}
	wg.Wait()
	stopped <- struct{}{}
}

func HandleWorkerContext(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case work := <-WorkPool:
			work.Working()
		case <-ctx.Done():
			log.Println("worker: caller has told us to stop")
			return
		}
	}

	return
}

func TestWorkerContext(t *testing.T) {
	WorkPool = make(chan work)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stoppedChan := make(chan struct{})
	go RunWorkerContext(ctx, stoppedChan)

	go func() {
		list := benchmarkList()
		for _, l := range list {
			WorkPool <- work{fmt.Sprint(l)}
		}
	}()

	// listen for C-c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println("main: received C-c - shutting down")

	// tell the goroutine to stop
	log.Println("main: telling workers to stop")
	cancel()
	// and wait for them to reply back
	<-stoppedChan
	log.Println("main: workers has told us they've finished")

	ctxTimeout, cancelTimeout := context.WithTimeout(ctx, 100*time.Microsecond)
	defer cancelTimeout()
	for {
		select {
		case work := <-WorkPool:
			work.SavingDB(ctx)
		case <-ctxTimeout.Done():
			log.Println("main: cann't wait any more")
			return
		}
	}

	return
}
