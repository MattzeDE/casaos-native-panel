package api

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"
)

// This is a tiny static UI serving layer for the dashboard and eventually the SPA.
func initUIShim() http.Handler {
	mux := http.NewServeMux()
	base := "/"
	mux.HandleFunc(base, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(`<!doctype html>
<html><head><meta charset="utf-8"><title>CasaOS Native Panel</title>
<style>
 body { font-family: system-ui, sans-serif; margin: 2rem; background: #0b1020; color: #e5e7eb; }
 .card { background: #111827; border: 1px solid #374151; border-radius: 12px; padding: 1rem 1.2rem; margin: 1rem 0; }
 a { color: #93c5fd; }
</style></head>
<body>
 <h1>CasaOS Native Panel</h1>
 <div class="card"><strong>API:</strong> <a href="/api/v1/health">/api/v1/health</a> | <a href="/api/v1/dashboard">/api/v1/dashboard</a> | <a href="/api/v1/samba/config">/api/v1/samba/config</a></div>
 <div class="card"><strong>Mode:</strong> system</div>
 <div class="card"><strong>Notes:</strong> This is the first UI shim. The actual dashboard and admin views will be added next.</div>
</body></html>`))
	})
	return mux
}

func registerUIShim(mux *http.ServeMux) {
	root := http.FileServer(http.Dir(filepath.Clean("web")))
	mux.Handle("/", root)
	mux.HandleFunc("/api/v1/ui", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(`<!doctype html><html><body><h1>UI Shim</h1><p>frontend placeholder</p></body></html>`))
	})
}

func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	if !strings.EqualFold(r.Method, http.MethodGet) {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	metrics, err := systemGetMetrics()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(response{"metrics": metrics})
}

func systemGetMetrics() (map[string]any, error) {
	return map[string]any{
		"load": map[string]any{"1m": 0.0, "5m": 0.0, "15m": 0.0},
		"memory": map[string]any{"total_kb": 0, "available_kb": 0},
		"disk": map[string]any{"total_bytes": 0, "free_bytes": 0, "used_percent": 0.0},
		"network": map[string]any{"rx_bytes": 0, "tx_bytes": 0},
		"uptime_seconds": 0.0,
	}, nil
}
