package psn

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

// ListPIDs reads all PIDs in '/proc'.
func ListPIDs() ([]int64, error) {
	ds, err := ioutil.ReadDir("/proc")
	if err != nil {
		return nil, err
	}

	pids := make([]int64, 0, len(ds))
	for _, f := range ds {
		if f.IsDir() && isInt(f.Name()) {
			id, err := strconv.ParseInt(f.Name(), 10, 64)
			if err != nil {
				return nil, err
			}
			pids = append(pids, id)
		}
	}
	return pids, nil
}

// ListProcFds reads '/proc/*/fd/*' to grab process IDs.
func ListProcFds() ([]string, error) {
	// returns the names of all files matching pattern
	// or nil if there is no matching file
	fs, err := filepath.Glob("/proc/[0-9]*/fd/[0-9]*")
	if err != nil {
		return nil, err
	}
	return fs, nil
}

func pidFromFd(s string) (int64, error) {
	// get 5261 from '/proc/5261/fd/69'
	return strconv.ParseInt(strings.Split(s, "/")[2], 10, 64)
}

// GetProgram returns the program name.
func GetProgram(pid int64) (string, error) {
	// Readlink needs root permission
	// return os.Readlink(fmt.Sprintf("/proc/%d/exe", pid))

	rs, err := rawProcStatus(pid)
	if err != nil {
		return "", err
	}
	return rs.Name, nil
}

// GetDevice returns the device name where dir is mounted.
// It parses '/etc/mtab'.
func GetDevice(mounted string) (string, error) {
	f, err := openToRead("/etc/mtab")
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		if len(txt) == 0 {
			continue
		}

		fields := strings.Fields(txt)
		if len(fields) < 2 {
			continue
		}

		dev := strings.TrimSpace(fields[0])
		at := strings.TrimSpace(fields[1])
		if mounted == at {
			return dev, nil
		}
	}

	return "", fmt.Errorf("no device found, mounted at %q", mounted)
}

var errNoDefaultRoute = fmt.Errorf("could not find default route")

func getDefaultRoute() (*syscall.NetlinkMessage, error) {
	dat, err := syscall.NetlinkRIB(syscall.RTM_GETROUTE, syscall.AF_UNSPEC)
	if err != nil {
		return nil, err
	}

	msgs, msgErr := syscall.ParseNetlinkMessage(dat)
	if msgErr != nil {
		return nil, msgErr
	}

	rtmsg := syscall.RtMsg{}
	for _, m := range msgs {
		if m.Header.Type != syscall.RTM_NEWROUTE {
			continue
		}
		buf := bytes.NewBuffer(m.Data[:syscall.SizeofRtMsg])
		if rerr := binary.Read(buf, binary.LittleEndian, &rtmsg); rerr != nil {
			continue
		}
		if rtmsg.Dst_len == 0 {
			// zero-length Dst_len implies default route
			return &m, nil
		}
	}

	return nil, errNoDefaultRoute
}

func getIface(idx uint32) (*syscall.NetlinkMessage, error) {
	dat, err := syscall.NetlinkRIB(syscall.RTM_GETADDR, syscall.AF_UNSPEC)
	if err != nil {
		return nil, err
	}

	msgs, msgErr := syscall.ParseNetlinkMessage(dat)
	if msgErr != nil {
		return nil, msgErr
	}

	ifaddrmsg := syscall.IfAddrmsg{}
	for _, m := range msgs {
		if m.Header.Type != syscall.RTM_NEWADDR {
			continue
		}
		buf := bytes.NewBuffer(m.Data[:syscall.SizeofIfAddrmsg])
		if rerr := binary.Read(buf, binary.LittleEndian, &ifaddrmsg); rerr != nil {
			continue
		}
		if ifaddrmsg.Index == idx {
			return &m, nil
		}
	}

	return nil, errNoDefaultRoute
}

var errNoDefaultInterface = fmt.Errorf("could not find default interface")

// GetDefaultInterface returns the default network interface
// (copied from https://github.com/coreos/etcd/blob/master/pkg/netutil/routes_linux.go).
func GetDefaultInterface() (string, error) {
	rmsg, rerr := getDefaultRoute()
	if rerr != nil {
		return "", rerr
	}

	_, oif, err := parsePREFSRC(rmsg)
	if err != nil {
		return "", err
	}

	ifmsg, ierr := getIface(oif)
	if ierr != nil {
		return "", ierr
	}

	attrs, aerr := syscall.ParseNetlinkRouteAttr(ifmsg)
	if aerr != nil {
		return "", aerr
	}

	for _, attr := range attrs {
		if attr.Attr.Type == syscall.IFLA_IFNAME {
			return string(attr.Value[:len(attr.Value)-1]), nil
		}
	}
	return "", errNoDefaultInterface
}

// parsePREFSRC returns preferred source address and output interface index (RTA_OIF).
func parsePREFSRC(m *syscall.NetlinkMessage) (host string, oif uint32, err error) {
	var attrs []syscall.NetlinkRouteAttr
	attrs, err = syscall.ParseNetlinkRouteAttr(m)
	if err != nil {
		return "", 0, err
	}

	for _, attr := range attrs {
		if attr.Attr.Type == syscall.RTA_PREFSRC {
			host = net.IP(attr.Value).String()
		}
		if attr.Attr.Type == syscall.RTA_OIF {
			oif = binary.LittleEndian.Uint32(attr.Value)
		}
		if host != "" && oif != uint32(0) {
			break
		}
	}

	if oif == 0 {
		err = errNoDefaultRoute
	}
	return
}
