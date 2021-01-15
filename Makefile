compile:
	echo "Compiling for most os's and platforms"
	echo "-> linux:"
	GOOS=linux GOARCH=386 go build -o bin/whgoxy-linux-386 ./cmd/whgoxy/main.go
	GOOS=linux GOARCH=amd64 go build -o bin/whgoxy-linux-amd64 ./cmd/whgoxy/main.go
	GOOS=linux GOARCH=arm go build -o bin/whgoxy-linux-arm ./cmd/whgoxy/main.go
	GOOS=linux GOARCH=arm64 go build -o bin/whgoxy-linux-arm64 ./cmd/whgoxy/main.go
	echo "-> darwin:"
	GOOS=darwin GOARCH=386 go build -o bin/whgoxy-darwin-386 ./cmd/whgoxy/main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/whgoxy-darwin-amd64 ./cmd/whgoxy/main.go
	echo "-> freebsd:"
	GOOS=freebsd GOARCH=arm go build -o bin/whgoxy-freebsd-arm ./cmd/whgoxy/main.go
	GOOS=freebsd GOARCH=amd64 go build -o bin/whgoxy-freebsd-amd64 ./cmd/whgoxy/main.go
	GOOS=freebsd GOARCH=386 go build -o bin/whgoxy-freebsd-386 ./cmd/whgoxy/main.go
	echo "-> windows:"
	GOOS=windows GOARCH=amd64 go build -o bin/whgoxy-windows-amd64.exe ./cmd/whgoxy/main.go
	GOOS=windows GOARCH=386 go build -o bin/whgoxy-windows-i386.exe ./cmd/whgoxy/main.go