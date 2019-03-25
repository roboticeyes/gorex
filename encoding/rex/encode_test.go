package rex

import (
	"bytes"
	"os"
	"testing"
)

func TestEncodingHeader(t *testing.T) {

	rexFile := File{
		Header: *CreateHeader(),
	}
	var buf bytes.Buffer
	e := NewEncoder(&buf)
	n, err := e.Encode(rexFile)
	if err != nil {
		t.Fatalf("Error during encoding %v", err)
	}
	if n != 86 {
		t.Fatalf("Header size does not match")
	}

	f, _ := os.Create("test.rex")
	f.Write(buf.Bytes())
	defer f.Close()
}
