default: lint format

lint:
  golangci-lint run

format:
  fd '.*\.go' | xargs -L1 go fmt
  fd '.*\.go' | xargs -L1 go fix

build:
  go build 'git.sr.ht/~ansipunk/weaver/cmd/weaver'
