package tmux

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

type Project struct {
	Name           string   `yaml:"name"`
	Root           string   `yaml:"root"`
	OnProjectStart []string `yaml:"on_project_start"`
	Windows        []Window `yaml:"windows"`
}

func StartProject(name string) {
	p, err := LoadProject(name)
	if err != nil {
		fmt.Println("Unable to load project:", err)
		os.Exit(1)
	}

	// Run startup commands
	if len(p.OnProjectStart) > 0 {
		for _, command := range p.OnProjectStart {
			args := strings.Fields(command)
			cmd := exec.Command(args[0], args[1:]...)
			cmd.Dir = p.GetRoot()
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println("Unable to run start command:", err, output)
				os.Exit(1)
			}
			fmt.Println(string(output))
		}
	}

	tmux := CreateTmux(false)

	if !sessionExists(name) {
		tmux.Run("new-session", "-d", "-s", name, "-n", p.Windows[0].Name, "-c", p.Windows[0].Root)

		for index, window := range p.Windows {
			if index > 0 {
				window.Create(tmux)
			}

			window.SendCommands(tmux)
		}

		p.Windows[0].Focus(tmux)
	}

	tmux.Attach(name)
}

// LoadProject loads and parses the config for the given project.
func LoadProject(name string) (*Project, error) {
	project := &Project{}

	home, _ := homedir.Dir()
	fileName := path.Join(home, ".workon", name+".yml")

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return project, err
	}

	err = yaml.Unmarshal(data, project)
	if len(project.Windows) < 1 {
		return project, errors.New("No windows defined")
	}

	rootPath := project.GetRoot()
	for index, window := range project.Windows {
		project.Windows[index].ID = fmt.Sprintf("%s:%d", project.Name, index)
		project.Windows[index].Root = filepath.Join(rootPath, window.Root)
	}

	return project, err
}

func (p *Project) GetRoot() string {
	rootPath, err := homedir.Expand(p.Root)
	if err != nil {
		fmt.Println("Unable to find root path")
	}
	return rootPath
}

func sessionExists(name string) bool {
	t := Tmux{}
	result, err := t.Exec("ls")
	if err != nil {
		return false
	}

	re := regexp.MustCompile(fmt.Sprintf("^%s:", name))
	return re.MatchString(string(result))
}
