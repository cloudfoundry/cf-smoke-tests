package cf_health_checks_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-cf-experimental/cf-test-helpers/cf"
	. "github.com/pivotal-cf-experimental/cf-test-helpers/generator"
	"github.com/vito/cmdtest"
	. "github.com/vito/cmdtest/matchers"
	"os"
)

var _ = Describe("Application", func() {
	BeforeEach(func() {
		os.Setenv("CF_COLOR", "false")
		AppName = RandomName()

		Expect(Cf("push", AppName, "-p", AppPath)).To(Say("App started"))
	})

	AfterEach(func() {
		Expect(Cf("delete", AppName, "-f")).To(Say("OK"))
	})

	Describe("pushing", func() {
		It("makes the app reachable via its bound route", func() {
			Eventually(Curling("/")).Should(Say("It just needed to be restarted!"))
		})
	})

	Describe("scaling", func() {
		BeforeEach(func() {
			Expect(Cf("scale", AppName, "-i", "2")).To(Say("OK"))
		})

		It("reports 2 instances", func() {
			Eventually(func() *cmdtest.Session {
				return Cf("app", AppName)
			}, 10).Should(Say("instances: 2/2"))
		})

		It("actually starts a second instance", func() {
			Eventually(Curling("/")).Should(Say("\"instance_index\":0"))
			Eventually(Curling("/")).Should(Say("\"instance_index\":1"))
		})
	})

	Describe("deleting", func() {
		BeforeEach(func() {
			Expect(Cf("delete", AppName, "-f")).To(Say("OK"))
		})

		It("removes the application", func() {
			Expect(Cf("app", AppName)).To(Say("not found"))
		})

		It("makes the app unreachable", func() {
			Eventually(Curling("/")).Should(Say("404"))
		})
	})
})
