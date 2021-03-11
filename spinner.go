package spinner

import (
	"fmt"
	"time"
)

type Spinner struct {
	Shape string
	tick  *time.Ticker
	stop  chan int
}

func New() *Spinner {
	return &Spinner{
		Shape: DEFAULT,
		tick:  time.NewTicker(100 * time.Millisecond),
		stop:  make(chan int),
	}
}

// during spinning make sure not to print anything to stdout
func (s *Spinner) Spin(format string, a ...interface{}) {
	s.tick.Reset(100 * time.Millisecond) // in case it was stopped before
	prepareTerminal(format, a...)

	go func() {
		for {
			for _, c := range s.Shape {
				select {
				case <-s.tick.C:
					fmt.Printf("%c\b", c)
				case <-s.stop:
					restoreCursor()
					return
				}
			}
		}
	}()
}

func (s *Spinner) Stop() {
	s.tick.Stop()
	s.stop <- 1
}

func prepareTerminal(format string, a ...interface{}) {
	fmt.Printf(format + " ", a...)
	fmt.Print("\x1b[?25l") // hide cursor
}

func restoreCursor() {
	// we first print a white character to overwrite the last
	// printed shape by spinner.
	fmt.Print(" \n\x1b[?25h")
}
