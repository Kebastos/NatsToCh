# NatsToCh
Simple data transfer application from Nats to ClickHouse.

[![Go Report Card](https://goreportcard.com/badge/github.com/Kebastos/NatsToCh)](https://goreportcard.com/report/github.com/Kebastos/NatsToCh)
[![test](https://github.com/Kebastos/NatsToCh/actions/workflows/go-test.yml/badge.svg?branch=main)](https://github.com/Kebastos/NatsToCh/actions/workflows/go-test.yml)
[![Go Coverage](https://github.com/Kebastos/NatsToCh/wiki/coverage.svg)](https://raw.githack.com/wiki/Kebastos/NatsToCh/coverage.html)
[![golangci-lint](https://github.com/Kebastos/NatsToCh/actions/workflows/golangci-lint.yaml/badge.svg?branch=main)](https://github.com/Kebastos/NatsToCh/actions/workflows/golangci-lint.yaml)


## Simple example
1. Run [Nats](https://hub.docker.com/_/nats)
2. Run [Clickhouse](https://hub.docker.com/r/clickhouse/clickhouse-server)
3. Create a new table with name - NatsMessages
    - Table must contain five required columns

      | Name     | DataType | Description                                                     |
      |----------|----------|-----------------------------------------------------------------|
      | Id       | String   | Unique UUID for every row                                       |
      | ClientId | String   | The name is taken from the configuration file, Nats:client_name |
      | Subject  | String   | Subject where the message was received from                     |
      | CreateDateTime | DateTime | Time of receipt                                           |
      | Content  | String   | Message content                                                 |

```clickhouse
create table if not exists NatsMessages
(
    Id             String,
    ClientId       String,
    Subject        String,
    CreateDateTime DateTime,
    Content        String
)
    engine = MergeTree PARTITION BY toStartOfWeek(CreateDateTime)
        ORDER BY (Subject, CreateDateTime)
        TTL CreateDateTime + toIntervalMonth(1)
        SETTINGS index_granularity = 8192;
```
   
4. Clone repository, build and run it with default configuration
```bash
git clone https://github.com/Kebastos/NatsToCh.git
cd ./natsToCh
go build .
./nats2ch
```


## Settings

In this section you can configure subjects.
```yaml
subjects: [
    {
      ## Name of subject
      name: "test",
      ## Name of queue
      ......
    },
    {
      .......
    },
]
```
You may configure two different ways for inserting data into the Clickhouse.
The first option (it is not recommended due to frequent inserts) - when ```use_buffer``` is ```false```.
In this case nats2ch will get messages and insert data into the Clickhouse.
```yaml
{
  ## Name of subject
  name: "test",
  ## Name of queue
  queue: "testQueue",
  ## Name of table for insert
  table_name: "testCash",
  ## Enable Buffer mode, when messages from Nats stacks in buffer.
  use_buffer: false,
}
```
The second option - when ```use_buffer``` is ```true```.
In this case you must configure section with name - ```buffer_config``` and when nats2ch will get messages and put it into buffer.
Buffer will drain after one of two events:
1. Buffers size greater than ```max_size``` option.
2. ```max_wait``` is expired. 

This can help you avoid problems with excessive memory usage.

For example:
```yaml
{
    ......
    use_buffer: true,
    buffer_config: {
      ## Max size of buffer for starting to drain. Must be greater than 0.
      max_size: 10,
      ## Max timeout of buffer for starting to drain. Must be greater than 0.
      max_wait: 60s,
}
```

