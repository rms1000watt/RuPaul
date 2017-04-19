## Rygen

Golang code generator

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
go get -u -v github.com/rms1000watt/rygen
```

### Usage

General Usage

```sh
# Edit examples/rygen.yml as needed
go run main.go generate -f examples/rygen.yml
```

Testing Purposes

```sh
clear && rm -rf out && go run main.go generate -f examples/rygen.yml
```

**WIP**