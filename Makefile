init: clean generate
	go mod tidy
	go mod vendor

compile-api:
	go tool oapi-codegen -config ./misc/cfg.yaml ./misc/api.yaml

clean:
	rm -rf generated
	rm -f ./app/cmd/wire_gen.go

generate:
	mkdir -p generated
	wire ./app/cmd/wire.go
	$(MAKE) compile-api

up:
	docker-compose up -d

wait:
	@echo "⏳ Waiting for services to be healthy..."
	@for service in postgres redis mongodb; do \
		echo "⏳ Waiting for $$service..."; \
		until [ "$$(docker inspect --format='{{.State.Health.Status}}' $$service)" = "healthy" ]; do \
			sleep 1; \
		done; \
		echo "✅ $$service is healthy."; \
	done

run: up wait
	go run ./app/cmd main