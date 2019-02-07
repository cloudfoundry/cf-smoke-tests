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
	var useExistingApp = testConfig.LoggingApp != ""
	var appName string

	Describe("cf logs", func() {
		AfterEach(func() {
			defer func() {
				if testConfig.Cleanup && !useExistingApp {
					Expect(cf.Cf("delete", appName, "-f", "-r").Wait(testConfig.GetDefaultTimeout())).To(Exit(0))
				}
			}()
			smoke.AppReport(appName, testConfig.GetDefaultTimeout())
		})

		Context("linux", func() {
			BeforeEach(func() {
				if !useExistingApp {
					appName = generator.PrefixedRandomName("SMOKES", "APP")
					Expect(cf.Cf("push", appName, "-b", "ruby_buildpack", "-p", smoke.SimpleRubyAppBitsPath, "-d", testConfig.AppsDomain).Wait(testConfig.GetPushTimeout())).To(Exit(0))
				}
			})

			It("can see app messages in the logs", func() {
				Eventually(func() *Session {
					appLogsSession := smoke.Logs(testConfig.UseLogCache, appName)
					Expect(appLogsSession.Wait(testConfig.GetDefaultTimeout())).To(Exit(0))

					return appLogsSession
				}, testConfig.GetDefaultTimeout()*5).Should(Say(`\[(App|APP).*/0\]`))
			})
		})

		Context("windows", func() {
			BeforeEach(func() {
				smoke.SkipIfNotWindows(testConfig)

				appName = generator.PrefixedRandomName("SMOKES", "APP")
				Expect(cf.Cf("push", appName, "-p", smoke.SimpleDotnetAppBitsPath, "-d", testConfig.AppsDomain, "-s", testConfig.GetWindowsStack(), "-b", "hwc_buildpack").Wait(testConfig.GetPushTimeout())).To(Exit(0))
			})

			It("can see app messages in the logs", func() {
				Eventually(func() *Session {
					appLogsSession := cf.Cf("logs", "--recent", appName)
					Expect(appLogsSession.Wait(testConfig.GetDefaultTimeout())).To(Exit(0))
					return appLogsSession
				}, testConfig.GetDefaultTimeout()*5).Should(Say(`\[(App|APP).*/0\]`))
			})
		})
	})
})
