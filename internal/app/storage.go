package app

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
	Token        string `json:"token,omitempty"`
	EmployeeID   string `json:"employee_id,omitempty"`
	EmployeeName string `json:"employee_name,omitempty"`
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
// Salva o nome do funcionário no keyring ou em um arquivo
func SaveEmployeeName(employeeName string) error {
	if err := keyring.Set(service, "employee_name", employeeName); err == nil {
		return nil
	}
	// fallback para arquivo
	return saveEmployeeNameToFile(employeeName)
}

// Carrega o nome do funcionário do keyring ou do arquivo
func LoadEmployeeName() (string, error) {
	if employeeName, err := keyring.Get(service, "employee_name"); err == nil {
		return employeeName, nil
	}
	return loadEmployeeNameFromFile()
}

func getConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "ttq", "config.json")
}

func loadConfigFromFile() (Config, error) {
	var cfg Config
	path := getConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func writeConfigToFile(cfg Config) error {
	data, _ := json.MarshalIndent(cfg, "", "  ")
	path := getConfigPath()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func saveTokenToFile(token string) error {
	cfg, _ := loadConfigFromFile() // se erro, começa vazio
	cfg.Token = token
	return writeConfigToFile(cfg)
}

func loadTokenFromFile() (string, error) {
	cfg, err := loadConfigFromFile()
	if err != nil || cfg.Token == "" {
		return "", errors.New("nenhum token encontrado")
	}
	return cfg.Token, nil
}

func saveEmployeeIDToFile(employeeID string) error {
	cfg, _ := loadConfigFromFile()
	cfg.EmployeeID = employeeID
	return writeConfigToFile(cfg)
}

func saveEmployeeNameToFile(employeeName string) error {
	cfg, _ := loadConfigFromFile()
	cfg.EmployeeName = employeeName
	return writeConfigToFile(cfg)
}

func loadEmployeeIDFromFile() (string, error) {
	cfg, err := loadConfigFromFile()
	if err != nil || cfg.EmployeeID == "" {
		return "", errors.New("nenhum ID de funcionário encontrado")
	}
	return cfg.EmployeeID, nil
}

func loadEmployeeNameFromFile() (string, error) {
	cfg, err := loadConfigFromFile()
	if err != nil || cfg.EmployeeName == "" {
		return "", errors.New("nenhum nome de funcionário encontrado")
	}
	return cfg.EmployeeName, nil
}
