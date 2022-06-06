.PHONY: run
run:
	go run cmd/main.go

.PHONY: docker-up
docker-up:
	docker-compose up --build

.PHONY: proto-gen
proto-gen:
	protoc --proto_path=proto proto/*.proto --go_out=proto/ --experimental_allow_proto3_optional
	protoc --proto_path=proto proto/*.proto --go-grpc_out=proto/ --experimental_allow_proto3_optional

.PHONY: proto-clean
proto-clean:
	rm proto/*.go

.PHONY: migrate
migrate:
	migrate -path ./migrations -database 'postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable' up

.PHONY: migrate-down
migrate-down:
	migrate -path ./migrations -database 'postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable' down

.DEFAULT_GOAL = run
