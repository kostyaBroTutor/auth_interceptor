# Auth Interceptor

This repository stores a implementation of the gRPC interceptor,
that check auth roles or permissions, and authenticate user. \
This repository created for the [habr post(TODO link)](example.com).

## Cloning repository

This repository contains submodules, 
so you need to clone it with `--recursive` flag:

```bash
git clone --recursive https://github.com/kostyaBroTutor/auth_interceptor.git
```
You can run `make generate` for regenerating protobuf files, mocks, and other.

## 

## Short description

This repository contains an implementation of a gRPC interceptor in Go 
for handling authorization and authentication in gRPC services. 
The primary function of this interceptor is 
to verify tokens present in the metadata of gRPC requests and 
determine if the request has appropriate permissions 
to access a particular service method. 

The interceptor uses an AuthClient to authenticate tokens and 
retrieve relevant information about the user. 
This user information, alongside required permissions for service methods, 
are cached within the interceptor to prevent unnecessary re-authentication. 
The interceptor uses gRPC reflection 
to fetch authentication requirements of service methods, 
providing an efficient mechanism 
for managing authentication and authorization in your gRPC services.

## Usage example 

To use this interceptor in your gRPC server, you must first initialize it:

```go
authClient := ... // Initialize your AuthClient
fullServiceName := "your.service.Name" // Your service name
reflectionClient := ... // Initialize your reflection client

interceptor := NewAuthUnaryServiceInterceptor(
	authClient,
	fullServiceName,
	reflectionClient,
)
```

Then, when setting up your gRPC server, you include the interceptor:

```go
grpcServer := grpc.NewServer(
	grpc.UnaryInterceptor(interceptor),
)
```

Now, any incoming requests to your gRPC server 
will first pass through the interceptor. 
If the request contains a valid token with the required permissions, 
it will be forwarded to the appropriate service method. 
If the request lacks a valid token or required permissions,
it will be rejected, and an error will be returned to the client. 

## Contributing

If you have suggestions for how this code could be improved, 
or want to report a bug, open an issue! 
Contributions are more than welcome.

## License

Distributed under the MIT License. See `LICENSE` for more information.

