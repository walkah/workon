package tmux

import (
	"strings"
)

type Window struct {
	Name     string   `yaml:"name"`
	Root     string   `yaml:"root,omitempty"`
	Commands []string `yaml:"commands,omitempty"`
	ID       string   `yaml:"-"`
	Panes    []Pane   `yaml:"panes,omitempty"`
}

func (w *Window) Create(tmux *Tmux) {
	tmux.Run("new-window", "-t", w.ID, "-n", w.Name, "-c", w.Root)

	for i, pane := range w.Panes {
		if i > 0 {
			err := tmux.SplitWindow(w.ID, pane.Type, w.Root)
			if err != nil {
				panic(err)
			}
		}

		if len(pane.Commands) > 0 {
			tmux.SendKeys(w.ID, strings.Join(pane.Commands, ";"))
		}
	}
}

func (w *Window) SendCommands(tmux *Tmux) {
	if len(w.Commands) > 0 {
		tmux.SendKeys(w.ID, strings.Join(w.Commands, ";"))
	}
}

func (w *Window) Focus(tmux *Tmux) {
	tmux.Run("select-window", "-t", w.ID)
}
