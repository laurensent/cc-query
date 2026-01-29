package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/term"
)

const (
	claudeOrangeFg = "\033[38;2;227;108;56m"
	dimOrangeFg    = "\033[38;2;50;24;12m"
	resetAttr      = "\033[0m"
)

type spinner struct {
	stop     chan struct{}
	done     chan struct{}
	stopOnce sync.Once
}

func getTermWidth() int {
	w, _, err := term.GetSize(int(os.Stderr.Fd()))
	if err != nil || w <= 0 {
		return 80
	}
	return w
}

func startSpinner() *spinner {
	s := &spinner{
		stop: make(chan struct{}),
		done: make(chan struct{}),
	}

	go func() {
		defer close(s.done)

		width := getTermWidth()
		segLen := width / 5
		if segLen < 6 {
			segLen = 6
		}
		maxPos := width - segLen

		const thinBar = "▔"

		// Insert a blank line at row 1, pushing existing content down.
		// This preserves all visible content — nothing is overwritten.
		fmt.Fprint(os.Stderr, "\033[s")       // save cursor
		fmt.Fprint(os.Stderr, "\033[1;1H")    // move to row 1
		fmt.Fprint(os.Stderr, "\033[L")       // insert line (content shifts down)
		fmt.Fprint(os.Stderr, "\033[u\033[B") // restore cursor + down 1 to compensate

		pos := 0
		dir := 1

		ticker := time.NewTicker(12 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-s.stop:
				// Delete the bar line, pulling content back up
				fmt.Fprint(os.Stderr, "\033[s")          // save cursor
				fmt.Fprint(os.Stderr, "\033[1;1H\033[M") // row 1 + delete line
				fmt.Fprint(os.Stderr, "\033[u\033[A")    // restore + up 1
				return
			case <-ticker.C:
				var buf strings.Builder
				buf.Grow(width * 24)
				for i := 0; i < width; i++ {
					if i >= pos && i < pos+segLen {
						buf.WriteString(claudeOrangeFg)
					} else {
						buf.WriteString(dimOrangeFg)
					}
					buf.WriteString(thinBar)
				}
				buf.WriteString(resetAttr)

				fmt.Fprintf(os.Stderr, "\033[s\033[1;1H%s\033[u", buf.String())

				pos += dir
				if pos >= maxPos {
					pos = maxPos
					dir = -1
				} else if pos <= 0 {
					pos = 0
					dir = 1
				}
			}
		}
	}()
	return s
}

func (s *spinner) Stop() {
	s.stopOnce.Do(func() {
		close(s.stop)
		select {
		case <-s.done:
		case <-time.After(500 * time.Millisecond):
		}
	})
}
