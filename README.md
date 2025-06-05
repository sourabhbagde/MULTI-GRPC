
# multi-grpc

This project demonstrates the ability to implement server-client pairs using the same gRPC protocol definition file (proto) but with different programming languages. The project includes two servers, one written in Python and the other in Go, and two clients, also written in Python and Go.

### Library Management System

library.proto file with message definitions and service methods is defined.

## Key Concepts
Unary RPC: Add a new book to the catalog (AddBook(BookRequest) returns (BookResponse))
Server Streaming RPC: Get a list of available books (ListAvailableBooks(EmptyRequest) returns (stream BookResponse))
Client Streaming RPC: Borrow multiple books in a single request (BorrowBooks(stream BookRequest) returns (BorrowStatus))
Bidirectional Streaming RPC: Real-time updates for book availability
in the library (LiveBookUpdates(stream BorrowRequest) returns (stream BookResponse))

## RFC Methods
BookRequest: Contains details for adding a book.
BookResponse: Confirmation or details after adding a book.
Book: Represents a book with attributes like title, author, etc.
BorrowRequest: Contains details for borrowing a book.
BorrowStatus: The response after borrowing books.
EmptyRequest: Used for server streaming when no input is needed.


## Defining Services with four methods
AddBook: Unary RPC to add a new book.
ListAvailableBooks: Server streaming RPC to list available books.
BorrowBooks: Client streaming RPC to borrow multiple books.
LiveBookUpdates: Bidirectional streaming RPC for real-time updates.

## Key Elements
message Book: Represents a book's information (title, author, etc.).
message BorrowRequest: For borrowing a book.
message BorrowStatus: For borrowing responses.
message EmptyRequest: For methods with no input.
service LibraryService: Define the service with all the RPCs.

### commmad to create python grpc files in python folder

python -m grpc_tools.protoc -I=../proto --python_out=. --grpc_python_out=. ../proto/library.proto
genereated files are
library_pb2_grpc.py
library_pb2.py

### command to run server-client files in python folder

docker build -t library-server -f Dockerfile.server .
docker build -t library-client -f Dockerfile.client .
python library-server.py (in server terminal)
python library-client.py (in another client terminal)

### command to installation for go

brew install go
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest (necessary plugins to generate Go code from .proto file)
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

export PATH="$PATH:$(go env GOPATH)/bin" (to set executable file for Go)
go mod init grpc-library (to initialize go mod)

### command to generation for proto file for go

go mod init server-client-implementation
go get google.golang.org/grpc
go mod tidy

protoc -I=. --go_out=. --go-grpc_out=. proto/library.proto

### to build docker files

docker build -t server -f Dockerfile.go-server .
docker build -t client -f Dockerfile.go-client .

### to create network and run in docker

docker network create library-network
docker run -d --name go-server --network library-network -p 50051:50051 server
docker run --rm --name client --network library-network client
