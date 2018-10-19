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
	"github.com/cloudfoundry-incubator/cf-test-helpers/workflowhelpers"
)

const (
	SHARED_ISOLATION_SEGMENT_GUID = "933b4c58-120b-499a-b85d-4b6fc9e2903b"
	binaryHi                      = "Hello from a binary"
	BINARY_APP_BITS_PATH          = "../../assets/binary"
)

var _ = Describe("RoutingIsolationSegments", func() {
	var appsDomain string
	var orgGuid, orgName string
	var spaceGuid, spaceName string
	var isoSpaceGuid, isoSpaceName string
	var isoSegGuid string
	var isoSegName, isoSegDomain string
	var testSetup *workflowhelpers.ReproducibleTestSuiteSetup
	var testConfig *smoke.Config
	var appName string

	BeforeEach(func() {
		// New up a organization since we will be assigning isolation segments.
		// This has a potential to cause other tests to fail if running in parallel mode.
		testConfig = smoke.GetConfig()
		if testConfig.EnableIsolationSegmentTests != true {
			Skip("Skipping because EnableIsolationSegmentTests flag is set to false")
		}
		testSetup = workflowhelpers.NewSmokeTestSuiteSetup(testConfig)
		testSetup.Setup()

		appsDomain = testConfig.GetAppsDomains()
		orgName = testSetup.RegularUserContext().Org
		orgGuid = GetOrgGuidFromName(orgName)
		spaceName = testSetup.RegularUserContext().Space
		spaceGuid = GetSpaceGuidFromName(spaceName)
		isoSpaceName = spaceName
		isoSpaceGuid = spaceGuid
		appName = generator.PrefixedRandomName("SMOKES", "APP")

		isoSegName = testConfig.GetIsolationSegmentName()
		isoSegDomain = testConfig.GetIsolationSegmentDomain()

		if testConfig.GetUseExistingOrganization() && testConfig.GetUseExistingSpace() {
			if !OrgEntitledToIsolationSegment(orgGuid, isoSegName) {
				Fail(fmt.Sprintf("Pre-existing org %s is not entitled to isolation segment %s", orgName, isoSegName))
			}
			isoSpaceName = testConfig.GetIsolationSegmentSpace()
			isoSpaceGuid = GetSpaceGuidFromName(isoSpaceName)
			if !IsolationSegmentAssignedToSpace(isoSegName, isoSpaceGuid) {
				Fail(fmt.Sprintf("No isolation segment assigned  to pre-existing space %s", isoSpaceName))
			}
		}

		session := cf.Cf("curl", fmt.Sprintf("/v3/organizations?names=%s", orgName))
		bytes := session.Wait(testConfig.GetDefaultTimeout()).Out.Contents()
		orgGuid = GetGuidFromResponse(bytes)
	})

	AfterEach(func() {
		if testConfig.Cleanup {
			Expect(cf.Cf("delete", appName, "-f", "-r").Wait(CF_TIMEOUT_IN_SECONDS)).To(Exit(0))
		}
		testSetup.Teardown()
	})

	Context("When an app is pushed to a space assigned the shared isolation segment", func() {
		BeforeEach(func() {
			if !testConfig.GetUseExistingOrganization() && !testConfig.GetUseExistingSpace() {
				workflowhelpers.AsUser(testSetup.AdminUserContext(), testSetup.ShortTimeout(), func() {
					EntitleOrgToIsolationSegment(orgGuid, SHARED_ISOLATION_SEGMENT_GUID)
					AssignIsolationSegmentToSpace(spaceGuid, SHARED_ISOLATION_SEGMENT_GUID)
				})
			}
			Eventually(cf.Cf(
				"push", appName,
				"-p", BINARY_APP_BITS_PATH,
				"-b", "binary_buildpack",
				"-d", appsDomain,
				"-c", "./app"),
				testConfig.GetPushTimeout()).Should(Exit(0))
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
			isoSegGuid = GetIsolationSegmentGuid(isoSegName)
			if !testConfig.GetUseExistingOrganization() {
				EntitleOrgToIsolationSegment(orgGuid, isoSegGuid)
			}

			if !testConfig.GetUseExistingSpace() {
				AssignIsolationSegmentToSpace(isoSpaceGuid, isoSegGuid)
			}
			appName = generator.PrefixedRandomName("SMOKES", "APP")
			Eventually(cf.Cf("target", "-s", isoSpaceName), testConfig.GetDefaultTimeout()).Should(Exit(0))
			Eventually(cf.Cf(
				"push", appName,
				"-p", BINARY_APP_BITS_PATH,
				"-b", "binary_buildpack",
				"-d", isoSegDomain,
				"-c", "./app"),
				testConfig.GetPushTimeout()).Should(Exit(0))
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
