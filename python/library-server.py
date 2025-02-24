'''
@authors : Sourabh Bagde, Sai Charan Challa
'''
import grpc
import library_pb2
import library_pb2_grpc
from concurrent import futures
import logging

# storing the books
books_stored = {}

class LibraryService(library_pb2_grpc.LibraryServiceServicer):
    # unary RPCs: The client sends a single request to the server and gets a single response back, just like a normal function call.
    # sends book request and get details as book response.
    def AddBook(self, request, context):
        book = request.book
        # a check to see if already book is added.
        if book.id in books_stored:
            return library_pb2.BookResponse(message="Sorry, this book is already stored", book=book)
        # if not found, add the requested book.
        else:
            books_stored[book.id] = book
            return library_pb2.BookResponse(message="Thanks, this book is added", book=book)
        
    # server streaming RPCs: sends empty request and response as stream available of list of books.
    def ListAvailableBooks(self, request, context):
        for book in books_stored.values():
            if book.is_available:         
                yield library_pb2.BookResponse(message="Available book", book=book)

    # client streaming RPCs: sends stream of the empty request and response as the borrow status.
    def BorrowBooks(self, request_iterator, context):
        for request in request_iterator:
            book_id = request.book_id
            # from the stream, with book_id, if book_id is stored and available as TRUE - mark it as FALSE
            if book_id in books_stored and books_stored[book_id].is_available:
                books_stored[book_id].is_available = False
            else:
                # returns failure with a message if not available.
                return library_pb2.BorrowStatus(status="Failed", message=" Book "+ book_id+ " is not available")
        # returns success with a message on book found as available and able to borrow.
        return library_pb2.BorrowStatus(status="Success", message=" Books borrowed successfully")
    
    # bidirectional streaming RPCs: sends stream of borrow request and response of stream book response.
    def LiveBookUpdates(self, request_iterator, context):
        for request in request_iterator:
            book_id = request.book_id
            # from request stream, checks if any book is already stored - responds with an updated one.
            if book_id in books_stored:
                yield library_pb2.BookResponse(message="Updated book status", book=books_stored[book_id])

def serve():
    # port is defined and assigned to value of server port running.
    port = '50051'
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    library_pb2_grpc.add_LibraryServiceServicer_to_server(LibraryService(), server)
    server.add_insecure_port('[::]:' + port)
    server.start()
    print("Library gRPC server started and is runnning on port " + port)
    server.wait_for_termination()


if __name__ == "__main__":
    logging.basicConfig()
    serve()
