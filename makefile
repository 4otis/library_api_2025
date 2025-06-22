init :
	sudo systemctl restart docker
	sudo docker compose up -d

run : 
	go run ./cmd/main.go

tests :
	go test -v ./test/...

docs :
	swag init -g ./cmd/main.go --parseDependency --parseInternal --parseDepth 2

clean : 
	rm -rf docs/