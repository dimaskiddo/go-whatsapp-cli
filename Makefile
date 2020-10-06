BUILD_CGO_ENABLED  := 0
SERVICE_NAME       := go-whatsapp-cli
REBASE_URL         := "github.com/dimaskiddo/go-whatsapp-cli"
COMMIT_MSG         := "update improvement"

.PHONY:

.SILENT:

init:
	make clean
	GO111MODULE=on go mod init

init-dist:
	mkdir -p dist

vendor:
	make clean
	GO111MODULE=on go mod vendor

release:
	make vendor
	make clean-dist
	goreleaser --snapshot --skip-publish --rm-dist
	echo "Release '$(SERVICE_NAME)' complete, please check dist directory."

publish:
	make vendor
	make clean-dist
	GITHUB_TOKEN=$(GITHUB_TOKEN) goreleaser --rm-dist
	echo "Publish '$(SERVICE_NAME)' complete, please check your repository releases."

build:
	make vendor
	CGO_ENABLED=$(BUILD_CGO_ENABLED) go build -ldflags="-s -w" -a -o $(SERVICE_NAME) cmd/main/main.go
	echo "Build '$(SERVICE_NAME)' complete."

run:
	go run cmd/main/main.go

clean-dist:
	rm -rf dist

clean:
	make clean-dist
	rm -rf vendor

commit:
	make vendor
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
