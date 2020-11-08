build: clean make-dir copy-dependent build-mac zip-mac build-linux zip-linux clean2
clean:
	rm -rf dist/*
clean2:
	rm -rf dist/work
make-dir:
	mkdir -p ./dist/{work,releases}
build-mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o dist/work
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/work
copy-dependent:
	cp -r ./lib ./CHANGELOG ./COMMANDS.md ./README.md ./LICENSE dist/work
zip-mac:
	zip -r ./dist/releases/deer-executor_darwin_amd64.zip ./dist/work/*
zip-linux:
	zip -r ./dist/releases/deer-executor_linux_amd64.zip ./dist/work/*