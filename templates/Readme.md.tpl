## {{.CommandLine.AppName}}

Description for {{.CommandLine.AppName}} goes here.

### Installation

```sh
go get -u -v {{.MainImportPath}}
```

### Usage

```sh
go run main.go
```

### Deploy

```sh
go build
docker build --rm -t --no-cache {{.DockerPath}}:{{.Version}} .
docker push {{.DockerPath}}:{{.Version}}
```
