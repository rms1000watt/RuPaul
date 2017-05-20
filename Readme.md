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
- `Dockerfile..tpl` will go to `/Dockerfile`

### Installation

```sh
go get -u -v github.com/rms1000watt/rupaul
go get -u -v github.com/magical/argon2
go get -u -v github.com/spf13/cobra
go get -u -v github.com/spf13/pflag
```

### Usage

General Usage

```sh
# Edit examples/rupaul.yml as needed
go run main.go generate -f examples/rupaul.yml
```

Testing Purposes

```sh
clear; rm -rf out; go run main.go generate -f examples/rupaul.yml
PROJECT_PATH=$(go env GOPATH)/src/github.com/rms1000watt/rupaul-test bash -c 'rm -rf $PROJECT_PATH && mkdir $PROJECT_PATH  && mkdir $PROJECT_PATH/certs && cp -r out/* $PROJECT_PATH && cp -r certs/* $PROJECT_PATH/certs && cp out/.gitignore $PROJECT_PATH/'

cd ../rupaul-test; clear; go run main.go serve

# Fails
curl -X POST -d '{"first_name":"Chet","middle_name":"Darf","last_name":"Star"}' localhost:8080/person
curl -X POST -d '{"first_name":"Chet","middle_name":"Darf","last_name":"Star","age":33}' localhost:8080/person
curl -X POST -d '{"first_name":"Chet","middle_name":"Darf","last_name":"Star","age":33,"account":123.123}' https://localhost:8080/person
curl -X POST -d '{"first_name":"Chet","middle_name":"Darf","last_name":"StarStarStarStarStarStarStarStarStarStarStarStarStar","age":33,"account":123.123,"password":"pASSword"}' https://localhost:8080/person

# Success
curl -X POST -d '{"first_name":"ChetChetChetChet","middle_name":"Darf","last_name":"Star","age":33,"account":123.123,"password":"pASSword"}' --insecure https://localhost:8080/person
curl -X POST -d '{"first_name":"ChetChetChetChet","middle_name":"Darf","last_name":"Star","age":33,"account":123.123,"password":"pASSword","planet":{"name":"earth"}}' --insecure https://localhost:8080/person
curl -X POST -d '{"first_name":"Chet","middle_name":"Darf","last_name":"Star","age":33,"account":123.123,"password":"pASSword","gossip":"hello world"}' --insecure https://localhost:8080/person
curl -X POST -d '{"first_name":"Chet","middle_name":"Darf","last_name":"Star","age":33,"account":123.123,"password":"pASSword","gossip":"hello world","grocery_list":["pigs","blanket"],"age_list":[3,5,1]}' --insecure https://localhost:8080/person
curl --insecure 'https://localhost:8080/person?first_name=ryan&last_name=smith&age=88'
curl --insecure 'https://localhost:8080/ticket?ticket_id=123123'
curl -X POST -d '{"ticket_id":"123123"}' --insecure https://localhost:8080/ticket
https://localhost:8080/person?first_name=ryan&last_name=smith&age=88
```

**WIP**

### TODO:

- [x] Handle array data types
- [x] Handle struct data types
- [] Handle nested structs
- [] Validation of structs and nested structs
- [] Transformation of structs and nested structs
- [] Docs, docs, docs
- [] Config sanitization
- [] Template cleanup (use functions to format var names instead of in-template)
- [] Connectors
- [] Move from fmt.Println() to log.Debug() or similar
- [] More and better examples
- [] ~~Consider moving helper/validate/transform/unmarshal logic to separate library~~ (Encryption/Decryption needs to be handled properly by the user)
- [] Govendor
