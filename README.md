# Redis Clone in Go

A high-performance Redis-compatible in-memory data store built from scratch in Go.

This project was developed to understand how Redis works internally by implementing its networking model, protocol parsing, data structures, persistence, transactions, Pub/Sub, and eviction policies without relying on external Redis libraries.

---

## Features

### Networking

- TCP Server
- Linux `epoll` based Event Loop
- Non-blocking sockets
- Single-threaded event-driven architecture
- Persistent client connections (Keep Alive)
- Per-client input buffers
- Multiple concurrent clients

---

### RESP Protocol

- RESP2 Parser
- Incremental parsing
- Partial packet handling
- Multiple commands in a single TCP packet (Pipelining)
- RESP serialization

---

### String Commands

- PING
- SET
- GET
- DEL
- EXISTS
- EXPIRE
- TTL
- PERSIST
- APPEND
- STRLEN
- GETSET
- INCR
- DECR
- INCRBY
- DECRBY
- MSET
- MGET

---

### List Commands

- LPUSH
- RPUSH
- LPOP
- RPOP
- LLEN
- LRANGE

---

### Hash Commands

- HSET
- HGET
- HDEL
- HEXISTS
- HLEN
- HGETALL
- HKEYS
- HVALS
- HMGET
- HSETNX
- HSTRLEN
- HINCRBY

---

### Set Commands

- SADD
- SREM
- SISMEMBER
- SCARD
- SMEMBERS

---

### Sorted Set Commands

Implemented using a custom Skip List.

- ZADD
- ZREM
- ZSCORE
- ZCARD
- ZRANGE
- ZREVRANGE
- ZRANK
- ZREVRANK
- ZCOUNT

---

### Transactions

- MULTI
- EXEC
- DISCARD

Features

- Atomic execution
- Queued commands
- Ordered responses
- Error propagation

---

### Publish / Subscribe

- SUBSCRIBE
- PUBLISH

Features

- Multiple channels
- Multiple subscribers
- Automatic cleanup on client disconnect

---

### Persistence

Append Only File (AOF)

- Command logging
- Startup replay
- Database reconstruction
- AOF Rewrite (BGREWRITEAOF)

---

### Memory Management

Simple Least Recently Used (LRU) eviction policy.

- Last access tracking
- Automatic eviction when memory limit is reached

---

## Architecture

```
                        Clients
                           │
                           ▼
                 TCP Listening Socket
                           │
                           ▼
                     Linux epoll
                  (Event Driven Loop)
                           │
        ┌──────────────────┴──────────────────┐
        ▼                                     ▼
 Per Client Buffer                     Per Client Buffer
        │                                     │
        └──────────────┬──────────────────────┘
                       ▼
             Incremental RESP Parser
                       │
                       ▼
              Command Dispatcher
                       │
      ┌────────────────┼────────────────┐
      ▼                ▼                ▼
 Strings           Collections      Transactions
(Lists/Hashes/
 Sets/ZSets)
      │                │                │
      └────────────────┼────────────────┘
                       ▼
               In-Memory Database
                       │
          ┌────────────┴────────────┐
          ▼                         ▼
      AOF Persistence         Pub/Sub Engine
                       │
                       ▼
                 RESP Serialization
                       │
                       ▼
                     Clients
```

---

## Project Structure

```
redis-clone/
│
├── cmd/
│
├── server/
│
├── networking/
│   ├── epoll.go
│   ├── server.go
│   └── client.go
│
├── protocol/
│   ├── parser.go
│   ├── serializer.go
│   └── resp.go
│
├── commands/
│   ├── string.go
│   ├── list.go
│   ├── hash.go
│   ├── set.go
│   ├── sortedset.go
│   ├── transaction.go
│   ├── pubsub.go
│   └── general.go
│
├── cache/
│
├── persistence/
│
├── skiplist/
│
├── pubsub/
│
└── main.go
```

---

## Performance

Benchmarked using `redis-benchmark`.

```
SET

100,000 requests
50 concurrent clients

≈ 65,000 requests/sec
Average Latency: 0.716 ms
```

```
GET

100,000 requests
50 concurrent clients

≈ 174,000 requests/sec
Average Latency: 0.214 ms
```

---

## Example Usage

```bash
go run .
```

Open another terminal.

```bash
redis-cli
```

```
PING

PONG
```

```
SET name Shivam

OK
```

```
GET name

"Shivam"
```

Lists

```
LPUSH fruits Apple Mango Banana
LRANGE fruits 0 -1
```

Hashes

```
HSET user name Shivam role Developer
HGETALL user
```

Sets

```
SADD skills Go Redis Linux
SMEMBERS skills
```

Sorted Sets

```
ZADD leaderboard 100 Alice 95 Bob
ZRANGE leaderboard 0 -1 WITHSCORES
```

Transactions

```
MULTI
SET count 1
INCR count
EXEC
```

Publish / Subscribe

Client 1

```
SUBSCRIBE news
```

Client 2

```
PUBLISH news "Hello World"
```

---

## Benchmark

```bash
redis-benchmark -p 6379 -t set -n 100000 -c 50

redis-benchmark -p 6379 -t get -n 100000 -c 50
```

---

## Technologies Used

- Go
- Linux epoll
- TCP Sockets
- RESP Protocol
- Skip Lists
- Hash Maps
- Event Driven Architecture

---

## Learning Outcomes

While building this project, I gained hands-on experience with:

- Low-level socket programming
- Event-driven server architecture
- Linux epoll
- RESP protocol internals
- Skip List implementation
- Redis persistence
- Transaction execution
- Pub/Sub messaging
- Memory eviction strategies
- Network protocol parsing
- Concurrent client handling
- High-performance server design

---

## Future Improvements

- Replication
- RDB Snapshots
- Cluster Support
- Lua Scripting
- Streams
- HyperLogLog
- Bitmaps
- ACL
- WATCH command
- Memory-based eviction policies

---

## License

MIT License