# AsyncAPI watcher

Builds an [asyncapi](https://www.asyncapi.com/) documentation for your microservices
communicating through [rabbitmq](https://www.rabbitmq.com/).

It listens to all published amqp messages and keeps an updated asyncapi
compliant documenation served at `/asyncapi`.

## Roadmap

- [ ] Add info and server sections to spec to make it valid asyncapi.
- [ ] Configure exchanges through env vars.
- [ ] Use postgres as database and configure through env var.
- [ ] Add CI with github actions.
- [ ] Build and publish docker image.
- [ ] Adapt to asyncapi 2.0.

## Running locally

Install dependencies
```
go mod download
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
