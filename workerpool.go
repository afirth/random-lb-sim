package main

import (
	"fmt"
	"time"
)

type WorkerPool struct {
	name              string
	workers           []*Worker
	stats             *WorkerPoolStats
	autoscaleCooldown time.Time
}

type WorkerPoolStats struct {
	failed         int
	success        int
	busy           int
	slots          int //in slots - "cores"
	workers        int
	errorRate      float64
	utilization    float64
	avgUtilization float64
	sampleCount    int
}

func NewWorkerPool(name string, workerCount int, workerSlots int) *WorkerPool {
	workers := make([]*Worker, workerCount)
	for i := range workers {
		workers[i] = NewWorker(workerSlots)
	}
	return &WorkerPool{
		name:              name,
		workers:           workers,
		stats:             &WorkerPoolStats{},
		autoscaleCooldown: time.Now(),
	}
}

// add or remove workers based on error rate
// totally dirty hack and you'll need to tweak it to stabilize quickly
// speed is a divisor, 1-10
func (wp *WorkerPool) Autoscale(targetErrorRate float64, workerSlots int, delay int) {
	cooldownDelayUp := time.Second * time.Duration(delay) //TODO
	cooldownDelayDown := time.Second * time.Duration(2*delay)
	stats := wp.Stats()
	if time.Now().After(wp.autoscaleCooldown) {
		if stats.errorRate > targetErrorRate {
			wp.autoscaleCooldown = time.Now().Add(cooldownDelayUp)
			wp.AddWorker(1, workerSlots)
		} else if stats.errorRate < targetErrorRate/1.1 {
			wp.autoscaleCooldown = time.Now().Add(cooldownDelayDown)
			wp.RemoveWorker(1)
		}
	}
}

// generate work at rps
// try to add the work to the worker pool by selecting a random worker
func (wp *WorkerPool) AddWorkRandomly(rps int, mean float64, stdDev float64) {
	interval := time.Millisecond * time.Duration(1000/rps) //TODO verify this is correct
	ticker := time.NewTicker(time.Duration(interval))
	for {
		select {
		case <-ticker.C:
			worker := wp.SelectRandom()
			worker.Render(NewRender(mean, stdDev))
		}
	}
}

func (wp *WorkerPool) AddWorker(count int, slots int) {
	fmt.Printf("adding %d workers\n", count)
	for i := 0; i < count; i++ {
		wp.workers = append(wp.workers, NewWorker(slots))
	}
}

// lost work is ignored //BUG kinda
func (wp *WorkerPool) RemoveWorker(count int) {
	fmt.Printf("removing %d workers\n", count)
	wp.workers = wp.workers[:len(wp.workers)-count]
}

// select a random worker from the pool
func (wp *WorkerPool) SelectRandom() *Worker {
	pick := randGen.Intn(len(wp.workers))
	return wp.workers[pick]
}

func (wp WorkerPool) PrintHeader() {
	fmt.Printf("Failed\tOK\tRunning\tState\n")
}

func (wp *WorkerPool) Print() {
	fmt.Print("\033[H\033[2J") //clear (*nix)
	wp.PrintHeader()
	for _, w := range wp.workers {
		w.Print()
	}
	stats := wp.Stats()
	fmt.Printf("Slots:\t")
	fmt.Printf("busy %d\t", stats.busy)
	fmt.Printf("total %d\t", stats.slots)
	fmt.Printf("%%Util %.2f%%\t", stats.utilization*100)
	fmt.Println()
	fmt.Printf("Jobs:\t")
	fmt.Printf("ok %d\t", stats.success)
	fmt.Printf("total %d\t", stats.failed+stats.success)
	fmt.Printf("failed %d\t", stats.failed)
	fmt.Println()
	fmt.Printf("Stats:\t")
	fmt.Printf("%%OK: %.2f%%\t", 100-stats.errorRate*100)
	fmt.Printf("avgUtil %.2f%%\t", stats.avgUtilization*100)
	fmt.Println()
	// fmt.Printf("sampleCount %d\n", stats.sampleCount)
}

func (wp *WorkerPool) Stats() (s WorkerPoolStats) {
	s = *wp.stats
	s.slots, s.busy = 0, 0 // these two are not cumulative
	for _, w := range wp.workers {
		s.slots += w.State.Slots
		s.busy += w.State.Busy
		wc := w.LatestStats()
		s.failed += wc.Failed
		s.success += wc.Success
	}
	s.errorRate = float64(s.failed) / float64(s.failed+s.success)
	s.utilization = float64(s.busy) / float64(s.slots)
	// calculate the overall avgUtilization so far
	s.sampleCount += 1
	if s.sampleCount == 1 {
		s.avgUtilization = s.utilization
	} else {
		s.avgUtilization = (s.avgUtilization*float64(s.sampleCount-1) + s.utilization) / float64(s.sampleCount)
	}
	wp.stats = &s
	return
}
