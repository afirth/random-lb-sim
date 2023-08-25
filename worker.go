package main

import (
	"fmt"
	"time"
)

// Worker is capable of performing parallel renders
// it keeps track of its own cumulative stats
type Worker struct {
	stats     *workerStats
	State     *workerState
	renders   []*Render
	startChan chan Render
	doneChan  chan int // index of the render which is done
}

type workerStats struct {
	cumulative WorkerCounters
	prev       WorkerCounters
}

type WorkerCounters struct {
	Failed, Success int
}

type workerState struct {
	Busy, Slots int
}

func NewWorker(slots int) *Worker {
	w := &Worker{
		renders:   make([]*Render, slots),
		startChan: make(chan Render),
		doneChan:  make(chan int),
		stats:     &workerStats{},
		State:     &workerState{Slots: slots},
	}
	go w.manage()
	return w
}

// substracts the previous stats from the current stats
// the workerPool keeps cumulative counts so they are not lost when a worker is removed
// the reads are not thread safe, but the writes are
func (w *Worker) LatestStats() WorkerCounters {
	if w.stats.prev == (WorkerCounters{}) {
		w.stats.prev = w.stats.cumulative
		return w.stats.cumulative
	} else {
		latest := WorkerCounters{
			Failed:  w.stats.cumulative.Failed - w.stats.prev.Failed,
			Success: w.stats.cumulative.Success - w.stats.prev.Success,
		}
		w.stats.prev = w.stats.cumulative
		return latest
	}
}

// the actual render is in manage() and the channel provides thread safety
func (w *Worker) Render(r Render) {
	w.startChan <- r
}

// listens to start and done channels
// fills or empties the slots
// manages the stats
func (w *Worker) manage() {
	for {
	Outer:
		select {
		// render is starting, use the first open slot
		case r := <-w.startChan:
			for i, render := range w.renders {
				if render == nil {
					w.renders[i] = &r
					w.State.Busy++
					go func() {
						// wait for the render to complete
						time.Sleep(r.duration)
						w.doneChan <- i
					}()
					break Outer // slot found
				}
			}
			// no slot found
			w.stats.cumulative.Failed++

			// render is complete, clear the slot
		case idx := <-w.doneChan:
			w.renders[idx] = nil
			w.State.Busy--
			w.stats.cumulative.Success++
		}
	}
}

func (w *Worker) Print() {
	fmt.Printf(" %d\t%d\t%d\t", w.stats.cumulative.Failed, w.stats.cumulative.Success, w.State.Busy)

	fmt.Print("|")
	for _, r := range w.renders {
		if r != nil {
			//print id as hexadecimal
			fmt.Printf("%02x|", r.id)
		} else {
			fmt.Print("  |")
		}
	}
	fmt.Println()
}
