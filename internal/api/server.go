package api

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "strings"
    "time"

    "github.com/lukelinux79/casaos-native-panel/internal/auth"
    "github.com/lukelinux79/casaos-native-panel/internal/runtime"
)

type Server struct {
    runtime runtime.Runtime
    authCfg auth.Config
}

type response map[string]any

func NewServer(r runtime.Runtime) *Server {
    return &Server{runtime: r, authCfg: auth.LoadFromEnv()}
}

func (s *Server) Handler() http.Handler {
    mux := http.NewServeMux()

    apiMux := http.NewServeMux()
    apiMux.HandleFunc("/api/v1/health", s.handleHealth)
    apiMux.HandleFunc("/api/v1/system/overview", s.handleOverview)
    apiMux.HandleFunc("/api/v1/services", s.handleServices)
    apiMux.HandleFunc("/api/v1/packages", s.handlePackages)
    apiMux.HandleFunc("/api/v1/users", s.handleUsers)
    apiMux.HandleFunc("/api/v1/groups", s.handleGroups)
    apiMux.HandleFunc("/api/v1/samba/config", s.handleSambaConfig)
    apiMux.HandleFunc("/api/v1/dashboard", s.handleDashboard)
    apiMux.HandleFunc("/api/v1/system/jobs", s.handleJobs)

    if s.authCfg.Enabled {
        mux.Handle("/api/", withBasicAuth(apiMux, s.authCfg))
    } else {
        mux.Handle("/api/", apiMux)
    }

    if _, err := os.Stat("web"); err == nil {
        fileServer := http.FileServer(http.Dir("web"))
        mux.Handle("/", fileServer)
    } else {
        mux.HandleFunc("/", s.handleRoot)
    }

    return mux
}

func withBasicAuth(next http.Handler, cfg auth.Config) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        user, pass, ok := r.BasicAuth()
        if !ok || user != cfg.Username || pass != cfg.Password {
            w.Header().Set("WWW-Authenticate", `Basic realm="CasaOS Native Panel"`)
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }
    _ = json.NewEncoder(w).Encode(response{"status": "ok", "service": "CasaOS Native Panel"})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(response{"status": "ok", "mode": s.runtime.Name()})
}

func (s *Server) handleOverview(w http.ResponseWriter, r *http.Request) {
    overview, err := s.runtime.GetOverview()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    _ = json.NewEncoder(w).Encode(response{"data": overview})
}

func (s *Server) handleServices(w http.ResponseWriter, r *http.Request) {
    services, err := s.runtime.GetServices()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    _ = json.NewEncoder(w).Encode(response{"services": services})
}

func (s *Server) handlePackages(w http.ResponseWriter, r *http.Request) {
    packages, err := s.runtime.GetPackages()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    _ = json.NewEncoder(w).Encode(response{"packages": packages})
}

func (s *Server) handleUsers(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        users, err := s.runtime.GetUsers()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        _ = json.NewEncoder(w).Encode(response{"users": users})
        return
    }

    if r.Method == http.MethodPost {
        var user runtime.User
        if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        if err := s.runtime.CreateUser(user); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.WriteHeader(http.StatusCreated)
        _ = json.NewEncoder(w).Encode(response{"status": "created", "user": user})
    }
}

func (s *Server) handleGroups(w http.ResponseWriter, r *http.Request) {
    groups, err := s.runtime.GetGroups()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    _ = json.NewEncoder(w).Encode(response{"groups": groups})
}

func (s *Server) handleSambaConfig(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        content, err := s.runtime.GetSambaConfig()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        _ = json.NewEncoder(w).Encode(response{"content": content})
        return
    }

    if r.Method == http.MethodPost || r.Method == http.MethodPut {
        var payload struct {
            Content string `json:"content"`
        }
        if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        if strings.TrimSpace(payload.Content) == "" {
            http.Error(w, "content is required", http.StatusBadRequest)
            return
        }
        if err := s.runtime.SaveSambaConfig(payload.Content); err != nil {
            status := http.StatusInternalServerError
            if strings.Contains(err.Error(), "testparm") || strings.Contains(err.Error(), "validation failed") {
                status = http.StatusBadRequest
            }
            w.WriteHeader(status)
            _ = json.NewEncoder(w).Encode(response{
                "error": err.Error(),
                "message": "Samba config validation failed and was not applied.",
            })
            return
        }
        w.WriteHeader(http.StatusAccepted)
        _ = json.NewEncoder(w).Encode(response{"status": "saved"})
    }
}

func (s *Server) handleJobs(w http.ResponseWriter, r *http.Request) {
    _ = json.NewEncoder(w).Encode(response{"jobs": []any{}})
}

func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    overview, err := s.runtime.GetOverview()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    services, err := s.runtime.GetServices()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    payload := response{
        "overview": overview,
        "services": services,
        "timestamp": time.Now().UTC().Format(time.RFC3339),
    }

    _ = json.NewEncoder(w).Encode(response{"data": payload})
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }
    _ = json.NewEncoder(w).Encode(response{"status": "ok", "service": "CasaOS Native Panel"})
}

func init() {
    // no-op placeholder to keep the file compile-safe under this minimal patch
}
