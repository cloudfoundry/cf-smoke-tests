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

	It("can see app messages in the logs", func() {
		if os.Getenv("CLEANUP_ENVIRONMENT") != "false" {
			Expect(Cf("push", AppName, "-p", AppPath)).To(Say("App started"))
		}

		Eventually(Cf("logs", "--recent", AppName)).Should(Say("[App/0]"))
	})
})
