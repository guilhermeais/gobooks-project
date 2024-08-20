package service

import "database/sql"

type Book struct {
	Id     int
	Title  string
	Author string
	Genre  string
}

type BookService struct {
	db *sql.DB
}

func NewBookService(db *sql.DB) *BookService {
	return &BookService{db: db}
}

func (bookService *BookService) Create(book *Book) error {
	result, err := bookService.db.Exec("INSERT INTO books (title, author, genre) VALUES (?, ?, ?)", book.Title, book.Author, book.Genre)
	if err != nil {
		return err
	}

	lastInsertId, err := result.LastInsertId()

	if err != nil {
		return err
	}

	book.Id = int(lastInsertId)

	return nil
}

func (bookService *BookService) GetAll() ([]Book, error) {
	results, err := bookService.db.Query("SELECT id, title, author, genre FROM books")

	if err != nil {
		return nil, err
	}

	var books []Book
	for results.Next() {
		var book Book
		err = results.Scan(&book.Id, &book.Title, &book.Author, &book.Genre)
		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	return books, nil
}

func (bookService *BookService) Get(id int) (*Book, error) {
	result := bookService.db.QueryRow("SELECT id, title, author, genre FROM books WHERE id = ?", id)

	var book Book
	err := result.Scan(&book.Id, &book.Title, &book.Author, &book.Genre)

	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (bookService *BookService) Update(book *Book) error {
	_, err := bookService.db.Exec("UPDATE books SET title = ?, author = ?, genre = ? WHERE id = ?", book.Title, book.Author, book.Genre, book.Id)

	if err != nil {
		return err
	}

	return nil
}

func (bookService *BookService) Delete(id int) error {
	_, err := bookService.db.Exec("DELETE FROM books WHERE id = ?", id)

	if err != nil {
		return err
	}

	return nil
}
