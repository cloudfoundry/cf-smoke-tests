package isolation_segments

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	V3_DEFAULT_MEMORY_LIMIT = "256"
	V3_JAVA_MEMORY_LIMIT    = "512"
	CF_TIMEOUT_IN_SECONDS   = 30
)

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
