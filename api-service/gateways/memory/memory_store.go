package memory

import (
	"exam-api/domain"
	"sync"
)

// This lines checks if Store implements domain.Storage
// It will fail at build time if not
var _ domain.Storage = (*Store)(nil)

type Store struct {
	products map[string]domain.Product
	// We are using a Read-Write Mutex here
	// This guarantees us when we lock and unlock it that either
	// At most one goroutine is writing in the map and none are reading or;
	// No goroutine is writing and any number are reading
	mu sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		products: make(map[string]domain.Product),
		mu:       sync.RWMutex{},
	}
}

func (s *Store) Save(product domain.Product) (string, bool, error) {
	// Lock - writer's lock
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.products[product.GetHash()]
	if ok {
		return product.GetHash(), true, nil
	}
	s.products[product.GetHash()] = product
	return product.GetHash(), false, nil
}

func (s *Store) Get(id string) (domain.Product, bool, error) {
	// RLock - reader's lock
	s.mu.RLock()
	defer s.mu.RUnlock()

	p, ok := s.products[id]
	if !ok {
		return domain.Product{}, false, nil
	}
	return p, true, nil
}

func (s *Store) Update(id string, diff domain.ProductDiff) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// check if id exists in map
	_, ok := s.products[id]

	// initialize new product
	newProduct := domain.Product{
		Name:         s.products[id].Name,
		Manufacturer: s.products[id].Manufacturer,
		Price:        diff.Diff.Price,
		Stock:        diff.Diff.Stock,
		Tags:         diff.Diff.Tags,
	}

	// update product
	s.products[id] = newProduct

	// return updated product
	return ok, nil
}

func (s *Store) Delete(id string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// check if id exists in map
	_, ok := s.products[id]

	// delete id from products
	delete(s.products, id)

	// return deleted product
	return ok, nil
}
