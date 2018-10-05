package finch

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	raven "github.com/getsentry/raven-go"
)

// Config is a type used for storing configuration information.
type Config map[string]interface{}

// Get pulls information from the config file or from an environment variable.
func (c Config) Get(key string) interface{} {
	if val, ok := c[key]; ok {
		return val
	}

	up := strings.ToUpper(key)

	if val := os.Getenv(up); val != "" {
		return val
	}

	return nil
}

// LoadConfig loads the saved config, if it exists.
//
// It looks for a FINCH_CONFIG environmental variable,
// before falling back to a file name config.json.
func LoadConfig() (*Config, error) {
	fileName := os.Getenv("FINCH_CONFIG")
	if fileName == "" {
		fileName = "config.json"
	}

	f, err := ioutil.ReadFile(fileName)
	if err != nil {
		return &Config{}, nil
	}

	var cfg Config
	json.Unmarshal(f, &cfg)

	return &cfg, nil
}

// Save saves the current Config struct.
//
// It uses the same file as LoadConfig.
func (c *Config) Save() error {
	b, err := json.Marshal(c)
	if err != nil {
		if sentryEnabled {
			raven.CaptureErrorAndWait(err, nil)
		}

		return err
	}

	fileName := os.Getenv("FINCH_CONFIG")
	if fileName == "" {
		fileName = "config.json"
	}

	return ioutil.WriteFile(fileName, b, 0600)
}
