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

// option defines configuration option
type option struct {
	configFile string
}

// Init initializes `config` from the default config file.
// use `WithConfigFile` to specify the location of the config file
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

// Option define an option for config package
type Option func(*option)

// WithConfigFile set `config` to use the given config file
func WithConfigFile(file string) Option {
	return func(opt *option) {
		opt.configFile = file
	}
}

// getDefaultConfigFile get default config file.
func getDefaultConfigFile() string {
	var (
		repoPath   = filepath.Join(os.Getenv("GOPATH"), "src/github.com/tokopedia/td-report-engine")
		configPath = filepath.Join(repoPath, "files/etc/td-report-engine/config.yaml")
	)

	_, err := os.Stat("/etc/td-report-engine/config.yaml")
	if !os.IsNotExist(err) {
		configPath = "/etc/td-report-engine/config.yaml"
	}

	return configPath
}

// Get config
func Get() *Config {
	if config == nil {
		config = &Config{}
	}
	return config
}
