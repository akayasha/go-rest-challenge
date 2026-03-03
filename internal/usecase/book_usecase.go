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

func (u *BookUsecase) Create(title, author string) (domain.Book, error) {
	if strings.TrimSpace(title) == "" || strings.TrimSpace(author) == "" {
		return domain.Book{}, errors.New("title and author required")
	}
	return u.repo.Create(domain.Book{Title: title, Author: author}), nil
}

func (u *BookUsecase) GetAll(author, pageStr, limitStr string) []domain.Book {
	books := u.repo.GetAll()

	// filter
	if author != "" {
		filtered := make([]domain.Book, 0, len(books))
		for _, b := range books {
			if strings.EqualFold(b.Author, author) {
				filtered = append(filtered, b)
			}
		}
		books = filtered
	}

	// pagination
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page > 0 && limit > 0 {
		start := (page - 1) * limit
		end := start + limit

		if start >= len(books) {
			return []domain.Book{}
		}
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

func (u *BookUsecase) Update(id int, title, author string) (domain.Book, error) {
	if title == "" || author == "" {
		return domain.Book{}, errors.New("invalid input")
	}
	return u.repo.Update(id, domain.Book{Title: title, Author: author})
}

func (u *BookUsecase) Delete(id int) error {
	return u.repo.Delete(id)
}
