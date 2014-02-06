package cf_health_checks_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf-experimental/cf-health-checks/config"
	"github.com/pivotal-cf-experimental/cf-test-helpers/runner"
	"github.com/vito/cmdtest"

	"testing"
)

func TestCfHealthChecks(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cf-Health-Checks Suite")
}

var IntegrationConfig = config.Load()

var AppName = ""

var AppPath = "./apps/ruby/simple"

func AppUri(endpoint string) string {
	return "http://" + AppName + "." + IntegrationConfig.AppsDomain + endpoint
}

func Curling(endpoint string) func() *cmdtest.Session {
	return func() *cmdtest.Session {
		return runner.Curl(AppUri(endpoint))
	}
}
