# AsyncAPI watcher

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
