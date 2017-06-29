package isolation_segments

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

const (
	V3_DEFAULT_MEMORY_LIMIT = "256"
	V3_JAVA_MEMORY_LIMIT    = "512"
	CF_TIMEOUT_IN_SECONDS   = 30
)

func AssignIsolationSegmentToSpace(spaceGuid, isoSegGuid string) {
	Eventually(cf.Cf("curl", fmt.Sprintf("/v3/spaces/%s/relationships/isolation_segment", spaceGuid),
		"-X",
		"PATCH",
		"-d",
		fmt.Sprintf(`{"data":{"guid":"%s"}}`, isoSegGuid)),
		CF_TIMEOUT_IN_SECONDS).Should(Exit(0))
}

func CreateIsolationSegment(name string) string {
	session := cf.Cf("curl", "/v3/isolation_segments", "-X", "POST", "-d", fmt.Sprintf(`{"name":"%s"}`, name))
	bytes := session.Wait(CF_TIMEOUT_IN_SECONDS).Out.Contents()

	var isolation_segment struct {
		Guid string `json:"guid"`
	}
	err := json.Unmarshal(bytes, &isolation_segment)
	Expect(err).ToNot(HaveOccurred())

	return isolation_segment.Guid
}

func EntitleOrgToIsolationSegment(orgGuid, isoSegGuid string) {
	Eventually(cf.Cf("curl",
		fmt.Sprintf("/v3/isolation_segments/%s/relationships/organizations", isoSegGuid),
		"-X",
		"POST",
		"-d",
		fmt.Sprintf(`{"data":[{ "guid":"%s" }]}`, orgGuid)),
		CF_TIMEOUT_IN_SECONDS).Should(Exit(0))
}

func GetGuidFromResponse(response []byte) string {
	type resource struct {
		Guid string `json:"guid"`
	}
	var GetResponse struct {
		Resources []resource `json:"resources"`
	}

	err := json.Unmarshal(response, &GetResponse)
	Expect(err).ToNot(HaveOccurred())

	if len(GetResponse.Resources) == 0 {
		Fail("No guid found for response")
	}

	return GetResponse.Resources[0].Guid
}

func GetIsolationSegmentGuid(name string) string {
	session := cf.Cf("curl", fmt.Sprintf("/v3/isolation_segments?names=%s", name))
	bytes := session.Wait(CF_TIMEOUT_IN_SECONDS).Out.Contents()
	return GetGuidFromResponse(bytes)
}

func GetIsolationSegmentGuidFromResponse(response []byte) string {
	type data struct {
		Guid string `json:"guid"`
	}
	var GetResponse struct {
		Data data `json:"data"`
	}

	err := json.Unmarshal(response, &GetResponse)
	Expect(err).ToNot(HaveOccurred())

	if (data{}) == GetResponse.Data {
		return ""
	}

	return GetResponse.Data.Guid
}

func IsolationSegmentExists(name string) bool {
	session := cf.Cf("curl", fmt.Sprintf("/v3/isolation_segments?names=%s", name))
	bytes := session.Wait(CF_TIMEOUT_IN_SECONDS).Out.Contents()
	type resource struct {
		Guid string `json:"guid"`
	}
	var GetResponse struct {
		Resources []resource `json:"resources"`
	}

	err := json.Unmarshal(bytes, &GetResponse)
	Expect(err).ToNot(HaveOccurred())
	return len(GetResponse.Resources) > 0
}

func OrgEntitledToIsolationSegment(orgGuid string, isoSegName string) bool {
	session := cf.Cf("curl", fmt.Sprintf("/v3/isolation_segments?names=%s&organization_guids=%s", isoSegName, orgGuid))
	bytes := session.Wait(CF_TIMEOUT_IN_SECONDS).Out.Contents()

	type resource struct {
		Guid string `json:"guid"`
	}
	var GetResponse struct {
		Resources []resource `json:"resources"`
	}

	err := json.Unmarshal(bytes, &GetResponse)
	Expect(err).ToNot(HaveOccurred())
	return len(GetResponse.Resources) > 0
}

func RevokeOrgEntitlementForIsolationSegment(orgGuid, isoSegGuid string) {
	Eventually(cf.Cf("curl",
		fmt.Sprintf("/v3/isolation_segments/%s/relationships/organizations/%s", isoSegGuid, orgGuid),
		"-X",
		"DELETE",
	), CF_TIMEOUT_IN_SECONDS).Should(Exit(0))
}

func SendRequestWithSpoofedHeader(host, domain string) *http.Response {
	req, _ := http.NewRequest("GET", fmt.Sprintf("http://wildcard-path.%s", domain), nil)
	req.Host = host

	resp, err := http.DefaultClient.Do(req)
	Expect(err).NotTo(HaveOccurred())
	return resp
}

func GetSpaceGuidFromName(spaceName string) string {
	session := cf.Cf("space", spaceName, "--guid")
	bytes := session.Wait(CF_TIMEOUT_IN_SECONDS).Out.Contents()
	return strings.TrimSpace(string(bytes))
}

func UnassignIsolationSegmentFromSpace(spaceGuid string) {
	Eventually(cf.Cf("curl", fmt.Sprintf("/v3/spaces/%s/relationships/isolation_segment", spaceGuid),
		"-X",
		"PATCH",
		"-d",
		`{"data":{"guid":null}}`),
		CF_TIMEOUT_IN_SECONDS).Should(Exit(0))
}
