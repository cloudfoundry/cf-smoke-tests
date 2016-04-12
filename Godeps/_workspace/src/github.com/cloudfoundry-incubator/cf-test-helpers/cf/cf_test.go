package cf_test

import (
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"

	. "github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	"github.com/cloudfoundry-incubator/cf-test-helpers/runner"
)

var _ = Describe("Cf", func() {
	It("sends the request to current CF target", func() {
		runner.CommandInterceptor = func(cmd *exec.Cmd) *exec.Cmd {
			Expect(cmd.Path).To(Equal(exec.Command("cf").Path))
			Expect(cmd.Args).To(Equal([]string{"cf", "apps"}))

			return exec.Command("bash", "-c", `exit 42`)
		}

		Eventually(Cf("apps"), 1*time.Second).Should(Exit(42))
	})
})
