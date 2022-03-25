package isolation_segments

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cloudfoundry/cf-smoke-tests/smoke"
	"github.com/cloudfoundry/cf-test-helpers/cf"
	"github.com/cloudfoundry/cf-test-helpers/generator"
	"github.com/cloudfoundry/cf-test-helpers/workflowhelpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

const (
	binaryHi          = "It just needed to be restarted!"
	binaryAppBitsPath = "../../assets/binary"
)

var (
	testConfig *smoke.Config
	testSetup  *workflowhelpers.ReproducibleTestSuiteSetup
)

var _ = Describe("RoutingIsolationSegments", func() {
	var appsDomain string
	var orgGUID, orgName string
	var spaceName string
	var isoSpaceGUID, isoSpaceName string
	var isoSegGUID string
	var isoSegName, isoSegDomain string
	var appName string

	BeforeEach(func() {
		if testConfig.EnableIsolationSegmentTests != true {
			Skip("Skipping because EnableIsolationSegmentTests flag is set to false")
		}

		appsDomain = testConfig.GetAppsDomains()

		orgName = testSetup.RegularUserContext().Org
		spaceName = testSetup.RegularUserContext().Space
		orgGUID = GetOrgGUIDFromName(orgName, testConfig.GetDefaultTimeout())

		isoSpaceName = testSetup.RegularUserContext().Space
		isoSpaceGUID = GetSpaceGUIDFromName(isoSpaceName, testConfig.GetDefaultTimeout())

		appName = generator.PrefixedRandomName("SMOKES", "APP")

		isoSegName = testConfig.GetIsolationSegmentName()
		isoSegDomain = testConfig.GetIsolationSegmentDomain()

		if testConfig.GetUseExistingOrganization() && testConfig.GetUseExistingSpace() {
			if !OrgEntitledToIsolationSegment(orgGUID, isoSegName, testConfig.GetDefaultTimeout()) {
				Fail(fmt.Sprintf("Pre-existing org %s is not entitled to isolation segment %s", orgName, isoSegName))
			}
			isoSpaceName = testConfig.GetIsolationSegmentSpace()
			isoSpaceGUID = GetSpaceGUIDFromName(isoSpaceName, testConfig.GetDefaultTimeout())
			if !IsolationSegmentAssignedToSpace(isoSpaceGUID, testConfig.GetDefaultTimeout()) {
				Fail(fmt.Sprintf("No isolation segment assigned  to pre-existing space %s", isoSpaceName))
			}
		}

		session := cf.Cf("curl", fmt.Sprintf("/v3/organizations?names=%s", orgName))
		bytes := session.Wait(testConfig.GetDefaultTimeout()).Out.Contents()
		orgGUID = GetGUIDFromResponse(bytes)
	})

	AfterEach(func() {
		if testConfig.Cleanup {
			Expect(cf.Cf("delete", appName, "-f", "-r").Wait(testConfig.GetDefaultTimeout())).To(Exit(0))
		}
	})

	Context("When an app is pushed to a space that has been assigned the shared isolation segment", func() {
		BeforeEach(func() {
			if testConfig.GetUseExistingOrganization() {
				Expect(orgDefaultIsolationSegmentIsShared(orgGUID, testConfig.GetDefaultTimeout())).To(BeTrue(), "Org's default isolation segment is not the shared isolation segment")
			}

			if testConfig.GetUseExistingSpace() {
				spaceSession := cf.Cf("space", testConfig.GetExistingSpace()).Wait(testConfig.GetDefaultTimeout())
				Expect(spaceSession).NotTo(Say(testConfig.GetIsolationSegmentName()), "Space should be assigned to the shared isolation segment")
			}

			Eventually(cf.Cf(
				"push", appName,
				"-p", binaryAppBitsPath,
				"-b", "binary_buildpack",
				"-m", "30M",
				"-k", "16M",
				"-c", "./app"),
				testConfig.GetPushTimeout()).Should(Exit(0))
		})

		It("is reachable from the shared router", func() {
			resp := SendRequestWithSpoofedHeader(fmt.Sprintf("%s.%s", appName, appsDomain), appsDomain, testConfig.SkipSSLValidation)
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(200))
			htmlData, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(htmlData)).To(ContainSubstring(binaryHi))
		})

		It("is not reachable from the isolation segment router", func() {
			//send a request to app in the shared domain, but through the isolation segment router
			resp := SendRequestWithSpoofedHeader(fmt.Sprintf("%s.%s", appName, appsDomain), isoSegDomain, testConfig.SkipSSLValidation)
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(404))
		})
	})

	Context("When an app is pushed to a space that has been assigned an Isolation Segment", func() {
		BeforeEach(func() {
			manifestPath := CreateManifestWithRoute(appName, isoSegDomain)

			CreateOrGetIsolationSegment(isoSegName, testConfig.GetDefaultTimeout())
			isoSegGUID = GetIsolationSegmentGUID(isoSegName, testConfig.GetDefaultTimeout())
			if !testConfig.GetUseExistingOrganization() {
				EntitleOrgToIsolationSegment(orgGUID, isoSegGUID, testConfig.GetDefaultTimeout())
			}

			if !testConfig.GetUseExistingSpace() {
				AssignIsolationSegmentToSpace(isoSpaceGUID, isoSegGUID, testConfig.GetDefaultTimeout())
			}
			Eventually(cf.Cf("target", "-s", isoSpaceName), testConfig.GetDefaultTimeout()).Should(Exit(0))
			Eventually(cf.Cf(
				"push", appName,
				"-p", binaryAppBitsPath,
				"-b", "binary_buildpack",
				"-m", "30M",
				"-k", "16M",
				"-f", manifestPath,
				"-c", "./app"),
				testConfig.GetPushTimeout()).Should(Exit(0))
		})

		AfterEach(func() {
			if !testConfig.GetUseExistingSpace() {
				ResetSpaceIsolationSegment(spaceName, isoSegName, testConfig.GetDefaultTimeout())
			}
			if !testConfig.GetUseExistingOrganization() {
				DisableOrgIsolationSegment(orgName, isoSegName, testConfig.GetDefaultTimeout())
			}
		})

		It("the app is reachable from the isolated router", func() {
			resp := SendRequestWithSpoofedHeader(fmt.Sprintf("%s.%s", appName, isoSegDomain), isoSegDomain, testConfig.SkipSSLValidation)
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(200))
			htmlData, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(htmlData)).To(ContainSubstring(binaryHi))
		})

		It("the app is not reachable from the shared router", func() {
			resp := SendRequestWithSpoofedHeader(fmt.Sprintf("%s.%s", appName, isoSegDomain), appsDomain, testConfig.SkipSSLValidation)
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(404))
		})
	})
})

func CreateManifestWithRoute(name string, domain string) string {
	file, err := ioutil.TempFile(os.TempDir(), "iso-segment-manifest-*.yml")
	Expect(err).NotTo(HaveOccurred())

	filePath := file.Name()

	_, err = file.Write([]byte(fmt.Sprintf("---\n"+
		"applications:\n"+
		"- name: %s\n"+
		"  routes:\n"+
		"  - route: %s.%s",
		name, name, domain)))
	Expect(err).NotTo(HaveOccurred())

	err = file.Close()
	Expect(err).NotTo(HaveOccurred())

	return filePath
}
