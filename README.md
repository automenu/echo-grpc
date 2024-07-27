# echo-grpc

The `echo-grpc` is a tiny gRPC service that echos what you start it with!

The service is built with [Connect][connect].
Its API uses [Protocol Buffers][proto], and the service
supports the [gRPC][grpc-protocol], [gRPC-Web][grpcweb-protocol], and [Connect protocol][connect-protocol].

## Getting started

### Prerequisites

- [Go][go]. For installation instructions, see Go’s [Getting Started][go-install] guide.
- [Buf CLI][buf-cli]. For installation instructions, see Buf's [Install the Buf CLI][buf-install] guide.

### Start the server

```sh
# Resolve Go dependencies
go mod tidy

# Start the server
make dev
```

### Make request

In a separate terminal window, in the root working directory, hit the API using the generated client:

```sh
make ping
```

If you have a recent version of cURL installed,
you can make a HTTP/1.1 POST request with a JSON payload.

```sh
curl \
    --header "Content-Type: application/json" \
    --data '{"message": "Hello, testing gRPC over HTTP..."}' \
    http://localhost:8080/echo.v1.EchoAPI/Echo
```

If you have a recent version of gRPCurl installed,
you can make gRPC requests without using a generated client.

```sh
grpcurl \
    -protoset <(buf build -o -) -plaintext \
    -d '{"message": "Hello, testing gRPC using gRPCurl..."}' \
    localhost:8080 echo.v1.EchoAPI/Echo
```

[buf-cli]: https://buf.build/docs/ecosystem/cli-overview
[buf-install]: https://buf.build/docs/installation
[connect]: https://github.com/connectrpc/connect-go
[connect-protocol]: https://connectrpc.com/docs/protocol
[go]: https://golang.org/
[go-install]: https://golang.org/doc/install
[grpc-protocol]: https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md
[grpcweb-protocol]: https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-WEB.md
[proto]: https://protobuf.dev/overview/
