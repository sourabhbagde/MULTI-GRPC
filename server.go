package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	// generated package (alias librarypb) for stub files.
	pb "server-client-implementation/proto"

	"google.golang.org/grpc"
)

// server implements the LibraryServiceServer interface.
type server struct {
	pb.UnimplementedLibraryServiceServer
	// using a mutex to protect concurrent access to books.
	books map[string]*pb.Book
	mu    sync.Mutex
}

func newServer() *server {
	return &server{
		books: make(map[string]*pb.Book),
	}
}

// unary RPCs: The client sends a single request to the server and gets a single response back, just like a normal function call.
// if a book already exists, it returns an error message; otherwise, it stores the book.
func (s *server) AddBook(ctx context.Context, req *pb.BookRequest) (*pb.BookResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	book := req.Book
	if _, exists := s.books[book.Id]; exists {
		return &pb.BookResponse{
			Message: "Sorry, this book is already stored",
			Book:    book,
		}, nil
	}

	s.books[book.Id] = book
	return &pb.BookResponse{
		Message: "Thanks, this book is added",
		Book:    book,
	}, nil
}

// ListAvailableBooks implements the server streaming RPC.
// this sends a stream of BookResponse messages for all available books.
// server streaming RPCs: sends empty request and response as stream available of list of books.
func (s *server) ListAvailableBooks(req *pb.EmptyRequest, stream pb.LibraryService_ListAvailableBooksServer) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, book := range s.books {
		if book.IsAvailable {
			res := &pb.BookResponse{
				Message: "Available book",
				Book:    book,
			}
			if err := stream.Send(res); err != nil {
				return err
			}
		}
	}
	return nil
}

// BorrowBooks implements the client streaming RPC.
// it processes a stream of BorrowRequest messages. For each request,
// if the requested book is available, it marks it as borrowed; if not, it returns a failure immediately.
// client streaming RPCs: sends stream of the empty request and response as the borrow status.
func (s *server) BorrowBooks(stream pb.LibraryService_BorrowBooksServer) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// When client has finished sending requests, return success.
			return stream.SendAndClose(&pb.BorrowStatus{
				Status:  "Success",
				Message: "Books borrowed successfully",
			})
		}
		if err != nil {
			return err
		}
		bookID := req.BookId
		if book, exists := s.books[bookID]; exists && book.IsAvailable {
			book.IsAvailable = false
		} else {
			// If any borrow request fails, immediately return failure.
			return stream.SendAndClose(&pb.BorrowStatus{
				Status:  "Failed",
				Message: "Book " + bookID + " is not available",
			})
		}
	}
}

// LiveBookUpdates implements the bidirectional streaming RPC.
// for each BorrowRequest received, if the corresponding book exists,
// it sends back an updated Book message (per the proto, LiveBookUpdates streams Book).
// bidirectional streaming RPCs: sends stream of borrow request and response of stream book response.
func (s *server) LiveBookUpdates(stream pb.LibraryService_LiveBookUpdatesServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		s.mu.Lock()
		book, exists := s.books[req.BookId]
		s.mu.Unlock()

		if exists {
			// Send the updated Book directly.
			if err := stream.Send(book); err != nil {
				return err
			}
		}
	}
}

func main() {
	// port is defined and assigned to value of server port running at 50052.
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterLibraryServiceServer(s, newServer())
	fmt.Println("Library gRPC server is running on port 50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
