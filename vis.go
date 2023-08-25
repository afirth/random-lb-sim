package main

// visualizes 10 workers capable of performing 10 parallel renders
// each worker is represented by a row, and each render is represented by an x

// each piece of incoming work has an associated duration which is a gaussian distribution durations of 1000ms, with a standard deviation of 200ms
// the work is generated at 10 rps

// the workers are capable of performing 10 renders in parallel
import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
)

var randGen = rand.New(rand.NewSource(42))

func main() {
	// create a new worker pool
	workerCount := 10
	workerSize := 4
	randomPool := NewWorkerPool("random", workerCount, workerSize)
	go inputLoop(randomPool, 50)
	for {
		randomPool.print()
		time.Sleep(200 * time.Millisecond)
	}
}

func inputLoop(pool WorkerPool, periodMillis int) {
	// generate work at 10rps
	// add the work to the worker pool
	interval := time.Millisecond * time.Duration(periodMillis)
	ticker := time.NewTicker(time.Duration(interval))
	defer ticker.Stop()

	// Start a goroutine to call periodicTask every n seconds
	go func() {
		for {
			select {
			case <-ticker.C:
				pool.AddRenderRandom() // discard failed renders for now
			}
		}
	}()
	select {} //don't stop
}

type Render struct {
	durationMillis time.Duration
}

func NewRender() Render {
	//generate random duration of normal distribution with mean 1000ms, std dev 200ms
	ms := int(math.Round(randGen.NormFloat64()*200 + 1000))
	return Render{
		durationMillis: time.Duration(ms) * time.Millisecond,
	}
}

type WorkerPool struct {
	name    string
	workers []*Worker
}

func NewWorkerPool(name string, workerCount int, workerSize int) WorkerPool {
	workers := make([]*Worker, workerCount)
	for i := range workers {
		workers[i] = NewWorker(workerSize)
	}
	return WorkerPool{
		name:    name,
		workers: workers,
	}
}
func (wp WorkerPool) print() {
	fmt.Print("\033[H\033[2J") //clear (*nix)
	failed, total := 0, 0
	for i, w := range wp.workers {
		fmt.Printf("worker %d:", i)
		w.print()
		fmt.Println()
		failed += w.failedRenders
		total += w.totalRenders
	}
	fmt.Printf("failed renders: %d\n", failed)
	fmt.Printf("total renders: %d\n", total)
}

// pick a random worker
// add the render to the worker
// if the worker is full, return an error
func (wp WorkerPool) AddRenderRandom() error {
	r := NewRender()
	pick := randGen.Intn(len(wp.workers))
	return wp.workers[pick].Render(r)
}

type Worker struct {
	renders       []*Render
	failedRenders int
	totalRenders  int
}

func NewWorker(size int) *Worker {
	return &Worker{
		renders:       make([]*Render, size),
		failedRenders: 0,
		totalRenders:  0,
	}
}
func (w *Worker) Render(r Render) error {
	// find the first open slot
	for i, render := range w.renders {
		if render == nil {
			w.renders[i] = &r
			go func() {
				// wait for the render to complete then clear the render
				time.Sleep(r.durationMillis)
				w.renders[i] = nil
				w.totalRenders++
			}()
			return nil
		}
	}
	w.failedRenders++
	return errors.New("no open slots")
}
func (w Worker) print() {
	// print the worker id
	// print an x for each render which is not null
	for _, r := range w.renders {
		if r != nil {
			fmt.Print("x")
		} else {
			fmt.Print(" ")
		}
	}
	fmt.Printf(" %d/%d", w.failedRenders, w.totalRenders)
}
