package tmux

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

type Tmux struct {
	BinPath string
	Debug   bool
}

func CreateTmux(debug bool) *Tmux {
	return &Tmux{Debug: debug}
}

func (t *Tmux) Exec(args ...string) ([]byte, error) {
	bin, err := t.getBinary()
	if err != nil {
		return []byte{}, err
	}
	if t.Debug {
		fmt.Println(bin, strings.Join(args, " "))
	}
	return exec.Command(bin, args...).CombinedOutput()
}

func (t *Tmux) Run(args ...string) {
	output, err := t.Exec(args...)
	if err != nil {
		fmt.Println(err, string(output))
	}
}

func (t *Tmux) Attach(name string) error {
	args := []string{}
	args = append(args, "-u", "attach-session", "-t", name)

	bin, err := t.getBinary()
	if err != nil {
		return err
	}
	err = syscall.Exec(bin, args, os.Environ())
	if err != nil {
		return err
	}
	return nil
}

func (t *Tmux) ListSessions() ([]string, error) {
	sessions := []string{}
	result, err := t.Exec("ls", "-F", "#{session_name}")
	if errors.Is(err, exec.ErrNotFound) {
		// No active sessions returns as an error.
		return sessions, err
	}

	lines := strings.Trim(string(result), "\n")
	fmt.Println("Sessions", lines)
	return strings.Split(lines, "\n"), nil
}

func (t *Tmux) KillSession(name string) error {
	_, err := t.Exec("kill-session", "-t", name)
	return err
}

func (t *Tmux) getBinary() (string, error) {
	if t.BinPath != "" {
		return t.BinPath, nil
	}

	tmux, err := exec.LookPath("tmux")
	if err != nil {
		return "", err
	}

	return tmux, nil
}
