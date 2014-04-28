package smoke

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/gbytes"
	"github.com/pivotal-cf-experimental/cf-test-helpers/cf"
	. "github.com/pivotal-cf-experimental/cf-test-helpers/generator"
)

var _ = Describe("Loggregator:", func() {
	var testConfig = GetConfig()
	var createTestApp = (testConfig.LoggingApp == "")
	var appName string

	BeforeEach(func() {
		if createTestApp {
			appName = RandomName()
			Eventually(cf.Cf("push", appName, "-p", SIMPLE_RUBY_APP_BITS_PATH), CF_PUSH_TIMEOUT_IN_SECONDS).Should(Exit(0))
		}  else {
			appName = testConfig.LoggingApp
		}
	})

	AfterEach(func() {
		if createTestApp {
			Eventually(cf.Cf("delete", appName, "-f"), CF_TIMEOUT_IN_SECONDS).Should(Exit(0))
		}
	})

	It("can see app messages in the logs", func() {
		appLogsSession := cf.Cf("logs", "--recent", appName)
		Eventually(appLogsSession, CF_TIMEOUT_IN_SECONDS).Should(Exit(0))
		Expect(appLogsSession).To(Say(`\[App/0\]`))
	})
})
