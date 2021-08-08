package main

import (
	"fmt"
	"time"
)

type FpsCounter struct {
	sum      float64
	counter  int
	divider  int
	fDivider float64
}

func NewFpsCounter(divider int) *FpsCounter {
	return &FpsCounter{
		divider:  divider,
		fDivider: float64(divider),
	}
}

func (f *FpsCounter) update(d time.Duration) {
	f.sum += d.Seconds()
	f.counter++
	if f.counter >= f.divider {
		fmt.Printf("FPS: %.2f \n", f.fDivider/f.sum)
		f.counter = 0
		f.sum = 0
	}
}
