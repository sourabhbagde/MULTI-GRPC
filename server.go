package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	pb "server-client-implementation/proto" // generated package (alias librarypb)

	"google.golang.org/grpc"
)

// server implements the LibraryServiceServer interface.
type server struct {
	pb.UnimplementedLibraryServiceServer
	// Using a mutex to protect concurrent access to books.
	books map[string]*pb.Book
	mu    sync.Mutex
}

func newServer() *server {
	return &server{
		books: make(map[string]*pb.Book),
	}
}

// AddBook implements the unary RPC.
// If a book already exists, it returns an error message; otherwise, it stores the book.
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
// It sends a stream of BookResponse messages for all available books.
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
// It processes a stream of BorrowRequest messages. For each request,
// if the requested book is available, it marks it as borrowed; if not, it returns a failure immediately.
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
// For each BorrowRequest received, if the corresponding book exists,
// it sends back an updated Book message (per the proto, LiveBookUpdates streams Book).
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
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterLibraryServiceServer(s, newServer())
	fmt.Println("Library gRPC server is running on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
