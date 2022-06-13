package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/cf-smoke-tests/smoke"
	"github.com/cloudfoundry/cf-test-helpers/workflowhelpers"

	. "github.com/onsi/ginkgo/v2"
	ginkgoconfig "github.com/onsi/ginkgo/v2/config"
	"github.com/onsi/ginkgo/v2/reporters"
	. "github.com/onsi/gomega"
)

func TestSmokeTests(t *testing.T) {
	RegisterFailHandler(Fail)

	testConfig := smoke.GetConfig()
	testSetup := workflowhelpers.NewSmokeTestSuiteSetup(testConfig)

	SynchronizedBeforeSuite(func() []byte {
		return nil
	}, func(data []byte) {
		testSetup.Setup()
	})

	SynchronizedAfterSuite(func() {
		testSetup.Teardown()
	}, func() {})

	rs := []Reporter{}

	if testConfig.ArtifactsDirectory != "" {
		os.Setenv("CF_TRACE", traceLogFilePath(testConfig))
		rs = append(rs, reporters.NewJUnitReporter(jUnitReportFilePath(testConfig)))
	}

	if testConfig.Reporter == "TeamCity" {
		rs = append(rs, reporters.NewTeamCityReporter(GinkgoWriter))
	}

	RunSpecsWithDefaultAndCustomReporters(t, "CF-Logging-Smoke-Tests", rs)
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
