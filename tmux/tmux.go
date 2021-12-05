package tmux

import (
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
	bin := t.getBinary()
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

func (t *Tmux) Attach(name string) {
	args := []string{}
	args = append(args, "-u", "attach-session", "-t", name)

	err := syscall.Exec(t.getBinary(), args, os.Environ())
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func (t *Tmux) ListSessions() []string {
	sessions := []string{}
	result, err := t.Exec("ls", "-F", "#{session_name}")
	if err != nil {
		// No active sessions returns as an error.
		return sessions
	}

	lines := strings.Trim(string(result), "\n")
	return strings.Split(lines, "\n")
}

func (t *Tmux) getBinary() string {
	if t.BinPath != "" {
		return t.BinPath
	}

	tmux, err := exec.LookPath("tmux")
	if err != nil {
		fmt.Println("Error:", err)
	}

	return tmux
}
