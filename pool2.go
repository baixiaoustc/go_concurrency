package go_concurrency

type Pool2 struct {
	channum   int
	workernum int
	Ch        chan bool
	Func      Handler
}

func (p Pool2) RunWorker(num int, in <-chan interface{}) {
	for i := 0; i < p.workernum; i++ {
		go p.FeedWorker(in)
	}
}

func (p Pool2) FeedWorker(in <-chan interface{}) {
	for n := range in {
		<-p.Ch
		go do(n, p.Ch)
	}
}

func do(n interface{}, ch chan bool) {
	defer func() {
		ch <- true
	}()

	if v, ok := n.(float64); ok {
		iobound(v)
	}
}
