
mac:
	go build -o remote-tail-mac main.go

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o remote-tail-linux	main.go

deploy:
	scp ./remote-tail-linux root@192.168.1.100:/usr/bin/remote-tail
	scp ./remote-tail-linux root@192.168.1.237:/usr/bin/remote-tail

deploy-local:
	cp ./remote-tail-mac /usr/local/bin/remote-tail

clean:
	rm -fr ./remote-tail-linux ./remote-tail-mac
