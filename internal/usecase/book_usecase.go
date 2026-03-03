package usecase

import (
	"errors"
	"go-rest-challenge/internal/domain"
	"go-rest-challenge/internal/repository"
	"strconv"
	"strings"
)

type BookUsecase struct {
	repo repository.BookRepository
}

func NewBookUsecase(r repository.BookRepository) *BookUsecase {
	return &BookUsecase{repo: r}
}

func (u *BookUsecase) Create(title, author string, year int) (domain.Book, error) {
	if strings.TrimSpace(title) == "" || strings.TrimSpace(author) == "" {
		return domain.Book{}, errors.New("title and author required")
	}
	return u.repo.Create(domain.Book{Title: title, Author: author, Year: year}), nil
}

func (u *BookUsecase) GetAll(author, pageStr, limitStr string) []domain.Book {
	allBooks := u.repo.GetAll()

	// 1. Initialize 'books' as an empty slice (not nil)
	books := []domain.Book{}

	// filter
	if author != "" {
		for _, b := range allBooks {
			if strings.EqualFold(b.Author, author) {
				books = append(books, b)
			}
		}
	} else {
		// If no filter, use all books
		books = allBooks
	}

	// Ensure we don't return nil if repo.GetAll() returned nil
	if books == nil {
		books = []domain.Book{}
	}

	// pagination
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page > 0 && limit > 0 {
		start := (page - 1) * limit
		if start < 0 {
			start = 0
		} // Guard against weird math

		if start >= len(books) {
			return []domain.Book{} // Returns []
		}

		end := start + limit
		if end > len(books) {
			end = len(books)
		}
		return books[start:end]
	}

	return books
}

func (u *BookUsecase) GetByID(id int) (domain.Book, error) {
	return u.repo.GetByID(id)
}

func (u *BookUsecase) Update(id int, title, author string, year int) (domain.Book, error) {
	if title == "" || author == "" {
		return domain.Book{}, errors.New("invalid input")
	}
	return u.repo.Update(id, domain.Book{Title: title, Author: author, Year: year})
}

func (u *BookUsecase) Delete(id int) error {
	return u.repo.Delete(id)
}
