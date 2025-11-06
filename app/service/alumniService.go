package service

import (
	"backendgo/app/model"
	"backendgo/app/repository"
	"github.com/gofiber/fiber/v2"
	"strconv"
)


// @Summary Ambil semua data alumni
// @Description Mengambil daftar lengkap semua alumni dari database (hanya bisa diakses user yang login)
// @Tags Alumni
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Berhasil mengambil data alumni"
// @Failure 401 {object} map[string]string "Token tidak valid atau tidak ditemukan"
// @Failure 500 {object} map[string]string "Terjadi kesalahan server"
// @Router /api/alumni [get]
func GetAllAlumniService(c *fiber.Ctx) error {
	data, err := repository.GetAllAlumni()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": data})
}


// @Summary Ambil data alumni berdasarkan ID
// @Description Mengambil data alumni tertentu berdasarkan ID-nya (hanya bisa diakses user yang login)
// @Tags Alumni
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID Alumni"
// @Success 200 {object} map[string]interface{} "Berhasil mengambil data alumni"
// @Failure 400 {object} map[string]string "ID tidak valid"
// @Failure 401 {object} map[string]string "Token tidak valid atau tidak ditemukan"
// @Failure 404 {object} map[string]string "Alumni tidak ditemukan"
// @Router /api/alumni/{id} [get]
func GetAlumniByIDService(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	data, err := repository.GetAlumniByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Alumni tidak ditemukan"})
	}

	return c.JSON(fiber.Map{"success": true, "data": data})
}


// @Summary Tambah alumni baru
// @Description Membuat data alumni baru dan menyimpannya ke database (hanya bisa diakses user yang login)
// @Tags Alumni
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body model.CreateAlumniRequest true "Data alumni baru"
// @Success 201 {object} map[string]interface{} "Alumni berhasil dibuat"
// @Failure 400 {object} map[string]string "Body request tidak valid"
// @Failure 401 {object} map[string]string "Token tidak valid atau tidak ditemukan"
// @Failure 500 {object} map[string]string "Gagal menyimpan ke database"
// @Router /api/alumni [post]
func CreateAlumniService(c *fiber.Ctx) error {
	var input model.CreateAlumniRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Body request tidak valid"})
	}

	data, err := repository.CreateAlumni(input)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "Alumni berhasil dibuat",
		"data":    data,
	})
}


// @Summary Update data alumni
// @Description Memperbarui data alumni berdasarkan ID (hanya bisa diakses user yang login)
// @Tags Alumni
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID Alumni"
// @Param body body model.UpdateAlumniRequest true "Data alumni yang diperbarui"
// @Success 200 {object} map[string]interface{} "Data alumni diperbarui"
// @Failure 400 {object} map[string]string "Body atau ID tidak valid"
// @Failure 401 {object} map[string]string "Token tidak valid atau tidak ditemukan"
// @Failure 500 {object} map[string]string "Gagal memperbarui data"
// @Router /api/alumni/{id} [put]
func UpdateAlumniService(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	var input model.UpdateAlumniRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Body request tidak valid"})
	}
	input.ID = id

	updated, err := repository.UpdateAlumni(input)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Data alumni diperbarui", "data": updated})
}


// @Summary Hapus data alumni
// @Description Menghapus data alumni berdasarkan ID (hanya bisa diakses user yang login)
// @Tags Alumni
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID Alumni"
// @Success 200 {object} map[string]string "Alumni berhasil dihapus"
// @Failure 400 {object} map[string]string "ID tidak valid"
// @Failure 401 {object} map[string]string "Token tidak valid atau tidak ditemukan"
// @Failure 500 {object} map[string]string "Gagal menghapus data alumni"
// @Router /api/alumni/{id} [delete]
func DeleteAlumniService(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	if err := repository.DeleteAlumni(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Alumni berhasil dihapus"})
}

// UpdateStatusKematianService godoc
// @Summary Update status kematian alumni
// @Description Mengubah status kematian alumni (hidup/wafat) (hanya bisa diakses user yang login)
// @Tags Alumni
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID Alumni"
// @Param body body model.UpdateStatusKematianRequest true "Status kematian alumni"
// @Success 200 {object} map[string]string "Status kematian diperbarui"
// @Failure 400 {object} map[string]string "Body atau ID tidak valid"
// @Failure 401 {object} map[string]string "Token tidak valid atau tidak ditemukan"
// @Failure 500 {object} map[string]string "Gagal memperbarui status kematian"
// @Router /api/alumni/{id}/kematian [put]
func UpdateStatusKematianService(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	var req model.UpdateStatusKematianRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Body tidak valid"})
	}

	if err := repository.UpdateStatusKematian(id, req.StatusKematian); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal update status kematian"})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Status kematian diperbarui"})
}



// @Summary Ambil data alumni dengan pagination dan pencarian
// @Description Mengambil daftar alumni dengan fitur pencarian, sorting, dan pagination (hanya bisa diakses user yang login)
// @Tags Alumni
// @Produce json
// @Security BearerAuth
// @Param page query int false "Nomor halaman (default 1)"
// @Param limit query int false "Jumlah data per halaman (default 10)"
// @Param sortBy query string false "Kolom pengurutan (default id)"
// @Param order query string false "Urutan (asc/desc)"
// @Param search query string false "Kata kunci pencarian"
// @Success 200 {object} model.AlumniResponse "Data alumni dengan pagination"
// @Failure 401 {object} map[string]string "Token tidak valid atau tidak ditemukan"
// @Failure 500 {object} map[string]string "Gagal mengambil data"
// @Router /api/alumni/pagination [get]
func GetAlumniWithPaginationService(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "id")
	order := c.Query("order", "asc")
	search := c.Query("search", "")

	offset := (page - 1) * limit

	data, err := repository.GetAlumniRepo(search, sortBy, order, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	total, _ := repository.CountAlumniRepo(search)

	response := model.AlumniResponse{
		Data: data,
		Meta: model.MetaInfo{
			Page:   page,
			Limit:  limit,
			Total:  total,
			Pages:  (total + limit - 1) / limit,
			SortBy: sortBy,
			Order:  order,
			Search: search,
		},
	}
	return c.JSON(response)
}
