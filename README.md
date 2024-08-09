# Overview

**Link** is an innovative Go program that simplifies NAT traversal and TCP connection forwarding. By seamlessly integrating server and client capabilities, Link bridges the gap between different network environments, ensuring reliable connectivity and uninterrupted data flow.

## Features

- **Unified Operation**: Link can operate as both a server and a client with a single executable. Simply configure the desired role using a URL, and the program will adapt to handle the specified connections. This flexibility streamlines deployment and reduces the need for multiple tools.

- **Auto Reconnection**: Link provides robust short-term reconnection capabilities. If either end of the connection experiences a disruption or dropout, the other end remains operational, ensuring uninterrupted service.

- **Connection Updates**: In scenarios where a client’s target service connection is interrupted and needs refreshing, Link supports real-time connection updates. The server can synchronize and provide the latest connection details, reducing downtime and maintaining connectivity.

## Usage

To run the program, provide a URL specifying the mode and connection addresses. The URL format is as follows:

```
server://linkAddr#targetAddr
client://linkAddr#targetAddr
```

### Server Mode

- `linkAddr`: The address for accepting client connections. For example, `:10101`.
- `targetAddr`: The address for listening to external connections. For example, `:10022`.

**Run as Server**

```bash
./link server://:10101#:10022
```

This command will listen for client connections on port `10101` , listen and forward data to port `10022`.

### Client Mode

- `linkAddr`: The address of the server to connect to. For example, `server_ip:10101`.
- `targetAddr`: The address of the target service to connect to. For example, `127.0.0.1:22`.

**Run as Client**

```bash
./link client://server_ip:10101#127.0.0.1:22
```

This command will establish link with `server_ip:10101` , connect and forward data to `127.0.0.1:22`.

## Container Usage

You can also run **Link** using a Docker container. The image is available at [ghcr.io/raymondragon/link](https://ghcr.io/raymondragon/link).

### Running with Docker

To run the container in server mode:

```bash
docker run --rm ghcr.io/raymondragon/link server://:10101#:10022
```

To run the container in client mode:

```bash
docker run --rm ghcr.io/raymondragon/link client://server_ip:10101#127.0.0.1:22
```

## License

This project is licensed under the MIT License. See the LICENSE file for details.
