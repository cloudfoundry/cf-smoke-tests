package smoke

import (
	"fmt"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

const (
	SimpleBinaryAppBitsPath = "../../assets/binary"
	SimpleDotnetAppBitsPath = "../../assets/dotnet_simple/Published"
)

func SkipIfNotWindows(testConfig *Config) {
	if !testConfig.EnableWindowsTests {
		Skip("Windows tests are disabled")
	}
}

func AppReport(appName string, timeout time.Duration) {
	Eventually(cf.Cf("app", appName, "--guid"), timeout).Should(Exit())
	Eventually(cf.Cf("logs", appName, "--recent"), timeout).Should(Exit())
}

func Logs(useLogCache bool, appName string) *Session {
	if useLogCache {
		return cf.Cf("tail", appName, "--lines", "125")
	}
	return cf.Cf("logs", "--recent", appName)
}

func Cf(args ...string) *Session {

	session := cf.Cf("version").Wait(GetConfig().GetDefaultTimeout())
	cfCliVersion := session.Buffer().Contents()

	if strings.HasPrefix(string(cfCliVersion), "cf version 6") {
		return cf.Cf(args...)
	} else if strings.HasPrefix(string(cfCliVersion), "cf version 7") {
		return cf.Cf(removeDomainParam(args...)...)
	}

	panic(fmt.Sprintf("Unsupported cf cli version: %s", cfCliVersion))
}

func removeDomainParam(args ...string) []string {
	var v7Args []string
	for i := 0; i < len(args); i++ {
		if args[i] == "-d" {
			i++
		} else {
			v7Args = append(v7Args, args[i])
		}
	}
	return v7Args
}
