package histutil

import (
	"testing"

	"src.elv.sh/pkg/store"
)

func TestMemStore_Cursor(t *testing.T) {
	s := NewMemStore("+ 0", "- 1", "+ 2")
	testCursorIteration(t, s.Cursor("+"), []store.Cmd{
		{Text: "+ 0", Seq: 0},
		{Text: "+ 2", Seq: 2},
	})
}

// Remaining methods tested along with HybridStore
