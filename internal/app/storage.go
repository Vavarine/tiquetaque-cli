package app

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/zalando/go-keyring"
)

const (
	service = "ttq"
	user    = "default"
)

type Config struct {
	Token        string `json:"token,omitempty"`
	EmployeeID   string `json:"employee_id,omitempty"`
	EmployeeName string `json:"employee_name,omitempty"`
}

func SaveToken(token string) error {
	if err := keyring.Set(service, user, token); err == nil {
		return nil
	}

	return saveTokenToFile(token)
}

func LoadToken() (string, error) {
	if token, err := keyring.Get(service, user); err == nil {
		return token, nil
	}
	return loadTokenFromFile()
}

func SaveEmployeeID(employeeID string) error {
	if err := keyring.Set(service, "employee_id", employeeID); err == nil {
		return nil
	}

	return saveEmployeeIDToFile(employeeID)
}

func LoadEmployeeID() (string, error) {
	if employeeID, err := keyring.Get(service, "employee_id"); err == nil {
		return employeeID, nil
	}
	return loadEmployeeIDFromFile()
}

func SaveEmployeeName(employeeName string) error {
	if err := keyring.Set(service, "employee_name", employeeName); err == nil {
		return nil
	}

	return saveEmployeeNameToFile(employeeName)
}

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
	cfg, _ := loadConfigFromFile()
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
