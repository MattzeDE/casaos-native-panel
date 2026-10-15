package runtime

import (
	"fmt"
	model "github.com/lukelinux79/casaos-native-panel/internal/model"
)

// Expose model types under the runtime package as aliases to avoid duplication.
// This lets other runtime files use Overview, Service, etc. while keeping
// the concrete definitions in internal/model.

type Overview = model.Overview
type Service = model.Service
type Package = model.Package
type User = model.User
type Group = model.Group
type Metrics = model.Metrics

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
	GetMetrics() (Metrics, error)
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
