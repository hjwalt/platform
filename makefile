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
	go build -o bin/flows

run:
	go run .

reset:
	./script/reset.sh

listen:
	./script/listen.sh

group-delete:
	./script/group-delete.sh
	
# --------------------

mocks: RUN
	mockgen -source=mock/interfaces.go -destination=mock/implementations.go -package=mock ;\
	mockgen -source=stateful_bun/connection.go -destination=mock/stateful_bun_connection.go -package=mock ;\

proto: RUN
	clang-format -style=file:.clang-format -i **/*.proto
	protoc  --proto_path=. --go_opt=paths=source_relative --go_out=. ./**/*.proto

# --------------------

RUN:
