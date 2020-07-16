package isolation_segments

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

const sharedIsolationSegmentGUID = "933b4c58-120b-499a-b85d-4b6fc9e2903b"

func AssignIsolationSegmentToSpace(spaceGUID, isoSegGUID string, timeout time.Duration) {
	Eventually(cf.Cf("curl", fmt.Sprintf("/v3/spaces/%s/relationships/isolation_segment", spaceGUID),
		"-X",
		"PATCH",
		"-d",
		fmt.Sprintf(`{"data":{"guid":"%s"}}`, isoSegGUID)),
		timeout).Should(Exit(0))
}

func EntitleOrgToIsolationSegment(orgGUID, isoSegGUID string, timeout time.Duration) {
	Eventually(cf.Cf("curl",
		fmt.Sprintf("/v3/isolation_segments/%s/relationships/organizations", isoSegGUID),
		"-X",
		"POST",
		"-d",
		fmt.Sprintf(`{"data":[{ "guid":"%s" }]}`, orgGUID)),
		timeout).Should(Exit(0))
}

func ResetSpaceIsolationSegment(spaceName, isoSegName string, timeout time.Duration) {
	Eventually(cf.Cf("reset-space-isolation-segment", spaceName), timeout).Should(Exit(0))
}

func DisableOrgIsolationSegment(orgName, isoSegName string, timeout time.Duration) {
	Eventually(cf.Cf("disable-org-isolation", orgName, isoSegName), timeout).Should(Exit(0))
}

func GetGUIDFromResponse(response []byte) string {
	type resource struct {
		GUID string `json:"guid"`
	}
	var GetResponse struct {
		Resources []resource `json:"resources"`
	}

	err := json.Unmarshal(response, &GetResponse)
	Expect(err).ToNot(HaveOccurred())

	if len(GetResponse.Resources) == 0 {
		Fail("No guid found for response")
	}

	return GetResponse.Resources[0].GUID
}

func CreateOrGetIsolationSegment(name string, timeout time.Duration) {
	isolationSegments := cf.Cf("isolation-segments").Wait(timeout)
	isolationSegmentsOutput := string(isolationSegments.Out.Contents())

	if !strings.Contains(isolationSegmentsOutput, name) {
		session := cf.Cf("create-isolation-segment", name)
		Eventually(session, timeout).Should(Exit(0))
	}
}

func GetIsolationSegmentGUID(name string, timeout time.Duration) string {
	session := cf.Cf("curl", fmt.Sprintf("/v3/isolation_segments?names=%s", name))
	bytes := session.Wait(timeout).Out.Contents()
	return GetGUIDFromResponse(bytes)
}

func DeleteIsolationSegment(name string, timeout time.Duration) {
	session := cf.Cf("delete-isolation-segment", name, "-f")
	Eventually(session, timeout).Should(Exit(0))
}

func OrgEntitledToIsolationSegment(orgGUID string, isoSegName string, timeout time.Duration) bool {
	session := cf.Cf("curl", fmt.Sprintf("/v3/isolation_segments?names=%s&organization_guids=%s", isoSegName, orgGUID))
	bytes := session.Wait(timeout).Out.Contents()

	type resource struct {
		GUID string `json:"guid"`
	}
	var GetResponse struct {
		Resources []resource `json:"resources"`
	}

	err := json.Unmarshal(bytes, &GetResponse)
	Expect(err).ToNot(HaveOccurred())
	return len(GetResponse.Resources) > 0
}

func IsolationSegmentAssignedToSpace(spaceGUID string, timeout time.Duration) bool {
	session := cf.Cf("curl", fmt.Sprintf("/v2/spaces/%s", spaceGUID))
	response := session.Wait(timeout).Out.Contents()
	type entity struct {
		GUID string `json:"isolation_segment_guid"`
	}
	var SpaceResponse struct {
		Entity entity `json:"entity"`
	}

	err := json.Unmarshal(response, &SpaceResponse)
	Expect(err).ToNot(HaveOccurred())

	return SpaceResponse.Entity.GUID != ""
}

func SendRequestWithSpoofedHeader(host, domain string, skipSslValidation bool) *http.Response {
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://wildcard-path.%s", domain), nil)
	req.Host = host

	transport := &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipSslValidation},
	}
	client := &http.Client{Transport: transport}
	resp, err := client.Do(req)
	Expect(err).NotTo(HaveOccurred())
	return resp
}

func GetSpaceGUIDFromName(spaceName string, timeout time.Duration) string {
	session := cf.Cf("space", spaceName, "--guid")
	bytes := session.Wait(timeout).Out.Contents()
	return strings.TrimSpace(string(bytes))
}

func GetOrgGUIDFromName(orgName string, timeout time.Duration) string {
	session := cf.Cf("org", orgName, "--guid")
	bytes := session.Wait(timeout).Out.Contents()
	return strings.TrimSpace(string(bytes))
}

func orgDefaultIsolationSegmentIsShared(orgGuid string, timeout time.Duration) bool {
	defaultIsolationSegmentsResponse := cf.Cf("curl", fmt.Sprintf("/v3/organizations/%s/relationships/default_isolation_segment", orgGuid)).Wait(timeout)

	var response struct {
		Data *struct {
			GUID string
		}
	}

	err := json.Unmarshal(defaultIsolationSegmentsResponse.Out.Contents(), &response)
	Expect(err).ToNot(HaveOccurred())

	if response.Data == nil {
		return true
	}

	if response.Data.GUID == sharedIsolationSegmentGUID {
		return true
	}

	return false
}
