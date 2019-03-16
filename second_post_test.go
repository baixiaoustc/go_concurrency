package go_concurrency

import (
	"log"
	"runtime"
	"sync"
	"testing"
)

func TestPipeline(t *testing.T) {
	c := gen(2, 3, 4)
	out := sq(c)

	for n := range out {
		log.Print(n)
	}
}

func TestPipelines(t *testing.T) {
	for n := range sq(sq(sq(gen(1, 2, 3, 4)))) {
		log.Print(n)
	}
}

func TestFan(t *testing.T) {
	c := gen(2, 3, 4)
	out1 := sq(c)
	out2 := sq(c)

	for n := range merge(out1, out2) {
		log.Print(n)
	}
}

func BenchmarkCPUSequential(b *testing.B) {
	for i := 0; i < b.N; i++ {
		list := benchmarkList()

		var sum float64
		for _, n := range list {
			sum += cpubound(n)
		}
	}
}

func BenchmarkCPUPipeline(b *testing.B) {
	for i := 0; i < b.N; i++ {
		list := benchmarkList()

		c := genChan(list)
		out := cpuChan(c)

		var sum float64
		for n := range out {
			sum += n
		}
	}
}

func BenchmarkCPUPipelineBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		list := benchmarkList()

		c := genChanBuffer(list)
		out := cpuChanBuffer(c)

		var sum float64
		for n := range out {
			sum += n
		}
	}
}

func BenchmarkCPUFan(b *testing.B) {
	for i := 0; i < b.N; i++ {
		list := benchmarkList()

		c := genChan(list)
		gonum := runtime.NumCPU() / 2
		outs := make([]<-chan float64, gonum)
		for i := 0; i < gonum; i++ {
			outs[i] = cpuChan(c)
		}

		var sum float64
		for n := range mergeChan(outs...) {
			sum += n
		}
	}
}

func BenchmarkCPUFanBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		list := benchmarkList()

		c := genChanBuffer(list)
		gonum := runtime.NumCPU() / 2
		outs := make([]<-chan float64, gonum)
		for i := 0; i < gonum; i++ {
			outs[i] = cpuChanBuffer(c)
		}

		var sum float64
		for n := range mergeChanBuffer(outs...) {
			sum += n
		}
	}
}

func BenchmarkCPUParallelize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		list := benchmarkList()
		gonum := runtime.NumCPU()

		var sum float64
		num := len(list)
		stride := num / gonum

		var wg sync.WaitGroup
		wg.Add(gonum)
		var mux sync.Mutex

		for g := 0; g < gonum; g++ {
			go func(g int) {
				start := g * stride
				end := start + stride
				if g == gonum-1 {
					end = num
				}

				var sumin float64
				for _, n := range list[start:end] {
					sumin += cpubound(n)
				}

				mux.Lock()
				sum += sumin
				mux.Unlock()

				wg.Done()
			}(g)
		}

		wg.Wait()
	}
}

func BenchmarkIOSequential(b *testing.B) {
	for i := 0; i < b.N; i++ {
		list := benchmarkList()

		var sum float64
		for _, n := range list {
			sum += iobound(n)
		}
	}
}

func BenchmarkIOPipeline(b *testing.B) {
	for i := 0; i < b.N; i++ {
		list := benchmarkList()

		c := genChan(list)
		out := ioChan(c)

		var sum float64
		for n := range out {
			sum += n
		}
	}
}
func BenchmarkIOPipelineBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		list := benchmarkList()

		c := genChanBuffer(list)
		out := ioChanBuffer(c)

		var sum float64
		for n := range out {
			sum += n
		}
	}
}

func BenchmarkIOFan(b *testing.B) {
	for i := 0; i < b.N; i++ {
		list := benchmarkList()

		c := genChan(list)
		gonum := runtime.NumCPU() / 2
		outs := make([]<-chan float64, gonum)
		for i := 0; i < gonum; i++ {
			outs[i] = ioChan(c)
		}

		var sum float64
		for n := range mergeChan(outs...) {
			sum += n
		}
	}
}

func BenchmarkIOFanBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		list := benchmarkList()

		c := genChanBuffer(list)
		gonum := runtime.NumCPU() / 2
		outs := make([]<-chan float64, gonum)
		for i := 0; i < gonum; i++ {
			outs[i] = ioChanBuffer(c)
		}

		var sum float64
		for n := range mergeChanBuffer(outs...) {
			sum += n
		}
	}
}

func BenchmarkIOParallelize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		list := benchmarkList()
		gonum := runtime.NumCPU()

		var sum float64
		num := len(list)
		stride := num / gonum

		var wg sync.WaitGroup
		wg.Add(gonum)
		var mux sync.Mutex

		for g := 0; g < gonum; g++ {
			go func(g int) {
				start := g * stride
				end := start + stride
				if g == gonum-1 {
					end = num
				}

				var sumin float64
				for _, n := range list[start:end] {
					sumin += iobound(n)
				}

				mux.Lock()
				sum += sumin
				mux.Unlock()

				wg.Done()
			}(g)
		}

		wg.Wait()
	}
}

func TestList(t *testing.T) {
	gonum := runtime.NumCPU()
	log.Println(gonum)
	log.Println(runtime.GOMAXPROCS(0))
	log.Println("end")
}

func TestLoop(t *testing.T) {
	go testBoring("here")
	i := 0
	list := make([]int, 0)
	for {
		list = append(list, 111)
		i++
	}
}

func TestCPUPool(t *testing.T) {
	channum := 100
	gonum := runtime.NumCPU()

	f := func(w *Work) {
		if v, ok := w.input.(float64); ok {
			log.Printf("%f, %f", v, cpubound(v))
		}
	}
	list := benchmarkList()
	c := genPoolChanBuffer(list)

	p := InitPool(channum, gonum, f)
	p.RunWorker()
	p.FeedWorker(c)
	p.Wait()
}

func BenchmarkCPUPool(b *testing.B) {
	channum := 100
	gonum := runtime.NumCPU()
	for i := 0; i < b.N; i++ {
		f := func(w *Work) {
			if v, ok := w.input.(float64); ok {
				cpubound(v)
			}
		}
		list := benchmarkList()
		c := genPoolChanBuffer(list)

		p := InitPool(channum, gonum, f)
		p.RunWorker()
		p.FeedWorker(c)
		p.Wait()
	}
}

func BenchmarkCPUPoolMin(b *testing.B) {
	channum := 100
	gonum := runtime.NumCPU() / 2
	for i := 0; i < b.N; i++ {
		f := func(w *Work) {
			if v, ok := w.input.(float64); ok {
				cpubound(v)
			}
		}
		list := benchmarkList()
		c := genPoolChanBuffer(list)

		p := InitPool(channum, gonum, f)
		p.RunWorker()
		p.FeedWorker(c)
		p.Wait()
	}
}

func BenchmarkCPUPoolMax(b *testing.B) {
	channum := 100
	gonum := 10 * runtime.NumCPU()
	for i := 0; i < b.N; i++ {
		f := func(w *Work) {
			if v, ok := w.input.(float64); ok {
				cpubound(v)
			}
		}
		list := benchmarkList()
		c := genPoolChanBuffer(list)

		p := InitPool(channum, gonum, f)
		p.RunWorker()
		p.FeedWorker(c)
		p.Wait()
	}
}

func BenchmarkIOPool(b *testing.B) {
	channum := 100
	gonum := runtime.NumCPU()
	for i := 0; i < b.N; i++ {
		f := func(w *Work) {
			if v, ok := w.input.(float64); ok {
				iobound(v)
			}
		}
		list := benchmarkList()
		c := genPoolChanBuffer(list)

		p := InitPool(channum, gonum, f)
		p.RunWorker()
		p.FeedWorker(c)
		p.Wait()
	}
}

func BenchmarkIOPoolMin(b *testing.B) {
	channum := 100
	gonum := runtime.NumCPU() / 2
	for i := 0; i < b.N; i++ {
		f := func(w *Work) {
			if v, ok := w.input.(float64); ok {
				iobound(v)
			}
		}
		list := benchmarkList()
		c := genPoolChanBuffer(list)

		p := InitPool(channum, gonum, f)
		p.RunWorker()
		p.FeedWorker(c)
		p.Wait()
	}
}

func BenchmarkIOPoolMax(b *testing.B) {
	channum := 100
	gonum := 10 * runtime.NumCPU()
	for i := 0; i < b.N; i++ {
		f := func(w *Work) {
			if v, ok := w.input.(float64); ok {
				iobound(v)
			}
		}
		list := benchmarkList()
		c := genPoolChanBuffer(list)

		p := InitPool(channum, gonum, f)
		p.RunWorker()
		p.FeedWorker(c)
		p.Wait()
	}
}
