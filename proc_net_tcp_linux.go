package psn

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

// TransportProtocol is tcp, tcp6.
type TransportProtocol int

const (
	TypeTCP TransportProtocol = iota
	TypeTCP6
)

func (tp TransportProtocol) String() string {
	switch tp {
	case TypeTCP:
		return "tcp"
	case TypeTCP6:
		return "tcp6"
	default:
		panic(fmt.Errorf("unknown %v", tp))
	}
}

// GetProcNetTCPByPID reads '/proc/$PID/net/tcp(6)' data.
func GetProcNetTCPByPID(pid int64, tp TransportProtocol) (ss []NetTCP, err error) {
	return parseProcNetTCPByPID(pid, tp)
}

type procNetColumnIndex int

const (
	proc_net_tcp_idx_sl procNetColumnIndex = iota
	proc_net_tcp_idx_local_address
	proc_net_tcp_idx_remote_address
	proc_net_tcp_idx_st
	proc_net_tcp_idx_tx_queue_rx_queue
	proc_net_tcp_idx_tr_tm_when
	proc_net_tcp_idx_retrnsmt
	proc_net_tcp_idx_uid
	proc_net_tcp_idx_timeout
	proc_net_tcp_idx_inode
)

var (
	// RPC_SHOW_SOCK
	// https://github.com/torvalds/linux/blob/master/include/trace/events/sunrpc.h
	netTCPStatus = map[string]string{
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

func parseProcNetTCPByPID(pid int64, tp TransportProtocol) ([]NetTCP, error) {
	fpath := fmt.Sprintf("/proc/%d/net/%s", pid, tp.String())
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
		if len(fs) < int(proc_net_tcp_idx_inode+1) {
			return nil, fmt.Errorf("not enough columns at %v", fs)
		}

		if first {
			if fs[0] != "sl" { // header
				return nil, fmt.Errorf("first line must be columns but got = %#q", fs)
			}
			first = false
			continue
		}

		row := make([]string, 10)
		copy(row, fs[:proc_net_tcp_idx_inode+1])

		rows = append(rows, row)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	var ipParse func(string) (string, int64, error)
	switch tp {
	case TypeTCP:
		ipParse = parseLittleEndianIpv4
	case TypeTCP6:
		ipParse = parseLittleEndianIpv6
	}

	nch, errc := make(chan NetTCP), make(chan error)
	for _, row := range rows {
		go func(row []string) {
			ns := NetTCP{}

			ns.Type = tp.String()

			sn, err := strconv.ParseUint(strings.Replace(row[proc_net_tcp_idx_sl], ":", "", -1), 10, 64)
			if err != nil {
				errc <- err
				return
			}
			ns.Sl = sn

			ns.LocalAddress = strings.TrimSpace(row[proc_net_tcp_idx_local_address])
			lp, lt, err := ipParse(row[proc_net_tcp_idx_local_address])
			if err != nil {
				errc <- err
				return
			}
			ns.LocalAddressParsedIPHost = strings.TrimSpace(lp)
			ns.LocalAddressParsedIPPort = lt

			ns.RemAddress = strings.TrimSpace(row[proc_net_tcp_idx_remote_address])
			rp, rt, err := ipParse(row[proc_net_tcp_idx_remote_address])
			if err != nil {
				errc <- err
				return
			}
			ns.RemAddressParsedIPHost = strings.TrimSpace(rp)
			ns.RemAddressParsedIPPort = rt

			ns.St = strings.TrimSpace(row[proc_net_tcp_idx_st])
			ns.StParsedStatus = strings.TrimSpace(netTCPStatus[row[proc_net_tcp_idx_st]])

			qs := strings.Split(row[proc_net_tcp_idx_tx_queue_rx_queue], ":")
			if len(qs) == 2 {
				ns.TxQueue = qs[0]
				ns.RxQueue = qs[1]
			}
			trs := strings.Split(row[proc_net_tcp_idx_tr_tm_when], ":")
			if len(trs) == 2 {
				ns.Tr = trs[0]
				ns.TmWhen = trs[1]
			}
			ns.Retrnsmt = row[proc_net_tcp_idx_retrnsmt]

			un, err := strconv.ParseUint(row[proc_net_tcp_idx_uid], 10, 64)
			if err != nil {
				errc <- err
				return
			}
			ns.Uid = un

			to, err := strconv.ParseUint(row[proc_net_tcp_idx_timeout], 10, 64)
			if err != nil {
				errc <- err
				return
			}
			ns.Timeout = to

			ns.Inode = strings.TrimSpace(row[proc_net_tcp_idx_inode])

			nch <- ns
		}(row)

	}

	nss := make([]NetTCP, 0, len(rows))
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
