run : 
	go run ./cmd/main.go

init :
	sudo sysytemctl restart docker
	sudo docker compose up -d

run_tests:
	go test -v ./test/library_test.go

docs:
	swag init -g ./cmd/main.go --parseDependency --parseInternal --parseDepth 2

clean : 
	rm -rf docs/