build:
	GOOS=linux GOARCH=amd64 go build cmd/server/main.go
deploy:
	scp ./main homeho:/home/homeho-backend