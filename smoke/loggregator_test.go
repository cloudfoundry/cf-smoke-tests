package smoke

import (
	"fmt"
	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	"github.com/cloudfoundry-incubator/cf-test-helpers/generator"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
)

var _ = Describe("Loggregator:", func() {
	var testConfig = GetConfig()
	var useExistingApp = (testConfig.LoggingApp != "")
	var appName string

	Describe("cf logs", func() {
		BeforeEach(func() {
			appName = testConfig.LoggingApp
			if !useExistingApp {
				appName = generator.RandomName()
				Expect(cf.Cf("push", appName, "-p", SIMPLE_RUBY_APP_BITS_PATH).Wait(CF_PUSH_TIMEOUT_IN_SECONDS)).To(Exit(0))
			}
		})

		AfterEach(func() {
			if !useExistingApp {
				Expect(cf.Cf("delete", appName, "-f").Wait(CF_TIMEOUT_IN_SECONDS)).To(Exit(0))
			}
		})

		It("can see app messages in the logs", func() {
			Eventually(func() *Session {
				appLogsSession := cf.Cf("logs", "--recent", appName)
				Expect(appLogsSession.Wait(CF_TIMEOUT_IN_SECONDS)).To(Exit(0))
				return appLogsSession
			}, 5).Should(Say(`\[App/0\]`))
		})
	})

	Describe("Syslog drains", func() {
		var drainListener *syslogDrainListener
		var serviceName string
		var appUrl string

		BeforeEach(func() {
			syslogDrainAddress := fmt.Sprintf("%s:%d", testConfig.SyslogIpAddress, testConfig.SyslogDrainPort)

			drainListener = &syslogDrainListener{port: testConfig.SyslogDrainPort}
			drainListener.StartListener()
			go drainListener.AcceptConnections()

			// verify listener is reachable via configured public IP
			var conn net.Conn

			var err error
			conn, err = net.Dial("tcp", syslogDrainAddress)
			Expect(err).ToNot(HaveOccurred())

			defer conn.Close()

			randomMessage := "random-message-" + generator.RandomName()
			_, err = conn.Write([]byte(randomMessage))
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() bool {
				return drainListener.DidReceive(randomMessage)
			}).Should(BeTrue())

			appName = generator.RandomName()
			appUrl = appName + "." + testConfig.AppsDomain
			Expect(cf.Cf("push", appName, "-p", SIMPLE_RUBY_APP_BITS_PATH).Wait(CF_PUSH_TIMEOUT_IN_SECONDS)).To(Exit(0))

			syslogDrainUrl := "syslog://" + syslogDrainAddress
			serviceName = "service-" + generator.RandomName()

			Expect(cf.Cf("cups", serviceName, "-l", syslogDrainUrl).Wait(CF_TIMEOUT_IN_SECONDS)).To(Exit(0))
			Expect(cf.Cf("bind-service", appName, serviceName).Wait(CF_TIMEOUT_IN_SECONDS)).To(Exit(0))
			Expect(cf.Cf("restage", appName).Wait(CF_PUSH_TIMEOUT_IN_SECONDS)).To(Exit(0))
		})

		AfterEach(func() {
			Expect(cf.Cf("delete", appName, "-f").Wait(CF_TIMEOUT_IN_SECONDS)).To(Exit(0))
			Expect(cf.Cf("delete-service", serviceName, "-f").Wait(CF_TIMEOUT_IN_SECONDS)).To(Exit(0))
			Expect(cf.Cf("delete-orphaned-routes", "-f").Wait(CF_PUSH_TIMEOUT_IN_SECONDS)).To(Exit(0))

			drainListener.Stop()
		})

		It("forwards app messages to registered syslog drains", func() {
			randomMessage := "random-message-" + generator.RandomName()
			http.Get("http://" + appUrl + "/log/" + randomMessage)

			Eventually(func() bool {
				return drainListener.DidReceive(randomMessage)
			}).Should(BeTrue())
		})
	})
})

type syslogDrainListener struct {
	sync.Mutex
	port             int
	listener         net.Listener
	receivedMessages string
}

func (s *syslogDrainListener) StartListener() {
	listenAddress := fmt.Sprintf(":%d", s.port)
	var err error
	s.listener, err = net.Listen("tcp", listenAddress)
	Expect(err).ToNot(HaveOccurred())
}

func (s *syslogDrainListener) AcceptConnections() {
	defer GinkgoRecover()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return
		}
		go s.handleConnection(conn)
	}
}

func (s *syslogDrainListener) Stop() {
	s.listener.Close()
}

func (s *syslogDrainListener) DidReceive(message string) bool {
	s.Lock()
	defer s.Unlock()

	return strings.Contains(s.receivedMessages, message)
}

func (s *syslogDrainListener) handleConnection(conn net.Conn) {
	defer GinkgoRecover()
	buffer := make([]byte, 65536)
	for {
		n, err := conn.Read(buffer)

		if err == io.EOF {
			return
		}
		Expect(err).ToNot(HaveOccurred())

		s.Lock()
		s.receivedMessages += string(buffer[0:n])
		s.Unlock()
	}
}
