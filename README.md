CF Smoke Tests
==============

Smoke tests are a suite of basic core functionality tests for Cloud Foundry.
They are suitable as an initial test against a new or updated deployment to
reveal fundamental problems with the system.

There are four tests in this suite, all of which are pretty simple:
1. `runtime`: Pushes an app and validates that HTTP requests are properly routed to the app.
2. `logging`: Pushes an app and validates that logs can be fetched for the app.
3. `etcd_cluster_check`: Pushes an app and validates that app instances cannot reach the IP addresses of the etcd cluster.
4. `isolation_segments`: Entitles an org to an isolation segment and pushes two apps,
  one to the isolation segment, and one to the shared segment.
  The test validates that isolated apps are only accessed via the isolated router,
  and that apps on the shared segment are only accessed via the shared router.

They are not intended to test more sophisticated functionality of Cloud Foundry
or to test administrator operations. The [CF Acceptance
Tests](https://github.com/cloudfoundry/cf-acceptance-tests) do perform this
more extensive testing, although they are designed to be run as part of a
development pipeline and not against production environments.


## Running the tests

### Set up your `go` environment

Set up your golang development environment, [per golang.org](http://golang.org/doc/install).

Make sure you have the following installed:
* [git](http://git-scm.com/)
* [mercurial](http://mercurial.selenic.com/) (for `go get`)
* [bazaar](http://bazaar.canonical.com/) ( for `go get`)
* [`cf` CLI](https://github.com/cloudfoundry/cli)
* [curl](http://curl.haxx.se/)

Check out a copy of `cf-smoke-tests` and make sure that it is added to your
`$GOPATH`.  The recommended way to do this is to run `go get -u -d
github.com/cloudfoundry/cf-smoke-tests`. You will receive a warning "no
buildable Go source files"; this can be ignored as there is no compilable go
code in the package.
(Alternatively, you can simply `cd` into the directory
and run `git pull`.)

### Test Setup

To run the CF Smoke tests, you will need:
- a running CF instance
- an environment variable `$CONFIG` which points to a `.json` file that
contains the application domain

Below is an example `integration_config.json`:
```json
{
  "suite_name"                      : "CF_SMOKE_TESTS",
  "api"                             : "api.bosh-lite.com",
  "apps_domain"                     : "bosh-lite.com",
  "user"                            : "non-admin",
  "password"                        : "super-secure",
  "org"                             : "CF-SMOKE-ORG",
  "space"                           : "CF-SMOKE-SPACE",
  "isolation_segment_space"         : "CF-SMOKE-ISOLATED-SPACE",
  "cleanup"                         : true,
  "use_existing_org"                : true,
  "use_existing_space"              : true,
  "logging_app"                     : "",
  "runtime_app"                     : "",
  "enable_windows_tests"            : false,
  "enable_etcd_cluster_check_tests" : false,
  "etcd_ip_address"                 : "",
  "backend"                         : "diego",
  "isolation_segment_name"          : "is1",
  "isolation_segment_domain"        : "is1.bosh-lite.com",
  "enable_isolation_segment_tests"  : true
}
```
**NOTE** Unless you supply an admin user, you _must_ use an existing space and org.
The tests will only pass if you have configured your environment in a way that allows isolation segments to be tested properly:
- If you do not provide an admin user, the smoke-tests `org` must be entitled to use the isolation segment and contain two spaces
    - The space that is referred to as `space` in the smoke-tests config must be assigned to the shared (i.e. not isolated) segment
    - The space that is referred to as `isolation_segment_space` in the smoke-tests-config must be assigned to the isolation segment
- Alternatively, you may provide an admin user and omit configuration of the existing spaces to allow smoke-tests to create the testing resources automatically.


If you are running the tests against bosh-lite or any other environment using
self-signed certificates, add

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

If you like to validate the security of your etcd cluster, set `enable_etcd_cluster_check_tests` to true and provide the `etcd_ip_address` to be the least restrictive IP that your etcd cluster has (private if that is the only IP etcd has, public otherwise)

If you like to run isolation segment test, set `enable_isolation_segment_tests` to true and provide values for `isolation_segment_name`, `isolation_segment_domain` and set `backend` to `diego`. Test setup assumes that isolation segment API resource with `isolation_segment_name` already exists. For more details on how to setup routing isolation segments, read this [document](https://docs.cloudfoundry.org/adminguide/routing-is.html)

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

## Contributing to Smoke Tests

### Guidelines
The goal of smoke tests
is to provide a small, simple set of tests
to verify basic deployment configuration.
As such, we have some guidelines
for contributing new tests to this suite.

#### 1. Creating API resources in the test
One basic rule for good test design is not to mock the object under test.
We can translate that idea to a suite like smoke tests in the following way:
If smoke tests exist to validate deployment configuration,
then smoke tests should not itself mutate deployment configuration.

There are, however, several resources
that can be defined as either deployment configuration or as API resources.
For example, shared app domains and isolation segments
are both resources that can be created via the API,
so it might be tempting to have a test create them in a `BeforeSuite`.
However, shared app domains and isolation segments really represent deployment configurations.
Accordingly, smoke tests should not create those resources as part of the test;
instead, it should validate (either implicitly or explicitly)
that those resources have already been created, and configured correctly.

Other API resources, like orgs and spaces
that exist simply to be able to push an app,
can absolutely be created as part of a test.

#### 2. Admin vs. Regular User workflows
There are two common user workflows
that can be validated in smoke tests.

1. **Regular user**:
Smoke tests are run with a configuration
that provides a non-admin user.
In this configuration,
tests should also expect
that the org and space used for the tests
have already been created.
This workflow is recommended for tests run against environments run by humans
-- in particular, production deployments.

2. **Admin user:**
Smoke tests are configured to run using admin credentials.
Given this configuration,
the tests may or may not use existing resources like orgs and spaces,
because an admin user can easily create them.
This configuration is recommended for tests run against environments created using automation tools,
for example, CI (continuous integration) environments on development teams.

### Dependency Management

Smoke Tests use [dep](https://github.com/golang/dep) to manage `go` dependencies.

All `go` packages required to run smoke tests are vendored into the `vendor/` directory.
