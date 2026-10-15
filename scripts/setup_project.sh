#!/usr/bin/env bash
set -euo pipefail

echo "[casaos-native-panel] Creating project skeleton..."
mkdir -p cmd/server
mkdir -p internal/api
mkdir -p internal/config
mkdir -p internal/runtime
mkdir -p internal/system
mkdir -p internal/samba
mkdir -p internal/users
mkdir -p internal/app
mkdir -p web/ui
mkdir -p scripts
mkdir -p docs

cat > Makefile <<'EOF'
.PHONY: deps build run test

APP_NAME=casaos-native-panel

all: build

deps:
	go mod download

build:
	go build -o bin/$(APP_NAME) ./cmd/server

run:
	go run ./cmd/server

test:
	go test ./...
EOF

cat > cmd/server/main.go <<'EOF'
package main

import "fmt"

func main() {
	fmt.Println("CasaOS Native Panel backend is starting...")
}
EOF

cat > internal/runtime/runtime.go <<'EOF'
package runtime

// Runtime is the abstraction layer for service execution.
// Docker mode and native system mode can both satisfy this interface.
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
	Hostname string
	OS       string
	Mode     string
}

type Service struct {
	Name    string
	Status  string
	Enabled bool
}

type Package struct {
	Name    string
	Version string
	State   string
}

type User struct {
	Username string
	HomeDir  string
	Shell    string
}

type Group struct {
	Name string
}
EOF

cat > internal/system/system.go <<'EOF'
package system

import "github.com/lukelinux79/casaos-native-panel/internal/runtime"

// SystemRuntime implements the native Debian system mode.
type SystemRuntime struct{}

func (s *SystemRuntime) Name() string { return "system" }

func (s *SystemRuntime) GetOverview() (runtime.Overview, error) {
	return runtime.Overview{Mode: "system"}, nil
}

func (s *SystemRuntime) GetServices() ([]runtime.Service, error) {
	return nil, nil
}

func (s *SystemRuntime) StartService(name string) error { return nil }
func (s *SystemRuntime) StopService(name string) error  { return nil }
func (s *SystemRuntime) RestartService(name string) error { return nil }

func (s *SystemRuntime) GetPackages() ([]runtime.Package, error) { return nil, nil }
func (s *SystemRuntime) InstallPackage(name string) error        { return nil }
func (s *SystemRuntime) RemovePackage(name string) error         { return nil }

func (s *SystemRuntime) GetUsers() ([]runtime.User, error) { return nil, nil }
func (s *SystemRuntime) CreateUser(u runtime.User) error  { return nil }

func (s *SystemRuntime) GetGroups() ([]runtime.Group, error) { return nil, nil }

func (s *SystemRuntime) GetSambaConfig() (string, error) { return "", nil }
func (s *SystemRuntime) SaveSambaConfig(content string) error { return nil }
EOF

cat > docs/architecture.md <<'EOF'
# Architecture

## Overview
This project aims to keep the CasaOS-friendly management experience while supporting native Debian 12 operation without Docker.

## Modes
- `system` mode: manage Debian services, apt packages, Samba, users, and logs directly
- `docker` mode: keep compatibility with container-based CasaOS workflows

## Runtime abstraction
The core of the project is a runtime interface that both Docker and system mode implement.

This lets the UI and API remain stable while backend implementation changes depending on the desired mode.

## High-level flow
1. User opens the panel
2. Frontend requests data from the backend API
3. Backend calls runtime adapter
4. Adapter performs system or Docker-specific work
5. Response is normalized back to the UI

## Native system operations
The system adapter should support:
- systemd service control
- apt package management
- users and groups
- Samba configuration editing with validation
- logs via journalctl

## Security assumptions
- Root privileges may be needed for system modifications
- All system-critical writes should use backup + validation
- TLS and web access protection should be configured in deployment
EOF

cat > scripts/install-system.sh <<'EOF'
#!/usr/bin/env bash
set -euo pipefail

echo "[install-system] Native Debian system mode setup"
echo "This script is a placeholder for the Debian 12 native installation flow."
echo "In a later iteration it will install the backend, configure systemd, and set up the UI."
EOF

chmod +x scripts/install-system.sh

cat > .gitignore <<'EOF'
/bin/
/.idea/
/.vscode/
*.log
.env
EOF

cat > .github/workflows/ci.yml <<'EOF'
name: CI

on:
  push:
    branches: [ main ]
  pull_request:

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - run: go mod download
      - run: go build ./cmd/server
EOF

echo "Project scaffold created successfully."
