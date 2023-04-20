#!/bin/bash
docker run -v $PWD:/home/cf-smoke-tests -w /home/cf-smoke-tests cloudfoundry/cf-deployment-concourse-tasks go build -o app
