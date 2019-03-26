package rex

import (
	"testing"
)

func TestHeader(t *testing.T) {
	h := CreateHeader()

	if h.version != 1 {
		t.Error("Wrong REX version")
	}
}
