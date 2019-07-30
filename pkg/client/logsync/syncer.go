package logsync

import (
	"bufio"
	"io"
	"strings"
	"sync"
	"time"
)

type queued struct {
	buffer  []string
	scanner *bufio.Scanner
	mutex   sync.Mutex
}

func (q *queued) ready() bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return len(q.buffer) > 0
}

func (q *queued) peek() string {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.buffer[0]
}

func (q *queued) pop() {
	q.mutex.Lock()
	q.buffer = q.buffer[1:]
	q.mutex.Unlock()
}

func getQueued(reader io.Reader, action func()) *queued {
	q := queued{
		scanner: bufio.NewScanner(reader),
	}

	go func() {
		for {
			if q.scanner.Scan() {
				q.mutex.Lock()
				q.buffer = append(q.buffer, q.scanner.Text())
				q.mutex.Unlock()
			} else { // what to do when errors
				break
			}
			action()
		}
	}()

	return &q
}

// Sync returns an io.Reader that synchronizes all the readers chronologically
func Sync(readers []io.Reader) io.Reader {
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

	action := func() {
		mutex.Lock()
		defer mutex.Unlock()
		if logs == 0 {
			if active == 0 {
				writer.Close()
			}
			return
		}
		recent := time.Time{}
		ind := -1
		for q := range qs {
			if qs[q].ready() {
				s := qs[q].peek()
				if i := strings.Index(s, " "); i != -1 {
					t, _ := time.Parse(time.RFC3339Nano, s[:i-1])
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
			writer.Write([]byte(s[strings.Index(s, " ")+1:] + "\n"))
			qs[ind].pop()
			logs--
		}
	}

	for i, reader := range readers {
		q := queued{
			scanner: bufio.NewScanner(reader),
		}

		go func() {
			for {
				if q.scanner.Scan() {
					mutex.Lock()
					q.mutex.Lock()
					q.buffer = append(q.buffer, q.scanner.Text())
					logs++
					q.mutex.Unlock()
					mutex.Unlock()
				} else { // what to do when errors
					activeLock.Lock()
					active--
					activeLock.Unlock()
					break
				}
				action()
			}
			action()
		}()

		qs[i] = &q
	}

	mutex.Unlock()

	return reader
}
