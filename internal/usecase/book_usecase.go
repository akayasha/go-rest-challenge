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
	if allBooks == nil {
		allBooks = []domain.Book{}
	}

	// Filter
	author = strings.TrimSpace(author)
	filtered := []domain.Book{}
	if author != "" {
		for _, b := range allBooks {
			if strings.EqualFold(strings.TrimSpace(b.Author), author) {
				filtered = append(filtered, b)
			}
		}
	} else {
		filtered = allBooks
	}

	// Pagination
	page := 1
	limit := len(filtered)

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	start := (page - 1) * limit
	if start >= len(filtered) {
		return []domain.Book{}
	}

	end := start + limit
	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[start:end]
}

func (u *BookUsecase) GetByID(id int) (domain.Book, error) {
	return u.repo.GetByID(id)
}

func (u *BookUsecase) Update(id int, title, author string, year int) (domain.Book, error) {
	// 1. Validate input first
	if strings.TrimSpace(title) == "" || strings.TrimSpace(author) == "" {
		return domain.Book{}, errors.New("invalid input")
	}

	// 2. Check if book exists
	_, err := u.repo.GetByID(id)
	if err != nil {
		return domain.Book{}, repository.ErrNotFound
	}

	// 3. Perform the update
	return u.repo.Update(id, domain.Book{Title: title, Author: author, Year: year})
}

func (u *BookUsecase) Delete(id int) error {
	return u.repo.Delete(id)
}
