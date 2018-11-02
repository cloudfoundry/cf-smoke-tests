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
where.exe ginkgo
if ($LASTEXITCODE -ne 0) {
  New-Item -Path "$env:GOPATH\src\github.com\onsi" -ItemType Directory -Force
  Copy-Item -Recurse -Force vendor\github.com\onsi\ginkgo "$env:GOPATH\src\github.com\onsi"
  go.exe install -v github.com\onsi\ginkgo\ginkgo
}

$env:CF_DIAL_TIMEOUT=11

ginkgo.exe -r --succinct -slowSpecThreshold=300 $args
exit $LASTEXITCODE
