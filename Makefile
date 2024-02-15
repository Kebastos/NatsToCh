.PHONY: run
run: ## Build and run server.
	go build -o bin/server -ldflags="-X 'main.Version=v1.0.0'"  main.go
	bin/server