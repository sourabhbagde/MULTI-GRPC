# multi-grpc

This project demonstrates the ability to implement server-client pairs using the same gRPC protocol definition file (proto) but with different programming languages. The project includes two servers, one written in Python and the other in Go, and two clients, also written in Python and Go.

Library Management System

library.proto file with message definitions and service methods is defined.

Unary RPC: Add a new book to the catalog (AddBook(BookRequest) returns (BookResponse))
Server Streaming RPC: Get a list of available books (ListAvailableBooks(EmptyRequest) returns (stream BookResponse))
Client Streaming RPC: Borrow multiple books in a single request (BorrowBooks(stream BookRequest) returns (BorrowStatus))
Bidirectional Streaming RPC: Real-time updates for book availability
in the library (LiveBookUpdates(stream BorrowRequest) returns (stream BookResponse))

Define the messages for books and other data types used in the RPC methods:

BookRequest: Contains details for adding a book.
BookResponse: Confirmation or details after adding a book.
Book: Represents a book with attributes like title, author, etc.
BorrowRequest: Contains details for borrowing a book.
BorrowStatus: The response after borrowing books.
EmptyRequest: Used for server streaming when no input is needed.

Define the service with four methods:

AddBook: Unary RPC to add a new book.
ListAvailableBooks: Server streaming RPC to list available books.
BorrowBooks: Client streaming RPC to borrow multiple books.
LiveBookUpdates: Bidirectional streaming RPC for real-time updates.

Key Elements:
message Book: Represents a book's information (title, author, etc.).
message BorrowRequest: For borrowing a book.
message BorrowStatus: For borrowing responses.
message EmptyRequest: For methods with no input.
service LibraryService: Define the service with all the RPCs.

# commmad to create python grpc files in python folder

python -m grpc_tools.protoc -I=../proto --python_out=. --grpc_python_out=. ../proto/library.proto
genereated files are
library_pb2_grpc.py
library_pb2.py

# command to run server-client files in python folder

python library-server.py (in server terminal)
python library-client.py (in another client terminal)
