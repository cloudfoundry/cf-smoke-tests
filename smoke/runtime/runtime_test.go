package runtime

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/cloudfoundry/cf-smoke-tests/smoke"
	"github.com/cloudfoundry/cf-test-helpers/cf"
	"github.com/cloudfoundry/cf-test-helpers/generator"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Runtime:", func() {
	var testConfig = smoke.GetConfig()
	var appName string
	var appURL string
	var expectedNullResponse string
	var manifestPath string

	BeforeEach(func() {
		appName = testConfig.RuntimeApp
		if appName == "" {
			appName = generator.PrefixedRandomName("SMOKES", "APP")
		}

		appURL = "https://" + appName + "." + testConfig.AppsDomain

		Eventually(func() error {
			var err error
			expectedNullResponse, err = getBodySkipSSL(testConfig.SkipSSLValidation, appURL)
			return err
		}, testConfig.GetDefaultTimeout()).Should(BeNil())

		manifestPath = createManifestWithRoute(appName, testConfig.AppsDomain)
	})

	AfterEach(func() {
		defer func() {
			if testConfig.Cleanup {
				Expect(cf.Cf("delete", appName, "-f", "-r").Wait(testConfig.GetDefaultTimeout())).To(Exit(0))
			}
		}()
		smoke.AppReport(appName, testConfig.GetDefaultTimeout())
	})

	Context("linux apps", func() {
		It("can be pushed, scaled and deleted", func() {

			Expect(cf.Cf("push",
				appName,
				"-b", "binary_buildpack",
				"-m", "30M",
				"-k", "16M",
				"-f", manifestPath,
				"-p", smoke.SimpleBinaryAppBitsPath).Wait(testConfig.GetPushTimeout())).To(Exit(0))

			runPushTests(appName, appURL, expectedNullResponse, testConfig)
		})
	})

	Context("windows apps", func() {
		It("can be pushed, scaled and deleted", func() {
			smoke.SkipIfNotWindows(testConfig)

			Expect(cf.Cf("push",
				appName,
				"-f", manifestPath,
				"-p", smoke.SimpleDotnetAppBitsPath,
				"-s", testConfig.GetWindowsStack(),
				"-b", "hwc_buildpack").Wait(testConfig.GetPushTimeout())).To(Exit(0))

			runPushTests(appName, appURL, expectedNullResponse, testConfig)
		})
	})
})

func runPushTests(appName, appURL, expectedNullResponse string, testConfig *smoke.Config) {
	Eventually(func() (string, error) {
		return getBodySkipSSL(testConfig.SkipSSLValidation, appURL)
	}, testConfig.GetDefaultTimeout()).Should(ContainSubstring("It just needed to be restarted!"))

	instances := 2
	maxAttempts := 120

	ExpectAppToScale(appName, instances, testConfig.GetScaleTimeout())

	ExpectAllAppInstancesToStart(appName, instances, maxAttempts, testConfig.GetAppStatusTimeout())

	ExpectAllAppInstancesToBeReachable(appURL, instances, maxAttempts)

	if testConfig.Cleanup {
		Expect(cf.Cf("delete", appName, "-f", "-r").Wait(testConfig.GetDefaultTimeout())).To(Exit(0))

		Eventually(func() (string, error) {
			return getBodySkipSSL(testConfig.SkipSSLValidation, appURL)
		}, testConfig.GetDefaultTimeout()).Should(ContainSubstring(string(expectedNullResponse)))
	}
}

func ExpectAppToScale(appName string, instances int, timeout time.Duration) {
	Expect(cf.Cf("scale", appName, "-i", strconv.Itoa(instances)).Wait(timeout)).To(Exit(0))
}

// Gets app status (up to maxAttempts) until all instances are up
func ExpectAllAppInstancesToStart(appName string, instances int, maxAttempts int, timeout time.Duration) {
	var found bool
	expectedOutput := regexp.MustCompile(fmt.Sprintf(`instances:\s+%d/%d`, instances, instances))

	outputMatchers := make([]*regexp.Regexp, instances)
	for i := 0; i < instances; i++ {
		outputMatchers[i] = regexp.MustCompile(fmt.Sprintf(`#%d\s+running`, i))
	}

	for i := 0; i < maxAttempts; i++ {
		session := cf.Cf("app", appName)
		Expect(session.Wait(timeout)).To(Exit(0))

		output := string(session.Out.Contents())
		found = expectedOutput.MatchString(output)

		if found {
			for _, matcher := range outputMatchers {
				matches := matcher.FindStringSubmatch(output)
				if matches == nil {
					found = false
					break
				}
			}
		}

		if found {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	Expect(found).To(BeTrue(), fmt.Sprintf("Wanted to see '%s' (all instances running) in %d attempts, but didn't", expectedOutput, maxAttempts))
}

// Curls the appURL (up to maxAttempts) until all instances have been seen
func ExpectAllAppInstancesToBeReachable(appURL string, instances int, maxAttempts int) {
	matcher := regexp.MustCompile(`instance[ _]index["]{0,1}:[ ]{0,1}(\d+)`)

	branchesSeen := make([]bool, instances)
	var sawAll bool
	var testConfig = smoke.GetConfig()
	for i := 0; i < maxAttempts; i++ {
		var output string
		Eventually(func() error {
			var err error
			output, err = getBodySkipSSL(testConfig.SkipSSLValidation, appURL)
			return err
		}, testConfig.GetDefaultTimeout()).Should(BeNil())

		matches := matcher.FindStringSubmatch(output)
		if matches == nil {
			Fail("Expected app curl output to include an instance_index; got " + output)
		}
		indexString := matches[1]
		index, err := strconv.Atoi(indexString)
		if err != nil {
			Fail("Failed to parse instance index value " + indexString)
		}
		branchesSeen[index] = true

		if allTrue(branchesSeen) {
			sawAll = true
			break
		}

		time.Sleep(time.Duration(5000/maxAttempts) * time.Millisecond)
	}

	Expect(sawAll).To(BeTrue(), fmt.Sprintf("Expected to hit all %d app instances in %d attempts, but didn't", instances, maxAttempts))
}

func allTrue(bools []bool) bool {
	for _, curr := range bools {
		if !curr {
			return false
		}
	}
	return true
}

func getBodySkipSSL(skip bool, url string) (string, error) {
	transport := &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skip},
	}
	client := &http.Client{Transport: transport}
	resp, err := client.Get(url)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func createManifestWithRoute(name string, domain string) string {
	file, err := ioutil.TempFile(os.TempDir(), "runtime-manifest-*.yml")
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
