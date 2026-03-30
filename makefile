run:
	air --build.cmd "go build -o app/api/main ./app/api" --build.bin "./app/api/main"

pretty:  
	make run | jq -R 'fromjson? | select(type == "object")'

upgrade:
	go get -u -v ./...
	go mod tidy
	go mod vendor

lint:
	CGO_ENABLED=0 go vet ./...
	staticcheck -checks=all ./...

dev-tools:
	go install github.com/cosmtrek/air@latest
	brew install jq
