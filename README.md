# What is gRPC?
- An RPC is a Remote Procedure Call.
- In your CLIENT code, it looks like you're just calling a function directly on the SERVER.

![alt text](https://grpc.io/img/landing-2.svg)

### How to get started?
- At the core of gRPC, we need to define the messages and services using [Protocol Buffers](https://developers.google.com/protocol-buffers)
- The rest of the gRPC code will be generated for us and we'll have to provide an implementation for it. 
- One `.proto` file works for over 12 programming languages (server and client), and allows us to use a framework that scales to millions of RPC per seconds.

### Why Protocol Buffers?
- Protocol Buffers are language agnostic.
- Code can be generated for pretty much any language.
- Data is binary and efficiently serialized (small payloads).
- Very convenient for transporting a lot of data.
- Allows for easy API evolution using rules.
