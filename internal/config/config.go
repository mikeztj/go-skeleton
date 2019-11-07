package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var (
	config *Config
)

const (
	envDevelopment = "development"
	envStaging     = "staging"
	envProduction  = "production"
)

type option struct {
	configFile string
}

// Init ...
func Init(opts ...Option) error {
	opt := &option{
		configFile: getDefaultConfigFile(),
	}
	for _, optFunc := range opts {
		optFunc(opt)
	}

	out, err := ioutil.ReadFile(opt.configFile)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(out, &config)
}

// Option ...
type Option func(*option)

// WithConfigFile ...
func WithConfigFile(file string) Option {
	return func(opt *option) {
		opt.configFile = file
	}
}

func getDefaultConfigFile() string {
	var (
		repoPath   = filepath.Join(os.Getenv("GOPATH"), "src/github.com/vilbert/go-skeleton")
		configPath = filepath.Join(repoPath, "files/etc/skeleton/skeleton.environment.yaml")
		env        = os.Getenv("ENV")
	)

	if env != "" {
		if env == envStaging {
			configPath = "files/etc/skeleton/skeleton.environment.yaml"
		}
	}
	return configPath
}

// Get ...
func Get() *Config {
	return config
}
