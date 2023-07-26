# AuthProxy

AuthProxy is a Go application that acts as an authentication proxy, allowing secure communication between clients and backend servers. It performs authentication checks by validating session cookies and authorizes requests based on specified cookie keys. The application supports WebSocket connections and HTTP requests, forwarding them to the appropriate backend server.

## Features

- WebSocket and HTTP support: AuthProxy can handle WebSocket connections as well as standard HTTP requests.
- Authentication and Authorization: It checks session cookies against specified cookie keys to ensure proper authorization.
- Cross-Origin Resource Sharing (CORS) support: AuthProxy sets the necessary CORS headers to allow cross-origin requests.

## Configuration

Before running AuthProxy, you need to set the following environment variables:

- `COOKIE_KEYS`: The list of cookie keys for authentication, separated by commas (e.g., `session1, session2`).
- `REDIS_ADDR`: The Redis server address (default: `127.0.0.1:6379`).
- `REDIS_PASSWORD`: The Redis server password (if applicable, default: empty string).
- `REDIS_DB`: The Redis database number (default: `0`).
- `REDIRECT_URL`: The URL to redirect when errors occur (e.g., `https://www.xxx.com`).
- `HEARTBEAT_TIME`: WebSocket connection heartbeat interval in seconds (default: `30` seconds).

## Building and Running

To build the AuthProxy application, use the following command:

```bash
go build -o auth-proxy ./path/to/your/AuthProxy
```

Ensure you have all the necessary dependencies installed.

To run AuthProxy, execute the following command:

~~~shell
./auth-proxy
~~~

By default, AuthProxy listens on port 80. You can modify the port by changing the `http.ListenAndServe` function in the `main` function.

## Contributing

If you'd like to contribute to AuthProxy, feel free to open a pull request or submit issues for bug fixes, improvements, or new features.

