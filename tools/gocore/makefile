build:
	statik -src=./ui/web -dest=./ui/
	GOARCH=wasm GOOS=js go build -o ui/web/app.wasm
	statik -src=./ui/web -dest=./ui/
	# go build
	# rm -rf ./ui/statik/*

run: build
	./goui