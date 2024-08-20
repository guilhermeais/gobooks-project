package web

import (
	"encoding/json"
	"gobooks/internal/service"
	"net/http"
	"strconv"
)

type BookHandlers struct {
	bookService *service.BookService
}

func (handler *BookHandlers) GetBooks(
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	books, err := handler.bookService.GetAll()

	if err != nil {
		http.Error(responseWriter, "Failed to get books", http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(responseWriter).Encode(books)
}

func (handler *BookHandlers) CreateBook(
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	var book service.Book
	err := json.NewDecoder(request.Body).Decode(&book)

	if err != nil {
		http.Error(responseWriter, "Failed to parse the payload", http.StatusBadRequest)
		return
	}

	err = handler.bookService.Create(&book)

	if err != nil {
		http.Error(responseWriter, "Failed to create book", http.StatusInternalServerError)
		return
	}

	responseWriter.WriteHeader(http.StatusCreated)
	json.NewEncoder(responseWriter).Encode(book)
}

func (handler *BookHandlers) GetBook(
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	idStr := request.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(responseWriter, "Invalid book id", http.StatusBadRequest)
		return
	}
	book, err := handler.bookService.Get(id)

	if err != nil {
		http.Error(responseWriter, "Failed to get books", http.StatusInternalServerError)
		return
	}

	if book == nil {
		http.Error(responseWriter, "Book not found", http.StatusNotFound)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(responseWriter).Encode(book)
}

func (handler *BookHandlers) UpdateBook(
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	idStr := request.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(responseWriter, "Invalid book id", http.StatusBadRequest)
		return
	}

	var bookUpdate service.Book

	if err = json.NewDecoder(request.Body).Decode(&bookUpdate); err != nil {
		http.Error(responseWriter, "Failed to parse the payload", http.StatusBadRequest)
		return
	}

	bookUpdate.Id = id

	if err = handler.bookService.Update(&bookUpdate); err != nil {
		http.Error(responseWriter, "Failed to update book", http.StatusInternalServerError)
		return
	}

	responseWriter.WriteHeader(http.StatusNoContent)
	json.NewEncoder(responseWriter).Encode(bookUpdate)
}

func (handler *BookHandlers) DeleteBook(
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	idStr := request.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(responseWriter, "Invalid book id", http.StatusBadRequest)
		return
	}

	if err = handler.bookService.Delete(id); err != nil {
		http.Error(responseWriter, "Failed to delete book", http.StatusInternalServerError)
		return
	}

	responseWriter.WriteHeader(http.StatusNoContent)
}
