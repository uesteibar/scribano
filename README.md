# Scribano

Builds an [asyncapi](https://www.asyncapi.com/) documentation for your microservices
communicating through [rabbitmq](https://www.rabbitmq.com/).

It listens to all published amqp messages and keeps an updated asyncapi
compliant documenation served at `/asyncapi`.

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
go run main.go -u https://raw.githubusercontent.com/uesteibar/scribano/master/fixtures/test/yaml_config.yml
```

## Running on docker

```
docker run -e PG_URL='postgresql://postgres:postgres@localhost:5433/asyncapi' uesteibar/scribano ./scribano -u https://raw.githubusercontent.com/uesteibar/scribano/master/fixtures/test/yaml_config.yml
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
