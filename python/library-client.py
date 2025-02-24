'''
@authors : Sourabh Bagde, Sai Charan Challa
'''
import grpc
import library_pb2
import library_pb2_grpc

# add a book
def add_book(stub, book_id, title, author, genre, is_available):
    book = library_pb2.Book(id=book_id, title=title, author=author, genre=genre, is_available=is_available)
    request = library_pb2.BookRequest(book=book)
    response = stub.AddBook(request)
    print("AddBook response " + response.message)

# list of available books
def list_available_books(stub):
    request = library_pb2.EmptyRequest()
    print("Available Books:")
    for book in stub.ListAvailableBooks(request):
        print(book)

# borrow books on stub request.
def borrow_books(stub):
    requests = [library_pb2.BorrowRequest(book_id="1"),library_pb2.BorrowRequest(book_id="2")]
    # for request in requests:
    #     response= stub.BorrowBooks([request])
    response = stub.BorrowBooks(iter(requests))
    print(f"BorrowBooks Response: {response.status} with message: {response.message}")

# 
def live_book_updates(stub):
    def request_iterator():
        requests= [library_pb2.BorrowRequest(book_id="1", user_id="Sourabh"),library_pb2.BorrowRequest(book_id="2", user_id="Saicharan")]
        for request in requests:
            yield request
    responses = stub.LiveBookUpdates(request_iterator())
    for response in responses:
        status = "Available" if response.is_available else "Not Available"
        print("Live Update "+response.title+" is now "+status)

def run():
    # connecting host to gRPC server at port 50051 as channel
    with grpc.insecure_channel('library-server:50051') as channel:
        # creating a stub for library_pb2_grpc (methods of LibraryService) 
        stub = library_pb2_grpc.LibraryServiceStub(channel)
        print("Client connected to gRPC server...\n")
        # calling all methods.
        # add_book(stub)
        # adding first book with parameters.
        add_book(stub, "1", "Book Title 1", "Sourabh Bagde", "Educational", True)
        # adding another book with parameters.
        add_book(stub, "2", "Book Title 2", "Sai Charan Challa", "Research", True)
        list_available_books(stub)
        borrow_books(stub)
        live_book_updates(stub)

if __name__ == '__main__':
    run()