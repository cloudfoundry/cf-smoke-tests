package smoke

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-cf-experimental/cf-test-helpers/cf"
	. "github.com/pivotal-cf-experimental/cf-test-helpers/generator"
	"github.com/vito/cmdtest"
	. "github.com/vito/cmdtest/matchers"
	"os"
)

var _ = Describe("Application Flow", func() {
	BeforeEach(func() {
		os.Setenv("CF_COLOR", "false")
		AppName = RandomName()
	})

	AfterEach(func() {
		Expect(Cf("delete", AppName, "-f")).To(Say("OK"))
	})

	It("can be pushed, scaled and deleted", func() {
		Expect(Cf("push", AppName, "-p", AppPath)).To(Say("App started"))
		Eventually(Curling("/")).Should(Say("It just needed to be restarted!"))

		Expect(Cf("scale", AppName, "-i", "2")).To(Say("OK"))
		Eventually(func() *cmdtest.Session {
			return Cf("app", AppName)
		}, 10).Should(Say("instances: 2/2"))
		Eventually(Curling("/")).Should(Say("\"instance_index\":0"))
		Eventually(Curling("/")).Should(Say("\"instance_index\":1"))

		Expect(Cf("delete", AppName, "-f")).To(Say("OK"))
		Expect(Cf("app", AppName)).To(Say("not found"))
		Eventually(Curling("/")).Should(Say("404"))
	})
})
