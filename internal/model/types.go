package model

// Shared data model used across runtime and system packages.

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

type Metrics struct {
	Load1          float64 `json:"load1"`
	Load5          float64 `json:"load5"`
	Load15         float64 `json:"load15"`
	MemTotalKB     int64   `json:"mem_total_kb"`
	MemAvailableKB int64   `json:"mem_available_kb"`
	DiskTotalBytes uint64  `json:"disk_total_bytes"`
	DiskFreeBytes  uint64  `json:"disk_free_bytes"`
	DiskUsedPct    float64 `json:"disk_used_percent"`
	NetRxBytes     uint64  `json:"net_rx_bytes"`
	NetTxBytes     uint64  `json:"net_tx_bytes"`
	UptimeSeconds  float64 `json:"uptime_seconds"`
}
