package logsync

import (
	"bufio"
	"io"
	"strings"
	"sync"
	"time"
)

// Helper struct to contain one log stream/buffer
type queued struct {
	buffer  []string       // All logs read from this stream (in order)
	scanner *bufio.Scanner // Stream
	mutex   sync.Mutex     // Buffer access mutex
}

// Returns whether more lines are available from this stream
func (q *queued) ready() bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return len(q.buffer) > 0
}

// Get the frontmost line
func (q *queued) peek() string {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.buffer[0]
}

// Remove the frontmost line
func (q *queued) pop() {
	q.mutex.Lock()
	q.buffer = q.buffer[1:]
	q.mutex.Unlock()
}

// Sync returns an io.Reader that synchronizes all the readers chronologically
func Sync(readers []io.Reader, processor func(string) string) io.Reader {
	qs := make([]*queued, len(readers))

	// Return
	reader, writer := io.Pipe()

	// number of active logs
	active := len(readers)
	var activeLock sync.RWMutex

	// Action mutex
	var mutex sync.Mutex
	mutex.Lock()
	logs := 0 // number of logs lines in all pods

	// There should be one of these called for each log
	action := func() {
		mutex.Lock()
		defer mutex.Unlock()
		if logs == 0 {
			activeLock.RLock()
			if active == 0 {
				writer.Close()
			}
			activeLock.RUnlock()
			return
		}
		recent := time.Time{}
		ind := -1
		for q := range qs {
			if qs[q].ready() {
				s := qs[q].peek()
				if i := strings.Index(s, " "); i != -1 {
					t, err := time.Parse(time.RFC3339Nano, s[:i])
					if err != nil { // Try non-nano
						t, _ = time.Parse(time.RFC3339, s[:i])
						// If also error... just continue downwards
					}
					if recent.IsZero() || t.Before(recent) {
						ind = q
						recent = t
					}
				} else { // hmm
					continue
				}
			}
		}
		if ind != -1 {
			s := qs[ind].peek()
			writer.Write([]byte(processor(s)))
			qs[ind].pop()
			logs--
		}
	}

	for i, r := range readers {
		q := &queued{
			scanner: bufio.NewScanner(r),
		}

		go func() {
			for q.scanner.Scan() {
				q.mutex.Lock()
				q.buffer = append(q.buffer, q.scanner.Text())
				logs++
				q.mutex.Unlock()
				go action()
			}
			activeLock.Lock()
			active--
			activeLock.Unlock()
			go action()
		}()

		qs[i] = q
	}

	mutex.Unlock()

	return reader
}
