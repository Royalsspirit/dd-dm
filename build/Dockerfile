FROM dd/base-golang

ENV threshold 10

RUN ["go", "get", "github.com/githubnemo/CompileDaemon"]

ENTRYPOINT CompileDaemon -build='go build -o ./bin/dd-monitoring ./cmd/cli' -command="./bin/dd-monitoring -threshold ${threshold}"
