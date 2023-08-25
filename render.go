package main

import (
	"math"
	"time"
)

type Render struct {
	id       int
	duration time.Duration
}

func NewRender(mean float64, stdDev float64) Render {
	// generate random distribution
	ms := int(math.Round(randGen.NormFloat64()*stdDev + mean))
	return Render{
		duration: time.Duration(ms) * time.Millisecond,
		id:       randGen.Intn(254), //to be rendered as %02x - not used for anything but visualization
	}
}
