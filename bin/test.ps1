$ErrorActionPreference = "Stop";
trap { $host.SetShouldExit(1) }

where.exe go
if ($LASTEXITCODE -ne 0) {
  echo "Go is not installed."
  exit 1
}

if (-not (Test-Path env:GOPATH)) {
  echo "GOPATH not specified"
  exit 1
}

$env:PATH="$env:GOPATH\bin;$env:PATH"

go.exe install github.com\onsi\ginkgo\v2\ginkgo

ginkgo.exe -r --succinct --poll-progress-after=300s $args
exit $LASTEXITCODE
