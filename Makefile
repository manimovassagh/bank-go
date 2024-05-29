db:
	docker compose up -d

build: db
	@go build -o bin/main .

run: build
	./bin/main
