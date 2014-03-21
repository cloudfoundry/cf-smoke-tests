package smoke

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-cf-experimental/cf-test-helpers/cf"
	. "github.com/pivotal-cf-experimental/cf-test-helpers/generator"
	. "github.com/vito/cmdtest/matchers"
	"os"
)

var _ = Describe("Loggregator:", func() {
	BeforeEach(func() {
		os.Setenv("CF_COLOR", "false")
		if os.Getenv("CLEANUP_ENVIRONMENT") == "false" {
			AppName = "smoke-test-app"
		}  else {
			AppName = RandomName()
		}
	})

	AfterEach(func() {
		if os.Getenv("CLEANUP_ENVIRONMENT") != "false" {
			Expect(Cf("delete", AppName, "-f")).To(Say("OK"))
		}
	})

	It("can see router requests in the logs", func() {
		if os.Getenv("CLEANUP_ENVIRONMENT") != "false" {
			Expect(Cf("push", AppName, "-p", AppPath)).To(Say("App started"))
		}

		Eventually(Curling("/")).Should(Say("It just needed to be restarted!"))

		// Curling multiple times because loggregator makes no guarantees about delivery of logs.
		Eventually(Curling("/")).Should(Say("Healthy"))
		Eventually(Cf("logs", "--recent", AppName)).Should(Say("[RTR]"))

		Eventually(Curling("/")).Should(Say("Healthy"))
		Eventually(Cf("logs", "--recent", AppName)).Should(Say("[App/0]"))
	})
})
