package main

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Web   WebConfig
	Debug DebugConfig
	ENV   string
}

type WebConfig struct {
	ReadTimeout     time.Duration `conf:"default:5s"`
	WriteTimeout    time.Duration `conf:"default:10s"`
	IdleTimeout     time.Duration `conf:"default:180s"`
	ShutdownTimeout time.Duration `conf:"default:20s"`
	APIHost         string        `conf:"default:0.0.0.0:3338"`
}

type DebugConfig struct {
	ReadTimeout     time.Duration `conf:"default:180s"`
	WriteTimeout    time.Duration `conf:"default:180s"`
	IdleTimeout     time.Duration `conf:"default:180s"`
	ShutdownTimeout time.Duration `conf:"default:20s"`
	APIHost         string        `conf:"default:0.0.0.0:3339"`
}

type SecretManagerConfig struct {
	ProjectID       string `conf:"default:887523945646"`
	SecretAppConfig string `conf:"default:go-base-structure"`
}

// ====================== //
// Hidden configurations  //
// ====================== //

type HiddenAppConfig struct {
	Postgres PostgresSecretConfig `yaml:"postgres"`
}

type PostgresSecretConfig struct {
	Host            string        `yaml:"host"`
	Port            string        `yaml:"port"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	Name            string        `yaml:"name"`
	MaxIdleConns    int           `yaml:"maxIdleConns"`
	ConnMaxIdleTime time.Duration `yaml:"connMaxIdleTime"`
	MaxOpenConns    int           `yaml:"maxOpenConns"`
	EnableTLS       bool          `yaml:"enableTLS"`
}

func readSecretConfig() (HiddenAppConfig, error) {
	var hiddenAppConfig HiddenAppConfig
	var cfg []byte
	var err error

	switch build {
	case "dev":
		cfg, err = os.ReadFile("./local_config.yaml")
		if err != nil {
			return HiddenAppConfig{}, fmt.Errorf("Error reading YAML file: %w\n", err)
		}
	}

	if err := yaml.Unmarshal(cfg, &hiddenAppConfig); err != nil {
		return HiddenAppConfig{}, fmt.Errorf("Error unmarshalling YAML file: %w\n", err)
	}

	return hiddenAppConfig, nil
}

// writeTempFile writes data to a temporary file and returns the file path.
func writeTempFile(prefix, data string) (string, error) {
	tmpFile, err := os.CreateTemp("", prefix)
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString(data); err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

func removeTempFile(file string) error {
	return os.Remove(file)
}
