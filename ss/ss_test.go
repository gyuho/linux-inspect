package ss

import (
	"fmt"
	"os"
	"testing"
)

func TestParseLittleEndianIpv4(t *testing.T) {
	var tests = map[string]string{
		"0101007F:0035": "127.0.1.1:53",
		"0100007F:0277": "127.0.0.1:631",
		"0100007F:049A": "127.0.0.1:1178",
		"0100007F:049B": "127.0.0.1:1179",
		"0100007F:049C": "127.0.0.1:1180",
		"0100007F:04FE": "127.0.0.1:1278",
		"0100007F:04FF": "127.0.0.1:1279",
		"0100007F:0500": "127.0.0.1:1280",
		"0100007F:0562": "127.0.0.1:1378",
		"0100007F:0563": "127.0.0.1:1379",
		"0100007F:0564": "127.0.0.1:1380",
		"0100007F:981F": "127.0.0.1:38943",
		"0100007F:B02D": "127.0.0.1:45101",
	}
	for k, v := range tests {
		host, port, err := parseLittleEndianIpv4(k)
		if err != nil {
			t.Error(err)
		}
		addr := host + port
		if addr != v {
			t.Errorf("got = %s, want = %s", addr, v)
		}
	}
}

func TestParseLittleEndianIpv6(t *testing.T) {
	var tests = map[string]string{
		"4506012691A700C165EB1DE1F912918C:4B72": "[8C91:12F9:E11D:EB65:C100:A791:2601:0645]:19314",
		"4506012691A700C165EB1DE1F912918C:F9D7": "[8C91:12F9:E11D:EB65:C100:A791:2601:0645]:63959",
		"4506012691A700C165EB1DE1F912918C:6251": "[8C91:12F9:E11D:EB65:C100:A791:2601:0645]:25169",
		"4506012691A700C165EB1DE1F912918C:BA85": "[8C91:12F9:E11D:EB65:C100:A791:2601:0645]:47749",
		"4506012691A700C165EB1DE1F912918C:4B82": "[8C91:12F9:E11D:EB65:C100:A791:2601:0645]:19330",
	}
	for k, v := range tests {
		host, port, err := parseLittleEndianIpv6(k)
		if err != nil {
			t.Error(err)
		}
		addr := "[" + host + "]" + port
		if addr != v {
			t.Errorf("got = %s, want = %s", addr, v)
		}
	}
}

func TestReadProcNetInternal(t *testing.T) {
	if _, err := readProcNet(TCP); err != nil {
		t.Error(err)
	}
	if _, err := readProcNet(TCP6); err != nil {
		t.Error(err)
	}
}

func TestReadProcFdInternal(t *testing.T) {
	if _, err := readProcFd(); err != nil {
		t.Error(err)
	}
}

func TestList(t *testing.T) {
	if _, err := List(nil); err != nil {
		t.Error(err)
	}
	if _, err := List(nil, TCP, TCP6); err != nil {
		t.Error(err)
	}
}

func TestWriteToTable(t *testing.T) {
	ps, err := List(nil, TCP, TCP6)
	if err != nil {
		t.Error(err)
	}
	WriteToTable(os.Stdout, ps...)
}

func TestFilterMatch(t *testing.T) {
	filter := Process{}
	filter.Program = "etcd"
	filter.PID = 1000
	p := Process{}
	p.Program = "etcd"
	p.PID = 1000
	if !filter.Match(p) {
		t.Errorf("got = false, want = true for %s", p)
	}
}

func TestListEtcd(t *testing.T) {
	filter := &Process{Program: "etcd"}
	fmt.Println("etcd filter:", filter)

	ps, err := List(filter, TCP, TCP6)
	if err != nil {
		t.Error(err)
	}
	WriteToTable(os.Stdout, ps...)

	pm := ListPorts(filter, TCP, TCP6)
	for pt := range pm {
		fmt.Printf("%9s is being used...\n", pt)
	}
}

func TestListPorts(t *testing.T) {
	pm := ListPorts(nil, TCP, TCP6)
	for pt := range pm {
		fmt.Printf("%9s is being used...\n", pt)
	}
}

func TestExistAddFree(t *testing.T) {
	ports := NewPorts()
	tp := ":1000"
	ports.Add(tp)
	if _, ok := ports.beingUsed[tp]; !ok {
		t.Errorf("%s should have been 'beingUsed'", tp)
	}
	if !ports.Exist(tp) {
		t.Errorf("%s should exist", tp)
	}
	ports.Free(tp)
	if _, ok := ports.beingUsed[tp]; ok {
		t.Errorf("%s should have been not 'beingUsed'", tp)
	}
}

func TestGetFreePort(t *testing.T) {
	pt, err := getFreePort(TCP, TCP6)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("Port %s is available for tcp.\n", pt)
}

func TestGetFreePorts(t *testing.T) {
	ps, err := GetFreePorts(3, TCP, TCP6)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("GetFreePorts: %v\n", ps)
}
