package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

type args struct {
	meanLatency     float64
	stdDevLatency   float64
	workerCount     int
	workerSlots     int
	targetErrorRate float64
	autoscaleDelay  int
	rps             int
}

func newArgs() *args {
	return &args{
		meanLatency:     1000,
		stdDevLatency:   1000,
		workerCount:     12,
		workerSlots:     15,
		targetErrorRate: 0.001, // 3 9's
		autoscaleDelay:  5,     // 1-10
		rps:             50,    //max 999 ;)
	}
}

func parseArgs() *args {
	a := newArgs()

	pflag.Float64VarP(&a.meanLatency, "meanLatency", "l", a.meanLatency, "Mean latency of render (ms)")
	pflag.Float64VarP(&a.stdDevLatency, "stdDevLatency", "d", a.stdDevLatency, "Standard deviation of latency (ms)")
	pflag.IntVarP(&a.workerCount, "workerCount", "w", a.workerCount, "Number of workers")
	pflag.IntVarP(&a.workerSlots, "workerSlots", "s", a.workerSlots, "Number of worker slots")
	pflag.Float64VarP(&a.targetErrorRate, "targetErrorRate", "e", a.targetErrorRate, "Target error rate (0.001 is 3 9's)")
	pflag.IntVarP(&a.autoscaleDelay, "autoscaleDelay", "a", a.autoscaleDelay, "Autoscale delay (seconds) - 1-10 is reasonable")
	pflag.IntVarP(&a.rps, "rps", "r", a.rps, "Requests per second to spawn (max 999)")

	pflag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		pflag.PrintDefaults()
	}

	pflag.Parse()
	return a
}
