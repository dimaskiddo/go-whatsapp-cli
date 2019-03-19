BUILD_OS           := linux
BUILD_OUTPUT       := gowa
REBASE_URL         := "github.com/dimaskiddo/go-whatsapp-cli"
COMMIT_MSG         := "update improvement"

.PHONY:

.SILENT:

init:
	make clean
	rm -f Gopkg.toml Gopkg.lock
	dep init -v

ensure:
	make clean
	dep ensure -v

compile:
	make ensure
	CGO_ENABLED=0 GOOS=$(BUILD_OS) go build -a -o ./build/$(BUILD_OUTPUT) *.go
	echo "Build complete please check build directory."

run:
	go run *.go

clean:
	rm -f ./build/$(BUILD_OUTPUT)
	rm -rf ./vendor

commit:
	make init
	make clean
	git add .
	git commit -am "$(COMMIT_MSG)"

rebase:
	rm -rf .git
	find . -type f -iname "*.go*" -exec sed -i '' -e "s%github.com/dimaskiddo/go-whatsapp-cli%$(REBASE_URL)%g" {} \;
	git init
	git remote add origin https://$(REBASE_URL).git

push:
	git push origin master

pull:
	git pull origin master
