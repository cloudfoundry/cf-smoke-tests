package smoke

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	SuiteName string `json:"suite_name"`

	Reporter string `json:"reporter"`

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
	EnableIsolationSegmentTests bool `json:"enable_isolation_segment_tests"`

	EtcdIpAddress string `json:"etcd_ip_address"`

	Backend string `json:"backend"`

	TimeoutScale           *float64 `json:"timeout_scale"`
	IsolationSegmentName   string   `json:"isolation_segment_name"`
	IsolationSegmentDomain string   `json:"isolation_segment_domain"`
}

func (c *Config) GetIsolationSegmentName() string {
	return c.IsolationSegmentName
}

func (c *Config) GetIsolationSegmentDomain() string {
	return c.IsolationSegmentDomain
}

func (c *Config) GetApiEndpoint() string {
	return c.ApiEndpoint
}

func (c *Config) GetConfigurableTestPassword() string {
	return c.Password
}

func (c *Config) GetPersistentAppOrg() string {
	return ""
}

func (c *Config) GetPersistentAppQuotaName() string {
	return ""
}

func (c *Config) GetPersistentAppSpace() string {
	return ""
}

func (c *Config) GetScaledTimeout(timeout time.Duration) time.Duration {
	return time.Duration(float64(timeout) * *c.TimeoutScale)
}

func (c *Config) GetAdminPassword() string {
	return c.Password
}

func (c *Config) GetExistingUser() string {
	return c.User
}

func (c *Config) GetExistingUserPassword() string {
	return c.Password
}

func (c *Config) GetShouldKeepUser() bool {
	return true
}

func (c *Config) GetUseExistingUser() bool {
	return true
}

func (c *Config) GetAdminUser() string {
	return c.User
}

func (c *Config) GetAppsDomains() string {
	return c.AppsDomain
}

func (c *Config) GetUseExistingOrganization() bool {
	return c.UseExistingOrg
}

func (c *Config) GetExistingOrganization() string {
	return c.Org
}

func (c *Config) GetExistingSpace() string {
	return c.Space
}

func (c *Config) GetUseExistingSpace() bool {
	return c.UseExistingSpace
}

func (c *Config) GetSkipSSLValidation() bool {
	return c.SkipSSLValidation
}

func (c *Config) GetNamePrefix() string {
	return "SMOKE"
}

func (c *Config) GetDefaultTimeout() int {
	return 30
}

func (c *Config) GetPushTimeout() int {
	return 300
}

func (c *Config) GetBackend() string {
	return c.Backend
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
	validateIsolationSegments(config)
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

	if config.UseExistingOrg && config.Org == "" {
		panic("missing configuration 'org'")
	}

	if config.UseExistingSpace && config.Space == "" {
		panic("missing configuration 'space'")
	}
}

func validateEtcdClusterCheckTests(config *Config) {
	if config.EnableEtcdClusterCheckTests == true && config.EtcdIpAddress == "" {
		panic("when etcd_cluster_check_tests is true, etcd_ip_address must be provided but it was not")
	}
}

func validateIsolationSegments(config *Config) {
	if !config.EnableIsolationSegmentTests {
		return
	}
	if config.GetBackend() != "diego" {
		panic("* Invalid Configuration: 'backend' must be set to 'diego' if 'enable_isolation_segment_tests' is true")
	}
	if config.GetIsolationSegmentName() == "" {
		panic("* Invalid configuration: 'isolation_segment_name' must be provided if 'enable_isolation_segment_tests' is true")
	}
	if config.GetIsolationSegmentDomain() == "" {
		panic("* Invalid configuration: 'isolation_segment_domain' must be provided if 'enable_isolation_segment_tests' is true")
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

	if config.TimeoutScale == nil {
		defaultTimeout := 1.0
		config.TimeoutScale = &defaultTimeout
	}
}

func configPath() string {
	path := os.Getenv("CONFIG")
	if path == "" {
		panic("Must set $CONFIG to point to an integration config .json file.")
	}

	return path
}
