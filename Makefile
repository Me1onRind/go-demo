http_run:
	go run ./cmd/http/main.go

supervisor_run:
	supervisord -c ./supervisor/supervisord.conf # /etc/rc.local  auto run in power on.

lint:
	golangci-lint run ./...
