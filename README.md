# CasaOS Native Panel

A simple, debian-first management panel based on CasaOS concepts but designed for native Debian 12 systems with optional Docker support.

## Goals
- Keep a clean and simple admin UI
- Support native Debian 12 x86_64 installation without Docker
- Support Docker mode as an optional compatibility layer
- Manage services, users, groups, packages, networks, Samba, and logs
- Keep the experience easy for everyday use and powerful for advanced users

## Project status
- Architecture planned
- Repository scaffold created
- Next: backend runtime abstraction and native system adapter

## High-level architecture

- Frontend: Vue-based UI that exposes a simple management dashboard
- Backend: Go API service
- Runtime abstraction layer: adapter for Docker and native system
- System operations:
  - systemd service control
  - apt package management
  - user/group management
  - Samba config editing
  - log viewing

## Default operating modes
- `system`: native Debian mode, no Docker required
- `docker`: CasaOS-style Docker container management mode

## Repository structure

```text
.
├── README.md
├── Makefile
├── .github/
│   └── workflows/
├── cmd/
│   └── server/
├── internal/
│   ├── api/
│   ├── config/
│   ├── runtime/
│   ├── system/
│   ├── samba/
│   ├── users/
│   └── app/
├── web/
│   └── ui/
├── scripts/
│   ├── install.sh
│   └── install-system.sh
├── docs/
│   └── architecture.md
└── go.mod
```

## Development workflow

```bash
# install Go dependencies
make deps

# build backend
make build

# run backend in development mode
make run
```

## License
MIT
