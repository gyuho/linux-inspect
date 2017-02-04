package psn

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/kr/pty"
)

// TopStream provides top command output stream.
type TopStream struct {
	cmd *exec.Cmd

	pmu sync.Mutex
	pt  *os.File

	// broadcast updates whenver available available
	wg                sync.WaitGroup
	rcond             *sync.Cond
	rmu               sync.RWMutex // protect results
	queue             []TopCommandRow
	pid2TopCommandRow map[int64]TopCommandRow
	err               error
	errc              chan error
}

// StartStream starts 'top' command stream.
func (cfg *TopConfig) StartStream() (*TopStream, error) {
	if err := cfg.createCmd(); err != nil {
		return nil, err
	}
	pt, err := pty.Start(cfg.cmd)
	if err != nil {
		return nil, err
	}

	str := &TopStream{
		cmd: cfg.cmd,

		pmu: sync.Mutex{},
		pt:  pt,

		wg:  sync.WaitGroup{},
		rmu: sync.RWMutex{},

		// pre-allocate
		queue:             make([]TopCommandRow, 0, 100),
		pid2TopCommandRow: make(map[int64]TopCommandRow, 500),
		err:               nil,
		errc:              make(chan error, 1),
	}
	str.rcond = sync.NewCond(&str.rmu)

	str.wg.Add(1)
	go str.enqueue()
	go str.dequeue()

	return str, nil
}

// ErrChan returns the error from stream.
func (str *TopStream) ErrChan() <-chan error {
	return str.errc
}

// feed new top results into the queue
func (str *TopStream) enqueue() {
	defer str.wg.Done()
	reader := bufio.NewReader(str.pt)
	for str.err == nil {
		// lock for pty
		str.pmu.Lock()
		data, _, lerr := reader.ReadLine()
		str.pmu.Unlock()

		data = bytes.TrimSpace(data)
		if topRowToSkip(data) {
			continue
		}
		line := string(data)

		// lock for results
		str.rmu.Lock()

		str.err = lerr
		if line == "" {
			str.rmu.Unlock()
			continue
		}

		row := strings.Fields(line)
		if len(row) != len(TopRowHeaders) {
			str.rmu.Unlock()
			continue
		}
		r, rerr := parseTopRow(row)
		if rerr != nil {
			str.err = rerr
		}

		str.queue = append(str.queue, r)
		if len(str.queue) == 1 {
			// we have a new output; signal!
			str.rcond.Signal()
		}
		str.rmu.Unlock()
	}

	// we got error; signal!
	str.rcond.Signal()
}

// dequeue polls from 'top' process.
// And signals error channel if any.
func (str *TopStream) dequeue() {
	str.rmu.Lock()
	for {
		// wait until there's output
		for len(str.queue) == 0 && str.err == nil {
			str.rcond.Wait()
		}

		// no output; should be error
		if len(str.queue) == 0 {
			break
		}

		row := str.queue[0]
		str.queue = str.queue[1:]

		str.pid2TopCommandRow[row.PID] = row
	}
	if expectedErr(str.err) {
		str.err = nil
	}
	str.rmu.Unlock()

	if str.err != nil {
		str.errc <- str.err
	}
}

func (str *TopStream) close(kill bool) (err error) {
	if str.cmd == nil {
		return str.err
	}
	if kill {
		str.cmd.Process.Kill()
	}

	err = str.cmd.Wait()

	str.pmu.Lock()
	str.pt.Close() // close file
	str.pmu.Unlock()

	str.wg.Wait()

	if err != nil {
		str.err = err
		if !kill && strings.Contains(err.Error(), "exit status") {
			str.err = nil // non-zero exit code
		} else if kill && expectedErr(err) {
			str.err = nil
		}
	}
	str.cmd = nil
	return str.err
}

func expectedErr(err error) bool {
	if err == nil {
		return true
	}
	es := err.Error()
	return strings.Contains(es, "signal:") || strings.Contains(es, "/dev/ptmx: input/output error")
}

// Stop kills the 'top' process and waits for it to exit.
func (str *TopStream) Stop() error {
	return str.close(true)
}

// Wait just waits for the 'top' process to exit.
func (str *TopStream) Wait() error {
	return str.close(false)
}

// Latest returns the latest top command outputs.
func (str *TopStream) Latest() map[int64]TopCommandRow {
	str.rmu.RLock()
	cm := make(map[int64]TopCommandRow, len(str.pid2TopCommandRow))
	for k, v := range str.pid2TopCommandRow {
		cm[k] = v
	}
	str.rmu.RUnlock()
	return cm
}
