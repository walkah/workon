package tmux

import (
	"testing"
)

func TestListSessions(t *testing.T) {
	m := CreateTmux(true)
	if _, err := m.ListSessions(); err != nil {
		t.Fatalf("ListSessions: %s", err)
	}
}
