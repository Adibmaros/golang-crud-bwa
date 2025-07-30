package main

import (
	"belajar-go-restapi-fiber/config"
	"belajar-go-restapi-fiber/entities"
	"log"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	db, err := config.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())
	validate := validator.New()

	// Health check endpoint untuk Railway
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Activity API is running!",
			"status":  "ok",
		})
	})

	// GET /activities
	app.Get("/activities", func(c *fiber.Ctx) error {
		rows, err := db.Query("SELECT id, title, category, description, activity_date, status, created_at, updated_at FROM activities")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		defer rows.Close()

		var activities []entities.Activity
		for rows.Next() {
			var activity entities.Activity
			err := rows.Scan(&activity.ID, &activity.Title, &activity.Category, &activity.Description, &activity.ActivityDate, &activity.Status, &activity.CreatedAt, &activity.UpdatedAt)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			activities = append(activities, activity)
		}

		// Check for errors during iteration
		if err := rows.Err(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"data": activities,
		})
	})

	// POST /activities
	app.Post("/activities", func(c *fiber.Ctx) error {
		var activity entities.Activity
		if err := c.BodyParser(&activity); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Validate the activity
		err = validate.Struct(&activity)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Validation failed: " + err.Error(),
			})
		}

		// Insert dengan timestamps untuk PostgreSQL
		_, err := db.Exec("INSERT INTO activities (title, category, description, activity_date, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)",
			activity.Title, activity.Category, activity.Description, activity.ActivityDate, activity.Status)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Activity created successfully",
		})
	})

	// PUT /activities/:id
	app.Put("/activities/:id", func(c *fiber.Ctx) error {
		calonId := c.Params("id")
		id, err := strconv.Atoi(calonId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid ID",
			})
		}
		var activity entities.Activity
		if err := c.BodyParser(&activity); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Validate the activity
		err = validate.Struct(&activity)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Validation failed: " + err.Error(),
			})
		}

		// Update dengan timestamps untuk PostgreSQL
		_, err = db.Exec("UPDATE activities SET title = $1, category = $2, description = $3, activity_date = $4, status = $5, updated_at = CURRENT_TIMESTAMP WHERE id = $6",
			activity.Title, activity.Category, activity.Description, activity.ActivityDate, activity.Status, id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message": "Activity updated successfully",
		})
	})

	// DELETE /activities/:id
	app.Delete("/activities/:id", func(c *fiber.Ctx) error {
		calonId := c.Params("id")
		id, err := strconv.Atoi(calonId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid ID",
			})
		}

		// Delete activity
		_, err = db.Exec("DELETE FROM activities WHERE id = $1", id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message": "Activity deleted successfully",
		})
	})

	// Get port from environment variable (untuk Railway)
	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}

	log.Printf("Server starting on port %s...", port)
	log.Fatal(app.Listen(":" + port))
}
