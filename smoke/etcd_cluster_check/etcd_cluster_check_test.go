package logging

import (
	"fmt"
	"net/url"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	"github.com/cloudfoundry-incubator/cf-test-helpers/generator"
	"github.com/cloudfoundry-incubator/cf-test-helpers/helpers"
	"github.com/cloudfoundry/cf-smoke-tests/smoke"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Etcd Cluster Check:", func() {
	var testConfig = smoke.GetConfig()
	var appName string

	Describe("etcd cluster checking", func() {
		BeforeEach(func() {
			if testConfig.EnableEtcdClusterCheckTests != true {
				Skip("Skipping because EnableEtcdClusterCheckTests flag set to false")
			}
			appName = generator.PrefixedRandomName("SMOKES", "APP")
			Eventually(cf.Cf(
				"push", appName,
				"-p", CURLER_RUBY_APP_BITS_PATH,
				"--no-start",
				"-b", "ruby_buildpack",
				"-d", testConfig.AppsDomain,
				"-i", "1"),
				CF_PUSH_TIMEOUT_IN_SECONDS,
			).Should(Exit(0))
			smoke.SetBackend(appName)
			Expect(cf.Cf("start", appName).Wait(CF_PUSH_TIMEOUT_IN_SECONDS)).To(Exit(0))
		})

		It("does not allow apps to directly modify etcd", func() {
			// Post to the app, which triggers a POST to etcd
			p := url.Values{}
			p.Add("host", testConfig.EtcdIpAddress)
			p.Add("path", "/v2/keys/foo")
			p.Add("port", "4001")
			p.Add("data", "value=updated_value")
			curlCmd := helpers.CurlSkipSSL(true, fmt.Sprintf("https://%s.%s/put/?%s", appName, testConfig.AppsDomain, p.Encode()))
			errorMsg := fmt.Sprintf("Connections from application containers to internal IP addresses such as etcd node <%s> were not rejected. Please review the documentation on Application Security Groups to disallow traffic from application containers to internal IP addresses: https://docs.cloudfoundry.org/adminguide/app-sec-groups.html", testConfig.EtcdIpAddress)
			Eventually(curlCmd, CF_TIMEOUT_IN_SECONDS).Should(Exit(0), errorMsg)
			Expect(string(curlCmd.Out.Contents())).To(ContainSubstring("Connection refused"), errorMsg)
		})

		AfterEach(func() {
			smoke.AppReport(appName, CF_TIMEOUT_IN_SECONDS)
			if testConfig.Cleanup {
				Expect(cf.Cf("delete", appName, "-f", "-r").Wait(CF_TIMEOUT_IN_SECONDS)).To(Exit(0))
			}
		})
	})
})
