build: clean make-dir build-mac build-linux copy-dependent
clean:
	rm -rf dist/* && rm -rf rice-box.go
make-dir:
	mkdir -p ./dist/{mac,linux}
build-mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o dist/mac
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/linux
copy-dependent:
	cp -r ./lib ./CHANGELOG ./COMMANDS.md ./README.md ./LICENSE dist/mac
	cp -r ./lib ./CHANGELOG ./COMMANDS.md ./README.md ./LICENSE dist/linux