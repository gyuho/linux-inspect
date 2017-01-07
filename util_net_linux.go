package psn

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"syscall"
)

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
