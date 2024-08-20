## Overview

**Link** is a powerful TCP connection management tool that simplifies NAT traversal, TCP forwarding and more. By seamlessly integrating three distinct running modes within a single binary file, Link bridges the gap between different network environments, redirecting services and handling TCP connections seamlessly, ensuring reliable network connectivity and ideal network environment. Also with highly integrated authorization handling, Link empowers you to efficiently manage user permissions and establish uninterrupted data flow, ensuring that sensitive resources remain protected while applications maintain high performance and responsiveness.

## Features

- **Unified Operation**: Link can function as a server, client, or broker, three roles from a single executable file.

- **Authorization Handling**: By IP address handling, Link ensures only authorized users gain access to sensitive resources.

- **In-Memory Certificate**: Provides a self-signed HTTPS certificate with a one-year validity, stored entirely in memory.

- **Auto Reconnection**: Providing robust short-term reconnection capabilities, ensuring uninterrupted service.

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

Note that only `server` and  `broker` mode support authorization Handling, which you can just add auth entry after `#`. For example:

```
server://linkAddr/targetAddr#authScheme//authAddr/secretPath
broker://linkAddr/targetAddr#authScheme//authAddr/secretPath
```

- **authScheme**: The option allows you to choose between using HTTP or HTTPS.
- **authAddr**: The server address and port designated for authorization handling.
- **secretPath**: The secret endpoint for processing authorization requests.

### Server Mode

- `linkAddr`: The address for accepting client connections. For example, `:10101`.
- `targetAddr`: The address for listening to external connections. For example, `:10022`.

**Run as Server**

```bash
./link server://:10101/:10022
```

- This command will listen for client connections on port `10101` , listen and forward data to port `10022`.

**Run as Server with authorization**

```bash
./link server://:10101/:10022#https://:8443/server
```

- The server handles authorization at `https://server_ip:8443/server`, on your visit and your IP logged.
- The server will listen for client connections on port `10101` , listen and forward data to port `10022`.

### Client Mode

- `linkAddr`: The address of the server to connect to. For example, `server_ip:10101`.
- `targetAddr`: The address of the target service to connect to. For example, `127.0.0.1:22`.

**Run as Client**

```bash
./link client://server_ip:10101/127.0.0.1:22
```

- This command will establish link with `server_ip:10101` , connect and forward data to `127.0.0.1:22`.

### Broker Mode

- `linkAddr`: The address for accepting client connections. For example, `:10101`.
- `targetAddr`: The address of the target service to connect to. For example, `127.0.0.1:22`.

**Run as Broker**

```bash
./link broker://:10101/127.0.0.1:22
```

- This command will listen for client connections on port `10101` , connect and forward data to `127.0.0.1:22`.

**Run as Broker with authorization**

```bash
./link broker://:10101/127.0.0.1:22#https://:8443/broker
```

- The server handles authorization at `https://server_ip:8443/broker`, on your visit and your IP logged.
- The server will listen for client connections on port `10101` , connect and forward data to `127.0.0.1:22`.

## Container Usage

You can also run **Link** using a Docker container. The image is available at [ghcr.io/raymondragon/link](https://ghcr.io/raymondragon/link).

To run the container in server mode with or without authorization:

```bash
docker run --rm ghcr.io/raymondragon/link server://:10101/:10022#https://:8443/server
```

```bash
docker run --rm ghcr.io/raymondragon/link server://:10101/:10022
```

To run the container in client mode:

```bash
docker run --rm ghcr.io/raymondragon/link client://server_ip:10101/127.0.0.1:22
```

To run the container in server mode with or without authorization:

```bash
docker run --rm ghcr.io/raymondragon/link broker://:10101/127.0.0.1:22#https://:8443/broker
```

```bash
docker run --rm ghcr.io/raymondragon/link broker://:10101/127.0.0.1:22
```

## License

This project is licensed under the MIT License. See the LICENSE file for details.
