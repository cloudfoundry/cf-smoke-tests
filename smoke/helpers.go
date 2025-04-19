package smoke

import (
	"time"

	"github.com/cloudfoundry/cf-test-helpers/v2/cf"

	. "github.com/onsi/ginkgo/v2"    //nolint:staticcheck
	. "github.com/onsi/gomega"       //nolint:staticcheck
	. "github.com/onsi/gomega/gexec" //nolint:staticcheck
)

const (
	SimpleBinaryAppBitsPath = "../../assets/binary"
	SimpleDotnetAppBitsPath = "../../assets/dotnet_simple/Published"
)

func SkipIfNotWindows(testConfig *Config) {
	if !testConfig.EnableWindowsTests {
		Skip("Windows tests are disabled")
	}
}

func AppReport(appName string, timeout time.Duration) {
	Eventually(cf.Cf("app", appName, "--guid"), timeout).Should(Exit())
	Eventually(cf.Cf("logs", appName, "--recent"), timeout).Should(Exit())
}

func Logs(useLogCache bool, appName string) *Session {
	if useLogCache {
		return cf.Cf("tail", appName, "--lines", "125")
	}
	return cf.Cf("logs", "--recent", appName)
}
