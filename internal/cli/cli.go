package cli

import (
	"fmt"
	"gobooks/internal/service"
	"os"
	"strconv"
	"time"
)

type BookCli struct {
	bookService *service.BookService
}

func NewBookCli(bookService *service.BookService) *BookCli {
	return &BookCli{bookService: bookService}
}

func (cli *BookCli) Run() {
	if len((os.Args)) < 2 {
		fmt.Println("Usage: books <command> [arguments]")
		return
	}

	command := os.Args[1]

	switch command {
	case "search":
		if len(os.Args) < 3 {
			fmt.Println("Usage: books search <bookTitle>")
			return
		}

		bookTitle := os.Args[2]
		cli.SearchBooks(bookTitle)
	case "simulate":
		if (len(os.Args)) < 3 {
			fmt.Println("Usage: books simulate <bookId> <bookId> <bookId> ...")
			return
		}

		bookIds := os.Args[2:]
		cli.SimulateReading(bookIds)
	}

}

func (cli *BookCli) SearchBooks(bookName string) {
	books, err := cli.bookService.SearchBooksByTitle(bookName)

	if err != nil {
		fmt.Println("Error searching books: ", err)
		return
	}

	fmt.Printf("Found %d books\n", len(books))
	for _, book := range books {
		fmt.Printf("Id: %d, Title: %s, Author: %s, Genre: %s\n", book.Id, book.Title, book.Author, book.Genre)
	}
}
func (cli *BookCli) SimulateReading(bookIdsString []string) {
	bookIds := []int{}
	for _, bookIdString := range bookIdsString {
		bookId, err := strconv.Atoi(bookIdString)
		if err != nil {
			fmt.Println("Invalid book id: ", bookIdString)
			continue
		}

		bookIds = append(bookIds, bookId)
	}

	results := cli.bookService.SimulateReadMultiple(bookIds, 2*time.Second)

	for _, result := range results {
		fmt.Println(result)
	}
}
