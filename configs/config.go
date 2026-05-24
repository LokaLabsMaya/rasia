package configs

import (
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents configuration variables
type Config struct {
	App     App     `yaml:"app"`
	Log     Log     `yaml:"log"`
	Http    Http    `yaml:"http"`
	Crypto  Crypto  `yaml:"crypto"`
	Secrets Secrets `yaml:"secrets"`
}

// Crypto holds the encryption key for secrets.
type Crypto struct {
	Key string `yaml:"key"`
}

// Secrets holds the SQLite configuration for the secrets manager.
type Secrets struct {
	DBPath string `yaml:"db_path"`
}

type App struct {
	Name string `yaml:"name"`
}

type Log struct {
	Level          uint32   `yaml:"level"`
	Format         string   `yaml:"format"`
	RedactedFields []string `yaml:"redacted_fields"`
}

type Http struct {
	Port                  string        `yaml:"port"`
	ReadTimeout           time.Duration `yaml:"read_timeout"`
	WriteTimeout          time.Duration `yaml:"write_timeout"`
	IdleTimeout           time.Duration `yaml:"idle_timeout"`
	StartTimeout          time.Duration `yaml:"start_timeout"`
	StopTimeout           time.Duration `yaml:"stop_timeout"`
	EnablePrintRoutes     bool          `yaml:"enable_print_routes"`
	DisableStartupMessage bool          `yaml:"disable_startup_message"`
	AllowOrigins          []string      `yaml:"allow_origins"`
}

// GetEnv returns the environment variable ENV
func (c *Config) GetEnv() Env {
	return Env(os.Getenv("ENV"))
}

// LoadConfig loads the configuration from the given file path
func LoadConfig(cfgFile string) *Config {

	var cfg Config

	// read file cfgFile
	data, err := os.ReadFile(cfgFile)
	if err != nil {
		log.Fatalf("read config error: %v", err)
	}

	// unmarshal yaml to config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("unmarshal yaml error: %v", err)
	}

	return &cfg
}
