# Go WhatsApp Implementation in Command-Line Interface (CLI)

This repository contains example of implementation [dimaskiddo/go-whatsapp](https://github.com/dimaskiddo/go-whatsapp) package. This example is using a codebase from [dimaskiddo/go-whatsapp-cli](https://github.com/dimaskiddo/go-whatsapp-cli).

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.
See deployment for notes on how to deploy the project on a live system.

### Prerequisites

Prequisites packages:
* Go (Go Programming Language)
* Dep (Go Dependencies Management Tool)
* Make (Automated Execution using Makefile)

Optional packages:
* GoReleaser (Go Automated Binaries Build)

### Installing

Below is the instructions to make this codebase running:
* Create a Go Workspace directory and export it as the extended GOPATH directory
```
cd <your_go_workspace_directory>
export GOPATH=$GOPATH:"`pwd`"
```
* Under the Go Workspace directory create a source directory
```
mkdir -p src/github.com/dimaskiddo/go-whatsapp-cli
```
* Move to the created directory and pull codebase
```
cd src/github.com/dimaskiddo/go-whatsapp-cli
git clone -b master https://github.com/dimaskiddo/go-whatsapp-cli.git .
```
* Run following command to pull dependecies package
```
make ensure
```
* Until this step you already can run this code by using this command
```
make run
```

## Running The Tests

Currently the test is not ready yet :)

## Deployment

To build this code to binaries for distribution purposes you can run following command:
```
make compile
```
The build result will shown in build directory

## Built With

* [Go](https://golang.org/) - Go Programming Languange
* [Dep](https://github.com/golang/dep) - Go Dependency Management Tool
* [GoReleaser](https://github.com/goreleaser/goreleaser) - Go Automated Binaries Build
* [Make](https://www.gnu.org/software/make/) - GNU Make Automated Execution

## Authors

* **Dimas Restu Hidayanto** - *Initial Work* - [DimasKiddo](https://github.com/dimaskiddo)

See also the list of [contributors](https://github.com/dimaskiddo/go-whatsapp-cli/contributors) who participated in this project

## Annotation

You can seek more information for the make command parameters in the [Makefile](https://raw.githubusercontent.com/dimaskiddo/go-whatsapp-cli/master/Makefile)
