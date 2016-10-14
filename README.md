CF Smoke Tests
==============

Smoke tests are a suite of basic core functionality tests for Cloud Foundry.
They are suitable as an initial test against a new or updated deployment to
reveal fundamental problems with the system. They are also safe to run
periodically against a production environment as they place minimal load on a
system and can be configured to require access only by a regular user. In
particular, Cloud Foundry operators can use these tests as a monitoring tool
against their running deployment.

They are not intended to test more sophisticated functionality of Cloud Foundry
or to test administrator operations. The [CF Acceptance
Tests](https://github.com/cloudfoundry/cf-acceptance-tests) do perform this
more extensive testing, although they are designed to be run as part of a
development pipeline and not against production environments.


## Running the tests

### Set up your `go` environment

Set up your golang development environment, [per golang.org](http://golang.org/doc/install).

You will probably also need the following SCM programs in order to `go get`
source code:
* [git](http://git-scm.com/)
* [mercurial](http://mercurial.selenic.com/)
* [bazaar](http://bazaar.canonical.com/)

See [Go CLI](https://github.com/cloudfoundry/cli) for instructions on
installing the go version of `cf`.

Make sure that [curl](http://curl.haxx.se/) is installed on your system.

Make sure that the go version of `cf` is accessible in your `$PATH`.

Check out a copy of `cf-smoke-tests` and make sure that it is added to your
`$GOPATH`.  The recommended way to do this is to run `go get -d
github.com/cloudfoundry/cf-smoke-tests`. You will receive a warning "no
buildable Go source files"; this can be ignored as there is no compilable go
code in the package.

### Test Setup

To run the CF Smoke tests, you will need:
- a running CF instance
- an environment variable `$CONFIG` which points to a `.json` file that
contains the application domain

Below is an example `integration_config.json`:
```json
{
  "suite_name"           : "CF_SMOKE_TESTS",
  "api"                  : "api.bosh-lite.com",
  "apps_domain"          : "bosh-lite.com",
  "user"                 : "admin",
  "password"             : "admin",
  "org"                  : "CF-SMOKE-ORG",
  "space"                : "CF-SMOKE-SPACE",
  "cleanup"              : true,
  "use_existing_org"     : false,
  "use_existing_space"   : false,
  "logging_app"          : "",
  "runtime_app"          : "",
  "enable_windows_tests" : false,
  "backend"              : "diego"
}
```

If you are running the tests with version newer than 6.0.2-0bba99f of the Go
CLI against bosh-lite or any other environment using self-signed certificates,
add

```
  "skip_ssl_validation": true
```

If you would like to preserve the organization, space, and app created during the
tests for debugging, add

```
  "cleanup": false
```

If you have deployed Windows cells, add

```
  "enable_windows_tests" : true
```

If you like to a specific backend, add (allowed diego, dea or empty (default))

```
  "backend" : "diego"
```


### Test Execution

To execute the tests, run:

```bash
./bin/test
```

Internally the `bin/test` script runs tests using [ginkgo](https://github.com/onsi/ginkgo).

Arguments, such as `-focus=`, `-nodes=`, etc., that are passed to the script are sent to `ginkgo`

For example, to execute tests in parallel across two processes one would run:

```bash
./bin/test -nodes=2
```

#### Seeing command-line output

To see verbose output from `cf`, use [ginkgo](https://github.com/onsi/ginkgo)'s `-v` flag.

```bash
./bin/test -v
```

#### Capturing CF CLI output

Set '`artifacts_directory`' in your `integration_config.json` (as shown below)
to store cf cli trace output. The output files will be saved inside the given
directory.

```
  "artifacts_directory": "/tmp/smoke-artifacts"
```

The following files may be created:

```bash
CF-TRACE-Smoke-1.txt
CF-TRACE-Smoke-2.txt
junit-Applications-1.xml
...
```

## Changing Smoke Tests

### Dependency Management

Smoke Tests use [gvt](https://github.com/FiloSottile/gvt) to manage `go` dependencies.

All `go` packages required to run smoke tests are vendored into the `vendor/` directory.

When making changes to the test suite that bring in additional `go` packages,
you should use the workflow described in the
[gvt documentation](https://github.com/FiloSottile/gvt#basic-usage).
