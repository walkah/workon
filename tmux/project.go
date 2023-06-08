package tmux

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

type Project struct {
	Name           string   `yaml:"name"`
	Root           string   `yaml:"root"`
	OnProjectStart []string `yaml:"on_project_start,omitempty"`
	OnProjectStop  []string `yaml:"on_project_stop,omitempty"`
	Windows        []Window `yaml:"windows"`
}

func StartProject(name string) {
	p, err := LoadProject(name)
	if err != nil {
		fmt.Println("Unable to load project:", err)
		os.Exit(1)
	}

	tmux := CreateTmux(false)

	if !sessionExists(name) {
		// Run startup commands
		p.RunCommands(p.OnProjectStart)

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

func StopProject(name string) {
	if !sessionExists(name) {
		return
	}

	t := Tmux{}
	t.KillSession(name)

	p, err := LoadProject(name)
	if err != nil {
		fmt.Println("Unable to load project:", err)
		os.Exit(1)
	}

	p.RunCommands(p.OnProjectStop)
}

func ListActiveProjects() ([]string, error) {
	activeProjects := []string{}

	projects, err := ListProjects()
	if err != nil {
		return activeProjects, err
	}

	for _, project := range projects {
		if sessionExists(project) {
			activeProjects = append(activeProjects, project)
		}
	}
	return activeProjects, nil
}

// ProjectList gets a list of
func ListProjects() ([]string, error) {
	configDir := getConfigDir()
	files, err := ioutil.ReadDir(configDir)
	if err != nil {
		return nil, err
	}

	projects := []string{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		ext := filepath.Ext(name)
		projects = append(projects, name[:len(name)-len(ext)])
	}

	return projects, nil
}

// LoadProject loads and parses the config for the given project.
func LoadProject(name string) (*Project, error) {
	project := &Project{}

	fileName := getConfigFilePath(name)

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return project, err
	}

	err = yaml.Unmarshal(data, project)
	if len(project.Windows) < 1 {
		return project, errors.New("no windows defined")
	}

	rootPath := project.GetRoot()
	for index, window := range project.Windows {
		project.Windows[index].ID = fmt.Sprintf("%s:%d", project.Name, index)
		project.Windows[index].Root = filepath.Join(rootPath, window.Root)
	}

	return project, err
}

func NewProject(name string) error {
	project := &Project{
		Name:           name,
		Root:           "~/",
		OnProjectStart: []string{""},
		Windows:        make([]Window, 3),
	}

	project.Windows[0] = Window{
		Name:     "shell",
		Commands: []string{""},
	}

	project.Windows[1] = Window{
		Name:     "server",
		Commands: []string{""},
	}

	project.Windows[2] = Window{
		Name:     "logs",
		Commands: []string{""},
	}

	project.Save()

	return EditProject(name)
}

func EditProject(name string) error {
	fileName := getConfigFilePath(name)

	_, err := os.Stat(fileName)
	if err != nil {
		return errors.New("config file does not exist")
	}

	editorName := os.Getenv("EDITOR")
	if editorName == "" {
		return errors.New("EDITOR variable not defined")
	}

	editor, err := exec.LookPath(editorName)
	if err != nil {
		return err
	}

	return syscall.Exec(editor, []string{editorName, fileName}, os.Environ())
}

func (p *Project) Save() error {
	fileName := getConfigFilePath(p.Name)

	_, err := os.Stat(fileName)
	if err == nil {
		return errors.New("config file already exists")
	}

	data, err := yaml.Marshal(p)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, data, 0644)
}

func (p *Project) GetRoot() string {
	rootPath, err := homedir.Expand(p.Root)
	if err != nil {
		fmt.Println("Unable to find root path")
	}
	return rootPath
}

func (p *Project) RunCommands(commands []string) {
	for _, command := range commands {
		if command == "" {
			continue
		}
		args := strings.Fields(command)
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = p.GetRoot()
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println("Unable to run command:", err)
			os.Exit(1)
		}
	}
}

func getConfigDir() string {
	home, _ := homedir.Dir()
	return path.Join(home, ".workon")
}

func getConfigFilePath(name string) string {
	return path.Join(getConfigDir(), name+".yml")
}

func sessionExists(name string) bool {
	t := Tmux{}

	sessions, err := t.ListSessions()
	if err != nil {
		return false
	}
	for _, s := range sessions {
		if s == name {
			return true
		}
	}
	return false
}
