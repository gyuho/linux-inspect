package psn

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	humanize "github.com/dustin/go-humanize"
)

type procNetDevColumnIndex int

const (
	proc_net_dev_idx_interface procNetDevColumnIndex = iota

	proc_net_dev_idx_receive_bytes
	proc_net_dev_idx_receive_packets
	proc_net_dev_idx_receive_errs
	proc_net_dev_idx_receive_drop
	proc_net_dev_idx_receive_fifo
	proc_net_dev_idx_receive_frame
	proc_net_dev_idx_receive_compressed
	proc_net_dev_idx_receive_multicast

	proc_net_dev_idx_transmit_bytes
	proc_net_dev_idx_transmit_packets
	proc_net_dev_idx_transmit_errs
	proc_net_dev_idx_transmit_drop
	proc_net_dev_idx_transmit_fifo
	proc_net_dev_idx_transmit_colls
	proc_net_dev_idx_transmit_carrier
)

// GetProcNetDev reads '/proc/net/dev'.
func GetProcNetDev() ([]NetDev, error) {
	f, err := openToRead("/proc/net/dev")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	header := true
	dss := []NetDev{}
	scanner := bufio.NewScanner(f)
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
		if len(ds) < int(proc_net_dev_idx_transmit_carrier+1) {
			return nil, fmt.Errorf("not enough columns at %v", ds)
		}

		d := NetDev{}

		d.Interface = strings.TrimSpace(ds[proc_net_dev_idx_interface])
		d.Interface = d.Interface[:len(d.Interface)-1] // remove ':' from 'wlp2s0:'

		mn, err := strconv.ParseUint(ds[proc_net_dev_idx_receive_bytes], 10, 64)
		if err != nil {
			return nil, err
		}
		d.ReceiveBytes = mn
		d.ReceiveBytesBytesN = mn
		d.ReceiveBytesParsedBytes = humanize.Bytes(mn)

		mn, err = strconv.ParseUint(ds[proc_net_dev_idx_transmit_bytes], 10, 64)
		if err != nil {
			return nil, err
		}
		d.TransmitBytes = mn
		d.TransmitBytesBytesN = mn
		d.TransmitBytesParsedBytes = humanize.Bytes(mn)

		mn, err = strconv.ParseUint(ds[proc_net_dev_idx_receive_packets], 10, 64)
		if err != nil {
			return nil, err
		}
		d.ReceivePackets = mn

		mn, err = strconv.ParseUint(ds[proc_net_dev_idx_receive_errs], 10, 64)
		if err != nil {
			return nil, err
		}
		d.ReceiveErrs = mn

		mn, err = strconv.ParseUint(ds[proc_net_dev_idx_receive_drop], 10, 64)
		if err != nil {
			return nil, err
		}
		d.ReceiveDrop = mn

		mn, err = strconv.ParseUint(ds[proc_net_dev_idx_receive_fifo], 10, 64)
		if err != nil {
			return nil, err
		}
		d.ReceiveFifo = mn

		mn, err = strconv.ParseUint(ds[proc_net_dev_idx_receive_frame], 10, 64)
		if err != nil {
			return nil, err
		}
		d.ReceiveFrame = mn

		mn, err = strconv.ParseUint(ds[proc_net_dev_idx_receive_compressed], 10, 64)
		if err != nil {
			return nil, err
		}
		d.ReceiveCompressed = mn

		mn, err = strconv.ParseUint(ds[proc_net_dev_idx_receive_multicast], 10, 64)
		if err != nil {
			return nil, err
		}
		d.ReceiveMulticast = mn

		mn, err = strconv.ParseUint(ds[proc_net_dev_idx_transmit_packets], 10, 64)
		if err != nil {
			return nil, err
		}
		d.TransmitPackets = mn

		mn, err = strconv.ParseUint(ds[proc_net_dev_idx_transmit_errs], 10, 64)
		if err != nil {
			return nil, err
		}
		d.TransmitErrs = mn

		mn, err = strconv.ParseUint(ds[proc_net_dev_idx_transmit_drop], 10, 64)
		if err != nil {
			return nil, err
		}
		d.TransmitDrop = mn

		mn, err = strconv.ParseUint(ds[proc_net_dev_idx_transmit_fifo], 10, 64)
		if err != nil {
			return nil, err
		}
		d.TransmitFifo = mn

		mn, err = strconv.ParseUint(ds[proc_net_dev_idx_transmit_colls], 10, 64)
		if err != nil {
			return nil, err
		}
		d.TransmitColls = mn

		mn, err = strconv.ParseUint(ds[proc_net_dev_idx_transmit_carrier], 10, 64)
		if err != nil {
			return nil, err
		}
		d.TransmitCarrier = mn

		dss = append(dss, d)
	}

	return dss, nil
}
