package runtime

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	ginkgoconfig "github.com/onsi/ginkgo/config"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	"github.com/cloudfoundry-incubator/cf-test-helpers/workflowhelpers"
	"github.com/cloudfoundry/cf-smoke-tests/smoke"
)

const (
	SIMPLE_RUBY_APP_BITS_PATH = "../../assets/ruby_simple"

	SIMPLE_DOTNET_APP_BITS_PATH = "../../assets/dotnet_simple/Published"

	CF_API_TIMEOUT = 1 * time.Minute

	// timeout for most cf cli calls
	CF_TIMEOUT_IN_SECONDS = 30

	// timeout for cf push cli calls
	CF_PUSH_TIMEOUT_IN_SECONDS = 300

	// timeout for cf scale cli calls
	CF_SCALE_TIMEOUT_IN_SECONDS = 120

	// timeout for cf app cli calls
	CF_APP_STATUS_TIMEOUT_IN_SECONDS = 120
)

func TestSmokeTests(t *testing.T) {
	testConfig := smoke.GetConfig()

	testUserContext := workflowhelpers.NewUserContext(
		testConfig.ApiEndpoint,
		testConfig.User,
		testConfig.Password,
		testConfig.Org,
		testConfig.Space,
		testConfig.SkipSSLValidation,
	)

	RegisterFailHandler(Fail)

	var originalCfHomeDir, currentCfHomeDir string

	SynchronizedBeforeSuite(func() []byte {
		originalCfHomeDir, currentCfHomeDir = workflowhelpers.InitiateUserContext(testUserContext, CF_API_TIMEOUT)

		if !testConfig.UseExistingOrg {
			Expect(cf.Cf("create-quota", quotaName(testConfig.Org), "-m", "10G", "-r", "10", "-s", "2").Wait(CF_TIMEOUT_IN_SECONDS)).To(Exit(0))
			Expect(cf.Cf("create-org", testConfig.Org).Wait(CF_TIMEOUT_IN_SECONDS)).To(Exit(0))
			Expect(cf.Cf("set-quota", testConfig.Org, quotaName(testConfig.Org)).Wait(CF_TIMEOUT_IN_SECONDS)).To(Exit(0))
		}

		Expect(cf.Cf("target", "-o", testConfig.Org).Wait(CF_TIMEOUT_IN_SECONDS)).To(Exit(0))

		if !testConfig.UseExistingSpace {
			Expect(cf.Cf("create-space", "-o", testConfig.Org, testConfig.Space).Wait(CF_TIMEOUT_IN_SECONDS)).To(Exit(0))
		}

		Expect(cf.Cf("target", "-s", testConfig.Space).Wait(CF_TIMEOUT_IN_SECONDS)).To(Exit(0))
		return []byte("")
	}, func(data []byte) {})

	SynchronizedAfterSuite(func() {
	}, func() {
		smoke.TestResourcesSummary(testConfig)

		if testConfig.Cleanup && !testConfig.UseExistingSpace {
			Expect(cf.Cf("delete-space", testConfig.Space, "-f").Wait(CF_TIMEOUT_IN_SECONDS)).To(Exit(0))
		}

		if testConfig.Cleanup && !testConfig.UseExistingOrg {
			Expect(cf.Cf("delete-org", testConfig.Org, "-f").Wait(CF_TIMEOUT_IN_SECONDS)).To(Exit(0))
			Expect(cf.Cf("delete-quota", quotaName(testConfig.Org), "-f").Wait(CF_TIMEOUT_IN_SECONDS)).To(Exit(0))
		}

		workflowhelpers.RestoreUserContext(testUserContext, CF_API_TIMEOUT, originalCfHomeDir, currentCfHomeDir)
	})

	rs := []Reporter{}

	if testConfig.ArtifactsDirectory != "" {
		os.Setenv("CF_TRACE", traceLogFilePath(testConfig))
		rs = append(rs, reporters.NewJUnitReporter(jUnitReportFilePath(testConfig)))
	}

	RunSpecsWithDefaultAndCustomReporters(t, "CF-Runtime-Smoke-Tests", rs)
}

func traceLogFilePath(testConfig *smoke.Config) string {
	return filepath.Join(testConfig.ArtifactsDirectory, fmt.Sprintf("CF-TRACE-%s-%d.txt", testConfig.SuiteName, ginkgoNode()))
}

func jUnitReportFilePath(testConfig *smoke.Config) string {
	return filepath.Join(testConfig.ArtifactsDirectory, fmt.Sprintf("junit-%s-%d.xml", testConfig.SuiteName, ginkgoNode()))
}

func ginkgoNode() int {
	return ginkgoconfig.GinkgoConfig.ParallelNode
}

func quotaName(prefix string) string {
	return prefix + "_QUOTA"
}
