package cf_test

import (
	"bytes"
	"os/exec"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	"github.com/cloudfoundry-incubator/cf-test-helpers/runner"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("CfAuth", func() {
	var callerOutput *bytes.Buffer
	var password string

	BeforeEach(func() {
		callerOutput = bytes.NewBuffer([]byte{})
		password = "superSecretPassword"

		GinkgoWriter = callerOutput
	})

	It("runs the cf auth command", func() {
		user := "myUser"

		runner.CommandInterceptor = func(cmd *exec.Cmd) *exec.Cmd {
			Expect(cmd.Path).To(Equal(exec.Command("cf").Path))
			Expect(cmd.Args).To(Equal([]string{
				"cf", "auth", user, password,
			}))

			return exec.Command("bash", "-c", "echo \"Authenticating...\nOK\"")
		}

		Eventually(cf.CfAuth(user, password)).Should(gbytes.Say("Authenticating...\nOK"))
	})

	It("does not expose the password", func() {
		user := "myUser"

		runner.CommandInterceptor = func(cmd *exec.Cmd) *exec.Cmd {
			Expect(cmd.Path).To(Equal(exec.Command("cf").Path))
			Expect(cmd.Args).To(Equal([]string{
				"cf", "auth", user, password,
			}))

			return exec.Command("bash", "-c", "echo \"Authenticating...\nOK\"")
		}

		cf.CfAuth(user, password).Wait()
		Expect(callerOutput.String()).NotTo(ContainSubstring(password))
		Expect(callerOutput.String()).NotTo(ContainSubstring("bash"))
		Expect(callerOutput.String()).NotTo(ContainSubstring("echo"))
		Expect(callerOutput.String()).To(ContainSubstring("REDACTED"))
	})
})
