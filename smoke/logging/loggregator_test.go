package logging

import (
	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	"github.com/cloudfoundry-incubator/cf-test-helpers/generator"
	"github.com/cloudfoundry/cf-smoke-tests/smoke"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Loggregator:", func() {
	var testConfig = smoke.GetConfig()
	var useExistingApp = (testConfig.LoggingApp != "")
	var appName string

	Describe("cf logs", func() {
		AfterEach(func() {
			smoke.AppReport(appName, CF_TIMEOUT_IN_SECONDS)
			if testConfig.Cleanup && !useExistingApp {
				Expect(cf.Cf("delete", appName, "-f", "-r").Wait(CF_TIMEOUT_IN_SECONDS)).To(Exit(0))
			}
		})

		Context("linux", func() {
			BeforeEach(func() {
				if !useExistingApp {
					appName = generator.PrefixedRandomName("SMOKES", "APP")
					Expect(cf.Cf("push", appName, "-b", "ruby_buildpack", "-p", SIMPLE_RUBY_APP_BITS_PATH, "-d", testConfig.AppsDomain, "--no-start").Wait(CF_PUSH_TIMEOUT_IN_SECONDS)).To(Exit(0))
					smoke.SetBackend(appName)
					Expect(cf.Cf("start", appName).Wait(CF_PUSH_TIMEOUT_IN_SECONDS)).To(Exit(0))
				}
			})

			It("can see app messages in the logs", func() {
				Eventually(func() *Session {
					appLogsSession := cf.Cf("logs", "--recent", appName)
					Expect(appLogsSession.Wait(CF_TIMEOUT_IN_SECONDS)).To(Exit(0))
					return appLogsSession
				}, CF_TIMEOUT_IN_SECONDS*5).Should(Say(`\[(App|APP).*/0\]`))
			})
		})

		Context("windows", func() {
			BeforeEach(func() {
				smoke.SkipIfWindows(testConfig)

				appName = generator.PrefixedRandomName("SMOKES", "APP")
				Expect(cf.Cf("push", appName, "-p", SIMPLE_DOTNET_APP_BITS_PATH, "-d", testConfig.AppsDomain, "-s", "windows2012R2", "-b", "hwc_buildpack", "--no-start").Wait(CF_PUSH_TIMEOUT_IN_SECONDS)).To(Exit(0))
				smoke.EnableDiego(appName)
				Expect(cf.Cf("start", appName).Wait(CF_PUSH_TIMEOUT_IN_SECONDS)).To(Exit(0))
			})

			It("can see app messages in the logs", func() {
				Eventually(func() *Session {
					appLogsSession := cf.Cf("logs", "--recent", appName)
					Expect(appLogsSession.Wait(CF_TIMEOUT_IN_SECONDS)).To(Exit(0))
					return appLogsSession
				}, CF_TIMEOUT_IN_SECONDS*5).Should(Say(`\[(App|APP).*/0\]`))
			})
		})
	})
})
