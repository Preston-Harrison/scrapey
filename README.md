# Forward Proxy Server

## Overview

This repository contains a forward proxy server that enables any computer operating within a network that allows only outbound traffic to act as a forward proxy. It achieves this by accepting proxy requests and forwarding raw bytes via TCP to a target host. The proxy server opens a TCP connection with the client, waits for a response containing a host, and then forwards bytes from the client to the target host.

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
./bin/server # Defaults to ports 5000 & 5001
```

```bash
# Run this on any computer with outbound traffic enabled
CHANNEL_KEY=some_long_secret SERVER_HOST=your_host:your_port ./bin/proxy
```

```go
// Make your request. Note at the moment, username is ignored. Password should be the channel key.
// This will change as it's dumb.
func main() {
	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(&url.URL{
				Scheme: "https", // Using http exposes channel key.
				User: url.UserPassword("", "your_channel_key"),
				Host: "your_host:your_port",
			}),
		},
	}
	response, err := httpClient.Get("https://example.com/")
}
```
