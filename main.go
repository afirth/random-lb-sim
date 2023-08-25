package main

// visualizes workers capable of performing parallel renders
// each worker is represented by a row, and each render is represented by an id

// each piece of incoming work has an associated duration which is a gaussian distribution
// if a render doesn't fit on the randomly assigned worker, it fails

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

var randGen = rand.New(rand.NewSource(42))

func main() {
	args := parseArgs()

	randomPool := NewWorkerPool("random", args.workerCount, args.workerSlots)
	go randomPool.AddWorkRandomly(args.rps, args.meanLatency, args.stdDevLatency)

	exitChan := spawnScreen()

Poller:
	for {
		select {
		case <-exitChan:
			break Poller

		case <-time.After(100 * time.Millisecond):
			randomPool.Autoscale(args.targetErrorRate, args.workerSlots, args.autoscaleDelay)
			randomPool.Print()
			fmt.Printf("inbound RPS: %d\n", args.rps)
		}
	}
	fmt.Print("\x1b[?1049l") // normal screen
	randomPool.Print()
}

func spawnScreen() (bye chan bool) {
	// start and clear the alt screen
	fmt.Print("\x1b[?1049h\x1b[2J\x1b[H")
	// break on enter (need to reset the screen)
	bye = make(chan bool)
	go func() {
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		bye <- true
	}()
	return
}
