build-consumer-image:
	docker build -t consumer -f ./build/Dockerfile .

build-sender:
	go build -o ./snd ./sender

test:
	go test ./...

lint:
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run --no-config --max-same-issues=100 \
		--issues-exit-code=1 \
		-v \
		--deadline=160s \
		--disable-all \
		--enable goconst \
		--enable golint \
		--enable gosimple \
		--enable maligned \
		--enable misspell \
		--enable ineffassign \
		--enable interfacer \
		--enable staticcheck \
		--enable structcheck \
		--enable unconvert \
		--enable varcheck \
		--enable gas \
		--enable unparam \
		--enable dogsled \
		--enable depguard \
		--enable gochecknoinits \
		--enable gofmt \
		--enable whitespace \
		--enable deadcode \
		--enable bodyclose \
		--enable stylecheck \
		--enable nakedret \
		./...

compose-up:
	docker-compose -f ./build/docker-compose.yml up

compose-down:
	docker-compose -f ./build/docker-compose.yml down