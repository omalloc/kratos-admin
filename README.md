# Kratos Admin

## Quick Start

```bash
make build
./bin/kratos-admin --conf ./configs
```

## Generate other auxiliary files by Makefile

```bash
# Download and update dependencies
make init
# Generate API files (include: pb.go, http, grpc, validate, swagger) by proto file
make api
# Generate all files
make all
```

## Docker

```bash
# build
docker build -t <your-docker-image-name> .

# run
docker run --rm -p 8000:8000 -p 9000:9000 -v </path/to/your/configs>:/data/conf <your-docker-image-name>
```

## Development

### with vscode debug

`.vscode/launch.json`

```jsonc
{
  // Use IntelliSense to learn about possible attributes.
  // Hover to view descriptions of existing attributes.
  // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch App",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "cwd": "${workspaceRoot}",
      "program": "${workspaceRoot}/cmd/kratos-admin/.",
      "args": ["--conf", "./configs/"]
    }
  ]
}
```

### protobuf usage

```bash
cd server
# Add a proto template
kratos proto add api/server/server.proto
# Generate the proto code
kratos proto client api/server/server.proto
# Generate the source code of service by proto file
kratos proto server api/server/server.proto -t internal/service

make api wire
```

### rename Package

```bash
$ export NEW_MODULE_NAME=example.com/you-are-project
$ export OLD_MODULE_NAME=github.com/omalloc/kratos-admin
$ go mod edit -module ${NEW_MODULE_NAME}
$ find . -type f -name '*.go' -exec sed -i -e 's,{OLD_MODULE_NAME},{NEW_MODULE_NAME},g' {};

go run cmd/kratos-admin/main.go cmd/kratos-admin/cobra.go cmd/kratos-admin/wire_gen.go --conf ./configs
```
