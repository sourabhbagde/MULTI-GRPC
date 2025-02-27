package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "server-client-implementation/proto" // generated package (alias librarypb)

	"google.golang.org/grpc"
)

func main() {
	// change here port to connect with server files in the python or go (server accordingly).
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewLibraryServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// adding a first book.
	book1 := &pb.Book{
		Id:          "1",
		Title:       "Book Title 1",
		Author:      "Sourabh Bagde",
		Genre:       "Educational",
		IsAvailable: true,
	}
	res, err := client.AddBook(ctx, &pb.BookRequest{Book: book1})
	if err != nil {
		log.Fatalf("Error calling AddBook: %v", err)
	}
	fmt.Println("AddBook response:", res.Message)

	// adding a second book.
	book2 := &pb.Book{
		Id:          "2",
		Title:       "Book Title 2",
		Author:      "Sai Charan Challa",
		Genre:       "Research",
		IsAvailable: true,
	}
	res, err = client.AddBook(ctx, &pb.BookRequest{Book: book2})
	if err != nil {
		log.Fatalf("Error calling AddBook: %v", err)
	}
	fmt.Println("AddBook response:", res.Message)

	// List available books (server streaming).
	stream, err := client.ListAvailableBooks(ctx, &pb.EmptyRequest{})
	if err != nil {
		log.Fatalf("Error calling ListAvailableBooks: %v", err)
	}
	fmt.Println("Available Books:")
	for {
		bookRes, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error receiving stream: %v", err)
		}
		fmt.Printf(" - %s by %s\n", bookRes.Book.Title, bookRes.Book.Author)
	}

	// Borrow books (client streaming).
	borrowStream, err := client.BorrowBooks(ctx)
	if err != nil {
		log.Fatalf("Error starting BorrowBooks stream: %v", err)
	}
	borrowRequests := []*pb.BorrowRequest{
		{BookId: "1", UserId: "User1"},
		{BookId: "2", UserId: "User2"},
	}
	for _, req := range borrowRequests {
		if err := borrowStream.Send(req); err != nil {
			log.Fatalf("Error sending borrow request: %v", err)
		}
	}
	borrowStatus, err := borrowStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving borrow status: %v", err)
	}
	fmt.Println("BorrowBooks response:", borrowStatus.Status, "-", borrowStatus.Message)

	// Live book updates (bidirectional streaming).
	liveStream, err := client.LiveBookUpdates(ctx)
	if err != nil {
		log.Fatalf("Error starting LiveBookUpdates stream: %v", err)
	}

	// launching a goroutine to receive live updates.
	go func() {
		for {
			// expecting a Book message.
			book, err := liveStream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receiving live update: %v", err)
			}
			status := "Available"
			if !book.IsAvailable {
				status = "Not Available"
			}
			fmt.Printf("Live Update: %s is now %s\n", book.Title, status)
		}
	}()

	// Send live update requests.
	liveRequests := []*pb.BorrowRequest{
		{BookId: "1", UserId: "Sourabh Bagde"},
		{BookId: "2", UserId: "Sai Charan Challa"},
	}
	for _, req := range liveRequests {
		if err := liveStream.Send(req); err != nil {
			log.Fatalf("Error sending live update request: %v", err)
		}
	}
	time.Sleep(2 * time.Second)
	liveStream.CloseSend()
}
