# kafka-proto-monitor

🛠️ A simple CLI debug utility for monitoring Kafka topics for the specific
messages and decode them using protobuf definitions.

Under the hood, `kafka-proto-monitor` relies on [protoreflect](https://github.com/jhump/protoreflect) in order to use raw (not compiled) proto files in runtime for unmarshaling Kafka messages.

## Installation

```bash
go install github.com/ri-nat/kafka-proto-monitor@latest
```

## Usage

```
Usage:
  kafka-proto-monitor [OPTIONS]

Application Options:
  -b, --broker=               Kafka broker URL (you can use this option multiple times) (default: localhost:9092)
  -t, --topic=                Kafka topic to read from (you can use this option multiple times)
  -a, --read-earliest         Read topic starting from the beginning (false by default)
  -r, --print-headers         Print message headers (false by default)
  -p, --proto-file=           Path to proto file
  -m, --proto-message=        Proto message to use
  -e, --proto-message-header= Name of Kafka message header, that contains message's proto name

Help Options:
  -h, --help                  Show this help message
```

## Examples

```bash
kafka-proto-monitor -p ./service.proto -m 'app.users.UserCreated' -t users -a -r
```

* Use `app.users.UserCreated` message from `service.proto`
* Connect to Kafka on `localhost:9092` (default)
* Read from topic `users`
* Read messages from the beginning of the topic (`-a`)
* Print message headers (`-r`)

---

```bash
kafka-proto-monitor -p ./service.proto -e proto-name -b kafka:9092 -t users
```

* Use all messages from `service.proto`
* Use `proto-name` header content as proto message name
* Connect to Kafka on `kafka:9092`
* Read from topic `users`
* Read only new messages (no `-a` flag)
* Do not print message headers (no `-r` flag)

## TODO

* Display help message on startup if no arguments provided
* Filter messages by header value
* Filter messages by parsed object attributes
* Construct proto name programmatically
* Be more verbose about message unmarshaling errors
* Use multiple proto files at once
* Filter printable attributes out
* Write more tests

## Licensing

This software is licensed under the MIT License. See [LICENSE](https://github.com/ri-nat/kafka-proto-monitor/blob/master/LICENSE) for the full license text.
