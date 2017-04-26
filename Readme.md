## RuPaul: The Golang Code Generator

RuPaul: The Golang Code Generator is the fiercest Golang starter kit you've encountered!

**WIP**

### Design

- The generator will never overwrite existing files
- All templates must have `[<DIR>. ...]<FILENAME>.<FILEEXTENSION>.tpl`

ie.

- `.gitignore.tpl` will go to `/.gitignore`
- `main.go.tpl` will go to `/main.go`
- `cmd.root.go.tpl` will go to `/cmd/root.go`

### Installation

```sh
go get -u -v github.com/rms1000watt/rupaul
```

### Usage

General Usage

```sh
# Edit examples/rupaul.yml as needed
go run main.go generate -f examples/rupaul.yml
```

Testing Purposes

```sh
clear && rm -rf out && go run main.go generate -f examples/rupaul.yml
PROJECT_PATH=$(go env GOPATH)/src/github.com/rms1000watt/rupaul-test bash -c 'rm -rf $PROJECT_PATH && mkdir $PROJECT_PATH && cp -r out/* $PROJECT_PATH'

cd ../rupaul-test && go run main.go serve

curl -X POST -d '{"first_name":"Chet","middle_name":"Darf","last_name":"Star"}' localhost:8080/person
curl -X POST -d '{"first_name":"Chet","middle_name":"Darf","last_name":"Star","age":33}' localhost:8080/person
curl -X POST -d '{"first_name":"Chet","middle_name":"Darf","last_name":"Star","age":33,"account":123.123}' localhost:8080/person
```

**WIP**