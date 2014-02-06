package cf_health_checks_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf-experimental/cf-test-helpers/runner"
	"github.com/vito/cmdtest"

	"testing"
)

func TestCfHealthChecks(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cf-Health-Checks Suite")
}

// var IntegrationConfig = config.Load()

var AppName = ""

var rubyAppPath = "./apps/ruby/simple"
var javaAppPath = "./apps/java/JavaTinyApp-1.1.war"

var AppPath = rubyAppPath

func AppUri(endpoint string) string {
	// TODO: Pull in IntegrationConfig from cf-acceptance-tests
	var AppsDomain = "10.244.0.34.xip.io"
	return "http://" + AppName + "." + AppsDomain + endpoint
}

func Curling(endpoint string) func() *cmdtest.Session {
	return func() *cmdtest.Session {
		return runner.Curl(AppUri(endpoint))
	}
}
