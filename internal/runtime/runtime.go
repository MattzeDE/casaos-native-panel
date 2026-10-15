package runtime

import "fmt"

// Runtime is the core abstraction used by the API layer.
type Runtime interface {
	Name() string
	GetOverview() (Overview, error)
	GetServices() ([]Service, error)
	StartService(name string) error
	StopService(name string) error
	RestartService(name string) error
	GetPackages() ([]Package, error)
	InstallPackage(name string) error
	RemovePackage(name string) error
	GetUsers() ([]User, error)
	CreateUser(u User) error
	GetGroups() ([]Group, error)
	GetSambaConfig() (string, error)
	SaveSambaConfig(content string) error
}

type Overview struct {
	Hostname   string `json:"hostname"`
	OS         string `json:"os"`
	Mode       string `json:"mode"`
	Kernel     string `json:"kernel"`
	DebianInfo string `json:"debian_info,omitempty"`
}

type Service struct {
	Name    string `json:"name"`
	State   string `json:"state"`
	Enabled bool   `json:"enabled"`
}

type Package struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	State   string `json:"state"`
}

type User struct {
	Username string   `json:"username"`
	HomeDir  string   `json:"home_dir,omitempty"`
	Shell    string   `json:"shell,omitempty"`
	Groups   []string `json:"groups,omitempty"`
}

type Group struct {
	Name    string   `json:"name"`
	Members []string `json:"members,omitempty"`
}

func New(mode string) (Runtime, error) {
	switch mode {
	case "system":
		return NewSystemRuntime(), nil
	case "docker":
		return NewDockerRuntime(), nil
	default:
		return nil, fmt.Errorf("unsupported runtime mode %q", mode)
	}
}
