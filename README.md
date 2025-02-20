# Golang Redis Clone

This project is a simple Redis clone written in Go, which aims to replicate the basic functionalities of Redis, including support for various commands and data structures. The server responds to commands sent via a TCP connection, making it compatible with the `redis-cli`.

## Features

- Basic commands for strings, hashes, and lists.
- Support for commands such as `SET`, `GET`, `DEL`, `EXISTS`, `HSET`, `HGET`, `HDEL`, `HEXISTS`, and more.
- Concurrent access to data structures using mutexes for thread safety.
- Lightweight and easy to use.

#
## Installation

To get started with the Redis Clone, follow these steps:

1. Clone the repository:
   ```bash
   git clone https://github.com/joaogabsoaresf/golang-github.com/joaogabsoaresf/golang-redis-clone.git
   cd github.com/joaogabsoaresf/golang-redis-clone
   ```

2. Ensure you have Go installed on your machine. You can download it from [https://go.dev/dl/](https://go.dev/dl/)

3. Run the server:
    ```bash
    go run main.go
    ```
4. The server will start listening on port `6371` (You can use whatever port you want, i used `6371` cause my `6379` was in used)

#
## Usage
You can interact wiht Redis clone using the `redis-cli -p ${PORT}` or any Redis client. Here are some example commands:

```bash
127.0.0.1:6371> SET key "value"
OK
127.0.0.1:6371> GET key
"value"
127.0.0.1:6371> EXISTS key
1
127.0.0.1:6371> DEL key
1
127.0.0.1:6371> EXISTS key
0

127.0.0.1:6371> HSET users u1 "Joao"
OK
127.0.0.1:6371> HGET users u1
"Joao"
127.0.0.1:6371> HEXISTS users u1
1
127.0.0.1:6371> HDEL users u1
1
```

#
## Commands Supported
- Strings
    - SET key value
    - GET key
    - DEL key key2...
    - EXISTS key
- Hashes - O(1) for Read and Write
    - HSET hash field value
    - HGET hash field
    - HDEL hash field
    - HEXISTS hash field

#
## Persistence

The Redis Clone supports data persistence by saving all commands in RESP format to a file. When the server starts, it reads this file to restore the previous state of the data in memory. This ensures that all items are retained even after the server is restarted.

### How It Works

- Each command executed by the server is logged in a designated file in RESP format.
- Upon startup, the server reads this file and executes the commands sequentially to populate the in-memory data structures.
- This mechanism allows for data recovery and ensures that the server can resume its operations with the same dataset after a restart.

### Usage

To enable persistence:
1. Ensure that the server has write permissions to the directory where it saves the RESP file.
2. The server automatically handles the reading and writing of commands to the file, so you don’t need to perform any additional setup.

This feature enhances the reliability of your Redis Clone
