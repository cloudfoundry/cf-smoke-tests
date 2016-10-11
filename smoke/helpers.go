package smoke

import (
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var DEFAULT_TIMEOUT = 30 * time.Second

func SetBackend(appName string) {
	config := GetConfig()
	if config.Backend == "diego" {
		EnableDiego(appName)
	} else if config.Backend == "dea" {
		DisableDiego(appName)
	}
}

func EnableDiego(appName string) {
	guid := GetAppGuid(appName)
	Eventually(cf.Cf("curl", "/v2/apps/"+guid, "-X", "PUT", "-d", `{"diego": true}`), DEFAULT_TIMEOUT).Should(Exit(0))
}

func DisableDiego(appName string) {
	guid := GetAppGuid(appName)
	Eventually(cf.Cf("curl", "/v2/apps/"+guid, "-X", "PUT", "-d", `{"diego": false}`), DEFAULT_TIMEOUT).Should(Exit(0))
}

func GetAppGuid(appName string) string {
	cfApp := cf.Cf("app", appName, "--guid")
	Eventually(cfApp, DEFAULT_TIMEOUT).Should(Exit(0))

	appGuid := strings.TrimSpace(string(cfApp.Out.Contents()))
	Expect(appGuid).NotTo(Equal(""))
	return appGuid
}

func SkipIfWindows(testConfig *Config) {
	if !testConfig.EnableWindowsTests {
		Skip("Windows tests are disabled")
	}
}

func AppReport(appName string, timeout time.Duration) {
	Eventually(cf.Cf("app", appName, "--guid"), timeout*time.Second).Should(Exit())
	Eventually(cf.Cf("logs", appName, "--recent"), timeout*time.Second).Should(Exit())
}

func TestResourcesSummary(testConfig *Config) {
	Expect(cf.Cf("org", testConfig.Org, "--guid").Wait(DEFAULT_TIMEOUT)).To(Exit(0))
	Expect(cf.Cf("space", testConfig.Space, "--guid").Wait(DEFAULT_TIMEOUT)).To(Exit(0))
}
