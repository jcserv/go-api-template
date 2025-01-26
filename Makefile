.PHONY: clean dev run/docker 

clean:
	rm main

dev:
	go build ./cmd/go-api-template/main.go && ./main

run/docker:
	docker compose up -d