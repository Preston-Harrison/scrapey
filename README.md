# Forward Proxy Server

## Overview

This repository contains a forward proxy server that enables any computer operating within a network that allows only outbound traffic to act as a forward proxy. It achieves this by having the proxy server open a TCP connection with the main server, and waiting for a request on the open connection. A client will proxy a request to the main server, which will use one of the open connections to send the request to the actual proxy.

## How it Works

The forward proxy server is divided into two main components: the server and the proxy.

### Server (/server)

The server component listens for incoming proxy requests from clients. It is responsible for:

- Accepting incoming proxy requests from clients.
- Extracting the target host from the client's request.
- Establishing a TCP connection with the target host.
- Forwarding raw bytes from the client to the target host and vice versa.

### Proxy (/proxy)

The proxy component initiates a TCP connection with the server and waits for the server to respond with a host. It is responsible for:

- Establishing a TCP connection with the server.
- Waiting for the server to provide the target host information.
- Forwarding raw bytes from the server to the target host.

### Configuration
You can configure the server's listening ports by modifying the `SERVER_PORT` and `PROXY_PORT` environment variables. Similarly, you can configure the proxy's connection details by modifying the `SERVER_HOST` environment variable.

### Usage
```bash
# Run this on server with publicly accessable IP & ports
./server # Defaults to ports 5000 & 5001
```

```bash
# Run this on any computer with outbound traffic enabled
AUTH_TOKEN=some_long_secret SERVER_HOST=your_host:your_port ./proxy
```

```go
// Make your request. Note at the moment, username is ignored. Password should be the auth token.
func main() {
	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(&url.URL{
				Scheme: "https", // Using http exposes the auth token.
				User: url.UserPassword("", "your_auth_token"),
				Host: "your_host:your_port",
			}),
		},
	}
	response, err := httpClient.Get("https://example.com/")
}
```

## Proxy-Server interractions
### Initial handshake
1. The proxy connects via TCP with the server.
2. The proxy sends an auth token in a single frame.
3. The server responds acknowledging the auth token is valid.
4. The connection remains open until the server is ready to send a request.

### Handling a request
*Note, only HTTPS is supported as it's easier to implement - encryped bytes are just streamed to/from the client and proxy*  
1. The client sends a request to the server.
2. The server checks the `Proxy-Authentication` header for the auth token.
3. The server finds an open proxy connection with a corresponding auth token, and aborts if there are none.
4. The server sends a single frame of the target host to the proxy.
5. The proxy responds if they were able to connect with the target.
6. The proxy opens a new connection with the server, going through the initial handshake process again in preparation for a new request.
6. The server copies bytes to/from the proxy and client, until either closes the connection.