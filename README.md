## Overview

**Link** is a powerful TCP connection management tool that simplifies NAT traversal, TCP forwarding and more. By seamlessly integrating three distinct running modes within a single binary file, Link bridges the gap between different network environments, redirecting services and handling TCP connections seamlessly, ensuring reliable network connectivity and ideal network environment.

## Features

- **Unified Operation**: Link can function as a server, client, or broker, three roles from a single executable file.

- **Auto Reconnection**: Link provides robust short-term reconnection capabilities, ensuring uninterrupted service.

- **Connection Updates**: In scenarios where connection is interrupted, Link supports real-time connection updates.

- **Service Forwarding**: Efficiently manage and redirect TCP connections from one service to entrypoints everywhere.

- **No External Dependencies**: Entirely built using Go's standard library, ensuring a lightweight and efficient solution.

## Usage

To run the program, provide a URL specifying the mode and connection addresses. The URL format is as follows:

```
server://linkAddr/targetAddr
client://linkAddr/targetAddr
broker://linkAddr/targetAddr
```

### Server Mode

- `linkAddr`: The address for accepting client connections. For example, `:10101`.
- `targetAddr`: The address for listening to external connections. For example, `:10022`.

**Run as Server**

```bash
./link server://:10101/:10022
```

This command will listen for client connections on port `10101` , listen and forward data to port `10022`.

### Client Mode

- `linkAddr`: The address of the server to connect to. For example, `server_ip:10101`.
- `targetAddr`: The address of the target service to connect to. For example, `127.0.0.1:22`.

**Run as Client**

```bash
./link client://server_ip:10101/127.0.0.1:22
```

This command will establish link with `server_ip:10101` , connect and forward data to `127.0.0.1:22`.

### Broker Mode

- `linkAddr`: The address for accepting client connections. For example, `:10101`.
- `targetAddr`: The address of the target service to connect to. For example, `127.0.0.1:22`.

**Run as Broker**

```bash
./link broker://:10101/127.0.0.1:22
```

This command will listen for client connections on port `10101` , connect and forward data to `127.0.0.1:22`.

## Container Usage

You can also run **Link** using a Docker container. The image is available at [ghcr.io/raymondragon/link](https://ghcr.io/raymondragon/link).

To run the container in server mode:

```bash
docker run --rm ghcr.io/raymondragon/link server://:10101/:10022
```

To run the container in client mode:

```bash
docker run --rm ghcr.io/raymondragon/link client://server_ip:10101/127.0.0.1:22
```

To run the container in broker mode:

```bash
docker run --rm ghcr.io/raymondragon/link broker://:10101/127.0.0.1:22
```

## License

This project is licensed under the MIT License. See the LICENSE file for details.
