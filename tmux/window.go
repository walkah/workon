package tmux

import (
	"strings"
)

type Window struct {
	Name     string   `yaml:"name"`
	Root     string   `yaml:"root,omitempty"`
	Commands []string `yaml:"commands,omitempty"`
	ID       string   `yaml:"-"`
}

func (w *Window) Create(tmux *Tmux) {
	tmux.Run("new-window", "-t", w.ID, "-n", w.Name, "-c", w.Root)
}

func (w *Window) SendCommands(tmux *Tmux) {
	if len(w.Commands) > 0 {
		tmux.Run("send-keys", "-t", w.ID, strings.Join(w.Commands, ";"))
		tmux.Run("send-keys", "-t", w.ID, "Enter")
	}
}

func (w *Window) Focus(tmux *Tmux) {
	tmux.Run("select-window", "-t", w.ID)
}
