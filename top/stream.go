package top

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/kr/pty"
)

// Stream provides top command output stream.
type Stream struct {
	cmd *exec.Cmd

	pmu sync.Mutex
	pt  *os.File

	// broadcast updates whenver available available
	wg      sync.WaitGroup
	rcond   *sync.Cond
	rmu     sync.RWMutex // protect results
	queue   []Row
	pid2Row map[int64]Row
	err     error
	errc    chan error

	// signal only once at initial, once the first line is ready
	readymu sync.Mutex
	ready   bool
	readyc  chan struct{}
}

// StartStream starts 'top' command stream.
func (cfg *Config) StartStream() (*Stream, error) {
	if err := cfg.createCmd(); err != nil {
		return nil, err
	}
	pt, err := pty.Start(cfg.cmd)
	if err != nil {
		return nil, err
	}

	str := &Stream{
		cmd: cfg.cmd,

		pmu: sync.Mutex{},
		pt:  pt,

		wg:  sync.WaitGroup{},
		rmu: sync.RWMutex{},

		// pre-allocate
		queue:   make([]Row, 0, 500),
		pid2Row: make(map[int64]Row, 500),
		err:     nil,
		errc:    make(chan error, 1),

		ready:  false,
		readyc: make(chan struct{}, 1),
	}
	str.rcond = sync.NewCond(&str.rmu)

	str.wg.Add(1)
	go str.enqueue()
	go str.dequeue()

	<-str.readyc
	return str, nil
}

// Stop kills the 'top' process and waits for it to exit.
func (str *Stream) Stop() error {
	return str.close(true)
}

// Wait just waits for the 'top' process to exit.
func (str *Stream) Wait() error {
	return str.close(false)
}

// ErrChan returns the error from stream.
func (str *Stream) ErrChan() <-chan error {
	return str.errc
}

// Latest returns the latest top command outputs.
func (str *Stream) Latest() map[int64]Row {
	str.rmu.RLock()
	cm := make(map[int64]Row, len(str.pid2Row))
	for k, v := range str.pid2Row {
		cm[k] = v
	}
	str.rmu.RUnlock()
	return cm
}

func (str *Stream) noError() (noErr bool) {
	str.rmu.RLock()
	noErr = str.err == nil
	str.rmu.RUnlock()
	return
}

// feed new top results into the queue
func (str *Stream) enqueue() {
	defer str.wg.Done()
	reader := bufio.NewReader(str.pt)
	for str.noError() {
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
		if len(row) != len(Headers) {
			str.rmu.Unlock()
			continue
		}

		r, rerr := parseRow(row)
		if rerr != nil {
			str.err = rerr
			str.rmu.Unlock()
			continue
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
func (str *Stream) dequeue() {
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

		str.pid2Row[row.PID] = row

		toc := false
		str.readymu.Lock()
		if !str.ready {
			str.ready = true
			toc = true
		}
		str.readymu.Unlock()
		if toc {
			close(str.readyc)
		}
	}
	if expectedErr(str.err) {
		str.err = nil
	}
	if str.err != nil {
		str.errc <- str.err
	}
	str.rmu.Unlock()
}

func (str *Stream) close(kill bool) (err error) {
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
		if !kill && strings.Contains(err.Error(), "exit status") {
			err = nil // non-zero exit code
		} else if kill && expectedErr(err) {
			err = nil
		}
	}
	str.cmd = nil
	return err
}

func expectedErr(err error) bool {
	if err == nil {
		return true
	}
	es := err.Error()
	return strings.Contains(es, "signal:") ||
		strings.Contains(es, "/dev/ptmx: input/output error") ||
		strings.Contains(es, "/dev/ptmx: file already closed")
}
