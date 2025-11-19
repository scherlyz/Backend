package test

import (
	"backendgo/app/service"
	"bytes"
	"net/http/httptest"
	"testing"
)


func TestGetAllAlumniService(t *testing.T) {
	app := setupApp()
	app.Get("/api/alumni", service.GetAllAlumniService)

	req := httptest.NewRequest("GET", "/api/alumni", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}


func TestGetAlumniByID_InvalidID(t *testing.T) {
	app := setupApp()
	app.Get("/api/alumni/:id", service.GetAlumniByIDService)

	req := httptest.NewRequest("GET", "/api/alumni/abc", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 400 {
		t.Errorf("Expected 400 invalid ID, got %d", resp.StatusCode)
	}
}


func TestCreateAlumni_InvalidBody(t *testing.T) {
	app := setupApp()
	app.Post("/api/alumni", service.CreateAlumniService)

	body := bytes.NewBuffer([]byte(`{ invalid json }`))
	req := httptest.NewRequest("POST", "/api/alumni", body)
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	if resp.StatusCode != 400 {
		t.Errorf("Expected 400 invalid body, got %d", resp.StatusCode)
	}
}


func TestUpdateAlumni_InvalidID(t *testing.T) {
	app := setupApp()
	app.Put("/api/alumni/:id", service.UpdateAlumniService)

	body := bytes.NewBuffer([]byte(`{"nama":"Test"}`))
	req := httptest.NewRequest("PUT", "/api/alumni/abc", body)
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	if resp.StatusCode != 400 {
		t.Errorf("Expected 400 invalid ID, got %d", resp.StatusCode)
	}
}

func TestUpdateAlumni_InvalidBody(t *testing.T) {
	app := setupApp()
	app.Put("/api/alumni/:id", service.UpdateAlumniService)

	body := bytes.NewBuffer([]byte(`{ invalid }`))
	req := httptest.NewRequest("PUT", "/api/alumni/1", body)
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	if resp.StatusCode != 400 {
		t.Errorf("Expected 400 invalid body, got %d", resp.StatusCode)
	}
}


func TestDeleteAlumni_InvalidID(t *testing.T) {
	app := setupApp()
	app.Delete("/api/alumni/:id", service.DeleteAlumniService)

	req := httptest.NewRequest("DELETE", "/api/alumni/xyz", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 400 {
		t.Errorf("Expected 400 invalid ID, got %d", resp.StatusCode)
	}
}


func TestUpdateStatusKematian_InvalidID(t *testing.T) {
	app := setupApp()
	app.Put("/api/alumni/:id/kematian", service.UpdateStatusKematianService)

	body := bytes.NewBuffer([]byte(`{"status_kematian": true}`))
	req := httptest.NewRequest("PUT", "/api/alumni/abc/kematian", body)
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	if resp.StatusCode != 400 {
		t.Errorf("Expected 400 invalid ID, got %d", resp.StatusCode)
	}
}

func TestUpdateStatusKematian_InvalidBody(t *testing.T) {
	app := setupApp()
	app.Put("/api/alumni/:id/kematian", service.UpdateStatusKematianService)

	body := bytes.NewBuffer([]byte(`{ invalid }`))
	req := httptest.NewRequest("PUT", "/api/alumni/1/kematian", body)
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	if resp.StatusCode != 400 {
		t.Errorf("Expected 400 invalid body, got %d", resp.StatusCode)
	}
}


func TestGetAlumniPagination(t *testing.T) {
	app := setupApp()
	app.Get("/api/alumni/pagination", service.GetAlumniWithPaginationService)

	req := httptest.NewRequest("GET", "/api/alumni/pagination?page=1&limit=10", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}
