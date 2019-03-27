package rex

import (
	"testing"
)

func TestHeader(t *testing.T) {
	h := CreateHeader()

	if h.Version != 1 {
		t.Error("Wrong REX version")
	}
}
