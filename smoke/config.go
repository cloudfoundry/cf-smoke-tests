package smoke

import (
	"encoding/json"
	"os"
	"time"
)

type Config struct {
	SuiteName string `json:"suite_name"`

	Reporter string `json:"reporter"`

	APIEndpoint string `json:"api"`
	AppsDomain  string `json:"apps_domain"`

	SkipSSLValidation bool `json:"skip_ssl_validation"`

	User         string `json:"user"`
	Password     string `json:"password"`
	Client       string `json:"client"`
	ClientSecret string `json:"client_secret"`

	Org   string `json:"org"`
	Space string `json:"space"`

	UseExistingOrg   bool `json:"use_existing_org"`
	UseExistingSpace bool `json:"use_existing_space"`

	UseLogCache bool `json:"use_log_cache"`

	// existing app names - if empty the space will be managed and a random app name will be used
	LoggingApp string `json:"logging_app"`
	RuntimeApp string `json:"runtime_app"`

	ArtifactsDirectory string `json:"artifacts_directory"`

	Cleanup bool `json:"cleanup"`

	EnableWindowsTests          bool   `json:"enable_windows_tests"`
	WindowsStack                string `json:"windows_stack"`
	EnableIsolationSegmentTests bool   `json:"enable_isolation_segment_tests"`

	TimeoutScale           *float64 `json:"timeout_scale"`
	IsolationSegmentName   string   `json:"isolation_segment_name"`
	IsolationSegmentDomain string   `json:"isolation_segment_domain"`
	IsolationSegmentSpace  string   `json:"isolation_segment_space"`
}

func (c *Config) GetIsolationSegmentName() string {
	return c.IsolationSegmentName
}

func (c *Config) GetIsolationSegmentDomain() string {
	return c.IsolationSegmentDomain
}

func (c *Config) GetIsolationSegmentSpace() string {
	return c.IsolationSegmentSpace
}

func (c *Config) GetApiEndpoint() string {
	return c.APIEndpoint
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

func (c *Config) GetExistingUser() string {
	return c.User
}

func (c *Config) GetExistingUserPassword() string {
	return c.Password
}

func (c *Config) GetExistingClient() string {
	return c.Client
}

func (c *Config) GetExistingClientSecret() string {
	return c.ClientSecret
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

func (c *Config) GetAdminPassword() string {
	return c.Password
}

func (c *Config) GetAdminClient() string {
	return c.Client
}

func (c *Config) GetAdminClientSecret() string {
	return c.ClientSecret
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

func (c *Config) GetDefaultTimeout() time.Duration {
	return 30 * time.Second
}

func (c *Config) GetPushTimeout() time.Duration {
	return 300 * time.Second
}

func (c *Config) GetScaleTimeout() time.Duration {
	return 120 * time.Second
}

func (c *Config) GetAppStatusTimeout() time.Duration {
	return 120 * time.Second
}

func (c *Config) GetWindowsStack() string {
	return c.WindowsStack
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
	loadConfigFromJSON(config)
	validateRequiredFields(config)
	validateIsolationSegments(config)
	return config
}

func newDefaultConfig() *Config {
	return &Config{
		UseExistingOrg:     false,
		UseExistingSpace:   false,
		Cleanup:            true,
		EnableWindowsTests: false,
		WindowsStack:       "windows2012R2",
	}
}

func validateRequiredFields(config *Config) {
	if config.SuiteName == "" {
		panic("missing configuration 'suite_name'")
	}

	if config.APIEndpoint == "" {
		panic("missing configuration 'api'")
	}

	if config.AppsDomain == "" {
		panic("missing configuration 'apps_domain'")
	}

	missingUserCredentials := config.User == "" || config.Password == ""
	missingClientCredentials := config.Client == "" || config.ClientSecret == ""

	if missingUserCredentials && missingClientCredentials {
		panic("missing configuration: you must provide either 'user'/'password' or 'client'/'client_secret'")
	}

	if config.UseExistingOrg && config.Org == "" {
		panic("missing configuration: you must provide 'org' if 'use_existing_org' is true")
	}

	if config.UseExistingSpace && config.Space == "" {
		panic("missing configuration: you must provide 'space' if 'use_existing_space' is true")
	}

	if config.UseExistingSpace && !config.UseExistingOrg {
		panic("missing configuration: 'use_existing_org' must be set if 'use_existing_space' is set")
	}
}

func validateIsolationSegments(config *Config) {
	if !config.EnableIsolationSegmentTests {
		return
	}
	if config.IsolationSegmentName == "" {
		panic("* Invalid configuration: 'isolation_segment_name' must be provided if 'enable_isolation_segment_tests' is true")
	}
	if config.IsolationSegmentDomain == "" {
		panic("* Invalid configuration: 'isolation_segment_domain' must be provided if 'enable_isolation_segment_tests' is true")
	}
	if config.UseExistingOrg && config.UseExistingSpace && config.IsolationSegmentSpace == "" {
		panic("* Invalid configuration: 'isolation_segment_space' must be provided if 'use_existing_org' and 'use_existing_space' are true")
	}
}

// Loads the config from json into the supplied config object
func loadConfigFromJSON(config *Config) {
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
