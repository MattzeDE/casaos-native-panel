package auth

import (
    "os"
)

type Config struct {
    Enabled  bool
    Username string
    Password string
}

func LoadFromEnv() Config {
    username := os.Getenv("CASAOS_ADMIN_USER")
    password := os.Getenv("CASAOS_ADMIN_PASSWORD")
    if username == "" && password == "" {
        return Config{Enabled: false}
    }
    if username == "" {
        username = "admin"
    }
    if password == "" {
        password = "admin"
    }
    return Config{Enabled: true, Username: username, Password: password}
}
