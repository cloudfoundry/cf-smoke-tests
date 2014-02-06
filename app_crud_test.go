package cf_health_checks_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-cf-experimental/cf-test-helpers/cf"
	. "github.com/pivotal-cf-experimental/cf-test-helpers/generator"
	. "github.com/vito/cmdtest/matchers"
)

var _ = Describe("Application", func() {
	// Behavior of original NYETs:

	// Clean up app from previous run
	// Create an app as a regular user
	// Clean up route from a previous run
	// Create a route as a regular user
	// Deploy an app
	// Start app (what's the difference from gcf push?)
	// Check if the app is routable

	// Scale the app
	// Check running instances
	// Check if the first and second instances are reachable

	// Delete the route
	// Delete the app
	// Check that the app's api is unavailable (?)
	// Check that the app's uri is unavailable

	// Monitoring around all of these things?

	// In this case, simplified to:

	// Push an app
	// Scale the app
	// Delete the app

	BeforeEach(func() {
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
