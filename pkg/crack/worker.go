package crack

import "github.com/lightningnetwork/lnd/aezeed"

type worker struct {
	id     int
	source chan *aezeed.Mnemonic
	done   chan bool
}

func (w *worker) Close() {
	w.done <- true // FIXME - may block for some reason
	close(w.source)
	close(w.done)
}

func newWorker(id int, password string, sink chan *aezeed.CipherSeed, workerDone chan int) *worker {
	worker := &worker{
		id:     id,
		source: make(chan *aezeed.Mnemonic),
		done:   make(chan bool),
	}
	go func() {
		for {
			select {
			case mnemonic := <-worker.source:
				// this is the actual work to be distributed across go routines
				seed, err := mnemonic.ToCipherSeed([]byte(password))
				if err != nil {
					sink <- nil
				} else {
					sink <- seed
				}
			case <-worker.done:
				// FIXME - if sink is blocking, this blocks as well. drain!
				workerDone <- worker.id
				return
			}
		}
	}()
	return worker
}

type WorkerPool struct {
	Workers []*worker
	Out     chan *aezeed.CipherSeed
	done    chan bool
}

// Consume reads mnemonics from channel and distributes them round-robin to its workers.
func (pool *WorkerPool) Consume(in <-chan *aezeed.Mnemonic) {
	w := 0
	for {
		select {
		case <-pool.done:
			return
		case mnemonic, ok := <-in:
			if !ok {
				return
			}
			pool.Workers[w].source <- mnemonic
			w = (w + 1) % len(pool.Workers)
		}
	}
}

// Close shuts down pool and all of its workers.
func (pool *WorkerPool) Close() {
	pool.done <- true // stops consumption
	// TODO - drain pool.Out and close pool.Out
	for _, worker := range pool.Workers {
		go worker.Close() // stops processing
	}
}

// NewWorkerPool creates a pool of size n to distribute mnemonic seed validation.
func NewWorkerPool(n int, password string) *WorkerPool {
	poolSink := make(chan *aezeed.CipherSeed)
	workerDone := make(chan int)
	workers := make([]*worker, n)
	for i := 0; i < n; i++ {
		workers[i] = newWorker(i, password, poolSink, workerDone)
	}
	pool := WorkerPool{
		Workers: workers,
		Out:     poolSink,
		done:    make(chan bool),
	}
	go func() {
		active := map[int]bool{}
		for i := 0; i < n; i++ {
			active[i] = true
		}
		for {
			select {
			case workerId := <-workerDone:
				delete(active, workerId)
				if len(active) == 0 {
					break
				}
			}
		}
	}()
	return &pool
}
