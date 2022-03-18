module github.com/cloudfoundry/cf-smoke-tests-release/src/smoke_tests

go 1.18

require (
	github.com/cloudfoundry-incubator/cf-test-helpers v1.0.1-0.20191216200933-cf8305784c93
	github.com/fsnotify/fsnotify v1.4.7
	github.com/hpcloud/tail v1.0.0
	github.com/onsi/ginkgo v1.9.0
	github.com/onsi/gomega v1.6.0
	golang.org/x/net v0.0.0-20190813141303-74dc4d7220e7
	golang.org/x/sys v0.0.0-20190813064441-fde4db37ae7a
	golang.org/x/text v0.3.2
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7
	gopkg.in/yaml.v2 v2.2.2
)

replace gopkg.in/fsnotify.v1 v1.4.7 => github.com/fsnotify/fsnotify v1.4.7
