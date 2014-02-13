package smoke

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-cf-experimental/cf-test-helpers/cf"
	. "github.com/pivotal-cf-experimental/cf-test-helpers/generator"
	. "github.com/vito/cmdtest/matchers"
	"os"
)

var _ = Describe("Logs", func() {
	BeforeEach(func() {
		os.Setenv("CF_COLOR", "false")
		AppName = RandomName()
	})

	AfterEach(func() {
		Expect(Cf("delete", AppName, "-f")).To(Say("OK"))
	})

	It("can see router requests in the logs", func() {
		Expect(Cf("push", AppName, "-p", AppPath)).To(Say("App started"))
		Eventually(Curling("/")).Should(Say("It just needed to be restarted!"))

		// Curling multiple times because loggregator makes no garauntees about delivery of logs.
		Eventually(Curling("/")).Should(Say("Healthy"))
		Eventually(Curling("/")).Should(Say("Healthy"))

		Eventually(Cf("logs", "--recent", AppName)).Should(Say("[RTR]"))
		Eventually(Cf("logs", "--recent", AppName)).Should(Say("[App/0]"))
	})
})
