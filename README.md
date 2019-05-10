# Scribano

[![CircleCI](https://circleci.com/gh/uesteibar/scribano/tree/master.svg?style=svg)](https://circleci.com/gh/uesteibar/scribano/tree/master)

Builds an [asyncapi](https://www.asyncapi.com/) documentation for your microservices
communicating through [rabbitmq](https://www.rabbitmq.com/).

It listens to all published amqp messages and keeps an updated asyncapi
compliant documentation served at `/asyncapi`.

Sample configuration can be found [here](https://github.com/uesteibar/scribano/blob/master/fixtures/test/yaml_config.yml).

## How does it work?

Scribano will subscribe to a set of configured RabbitMQ exchanges, and will start
consuming events from there.

When a message is consumed, scribano will infer a structure for the content.
For example, after consuming a message on the `some.key` routing key with the following payload:
```json
{
  "name": "infer type",
  "age": 27,
  "grade": 9.5,
  "canDrive": false,
  "birthDate": "1991-08-29",
  "lastLogin": "2015-06-10T13:23:30-08:00",
  "address": null,
  "emptyHash": {},
  "fines": [],
  "emptyHashes": [{}],
  "matrix": [
    [1, 2, 3],
    [3, 2, 1]
  ],
  "friends": [
    { "name": "pepe" },
    { "name": "gotera" }
  ],
  "car": {
    "brand": "mercedes"
  }
}
```

The following spec will be served on `/asyncapi`

```
{
  "asyncapi": "1.0.0",
  "info": {
    "title": "",
    "version": ""
  },
  "topics": {
    "some.key": {
      "publish": {
        "$ref": "#/components/messages/SomeKey"
      },
      "x-exchange": "/my-exchange"
    }
  },
  "components": {
    "messages": {
      "SomeKey": {
        "payload": {
          "type": "object",
          "properties": {
            "name": {
              "type": "string"
            },
            "age": {
              "type": "integer"
            },
            "grade": {
              "type": "number"
            },
            "birthDate": {
              "type": "string",
              "format": "date"
            },
            "lastLogin": {
              "type": "string",
              "format": "date-time"
            },
            "address": {
              "type": "string"
            },
            "emptyHash": {
              "type": "object"
            },
            "fines": {
              "type": "array",
              "items": {
                "type": "string"
              }
            },
            "emptyHashes": {
              "type": "array",
              "items": {
                "type": "object"
              }
            },
            "matrix": {
              "type": "array",
              "items": {
                "type": "array",
                "items": {
                  "type": "integer"
                }
              }
            },
            "friends": {
              "type": "array",
              "items": {
                "type": "object",
                "properties": {
                  "birthDate": {
                    "type": "string",
                    "format": "date"
                  },
                  "name": {
                    "type": "string"
                  }
                }
              }
            },
            "car": {
              "type": "object",
              "properties": {
                "brand": {
                  "type": "string"
                }
              }
            }
          }
        }
      }
    }
  }
}
```

### Type inference

|            value            | inferred type |   format  |                                  notes                                 |
|:---------------------------:|:-------------:|:---------:|:----------------------------------------------------------------------:|
|        "some string"        |     string    |           |                                                                        |
|         "1991-08-29"        |     string    |    date   | Format must match exactly, otherwise it is considered a regular string |
| "2015-06-10T13:23:30-08:00" |     string    | date-time | Format must match exactly, otherwise it is considered a regular string |
|         true / false        |    boolean    |           |                                                                        |
|            [...]            |     array     |           | Type for the array values is inferred by checking on the first element |
|            {...}            |     object    |           |     Types for the fields inside the object are checked recursively     |
|             null            |     string    |           |                                                                        |

### Optional fields

Fields are set as optional using the custom `x-optional: true` attribute.

When a fields is received in a message and then not received in a
subsequent message, it is considered optional as the presence of that
field cannot be guaranteed.

When a field is received in a message and was never previously received,
it is also considered optional.


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
