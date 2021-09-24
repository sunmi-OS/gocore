build:
	GOARCH=wasm GOOS=js go build -o ui/web/app.wasm
	statik -src=./ui/web -dest=./ui/
	go build

run: build
	./gocore