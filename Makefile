#NEED_TEST_DIR=$(shell go list ./...|grep -v -e "github.com/Me1onRind/go-demo/\(internal/usecase\|protocol\|internal/infrastructure/logger\|internal/infrastructure/client\)")

http_run:
	go run ./cmd/http/main.go

supervisor_run:
	mkdir -p ./supervisor/log
	mkdir -p ./supervisor/run
	supervisord -c ./supervisor/supervisord.conf # /etc/rc.local  auto run in power on.

lint:
	golangci-lint run -v ./...

mock:
	mockgen --source ./internal/infrastructure/client/etcd/etcd.go --destination ./internal/infrastructure/client/etcd/mock_etcd.go --package etcd
	mockgen --source ./internal/domain/iddm/iddm.go --destination ./internal/domain/iddm/mock_iddm.go --package iddm

test:
	@skip_io=1 go test -cover ./...

test_cover:
	@skip_io=1 go test ./... -coverprofile=/tmp/go_test.out
	go tool cover -html=/tmp/go_test.out -o=/root/share/coverage.html

test_cover_gui:
	@skip_io=1 go test ./... -coverprofile=/tmp/go_test.out
	go tool cover -html=/tmp/go_test.out

generate:
	go generate ./...
