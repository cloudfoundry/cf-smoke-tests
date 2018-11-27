package isolation_segments

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/cloudfoundry-incubator/cf-test-helpers/workflowhelpers"
	"github.com/cloudfoundry/cf-smoke-tests/smoke"
	ginkgoconfig "github.com/onsi/ginkgo/config"
	"github.com/onsi/ginkgo/reporters"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	testSetup  *workflowhelpers.ReproducibleTestSuiteSetup
	testConfig *smoke.Config
)

func TestSmokeTests(t *testing.T) {
	RegisterFailHandler(Fail)

	testConfig = smoke.GetConfig()
	testSetup = workflowhelpers.NewSmokeTestSuiteSetup(testConfig)
	rs := []Reporter{}

	if testConfig.ArtifactsDirectory != "" {
		os.Setenv("CF_TRACE", traceLogFilePath(testConfig))
		rs = append(rs, reporters.NewJUnitReporter(jUnitReportFilePath(testConfig)))
	}

	if testConfig.Reporter == "TeamCity" {
		rs = append(rs, reporters.NewTeamCityReporter(GinkgoWriter))
	}

	RunSpecsWithDefaultAndCustomReporters(t, "CF-Isolation-Segment-Smoke-Tests", rs)
}

var _ = SynchronizedBeforeSuite(func() []byte {
	testSetup.Setup()
	return nil
}, func(data []byte) {})

var _ = SynchronizedAfterSuite(func() {
}, func() {
	testSetup.Teardown()
})

func traceLogFilePath(testConfig *smoke.Config) string {
	return filepath.Join(testConfig.ArtifactsDirectory, fmt.Sprintf("CF-TRACE-%s-%d.txt", testConfig.SuiteName, ginkgoNode()))
}

func jUnitReportFilePath(testConfig *smoke.Config) string {
	return filepath.Join(testConfig.ArtifactsDirectory, fmt.Sprintf("junit-%s-%d.xml", testConfig.SuiteName, ginkgoNode()))
}

func ginkgoNode() int {
	return ginkgoconfig.GinkgoConfig.ParallelNode
}
