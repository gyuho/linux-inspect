package ss

import (
	"fmt"
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

func TestListProcess(t *testing.T) {
	ps4, err := ListProcess(TCP)
	if err != nil {
		t.Error(err)
	}
	for _, p := range ps4 {
		fmt.Println(p)
	}
	ps6, err := ListProcess(TCP6)
	if err != nil {
		t.Error(err)
	}
	for _, p := range ps6 {
		fmt.Println(p)
	}
}
