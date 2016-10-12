FROM golang:1.7

RUN go get -d github.com/cloudfoundry/cf-smoke-tests || echo ""

RUN curl -L "https://cli.run.pivotal.io/stable?release=linux64-binary&source=github" | tar -zx
RUN mv cf /bin/cf

WORKDIR /go
RUN mkdir -p src/github.com/onsi

WORKDIR /go/src/github.com/cloudfoundry/cf-smoke-tests
RUN cp -R vendor/github.com/onsi/ginkgo /go/src/github.com/onsi/ginkgo
RUN go install -v github.com/onsi/ginkgo/ginkgo

CMD echo $CONFIG_CONTENT > /go/runtime_config.json ; CONFIG=/go/runtime_config.json bin/test -v


# Sample invocation:
# docker run -e "CONFIG_CONTENT=$(cat ../runtime-config.json)" <image id>
