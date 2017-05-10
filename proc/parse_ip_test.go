package proc

import "testing"

func TestParseIP(t *testing.T) {
	ipv6, port, err := parseLittleEndianIpv6("4506012640B600C10C1136C5C1EB0C75:B0BA")
	if err != nil {
		t.Fatal(err)
	}
	if ipv6 != "750C:EBC1:C536:110C:C100:B640:2601:0645" {
		t.Fatalf("ipv6 expected '750C:EBC1:C536:110C:C100:B640:2601:0645', got %s", ipv6)
	}
	if port != 45242 {
		t.Fatalf("port expected '45242', got %d", port)
	}

	ipv4, port, err := parseLittleEndianIpv4("0101007F:0035")
	if err != nil {
		t.Fatal(err)
	}
	if ipv4 != "127.0.1.1" {
		t.Fatalf("ipv4 expected '127.0.1.1', got %s", ipv4)
	}
	if port != 53 {
		t.Fatalf("port expected '53', got %d", port)
	}
}
