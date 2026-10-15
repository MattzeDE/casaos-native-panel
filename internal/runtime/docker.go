package runtime

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// DockerRuntime is a compatibility layer for Docker mode.
type DockerRuntime struct{}

func NewDockerRuntime() Runtime {
	return &DockerRuntime{}
}

func (d *DockerRuntime) Name() string { return "docker" }

func (d *DockerRuntime) GetOverview() (Overview, error) {
	return Overview{Mode: "docker", OS: detectOSName()}, nil
}

func (d *DockerRuntime) GetServices() ([]Service, error) { return nil, nil }
func (d *DockerRuntime) StartService(name string) error { return nil }
func (d *DockerRuntime) StopService(name string) error  { return nil }
func (d *DockerRuntime) RestartService(name string) error { return nil }
func (d *DockerRuntime) GetPackages() ([]Package, error) { return nil, nil }
func (d *DockerRuntime) InstallPackage(name string) error { return fmt.Errorf("docker mode does not install apt packages directly; use container image workflow") }
func (d *DockerRuntime) RemovePackage(name string) error  { return fmt.Errorf("docker mode does not manage apt packages directly") }
func (d *DockerRuntime) GetUsers() ([]User, error) { return nil, nil }
func (d *DockerRuntime) CreateUser(u User) error     { return fmt.Errorf("user management is handled by the system runtime in docker mode") }
func (d *DockerRuntime) GetGroups() ([]Group, error) { return nil, nil }
func (d *DockerRuntime) GetSambaConfig() (string, error) { return "# Docker mode has no direct Samba config in this compatibility layer.\n", nil }
func (d *DockerRuntime) SaveSambaConfig(content string) error { return fmt.Errorf("docker mode does not persist Samba config directly") }

func detectOSName() string {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil { return "unknown" }
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			return strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), "\"")
		}
	}
	return "unknown"
}

func shellOut(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil { return "", err }
	return strings.TrimSpace(string(out)), nil
}

func (d *DockerRuntime) GetMetrics() (Metrics, error) {
	return Metrics{}, nil
}
