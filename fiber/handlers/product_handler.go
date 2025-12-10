package handlers

import (
	"fiber-colledge-done/models"
	"fiber-colledge-done/storage"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	storage *storage.JSONStorage
}

func NewProductHandler(storage *storage.JSONStorage) *ProductHandler {
	return &ProductHandler{storage: storage}
}

type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=1"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Category    string  `json:"category"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=1"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Category    string  `json:"category"`
}

func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
	products, err := h.storage.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get products",
		})
	}

	return c.JSON(fiber.Map{
		"products": products,
		"count":    len(products),
	})
}

func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	product, err := h.storage.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get product",
		})
	}

	if product == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	return c.JSON(product)
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var req CreateProductRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Простая валидация
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Product name is required",
		})
	}

	if req.Price <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Product price must be greater than 0",
		})
	}

	now := time.Now()
	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := h.storage.Create(product); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create product",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(product)
}

func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	var req UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Простая валидация
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Product name is required",
		})
	}

	if req.Price <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Product price must be greater than 0",
		})
	}

	existingProduct, err := h.storage.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get product",
		})
	}

	if existingProduct == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	updatedProduct := &models.Product{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		CreatedAt:   existingProduct.CreatedAt,
		UpdatedAt:   time.Now(),
	}

	if err := h.storage.Update(id, updatedProduct); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update product",
		})
	}

	return c.JSON(updatedProduct)
}

func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	product, err := h.storage.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get product",
		})
	}

	if product == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	if err := h.storage.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete product",
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
