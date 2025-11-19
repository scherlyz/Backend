package test

import (
    serviceMongo "backendgo/app/serviceMongo"
    "backendgo/database"

    "bytes"
    "mime/multipart"
    "net/http/httptest"
    "testing"

    "github.com/gofiber/fiber/v2"
)


func createMultipartFile(field, filename string, content []byte) (*bytes.Buffer, string) {
    body := new(bytes.Buffer)
    writer := multipart.NewWriter(body)

    part, _ := writer.CreateFormFile(field, filename)
    part.Write(content)

    writer.Close()
    return body, writer.FormDataContentType()
}


func TestUploadFile_NoFile(t *testing.T) {
    
    database.MongoDB = nil

    app := setupApp()
    app.Post("/api/files/upload", func(c *fiber.Ctx) error {
        c.Locals("role", "user")
        c.Locals("user_id", 10)
        return serviceMongo.UploadFile(c)
    })

    req := httptest.NewRequest("POST", "/api/files/upload", nil)
    resp, _ := app.Test(req)

    if resp.StatusCode != 500 {
        t.Errorf("Expected 500 because MongoDB nil, got %d", resp.StatusCode)
    }
}

func TestUploadFile_WrongFormat(t *testing.T) {
    database.MongoDB = nil

    app := setupApp()
    app.Post("/api/files/upload", func(c *fiber.Ctx) error {
        c.Locals("role", "user")
        c.Locals("user_id", 10)
        return serviceMongo.UploadFile(c)
    })

    body, contentType := createMultipartFile("file", "test.txt", []byte("hai"))

    req := httptest.NewRequest("POST", "/api/files/upload", body)
    req.Header.Set("Content-Type", contentType)

    resp, _ := app.Test(req)

    if resp.StatusCode != 500 {
        t.Errorf("Expected 500 because MongoDB nil, got %d", resp.StatusCode)
    }
}


func TestGetAllFiles_NoMongo(t *testing.T) {
    database.MongoDB = nil

    app := setupApp()
    app.Get("/api/files", func(c *fiber.Ctx) error {
        c.Locals("role", "user")
        c.Locals("user_id", 10)
        return serviceMongo.GetAllFiles(c)
    })

    req := httptest.NewRequest("GET", "/api/files", nil)
    resp, _ := app.Test(req)

    if resp.StatusCode != 500 {
        t.Errorf("Expected 500, got %d", resp.StatusCode)
    }
}


func TestGetFileByID_InvalidID(t *testing.T) {
    database.MongoDB = nil

    app := setupApp()
    app.Get("/api/files/:id", func(c *fiber.Ctx) error {
        c.Locals("role", "user")
        c.Locals("user_id", 10)
        return serviceMongo.GetFileByID(c)
    })

    req := httptest.NewRequest("GET", "/api/files/invalid-hex", nil)
    resp, _ := app.Test(req)

    // karena MongoDB nil → return langsung 500
    if resp.StatusCode != 500 {
        t.Errorf("Expected 500 because MongoDB nil, got %d", resp.StatusCode)
    }
}


func TestDeleteFile_InvalidID(t *testing.T) {
    database.MongoDB = nil

    app := setupApp()
    app.Delete("/api/files/:id", func(c *fiber.Ctx) error {
        c.Locals("role", "user")
        c.Locals("user_id", 10)
        return serviceMongo.DeleteFile(c)
    })

    req := httptest.NewRequest("DELETE", "/api/files/salahID", nil)
    resp, _ := app.Test(req)

    // service akan return 500 jika MongoDB nil
    if resp.StatusCode != 500 {
        t.Errorf("Expected 500 because MongoDB nil, got %d", resp.StatusCode)
    }
}
