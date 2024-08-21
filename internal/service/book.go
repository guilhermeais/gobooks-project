package service

import (
	"database/sql"
	"fmt"
	"time"
)

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
		fmt.Println(err)
		return err
	}

	lastInsertId, err := result.LastInsertId()

	if err != nil {
		fmt.Println(err)
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

	defer results.Close()

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

func (bookService *BookService) SearchBooksByTitle(title string) ([]Book, error) {
	results, err := bookService.db.Query("SELECT id, title, author, genre FROM books WHERE title LIKE ?", "%"+title+"%")
	if err != nil {
		return nil, err
	}

	defer results.Close()

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

func (bookService *BookService) SimulateRead(bookId int, duration time.Duration, result chan string) {
	now := time.Now()

	book, err := bookService.Get(bookId)

	queryDuration := int(time.Since(now).Milliseconds())

	if err != nil {
		result <- fmt.Sprintf("Error on reading book %d: %s", bookId, err)
		return
	}

	if book == nil {
		result <- fmt.Sprintf("Book %d not found", bookId)
		return
	}

	time.Sleep(duration)
	result <- fmt.Sprintf("Book %d read in %d seconds (query took %dms)", bookId, duration, queryDuration)
}

func (bookService *BookService) SimulateReadMultiple(bookIds []int, duration time.Duration) []string {
	readChannel := make(chan string, len(bookIds))

	now := time.Now()
	for _, bookId := range bookIds {
		go bookService.SimulateRead(bookId, duration, readChannel)
	}

	var results []string

	for i := 0; i < len(bookIds); i++ {
		results = append(results, <-readChannel)
	}

	fmt.Printf("All books read in %dms\n", time.Since(now).Milliseconds())

	return results
}
