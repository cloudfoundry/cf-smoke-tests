module github.com/cloudfoundry/cf-smoke-tests-release/src/smoke_tests

go 1.18

require (
	github.com/cloudfoundry-incubator/cf-test-helpers v1.0.1-0.20191216200933-cf8305784c93
	github.com/cloudfoundry/cf-smoke-tests v0.0.0-20200921182851-bc7c19c52112
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.18.1
)

require (
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	golang.org/x/net v0.0.0-20210428140749-89ef3d95e781 // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	golang.org/x/text v0.3.6 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace gopkg.in/fsnotify.v1 v1.4.7 => github.com/fsnotify/fsnotify v1.4.7
