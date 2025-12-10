package storage

import (
	"encoding/json"
	"fiber-colledge-done/models"
	"os"
	"sync"
)

type JSONStorage struct {
	filename string
	mu       sync.RWMutex
	Products map[int]models.Product `json:"products"`
	NextID   int                    `json:"next_id"`
}

func NewJSONStorage(filename string) (*JSONStorage, error) {
	storage := &JSONStorage{
		filename: filename,
		Products: make(map[int]models.Product),
		NextID:   1,
	}

	// Загружаем данные из файла, если он существует
	if err := storage.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return storage, nil
}

func (s *JSONStorage) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, s)
}

func (s *JSONStorage) save() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filename, data, 0644)
}

func (s *JSONStorage) GetAll() ([]models.Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	products := make([]models.Product, 0, len(s.Products))
	for _, product := range s.Products {
		products = append(products, product)
	}

	return products, nil
}

func (s *JSONStorage) GetByID(id int) (*models.Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	product, exists := s.Products[id]
	if !exists {
		return nil, nil
	}

	return &product, nil
}

func (s *JSONStorage) Create(product *models.Product) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	product.ID = s.NextID
	s.NextID++
	s.Products[product.ID] = *product

	return s.save()
}

func (s *JSONStorage) Update(id int, product *models.Product) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.Products[id]; !exists {
		return nil // Продукт не найден
	}

	product.ID = id
	s.Products[id] = *product

	return s.save()
}

func (s *JSONStorage) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.Products[id]; !exists {
		return nil // Продукт не найден
	}

	delete(s.Products, id)
	return s.save()
}
