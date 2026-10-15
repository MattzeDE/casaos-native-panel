package runtime

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	sys "github.com/lukelinux79/casaos-native-panel/internal/system"
)

// SystemRuntime implements native Debian system operations.
type SystemRuntime struct{}

func NewSystemRuntime() Runtime {
	return &SystemRuntime{}
}

func (s *SystemRuntime) Name() string { return "system" }

func (s *SystemRuntime) GetOverview() (Overview, error) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	kernel, _ := shellOut("uname", "-r")
	osName := detectOSName()
	debianInfo, _ := shellOut("cat", "/etc/debian_version")

	return Overview{
		Hostname:   hostname,
		OS:         osName,
		Mode:       "system",
		Kernel:     kernel,
		DebianInfo: debianInfo,
	}, nil
}

func (s *SystemRuntime) GetServices() ([]Service, error) {
	out, err := shellOut("systemctl", "list-unit-files", "--type=service", "--no-pager", "--plain")
	if err != nil {
		return []Service{
			{Name: "ssh", State: "active", Enabled: true},
			{Name: "smbd", State: "inactive", Enabled: false},
		}, nil
	}

	lines := strings.Split(out, "\n")
	result := make([]Service, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.Contains(trimmed, "UNIT FILE") || strings.Contains(trimmed, "STATE") {
			continue
		}
		parts := strings.Fields(trimmed)
		if len(parts) < 2 {
			continue
		}
		name := parts[0]
		if !strings.HasSuffix(name, ".service") {
			continue
		}
		result = append(result, Service{
			Name:    strings.TrimSuffix(name, ".service"),
			State:   "unknown",
			Enabled: strings.Contains(strings.ToLower(trimmed), "enabled"),
		})
	}
	if len(result) == 0 {
		return []Service{{Name: "ssh", State: "active", Enabled: true}}, nil
	}
	return result, nil
}

func (s *SystemRuntime) StartService(name string) error {
	return exec.Command("systemctl", "start", name).Run()
}

func (s *SystemRuntime) StopService(name string) error {
	return exec.Command("systemctl", "stop", name).Run()
}

func (s *SystemRuntime) RestartService(name string) error {
	return exec.Command("systemctl", "restart", name).Run()
}

func (s *SystemRuntime) GetPackages() ([]Package, error) {
	out, err := shellOut("dpkg-query", "-W", "-f=${Package}\t${Version}\t${Status}\n")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(out, "\n")
	result := make([]Package, 0, len(lines))
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		result = append(result, Package{
			Name:    fields[0],
			Version: fields[1],
			State:   fields[2],
		})
	}
	return result, nil
}

func (s *SystemRuntime) InstallPackage(name string) error {
	cmd := exec.Command("apt-get", "update")
	if err := cmd.Run(); err != nil {
		return err
	}
	return exec.Command("apt-get", "install", "-y", name).Run()
}

func (s *SystemRuntime) RemovePackage(name string) error {
	return exec.Command("apt-get", "remove", "-y", name).Run()
}

func (s *SystemRuntime) GetUsers() ([]User, error) {
	out, err := shellOut("getent", "passwd")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(out, "\n")
	result := make([]User, 0, len(lines))
	for _, line := range lines {
		fields := strings.Split(line, ":")
		if len(fields) < 7 {
			continue
		}
		result = append(result, User{
			Username: fields[0],
			HomeDir:  fields[5],
			Shell:    fields[6],
		})
	}
	return result, nil
}

func (s *SystemRuntime) CreateUser(u User) error {
	if u.Username == "" {
		return fmt.Errorf("username is required")
	}
	args := []string{"-m", "-s", u.Shell}
	if u.Shell == "" {
		args = []string{"-m"}
	}
	return exec.Command("useradd", append(args, u.Username)...).Run()
}

func (s *SystemRuntime) GetGroups() ([]Group, error) {
	out, err := shellOut("getent", "group")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(out, "\n")
	result := make([]Group, 0, len(lines))
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) < 3 {
			continue
		}
		members := []string{}
		if parts[3] != "" {
			members = strings.Split(parts[3], ",")
		}
		result = append(result, Group{Name: parts[0], Members: members})
	}
	return result, nil
}

func (s *SystemRuntime) GetSambaConfig() (string, error) {
	content, err := os.ReadFile("/etc/samba/smb.conf")
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (s *SystemRuntime) SaveSambaConfig(content string) error {
	if _, err := os.Stat("/etc/samba"); err != nil {
		if mkdirErr := os.MkdirAll("/etc/samba", 0o755); mkdirErr != nil {
			return mkdirErr
		}
	}

	if err := os.WriteFile("/etc/samba/smb.conf", []byte(content), 0o644); err != nil {
		return err
	}
	return nil
}

func readLines(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

// GetMetrics collects system metrics and returns them as runtime.Metrics
func (s *SystemRuntime) GetMetrics() (Metrics, error) {
	m, err := sys.GetMetrics()
	if err != nil {
		return Metrics{}, err
	}
	return Metrics{
		Load1:          m.Load1,
		Load5:          m.Load5,
		Load15:         m.Load15,
		MemTotalKB:     m.MemTotalKB,
		MemAvailableKB: m.MemAvailableKB,
		DiskTotalBytes: m.DiskTotalBytes,
		DiskFreeBytes:  m.DiskFreeBytes,
		DiskUsedPct:    m.DiskUsedPct,
		NetRxBytes:     m.NetRxBytes,
		NetTxBytes:     m.NetTxBytes,
		UptimeSeconds:  m.UptimeSeconds,
	}, nil
}
