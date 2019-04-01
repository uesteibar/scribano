# Scribano

Builds an [asyncapi](https://www.asyncapi.com/) documentation for your microservices
communicating through [rabbitmq](https://www.rabbitmq.com/).

It listens to all published amqp messages and keeps an updated asyncapi
compliant documenation served at `/asyncapi`.

## Roadmap

- [x] Add info and server sections to spec to make it valid asyncapi.
- [x] Extract configuration to file.
- [x] Support consuming from multiple configurable exchanges.
- [x] Allow loading configuration from url instead of local file.
- [x] Use postgres as database.
- [ ] Add exchange information to each topic via x-exchange attribute.
- [ ] Add CI.
- [ ] Build and publish docker image on push to master.

## Running locally

Install dependencies
```
go mod download
```

Load environmental variables
```
source .env.dev
```

```
go run main.go -f fixtures/test/yaml_config.yml
```

You can also fetch the configuration from a url
```
go run main.go -u https://raw.githubusercontent.com/uesteibar/asyncapi-watcher/master/fixtures/test/yaml_config.yml
```

### Running tests

Start the rabbitmq server

```
docker-compose up -d
```

Run the tests

```
make test
```
