MODULE=github.com/hjwalt/platform

test:
	go test ./... -cover -coverprofile cover.out
	
testv:
	go test ./... -cover -coverprofile cover.out -v

cov: test
	go tool cover -func cover.out

htmlcov: test
	go tool cover -html cover.out -o cover.html

# --------------------

tidy:
	go mod tidy
	go fmt ./...

update:
	go get -u ./...
	go mod tidy
	go fmt ./...

# --------------------

build:
	go build -o bin/platform

run:
	go run ./main.go

reset:
	./script/reset.sh

# --------------------

mocks: RUN
	mockgen -source=flows/test_helper/interfaces.go -destination=flows/test_helper/implementations.go -package=test_helper ;\
	mockgen -source=flows/runtime_bun/connection.go -destination=flows/test_helper/stateful_bun_connection.go -package=test_helper ;\

proto: RUN
	clang-format -style=file:.clang-format -i **/*.proto
	protoc  --proto_path=. --go_opt=paths=source_relative --go_out=. ./**/*.proto


# --------------------

up: RUN
	docker compose up -d

down: RUN
	docker compose down


# --------------------

RUN:
