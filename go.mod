module github.com/cloudfoundry/cf-smoke-tests-release/src/smoke_tests

go 1.18

require (
	github.com/cloudfoundry-incubator/cf-test-helpers v1.0.1-0.20191216200933-cf8305784c93
	github.com/cloudfoundry/cf-smoke-tests v0.0.0-20200921182851-bc7c19c52112
	github.com/onsi/ginkgo v1.9.0
	github.com/onsi/gomega v1.6.0
)

require (
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/hpcloud/tail v1.0.0 // indirect
	golang.org/x/net v0.0.0-20190813141303-74dc4d7220e7 // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	golang.org/x/text v0.3.2 // indirect
	gopkg.in/fsnotify.v1 v1.4.7 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v2 v2.2.2 // indirect
)

replace gopkg.in/fsnotify.v1 v1.4.7 => github.com/fsnotify/fsnotify v1.4.7
