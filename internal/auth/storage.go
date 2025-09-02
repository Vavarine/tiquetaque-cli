package auth

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/zalando/go-keyring"
)

const (
	service = "ttq"     // nome que aparece no keyring
	user    = "default" // chave para guardar o token
)

// Config usado no fallback de arquivo
type Config struct {
	Token string `json:"token"`
}

// SaveToken tenta salvar no keyring, se falhar salva em ~/.config
func SaveToken(token string) error {
	if err := keyring.Set(service, user, token); err == nil {
		return nil
	}
	// fallback para arquivo
	return saveTokenToFile(token)
}

// LoadToken tenta carregar do keyring, se falhar tenta do arquivo
func LoadToken() (string, error) {
	if token, err := keyring.Get(service, user); err == nil {
		return token, nil
	}
	return loadTokenFromFile()
}

// SaveEmployeeID salva o ID do funcionário no keyring ou em um arquivo
func SaveEmployeeID(employeeID string) error {
	if err := keyring.Set(service, "employee_id", employeeID); err == nil {
		return nil
	}
	// fallback para arquivo
	return saveEmployeeIDToFile(employeeID)
}

// LoadEmployeeID tenta carregar o ID do funcionário do keyring, se falhar tenta do arquivo
func LoadEmployeeID() (string, error) {
	if employeeID, err := keyring.Get(service, "employee_id"); err == nil {
		return employeeID, nil
	}
	return loadEmployeeIDFromFile()
}

// ==== fallback em arquivo ~/.config/minhacli/config.json ====

func getConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "minhacli", "config.json")
}

func saveTokenToFile(token string) error {
	cfg := Config{Token: token}
	data, _ := json.MarshalIndent(cfg, "", "  ")
	path := getConfigPath()
	os.MkdirAll(filepath.Dir(path), 0700)
	return os.WriteFile(path, data, 0600)
}

func loadTokenFromFile() (string, error) {
	data, err := os.ReadFile(getConfigPath())
	if err != nil {
		return "", errors.New("nenhum token encontrado")
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return "", err
	}
	return cfg.Token, nil
}

func saveEmployeeIDToFile(employeeID string) error {
	cfg := Config{Token: employeeID}
	data, _ := json.MarshalIndent(cfg, "", "  ")
	path := getConfigPath()
	os.MkdirAll(filepath.Dir(path), 0700)
	return os.WriteFile(path, data, 0600)
}

func loadEmployeeIDFromFile() (string, error) {
	data, err := os.ReadFile(getConfigPath())
	if err != nil {
		return "", errors.New("nenhum ID de funcionário encontrado")
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return "", err
	}
	return cfg.Token, nil
}
