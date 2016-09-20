
mac:
	go build -o remote-tail-mac main.go

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o remote-tail-linux	main.go
