run : 
	go run ./cmd/main.go

docs:
	swag init -g ./cmd/main.go --parseDependency --parseInternal --parseDepth 2

clean : 
	rm -rf docs/