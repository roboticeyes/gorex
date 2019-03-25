package rex

import (
	"bytes"
	"os"
	"testing"
)

func TestEncodingHeader(t *testing.T) {

	rexFile := File{}

	var buf bytes.Buffer
	e := NewEncoder(&buf)
	n, err := e.Encode(rexFile)
	if err != nil {
		t.Fatalf("Error during encoding %v", err)
	}
	if n != 86 {
		t.Fatalf("Header size does not match")
	}
}

func TestEncodingPointList(t *testing.T) {

	pl := PointList{}
	pl.Points = append(pl.Points, Point{0.0, 0.0, 0.0})
	pl.Points = append(pl.Points, Point{1.0, 1.0, 0.0})
	pl.Points = append(pl.Points, Point{0.0, 1.0, 1.0})
	pl.Points = append(pl.Points, Point{0.0, 1.0, 1.0})

	rexFile := File{}
	rexFile.PointList = append(rexFile.PointList, pl)

	var buf bytes.Buffer
	e := NewEncoder(&buf)
	_, err := e.Encode(rexFile)
	if err != nil {
		t.Fatalf("Error during encoding %v", err)
	}

	f, _ := os.Create("test.rex")
	f.Write(buf.Bytes())
	defer f.Close()
}
