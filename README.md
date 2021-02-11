# What is gRPC?
- An RPC is a Remote Procedure Call.
- In your CLIENT code, it looks like you're just calling a function directly on the SERVER.

![alt text](https://grpc.io/img/landing-2.svg)

### How to get started?
- At the core of gRPC, we need to define the messages and services using [Protocol Buffers](https://developers.google.com/protocol-buffers).
- The rest of the gRPC code will be generated for us and we'll have to provide an implementation for it. 
- One `.proto` file works for over 12 programming languages (server and client), and allows us to use a framework that scales to millions of RPC per seconds.

### Why Protocol Buffers?
- Protocol Buffers are language agnostic.
- Code can be generated for pretty much any language.
- Data is binary and efficiently serialized (small payloads).
- Very convenient for transporting a lot of data.
- Allows for easy API evolution using rules.

### Protocol Buffers role in gRPC
- Protocol Buffers is used to define the:
    * Messages (data, Response and Request).
    * Service (service name and RPC endpoints).

### What is HTTP/2?
- gRPC leverages HTTP/2 as a backbone for communications.
- HTTP/2 is the newer standard for internet communcations that address common pitfall of HTTP/1.1 on modern web pages.

### How HTTP/1.1 works
- HTTP/1.1 opens a new TCP connection to a server at each request.
- It does not compress headers (which are plaintext).
- It only works with Response / Request mechanism (no server push).
- Was originally omposed of two commands:
    * GET to ask for content.
    * POST to send content.

### How HTTP 2 works
- HTTP 2 supports multiplexing
    * THe client and server can push messages in parallel over the same TCP connection.
    * This greatly reduces latency.
- HTTP 2 supports server push
    * Servers can push streams (multiple messages) for one request from the client.
    * This saves round trips (latency).
- HTTP 2 supports headers compression.
- HTTP 2 is binary.
- HTTP 2 is secure (SSL is not required but recommended by default).

### 4 Types of API in gRPC

![alt text](https://i.ibb.co/9sn4Yxn/apiTypes.png)

- Unary is what a traditional API looks like (HTTP REST).
- HTTP 2 as we've seen, enables APIs to now have streaming capabilities.
- The server and the client can push multiple messages as part of one request.
- In gRPC it's very easy to define these APIs.

### What is an Unary API?
- Unary RPC calls are the basic Response / Reques.
- The client sends one message to the server and recieves one response from the server.
- Unary calls are very well suited for small data.
- Start with Unary when writing APIs and use streaming API if performance is an issue.

### What is a Server Streaming API?
- Server Streaming RPC API is a new kind of API enabled thanks to HTTP 2.
- The client will send one message to the server and will receive many responses from the server.
- Streaming Servers are well suited
   * When the server needs to send a lot of data (big data).
   * When the server needs to `PUSH` data to the client continuously without having the client request all over again.
- In gPRC, Server Streaming Calls are defined using the keyword `stream`.
- As for each RPC call, we have to define a Request message and a Response message..

### What is a Client Streaming API?
- Server Streaming RPC API is a new kind of API enabled thanks to HTTP 2.
- The client will send many messages to the server and will receive one response from the server.
- Streaming Client are well suited
   * When the client needs to send a lot of data (big data).
   * When the server processing is expensive and should happen as the client sends data.
   * When the client needs to `PUSH` data to the server without really expecting a response.
- In gPRC, Client Streaming Calls are defined using the keyword `stream`.
- As for each RPC call, we have to define a Request message and a Response message.

### What is a Bi Directional Streaming API?
- Bi Directional Streaming RPC API is a new kind of API enabled thanks to HTTP 2.
- The client will send many messages to the server and will receive many responses from the server.
- The number of requests and responses does not have to much.
- Bi Directional Streaming RPC is well suited
   * When the client and the server needs to send a lot of data asynchronously.
   * `Chat` protocol.
   * Long running connections.
- In gPRC, Bi Directional Streaming API Calls are defined using the keyword `stream`, twice.
- As for each RPC call, we have to define a Request message and a Response message.

### Scalability in gRPC
- gRPC Servers are asynchronous by default. This means they do not block threads on request. Therefore each gRPC server can serve millions of requests in parallel.
- gRPC Clients can be asynchronous or synchronous (blocking). The client decides which model works best for the performance needs.
- gRPC Clients can perform client side load balancing.

### Security in gRPC
- By default gRPC strongly advocates for the use of SSL in any API.
- Each language will provide an API to load gRPC with the required certificates and provide encryption capability out of the box.
- Additionally using Interceptors, we can also provide authntication.

### gRPC vs REST

| gRPC | REST  |
|---|---|
| Protocol Buffers - smaller, faster  | JSON - text based, slower, bigger  |
| HTTP 2 (low latency)  | HTTP1.1 (higher latency)  |
| Bidirectional & Async  | Client => Server requests only  |
| Stream Support  | Response / Request support only  |
| API Oriented (no constraints - free design)  | CRUD Oriented / POST GET PUT DELETE  |
| Code Generation through Protocol Buffers in any language  | Code Generation though OpenAPI / Swagger (add-on)  |
| RPC based - gRPC does the plumbing  | HTTP verb based - we have to write the plumbing or use 3rd party library  |

### Why use gRPC
- Easy code definition.
- Uses a modern, low latency HTTP 2 transport mechanism.
- SSL security is built in.
- Support for streaming APIs for maximum performance.
- gRPC is API oriented, instead of Resource Oriented like REST.

### Dependencies Setup
- `go get -u google.golang.org/grpc`
- `go install google.golang.org/protobuf/cmd/protoc-gen-go`

### Generate Go code through `.proto` file
```
protoc greet/greetpb/greet.proto --go_out=plugins=grpc:.
protoc calculator/calculatorpb/calculator.proto --go_out=plugins=grpc:.
```

### gRPC Error Handling
- `[inline link](https://avi.im/grpc-errors/)`
- `[inline link](https://grpc.io/docs/guides/error/)`

### What is gRPC Deadlines?
- Deadlines allow gRPC clients to specify how long they are willing to wait for an RPC to complete before the RPC is returned with the error DEADLINE_EXCEEDED.
- The gRPC documentation recommends you set a deadline for all client RPC calls.
- The server should check if the deadling has exceeded and calcel the work it is doing.
- Deadlines are propagated across if gRPC calls are chained. 
   * A => B => C  (Deadline for A is passed to B and then passed to C)
- More info: `[inline link](https://grpc.io/blog/deadlines/)`

### SSL Encryption in gRPC
- In production, gRPC calls should be running with encryption enabled.
- This is done by generating SSL certificates.
- SSL allows communication to be secure end-to-end and ensuring no Man in the middle attack can be performed.
![alt text](https://i.ibb.co/127Yvzg/gRPCencr.png)

#### Hands on SSL Encryption in gRPC with Go
- Setup a certificate authority.
- Setup a server certificate.
   * Setup the server to use TLS.
- Sign a server certificate.
   * Setup the client to connect securely over TLS.
- More info: `[inline link](https://grpc.io/docs/guides/auth/)`

#### How SSL works
- When you communicate over the internet, your data is visible by all the servers that transfer your packet.
- Any router in the middle can view the packets you're sending using PLAINTEXT.
![alt text](https://i.ibb.co/3NssRWC/http.png)

- SSL allows clients and servers to encrypt packets.
![alt text](https://i.ibb.co/dPwZpyY/https.png)



























