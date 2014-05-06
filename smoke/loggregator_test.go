package smoke

import (
	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	"github.com/cloudfoundry-incubator/cf-test-helpers/generator"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Loggregator:", func() {
	var testConfig = GetConfig()
	var useExistingApp = (testConfig.LoggingApp != "")
	var appName string

	BeforeEach(func() {
		appName = testConfig.LoggingApp
		if !useExistingApp {
			appName = generator.RandomName()
			Eventually(cf.Cf("push", appName, "-p", SIMPLE_RUBY_APP_BITS_PATH), CF_PUSH_TIMEOUT_IN_SECONDS).Should(Exit(0))
		}
	})

	AfterEach(func() {
		if !useExistingApp {
			Eventually(cf.Cf("delete", appName, "-f"), CF_TIMEOUT_IN_SECONDS).Should(Exit(0))
		}
	})

	It("can see app messages in the logs", func() {
		appLogsSession := cf.Cf("logs", "--recent", appName)
		Eventually(appLogsSession, CF_TIMEOUT_IN_SECONDS).Should(Exit(0))
		Expect(appLogsSession).To(Say(`\[App/0\]`))
	})
})
