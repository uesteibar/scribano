# Thoth

## Running locally

Install [govendor](https://github.com/kardianos/govendor)
```
go get github.com/kardianos/govendor
```

Install dependencies
```
govendor sync
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
