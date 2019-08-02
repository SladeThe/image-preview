test:
	go test ./...

run-dev-docker: test
	rm -rf ./build/bin/*
	mkdir -p ./build/bin
	go build -o ./build/bin/image-preview ./cmd/image-preview/
	chmod +x ./build/bin/image-preview
	docker-compose -p image-preview -f ./build/package/docker/dev/docker-compose.yml up
