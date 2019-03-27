package rex

import (
	"bufio"
	"fmt"
	"os"
	"testing"
	// "github.com/go-gl/mathgl/mgl32"
)

func TestDecodingHeader(t *testing.T) {

	file, err := os.Open("mesh.rex")
	defer file.Close()
	if err != nil {
		t.Fatalf("Cannot read testfile: %v", err)
	}

	r := bufio.NewReader(file)

	d := NewDecoder(r)
	header, rex, err := d.Decode()

	if err != nil {
		t.Fatalf("TEST ERROR: %v", err)
	}
	if rex == nil {
		fmt.Println("rex file is nil")
	}

	fmt.Println(header)
}
