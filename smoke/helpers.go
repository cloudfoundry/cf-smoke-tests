package smoke

import (
	"time"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var DEFAULT_TIMEOUT = 30 * time.Second

func SkipIfWindows(testConfig *Config) {
	if !testConfig.EnableWindowsTests {
		Skip("Windows tests are disabled")
	}
}

func AppReport(appName string, timeout time.Duration) {
	Eventually(cf.Cf("app", appName, "--guid"), timeout*time.Second).Should(Exit())
	Eventually(cf.Cf("logs", appName, "--recent"), timeout*time.Second).Should(Exit())
}

func Logs(useLogCache bool, appName string) *Session {
	if useLogCache {
		return cf.Cf("tail", appName, "--lines", "125")
	}
	return cf.Cf("logs", "--recent", appName)
}
