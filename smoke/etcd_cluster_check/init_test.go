package logging

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

	"github.com/cloudfoundry-incubator/cf-test-helpers/workflowhelpers"
	"github.com/cloudfoundry/cf-smoke-tests/smoke"
)

const (
	CURLER_RUBY_APP_BITS_PATH = "../../assets/curler"

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
	var testSetup *workflowhelpers.ReproducibleTestSuiteSetup
	testConfig := smoke.GetConfig()

	RegisterFailHandler(Fail)

	SynchronizedBeforeSuite(func() []byte {
		return nil
	}, func(data []byte) {
		testSetup = workflowhelpers.NewSmokeTestSuiteSetup(testConfig)
		testSetup.Setup()
	})

	AfterSuite(func() {
		testSetup.Teardown()
	})

	rs := []Reporter{}

	if testConfig.ArtifactsDirectory != "" {
		os.Setenv("CF_TRACE", traceLogFilePath(testConfig))
		rs = append(rs, reporters.NewJUnitReporter(jUnitReportFilePath(testConfig)))
	}

	RunSpecsWithDefaultAndCustomReporters(t, "CF-EtcdClusterCheck-Smoke-Tests", rs)
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
