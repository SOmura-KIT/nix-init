.PHONY: build

build: 
	go build -o bin/nix-init cmd/nix-init/main.go

clean:
	rm -rf bin/kit
