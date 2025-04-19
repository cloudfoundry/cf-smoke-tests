package runtime

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/cf-smoke-tests/smoke"
	"github.com/cloudfoundry/cf-test-helpers/v2/workflowhelpers"

	. "github.com/onsi/ginkgo/v2"
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

	_, rc := GinkgoConfiguration()

	if testConfig.ArtifactsDirectory != "" {
		err := os.Setenv("CF_TRACE", traceLogFilePath(testConfig))
		Expect(err).ToNot(HaveOccurred())
		rc.JUnitReport = jUnitReportFilePath(testConfig)
	}

	RunSpecs(t, "CF-Runtime-Smoke-Tests", rc)
}

func traceLogFilePath(testConfig *smoke.Config) string {
	return filepath.Join(testConfig.ArtifactsDirectory, fmt.Sprintf("CF-TRACE-%s-%d.txt", testConfig.SuiteName, GinkgoParallelProcess()))
}

func jUnitReportFilePath(testConfig *smoke.Config) string {
	return filepath.Join(testConfig.ArtifactsDirectory, fmt.Sprintf("junit-%s-%d.xml", testConfig.SuiteName, GinkgoParallelProcess()))
}
