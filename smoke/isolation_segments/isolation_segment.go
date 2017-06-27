package isolation_segments

import (
	"fmt"
	"io/ioutil"

	"github.com/cloudfoundry/cf-smoke-tests/smoke"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	"github.com/cloudfoundry-incubator/cf-test-helpers/generator"
)

const (
	SHARED_ISOLATION_SEGMENT_GUID = "933b4c58-120b-499a-b85d-4b6fc9e2903b"
	binaryHi                      = "Hello from a binary"
	SHARED_ISOLATION_SEGMENT_NAME = "shared"
	BINARY_APP_BITS_PATH          = "../../assets/binary"
)

var _ = Describe("RoutingIsolationSegments", func() {
	var appsDomain string
	var orgGuid, orgName string
	var spaceGuid, spaceName string
	var isoSpaceGuid, isoSpaceName string
	var isoSegName, isoSegDomain string
	var testConfig *smoke.Config

	BeforeEach(func() {

		// New up a organization since we will be assigning isolation segments.
		// This has a potential to cause other tests to fail if running in parallel mode.
		testConfig = smoke.GetConfig()
		appsDomain = testConfig.GetAppsDomains()

		if testConfig.EnableIsolationSegmentTests != true {
			Skip("Skipping because EnableIsolationSegmentTests flag is set to false")
		}
		orgName = testConfig.GetExistingOrganization()
		spaceName = testConfig.GetExistingSpace()
		spaceGuid = GetSpaceGuidFromName(spaceName)

		isoSpaceName = testConfig.GetExistingIsoSpace()
		isoSpaceGuid = GetSpaceGuidFromName(isoSpaceName)

		isoSegName = testConfig.GetIsolationSegmentName()
		isoSegDomain = testConfig.GetIsolationSegmentDomain()
		Expect(IsolationSegmentExists(isoSegName)).To(BeTrue())

		session := cf.Cf("curl", fmt.Sprintf("/v3/organizations?names=%s", orgName))
		bytes := session.Wait(testConfig.GetDefaultTimeout()).Out.Contents()
		orgGuid = GetGuidFromResponse(bytes)
	})

	Context("When an app is pushed to a space assigned the shared isolation segment", func() {
		var appName string

		BeforeEach(func() {
			Eventually(cf.Cf("t", "-o", testConfig.GetExistingOrganization(), "-s", testConfig.GetExistingSpace()),
				testConfig.GetPushTimeout()).Should(Exit(0))

			appName = generator.PrefixedRandomName("SMOKES", "APP")
			Eventually(cf.Cf(
				"push", appName,
				"-p", BINARY_APP_BITS_PATH,
				"--no-start",
				"-b", "binary_buildpack",
				"-d", appsDomain,
				"-c", "./app"),
				testConfig.GetPushTimeout()).Should(Exit(0))
			smoke.SetBackend(appName)
			Eventually(cf.Cf("start", appName), testConfig.GetDefaultTimeout()).Should(Exit(0))
		})

		AfterEach(func() {
			smoke.AppReport(appName, CF_TIMEOUT_IN_SECONDS)
			if testConfig.Cleanup {
				Eventually(cf.Cf("delete", "-f", "-r", appName), testConfig.GetDefaultTimeout()).Should(Exit(0))
			}
		})

		It("is reachable from the shared router", func() {
			resp := SendRequestWithSpoofedHeader(fmt.Sprintf("%s.%s", appName, appsDomain), appsDomain)
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(200))
			htmlData, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(htmlData)).To(ContainSubstring(binaryHi))
		})

		It("is not reachable from the isolation segment router", func() {
			//send a request to app in the shared domain, but through the isolation segment router
			resp := SendRequestWithSpoofedHeader(fmt.Sprintf("%s.%s", appName, appsDomain), isoSegDomain)
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(404))
		})
	})

	Context("When an app is pushed to a space that has been assigned an Isolation Segment", func() {
		var appName string

		BeforeEach(func() {
			Eventually(cf.Cf("t", "-o", testConfig.GetExistingOrganization(), "-s", testConfig.GetExistingIsoSpace()),
				testConfig.GetPushTimeout()).Should(Exit(0))

			appName = generator.PrefixedRandomName("SMOKES", "APP")
			Eventually(cf.Cf(
				"push", appName,
				"-p", BINARY_APP_BITS_PATH,
				"--no-start",
				"-b", "binary_buildpack",
				"-d", isoSegDomain,
				"-c", "./app"),
				testConfig.GetPushTimeout()).Should(Exit(0))
			smoke.SetBackend(appName)
			Eventually(cf.Cf("start", appName), testConfig.GetDefaultTimeout()).Should(Exit(0))
		})

		AfterEach(func() {
			smoke.AppReport(appName, CF_TIMEOUT_IN_SECONDS)
			if testConfig.Cleanup {
				Eventually(cf.Cf("delete", "-f", "-r", appName), testConfig.GetDefaultTimeout()).Should(Exit(0))
			}
		})
		It("the app is reachable from the isolated router", func() {
			resp := SendRequestWithSpoofedHeader(fmt.Sprintf("%s.%s", appName, isoSegDomain), isoSegDomain)
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(200))
			htmlData, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(htmlData)).To(ContainSubstring(binaryHi))
		})

		It("the app is not reachable from the shared router", func() {
			resp := SendRequestWithSpoofedHeader(fmt.Sprintf("%s.%s", appName, isoSegDomain), appsDomain)
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(404))
		})
	})
})
