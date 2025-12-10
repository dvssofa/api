package main

import (
	"fiber-colledge-done/handlers"
	"fiber-colledge-done/storage"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Инициализация хранилища
	storage, err := storage.NewJSONStorage("products.json")
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Инициализация обработчиков
	productHandler := handlers.NewProductHandler(storage)

	// Создание Fiber приложения
	app := fiber.New(fiber.Config{
		AppName: "Product API",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())

	// Маршруты
	api := app.Group("/api/v1")

	// Продукты
	products := api.Group("/products")
	{
		products.Get("/", productHandler.GetAllProducts)
		products.Get("/:id", productHandler.GetProduct)
		products.Post("/", productHandler.CreateProduct)
		products.Put("/:id", productHandler.UpdateProduct)
		products.Delete("/:id", productHandler.DeleteProduct)
	}

	// Корневой маршрут
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Product API is running",
			"version": "1.0.0",
		})
	})

	// Обработка 404
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Endpoint not found",
		})
	})

	// Запуск сервера
	log.Println("Server starting on :3000")
	log.Fatal(app.Listen(":3000"))
}
