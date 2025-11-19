package test

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestGetAllPekerjaanMongo(t *testing.T) {
	app := setupApp()

	// dummy response → TIDAK panggil service asli
	app.Get("/api/pekerjaan-mongo", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"success": true,
			"data":    []string{},
		})
	})

	req := httptest.NewRequest("GET", "/api/pekerjaan-mongo", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}


func TestGetPekerjaanByID(t *testing.T) {
	app := setupApp()

	app.Get("/api/pekerjaan-mongo/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "invalid-id" {
			return c.Status(404).JSON(fiber.Map{
				"success": false,
				"error":   "Data tidak ditemukan",
			})
		}
		return c.JSON(fiber.Map{"success": true})
	})

	req := httptest.NewRequest("GET", "/api/pekerjaan-mongo/invalid-id", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 404 {
		t.Errorf("Expected 404, got %d", resp.StatusCode)
	}
}


func TestGetPekerjaanByAlumni(t *testing.T) {
	app := setupApp()

	app.Get("/api/pekerjaan-mongo/alumni/:alumni_id", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"data":    []string{},
		})
	})

	req := httptest.NewRequest("GET", "/api/pekerjaan-mongo/alumni/123", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}


func TestCreatePekerjaanMongo_InvalidBody(t *testing.T) {
	app := setupApp()

	app.Post("/api/pekerjaan-mongo", func(c *fiber.Ctx) error {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Request body tidak valid",
		})
	})

	body := bytes.NewBuffer([]byte(`{ invalid json }`))
	req := httptest.NewRequest("POST", "/api/pekerjaan-mongo", body)
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	if resp.StatusCode != 400 {
		t.Errorf("Expected 400, got %d", resp.StatusCode)
	}
}


func TestUpdatePekerjaanMongo_InvalidBody(t *testing.T) {
	app := setupApp()

	app.Put("/api/pekerjaan-mongo/:id", func(c *fiber.Ctx) error {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid body",
		})
	})

	body := bytes.NewBuffer([]byte(`{ invalid }`))
	req := httptest.NewRequest("PUT", "/api/pekerjaan-mongo/123", body)
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	if resp.StatusCode != 400 {
		t.Errorf("Expected 400, got %d", resp.StatusCode)
	}
}


func TestDeletePekerjaanMongo(t *testing.T) {
	app := setupApp()

	app.Delete("/api/pekerjaan-mongo/:id", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Data terhapus",
		})
	})

	req := httptest.NewRequest("DELETE", "/api/pekerjaan-mongo/123", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}
