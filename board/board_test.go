package board

import (
	"testing"
)

func Test_BoardToString(t *testing.T) {
	b := NewBoard()
	if len(b.ToString()) != 64 {
		t.Fail()
	}
}
