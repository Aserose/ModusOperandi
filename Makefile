build:
	go build -o ./.bin/app.exe -ldflags -H=windowsgui cmd/main.go

run: build
	./.bin/app