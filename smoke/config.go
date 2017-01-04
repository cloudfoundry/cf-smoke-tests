package smoke

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	SuiteName string `json:"suite_name"`

	ApiEndpoint string `json:"api"`
	AppsDomain  string `json:"apps_domain"`

	SkipSSLValidation bool `json:"skip_ssl_validation"`

	User     string `json:"user"`
	Password string `json:"password"`

	Org   string `json:"org"`
	Space string `json:"space"`

	UseExistingOrg   bool `json:"use_existing_org"`
	UseExistingSpace bool `json:"use_existing_space"`

	// existing app names - if empty the space will be managed and a random app name will be used
	LoggingApp string `json:"logging_app"`
	RuntimeApp string `json:"runtime_app"`

	ArtifactsDirectory string `json:"artifacts_directory"`

	SyslogDrainPort int    `json:"syslog_drain_port"`
	SyslogIpAddress string `json:"syslog_ip_address"`

	Cleanup bool `json:"cleanup"`

	EnableWindowsTests          bool `json:"enable_windows_tests"`
	EnableEtcdClusterCheckTests bool `json:"enable_etcd_cluster_check_tests"`

	EtcdIpAddress string `json:"etcd_ip_address"`

	Backend string `json:"backend"`
}

// singleton cache
var cachedConfig *Config

func GetConfig() *Config {
	if cachedConfig == nil {
		cachedConfig = loadConfig()
	}
	return cachedConfig
}

func loadConfig() *Config {
	config := newDefaultConfig()
	loadConfigFromJson(config)
	validateRequiredFields(config)
	validateEtcdClusterCheckTests(config)
	return config
}

func newDefaultConfig() *Config {
	return &Config{
		ArtifactsDirectory:          filepath.Join("..", "results"),
		UseExistingOrg:              false,
		UseExistingSpace:            false,
		Cleanup:                     true,
		EnableWindowsTests:          false,
		EnableEtcdClusterCheckTests: false,
		EtcdIpAddress:               "",
	}
}

func validateRequiredFields(config *Config) {
	if config.SuiteName == "" {
		panic("missing configuration 'suite_name'")
	}

	if config.ApiEndpoint == "" {
		panic("missing configuration 'api'")
	}

	if config.AppsDomain == "" {
		panic("missing configuration 'apps_domain'")
	}

	if config.User == "" {
		panic("missing configuration 'user'")
	}

	if config.Password == "" {
		panic("missing configuration 'password'")
	}

	if config.Org == "" {
		panic("missing configuration 'org'")
	}

	if config.Space == "" {
		panic("missing configuration 'space'")
	}
}

func validateEtcdClusterCheckTests(config *Config) {
	if config.EnableEtcdClusterCheckTests == true && config.EtcdIpAddress == "" {
		panic("when etcd_cluster_check_tests is true, etcd_ip_address must be provided but it was not")
	}
}

// Loads the config from json into the supplied config object
func loadConfigFromJson(config *Config) {
	path := configPath()

	configFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(config)
	if err != nil {
		panic(err)
	}
}

func configPath() string {
	path := os.Getenv("CONFIG")
	if path == "" {
		panic("Must set $CONFIG to point to an integration config .json file.")
	}

	return path
}
