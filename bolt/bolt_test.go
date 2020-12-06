package bolt

import (
	"testing"
)

func TestParsingEmptyTinyMap(t *testing.T) {
	msg := []byte{0xa0}
	m, n, err := ParseTinyMap(msg)
	if err != nil {
		t.Fatal(err)
	}
	if len(m) != 0 {
		t.Fatal("expected zero length map")
	}
	if n != 1 {
		t.Fatal("expected n=1, got", n)
	}
}

func TestParsingTinyInt(t *testing.T) {
	val, err := ParseTinyInt(0x0a)
	if err != nil {
		t.Fatal(err)
	}
	if val != 10 {
		t.Fatal("expected 10, got", val)
	}

	val, err = ParseTinyInt(0x69)
	if err != nil {
		t.Fatal(err)
	}
	if val != 105 {
		t.Fatal("expected 105, got", val)
	}

	val, err = ParseTinyInt(0x81)
	if err == nil {
		t.Fatal("expected to fail parsing, value is a tiny-string and not tiny-int!")
	}
}

func TestParsingInt(t *testing.T) {
	type test struct {
		buf          []byte
		expectedVal  int
		expectedSize int
	}

	tests := []test{
		test{[]byte{0xc8, 0x45}, 69, 2},
		test{[]byte{0xc8, 0xbb}, -69, 2},
		test{[]byte{0xc9, 0xfa, 0xc7}, -1337, 3},
		test{[]byte{0xc9, 0x14, 0x08}, 5128, 3},
		test{[]byte{0xca, 0x6b, 0x4b, 0xb4, 0x40}, 1800123456, 5},
		test{[]byte{0xca, 0xff, 0xfe, 0x1d, 0xc0}, -123456, 5},
		test{[]byte{0xcb, 0xff, 0xff, 0xa5, 0x0c, 0xef, 0x85, 0xc0, 0x01},
			-99999999999999, 9},
		test{[]byte{0xcb, 0x00, 0x00, 0x5a, 0xf3, 0x10, 0x7a, 0x3f, 0xff},
			99999999999999, 9},
	}

	for _, test := range tests {
		val, n, err := ParseInt(test.buf)
		if err != nil {
			t.Fatalf("failed test %#v: %s\n", test, err)
		}
		if val != test.expectedVal || n != test.expectedSize {
			t.Fatalf("expected (%d, %d), got (%d, %d)\n",
				test.expectedVal, test.expectedSize,
				val, n)
		}
	}
}

func TestParsingTinyString(t *testing.T) {
	val, n, err := ParseTinyString([]byte{0x87,
		// "address"
		0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
		// extra noise
		0xff, 0xff})
	if err != nil {
		t.Fatal(err)
	}

	if val != "address" {
		t.Fatal("expected 'address', got", val)
	}

	if n != len("address")+1 {
		t.Fatal("expected length of len('address')+1, got", n)
	}

	val, n, err = ParseTinyString([]byte{0x80})
	if err != nil {
		t.Fatal(err)
	}

	if val != "" {
		t.Fatal("expected '', got", val)
	}

	if n != 1 {
		t.Fatal("expected length of 1, got", n)
	}

	_, _, err = ParseTinyString([]byte{0xa1})
	if err == nil {
		t.Fatal("expected error for invalid tiny-string!")
	}
}

func TestParsingString(t *testing.T) {
	// 'neo4j-python/4.2.0 Python/3.8.6-final-0 (openbsd6)'
	msg := []byte{0xd0,
		0x32, 0x6e, 0x65, 0x6f, 0x34, 0x6a, 0x2d, 0x70,
		0x79, 0x74, 0x68, 0x6f, 0x6e, 0x2f, 0x34, 0x2e,
		0x32, 0x2e, 0x30, 0x20, 0x50, 0x79, 0x74, 0x68,
		0x6f, 0x6e, 0x2f, 0x33, 0x2e, 0x38, 0x2e, 0x36,
		0x2d, 0x66, 0x69, 0x6e, 0x61, 0x6c, 0x2d, 0x30,
		0x20, 0x28, 0x6f, 0x70, 0x65, 0x6e, 0x62, 0x73,
		0x64, 0x36, 0x29,
		// noise
		0xff, 0xff, 0x00}

	val, n, err := ParseString(msg)
	if err != nil {
		t.Fatal(err)
	}
	if val != "neo4j-python/4.2.0 Python/3.8.6-final-0 (openbsd6)" {
		t.Fatal("expected 'neo4j-python/4.2.0 Python/3.8.6-final-0 (openbsd6)', got", val)
	}
	if n != (2 + 0x32) {
		t.Fatal("expected 2 + 0x32 for length, got", n)
	}
}

func TestParsingTinymap(t *testing.T) {
	short := []byte{0xa1,
		// 'address'
		0x87, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
		// 'localhost:8888'
		0x8e, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x68, 0x6f,
		0x73, 0x74, 0x3a, 0x38, 0x38, 0x38, 0x38,
	}
	big := []byte{0xa5,
		// 'user_agent'
		0x8a, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x61, 0x67, 0x65, 0x6e, 0x74,
		// 'neo4j-python/4.2.0 Python/3.8.6-final-0 (openbsd6)'
		0xd0, 0x32, 0x6e, 0x65, 0x6f, 0x34, 0x6a, 0x2d, 0x70, 0x79, 0x74, 0x68, 0x6f, 0x6e, 0x2f, 0x34, 0x2e, 0x32, 0x2e, 0x30, 0x20, 0x50, 0x79, 0x74, 0x68, 0x6f, 0x6e, 0x2f, 0x33, 0x2e, 0x38, 0x2e, 0x36, 0x2d, 0x66, 0x69, 0x6e, 0x61, 0x6c, 0x2d, 0x30, 0x20, 0x28, 0x6f, 0x70, 0x65, 0x6e, 0x62, 0x73, 0x64, 0x36, 0x29, 0x87, 0x72, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67,
		// short, above
		0xa1, 0x87, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x8e, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x68, 0x6f, 0x73, 0x74, 0x3a, 0x38, 0x38, 0x38, 0x38,
		// 'scheme'
		0x86, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x65,
		// 'basic'
		0x85, 0x62, 0x61, 0x73, 0x69, 0x63,
		// 'principal'
		0x89, 0x70, 0x72, 0x69, 0x6e, 0x63, 0x69, 0x70, 0x61, 0x6c,
		// 'neo4j'
		0x85, 0x6e, 0x65, 0x6f, 0x34, 0x6a,
		// 'credentials'
		0x8b, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x73,
		// 'password'
		0x88, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64,
	}

	shortMap, n, err := ParseTinyMap(short)
	if err != nil {
		t.Fatalf("failed to parse short tinymap: %v", err)
	}
	if int(n) != len(short) {
		t.Fatalf("expected len of %d, got %d", len(short), n)
	}
	addrVal, found := shortMap["address"]
	if !found {
		t.Fatalf("failed to find 'address' in short tinymap")
	}
	addr, ok := addrVal.(string)
	if !ok {
		t.Fatalf("expected address to be a string, got %v", addr)
	}
	if addr != "localhost:8888" {
		t.Fatalf("expected address of 'localhost:8888', got %s", addr)
	}

	bigMap, _, err := ParseTinyMap(big)
	if err != nil {
		t.Fatalf("failed to parse big tinymap: %v", err)
	}
	agentVal, found := bigMap["user_agent"]
	if !found {
		t.Fatalf("failed ot find 'user_agent' in big tinymap")
	}
	agent, ok := agentVal.(string)
	if !ok {
		t.Fatalf("expected agent to be a string, got %v", agent)
	}
	if agent != "neo4j-python/4.2.0 Python/3.8.6-final-0 (openbsd6)" {
		t.Fatalf("got unexpected agent value: %s", agent)
	}
}

func TestTinyMapWithTinyArray(t *testing.T) {
	msg := []byte{0xa1,
		// 'fields'
		0x86, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73,
		// ['1'] (tiny array with a string of '1')
		0x91, 0x81, 0x31}
	m, _, err := ParseTinyMap(msg)
	if err != nil {
		t.Fatal(err)
	}

	val, found := m["fields"]
	if !found {
		t.Fatal("expected to find 'fields' in tinymap")
	}
	array, ok := val.([]interface{})
	if !ok {
		t.Fatal("expected value to be []interface{}")
	}
	first, ok := array[0].(string)
	if !ok {
		t.Fatal("expected array[0] to be a string")
	}
	if first != "1" {
		t.Fatal("expected array[0] to be '1', was", array[0])
	}
}

func TestTinyMapTinyInt(t *testing.T) {
	msg := []byte{0x0, 0x16,
		// success msg
		0xb1, 0x70,
		// tiny map with 2 entries
		0xa2,
		// "t_first"
		0x87, 0x74, 0x5f, 0x66, 0x69, 0x72, 0x73, 0x74,
		// 8
		0x08,
		// 'fields'
		0x86, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73,
		// '1'
		0x81, 0x31}

	m, _, err := ParseTinyMap(msg[4:])
	if err != nil {
		t.Fatal(err)
	}
	if len(m) != 2 {
		t.Fatal("expected to find 2 fields, got:", len(m))
	}
	val, found := m["t_first"]
	if !found {
		t.Fatal("expected to find a field called 't_first'")
	}
	i, ok := val.(int)
	if !ok {
		t.Fatal("expected value to be an int")
	}
	if i != 8 {
		t.Fatal("expected value to be 8, got:", val)
	}
}

func TestParsingSuccessMsg(t *testing.T) {
	msg := []byte{0xa3,
		0x87, 0x74, 0x5f, 0x66, 0x69, 0x72, 0x73, 0x74,
		0x4,
		0x86, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73,
		0x91,
		0x81, 0x78,
		0x83, 0x71, 0x69, 0x64,
		0x80}

	_, _, err := ParseTinyMap(msg)

	if err != nil {
		t.Fatal(err)
	}
	/*
		msg = []byte{0xa3, 0x87, 0x74, 0x5f, 0x66, 0x69, 0x72, 0x73, 0x74, 0x5, 0x86, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x91, 0x81, 0x78, 0x83, 0x71, 0x69, 0x64}

		val, _, err := ParseTinyMap(msg)
		if err != nil {
			t.Fatal(err)
		}
		if len(val) != 3 {
			t.Fatal("expected 3 map entries, saw", len(val))
		}
	*/
}
