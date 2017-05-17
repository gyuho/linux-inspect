package proc

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/gyuho/linux-inspect/pkg/fileutil"

	humanize "github.com/dustin/go-humanize"
)

type netDevColumnIndex int

const (
	net_dev_idx_interface netDevColumnIndex = iota

	net_dev_idx_receive_bytes
	net_dev_idx_receive_packets
	net_dev_idx_receive_errs
	net_dev_idx_receive_drop
	net_dev_idx_receive_fifo
	net_dev_idx_receive_frame
	net_dev_idx_receive_compressed
	net_dev_idx_receive_multicast

	net_dev_idx_transmit_bytes
	net_dev_idx_transmit_packets
	net_dev_idx_transmit_errs
	net_dev_idx_transmit_drop
	net_dev_idx_transmit_fifo
	net_dev_idx_transmit_colls
	net_dev_idx_transmit_carrier
)

// GetNetDev reads '/proc/net/dev'.
func GetNetDev() (nds []NetDev, err error) {
	var d []byte
	d, err = readNetDev()
	if err != nil {
		return nil, err
	}

	header := true
	scanner := bufio.NewScanner(bytes.NewReader(d))
	for scanner.Scan() {
		txt := scanner.Text()
		if len(txt) == 0 {
			continue
		}

		ds := strings.Fields(strings.TrimSpace(txt))
		if header {
			if strings.HasPrefix(ds[0], "Inter") {
				continue
			}
			if strings.HasSuffix(ds[0], "face") {
				header = false
				continue
			}
		}
		if len(ds) < int(net_dev_idx_transmit_carrier+1) {
			return nil, fmt.Errorf("not enough columns at %v", ds)
		}

		d := NetDev{}

		d.Interface = strings.TrimSpace(ds[net_dev_idx_interface])
		d.Interface = d.Interface[:len(d.Interface)-1] // remove ':' from 'wlp2s0:'

		mn, err := strconv.ParseUint(ds[net_dev_idx_receive_bytes], 10, 64)
		if err != nil {
			return nil, err
		}
		d.ReceiveBytes = mn
		d.ReceiveBytesBytesN = mn
		d.ReceiveBytesParsedBytes = humanize.Bytes(mn)

		mn, err = strconv.ParseUint(ds[net_dev_idx_transmit_bytes], 10, 64)
		if err != nil {
			return nil, err
		}
		d.TransmitBytes = mn
		d.TransmitBytesBytesN = mn
		d.TransmitBytesParsedBytes = humanize.Bytes(mn)

		mn, err = strconv.ParseUint(ds[net_dev_idx_receive_packets], 10, 64)
		if err != nil {
			return nil, err
		}
		d.ReceivePackets = mn

		mn, err = strconv.ParseUint(ds[net_dev_idx_receive_errs], 10, 64)
		if err != nil {
			return nil, err
		}
		d.ReceiveErrs = mn

		mn, err = strconv.ParseUint(ds[net_dev_idx_receive_drop], 10, 64)
		if err != nil {
			return nil, err
		}
		d.ReceiveDrop = mn

		mn, err = strconv.ParseUint(ds[net_dev_idx_receive_fifo], 10, 64)
		if err != nil {
			return nil, err
		}
		d.ReceiveFifo = mn

		mn, err = strconv.ParseUint(ds[net_dev_idx_receive_frame], 10, 64)
		if err != nil {
			return nil, err
		}
		d.ReceiveFrame = mn

		mn, err = strconv.ParseUint(ds[net_dev_idx_receive_compressed], 10, 64)
		if err != nil {
			return nil, err
		}
		d.ReceiveCompressed = mn

		mn, err = strconv.ParseUint(ds[net_dev_idx_receive_multicast], 10, 64)
		if err != nil {
			return nil, err
		}
		d.ReceiveMulticast = mn

		mn, err = strconv.ParseUint(ds[net_dev_idx_transmit_packets], 10, 64)
		if err != nil {
			return nil, err
		}
		d.TransmitPackets = mn

		mn, err = strconv.ParseUint(ds[net_dev_idx_transmit_errs], 10, 64)
		if err != nil {
			return nil, err
		}
		d.TransmitErrs = mn

		mn, err = strconv.ParseUint(ds[net_dev_idx_transmit_drop], 10, 64)
		if err != nil {
			return nil, err
		}
		d.TransmitDrop = mn

		mn, err = strconv.ParseUint(ds[net_dev_idx_transmit_fifo], 10, 64)
		if err != nil {
			return nil, err
		}
		d.TransmitFifo = mn

		mn, err = strconv.ParseUint(ds[net_dev_idx_transmit_colls], 10, 64)
		if err != nil {
			return nil, err
		}
		d.TransmitColls = mn

		mn, err = strconv.ParseUint(ds[net_dev_idx_transmit_carrier], 10, 64)
		if err != nil {
			return nil, err
		}
		d.TransmitCarrier = mn

		nds = append(nds, d)
	}

	return nds, nil
}

func readNetDev() ([]byte, error) {
	f, err := fileutil.OpenToRead("/proc/net/dev")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}
