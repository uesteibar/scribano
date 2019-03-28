# AsyncAPI watcher

Builds an [asyncapi](https://www.asyncapi.com/) documentation for your microservices
communicating through [rabbitmq](https://www.rabbitmq.com/).

It listens to all published amqp messages and keeps an updated asyncapi
compliant documenation served at `/asyncapi`.

## Roadmap

- [x] Add info and server sections to spec to make it valid asyncapi.
- [x] Extract configuration to file.
- [x] Support consuming from multiple configurable exchanges.
- [ ] Allow loading configuration from url instead of local file.
- [ ] Use postgres as database.
- [ ] Add CI with github actions.
- [ ] Build and publish docker image.

## Running locally

Install dependencies
```
go mod download
```

```
go run main.go -f fixtures/test/yaml_config.yml
```

### Running tests

Start the rabbitmq server

```
docker-compose up -d
```

Run the tests recursively for all subpackages

```
go test ./...
```
