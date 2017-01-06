package psn

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// TransportProtocol is tcp, tcp6.
type TransportProtocol int

const (
	TCP TransportProtocol = iota
	TCP6
)

// GetNetstat reads /proc/$PID/net/tcp(6) data.
func GetNetstat(pid int64, tp TransportProtocol) (ss []NetStat, err error) {
	for i := 0; i < 5; i++ {
		ss, err = parseProcNetStat(pid, tp)
		if err == nil {
			return ss, nil
		}
		log.Println("retrying;", err)
		time.Sleep(5 * time.Millisecond)
	}
	return
}

type colIdx int

const (
	// http://www.onlamp.com/pub/a/linux/2000/11/16/LinuxAdmin.html
	// [`sl` `local_address` `rem_address` `st` `tx_queue` `rx_queue` `tr`
	// `tm->when` `retrnsmt` `uid` `timeout` `inode`]
	idx_sl colIdx = iota
	idx_local_address
	idx_remote_address
	idx_st
	idx_tx_queue_rx_queue
	idx_tr_tm_when
	idx_retrnsmt
	idx_uid
	idx_timeout
	idx_inode
)

var (
	// RPC_SHOW_SOCK
	// https://github.com/torvalds/linux/blob/master/include/trace/events/sunrpc.h
	netstatStatus = map[string]string{
		"01": "ESTABLISHED",
		"02": "SYN_SENT",
		"03": "SYN_RECV",
		"04": "FIN_WAIT1",
		"05": "FIN_WAIT2",
		"06": "TIME_WAIT",
		"07": "CLOSE",
		"08": "CLOSE_WAIT",
		"09": "LAST_ACK",
		"0A": "LISTEN",
		"0B": "CLOSING",
	}
)

func parseProcNetStat(pid int64, tp TransportProtocol) ([]NetStat, error) {
	fpath := fmt.Sprintf("/proc/%d/net/tcp", pid)
	if tp == TCP6 {
		fpath += "6"
	}
	f, err := openToRead(fpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows := [][]string{}

	first := true
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		if len(txt) == 0 {
			continue
		}

		fs := strings.Fields(txt)
		if first {
			if fs[0] != "sl" { // header
				return nil, fmt.Errorf("first line must be columns but got = %#q", fs)
			}
			first = false
			continue
		}

		row := make([]string, 10)
		copy(row, fs[:idx_inode+1])

		rows = append(rows, row)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	var ipParse func(string) (string, int64, error)
	switch tp {
	case TCP:
		ipParse = parseLittleEndianIpv4
	case TCP6:
		ipParse = parseLittleEndianIpv6
	}

	nch, errc := make(chan NetStat), make(chan error)
	for _, row := range rows {
		go func(row []string) {
			ns := NetStat{}

			ns.Protocol = "tcp"
			if tp == TCP6 {
				ns.Protocol += "6"
			}

			sn, err := strconv.ParseUint(strings.Replace(row[idx_sl], ":", "", -1), 10, 64)
			if err != nil {
				errc <- err
				return
			}
			ns.Sl = sn

			ns.LocalAddress = strings.TrimSpace(row[idx_local_address])
			lp, lt, err := ipParse(row[idx_local_address])
			if err != nil {
				errc <- err
				return
			}
			ns.LocalAddressParsedIPHost = strings.TrimSpace(lp)
			ns.LocalAddressParsedIPPort = lt

			ns.RemAddress = strings.TrimSpace(row[idx_remote_address])
			rp, rt, err := ipParse(row[idx_remote_address])
			if err != nil {
				errc <- err
				return
			}
			ns.RemAddressParsedIPHost = strings.TrimSpace(rp)
			ns.RemAddressParsedIPPort = rt

			ns.St = strings.TrimSpace(row[idx_st])
			ns.StParsedStatus = strings.TrimSpace(netstatStatus[row[idx_st]])

			qs := strings.Split(row[idx_tx_queue_rx_queue], ":")
			if len(qs) == 2 {
				ns.TxQueue = qs[0]
				ns.RxQueue = qs[1]
			}
			trs := strings.Split(row[idx_tr_tm_when], ":")
			if len(trs) == 2 {
				ns.Tr = trs[0]
				ns.TmWhen = trs[1]
			}
			ns.Retrnsmt = row[idx_retrnsmt]

			un, err := strconv.ParseUint(row[idx_uid], 10, 64)
			if err != nil {
				errc <- err
				return
			}
			ns.Uid = un

			to, err := strconv.ParseUint(row[idx_timeout], 10, 64)
			if err != nil {
				errc <- err
				return
			}
			ns.Timeout = to

			ns.Inode = strings.TrimSpace(row[idx_inode])

			nch <- ns
		}(row)

	}

	nss := make([]NetStat, 0, len(rows))
	cn, limit := 0, len(rows)
	for cn != limit {
		select {
		case err := <-errc:
			return nil, err
		case p := <-nch:
			nss = append(nss, p)
			cn++
		}
	}
	close(nch)
	close(errc)

	return nss, nil
}

// SearchInode finds the matching process to the given inode.
func SearchInode(fds []string, inode string) (pid int64) {
	var mu sync.RWMutex

	var wg sync.WaitGroup
	wg.Add(len(fds))
	for _, fd := range fds {
		go func(fdpath string) {
			defer wg.Done()

			mu.RLock()
			done := pid != 0
			mu.RUnlock()
			if done {
				return
			}

			// '/proc/[pid]/fd' contains type:[inode]
			sym, err := os.Readlink(fdpath)
			if err != nil {
				return
			}
			if !strings.Contains(strings.TrimSpace(sym), inode) {
				return
			}

			pd, err := pidFromFd(fdpath)
			if err != nil {
				return
			}
			mu.Lock()
			pid = pd
			mu.Unlock()
		}(fd)
	}
	wg.Wait()

	if pid == 0 {
		pid = -1
	}
	return
}
