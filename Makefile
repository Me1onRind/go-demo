NEED_TEST_DIR=$(shell go list ./...|grep -v -e "github.com/Me1onRind/go-demo/\(internal/usecase\|protocol\|internal/infrastructure/logger\)")

http_run:
	go run ./cmd/http/main.go

supervisor_run:
	supervisord -c ./supervisor/supervisord.conf # /etc/rc.local  auto run in power on.

lint:
	golangci-lint run ./...

mock:
	mockgen --source ./internal/infrastructure/client/etcd/etcd.go --destination ./internal/infrastructure/client/etcd/mock_etcd.go --package etcd

test:
	@go test $(NEED_TEST_DIR)
