package repository

import (
	"errors"
	"go-rest-challenge/internal/domain"
	"sync"
)

type BookRepository interface {
	Create(book domain.Book) domain.Book
	GetAll() []domain.Book
	GetByID(id int) (domain.Book, error)
	Update(id int, book domain.Book) (domain.Book, error)
	Delete(id int) error
}

type MemoryBookRepository struct {
	mu     sync.RWMutex
	books  map[int]domain.Book
	nextID int
}

func NewMemoryBookRepository(capacity int) *MemoryBookRepository {
	return &MemoryBookRepository{
		books:  make(map[int]domain.Book, capacity),
		nextID: 1,
	}
}

func (r *MemoryBookRepository) Create(book domain.Book) domain.Book {
	r.mu.Lock()
	book.ID = r.nextID
	r.books[r.nextID] = book
	r.nextID++
	r.mu.Unlock()
	return book
}

func (r *MemoryBookRepository) GetAll() []domain.Book {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]domain.Book, 0, len(r.books))
	for _, b := range r.books {
		result = append(result, b)
	}
	return result
}

func (r *MemoryBookRepository) GetByID(id int) (domain.Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	book, ok := r.books[id]
	if !ok {
		return domain.Book{}, errors.New("not found")
	}
	return book, nil
}

func (r *MemoryBookRepository) Update(id int, book domain.Book) (domain.Book, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.books[id]; !ok {
		return domain.Book{}, errors.New("not found")
	}

	book.ID = id
	r.books[id] = book
	return book, nil
}

func (r *MemoryBookRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.books[id]; !ok {
		return errors.New("not found")
	}

	delete(r.books, id)
	return nil
}
