package proc

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/gyuho/linux-inspect/pkg/fileutil"

	"bytes"
)

// GetNetTCPByPID reads '/proc/$PID/net/tcp(6)' data.
func GetNetTCPByPID(pid int64, tp TransportProtocol) ([]NetTCP, error) {
	d, err := readNetTCP(pid, tp)
	if err != nil {
		return nil, err
	}

	var ipParse func(string) (string, int64, error)
	switch tp {
	case TypeTCP:
		ipParse = parseLittleEndianIpv4
	case TypeTCP6:
		ipParse = parseLittleEndianIpv6
	}
	return parseNetTCP(d, ipParse, tp.String())
}

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
		panic(fmt.Errorf("unknown transport protocol %d", tp))
	}
}

type netColumnIndex int

const (
	net_tcp_idx_sl netColumnIndex = iota
	net_tcp_idx_local_address
	net_tcp_idx_remote_address
	net_tcp_idx_st
	net_tcp_idx_tx_queue_rx_queue
	net_tcp_idx_tr_tm_when
	net_tcp_idx_retrnsmt
	net_tcp_idx_uid
	net_tcp_idx_timeout
	net_tcp_idx_inode
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

func parseNetTCP(d []byte, ipParse func(string) (string, int64, error), ipType string) ([]NetTCP, error) {
	rows := [][]string{}

	first := true
	scanner := bufio.NewScanner(bytes.NewReader(d))
	for scanner.Scan() {
		txt := scanner.Text()
		if len(txt) == 0 {
			continue
		}

		fs := strings.Fields(txt)
		if len(fs) < int(net_tcp_idx_inode+1) {
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
		copy(row, fs[:net_tcp_idx_inode+1])

		rows = append(rows, row)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	nch, errc := make(chan NetTCP), make(chan error)
	for _, row := range rows {
		go func(row []string) {
			np := NetTCP{}

			np.Type = ipType

			sn, err := strconv.ParseUint(strings.Replace(row[net_tcp_idx_sl], ":", "", -1), 10, 64)
			if err != nil {
				errc <- err
				return
			}
			np.Sl = sn

			np.LocalAddress = strings.TrimSpace(row[net_tcp_idx_local_address])
			lp, lt, err := ipParse(row[net_tcp_idx_local_address])
			if err != nil {
				errc <- err
				return
			}
			np.LocalAddressParsedIPHost = strings.TrimSpace(lp)
			np.LocalAddressParsedIPPort = lt

			np.RemAddress = strings.TrimSpace(row[net_tcp_idx_remote_address])
			rp, rt, err := ipParse(row[net_tcp_idx_remote_address])
			if err != nil {
				errc <- err
				return
			}
			np.RemAddressParsedIPHost = strings.TrimSpace(rp)
			np.RemAddressParsedIPPort = rt

			np.St = strings.TrimSpace(row[net_tcp_idx_st])
			np.StParsedStatus = strings.TrimSpace(netTCPStatus[row[net_tcp_idx_st]])

			qs := strings.Split(row[net_tcp_idx_tx_queue_rx_queue], ":")
			if len(qs) == 2 {
				np.TxQueue = qs[0]
				np.RxQueue = qs[1]
			}
			trs := strings.Split(row[net_tcp_idx_tr_tm_when], ":")
			if len(trs) == 2 {
				np.Tr = trs[0]
				np.TmWhen = trs[1]
			}
			np.Retrnsmt = row[net_tcp_idx_retrnsmt]

			un, err := strconv.ParseUint(row[net_tcp_idx_uid], 10, 64)
			if err != nil {
				errc <- err
				return
			}
			np.Uid = un

			to, err := strconv.ParseUint(row[net_tcp_idx_timeout], 10, 64)
			if err != nil {
				errc <- err
				return
			}
			np.Timeout = to

			np.Inode = strings.TrimSpace(row[net_tcp_idx_inode])

			nch <- np
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

func readNetTCP(pid int64, tp TransportProtocol) ([]byte, error) {
	fpath := fmt.Sprintf("/proc/%d/net/%s", pid, tp.String())
	f, err := fileutil.OpenToRead(fpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}
