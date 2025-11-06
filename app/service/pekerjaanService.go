package service

import (
	"backendgo/app/model"
	"backendgo/app/repository"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)


// GetAllPekerjaanService godoc
// @Summary Ambil semua data pekerjaan
// @Description Mengambil semua data pekerjaan dari database (hanya bisa diakses user yang login)
// @Tags Pekerjaan
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{} "Berhasil mengambil data pekerjaan"
// @Failure 401 {object} map[string]string "Token tidak valid atau tidak ditemukan"
// @Failure 500 {object} map[string]string "Terjadi kesalahan server"
// @Router /api/pekerjaan [get]
func GetAllPekerjaanService(c *fiber.Ctx) error {
	data, err := repository.GetAllPekerjaan()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": data})
}


// GetPekerjaanByIDService godoc
// @Summary Ambil pekerjaan berdasarkan ID
// @Description Mengambil detail pekerjaan berdasarkan ID (hanya bisa diakses user yang login)
// @Tags Pekerjaan
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID Pekerjaan"
// @Success 200 {object} map[string]interface{} "Data pekerjaan ditemukan"
// @Failure 400 {object} map[string]string "ID tidak valid"
// @Failure 401 {object} map[string]string "Token tidak valid atau tidak ditemukan"
// @Failure 404 {object} map[string]string "Data pekerjaan tidak ditemukan"
// @Router /api/pekerjaan/{id} [get]
func GetPekerjaanByIDService(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	data, err := repository.GetPekerjaanByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Pekerjaan tidak ditemukan"})
	}
	return c.JSON(fiber.Map{"success": true, "data": data})
}


// GetPekerjaanByAlumniIDService godoc
// @Summary Ambil pekerjaan berdasarkan ID Alumni
// @Description Mengambil semua data pekerjaan milik alumni tertentu (hanya bisa diakses user yang login)
// @Tags Pekerjaan
// @Security BearerAuth
// @Produce json
// @Param alumni_id path int true "ID Alumni"
// @Success 200 {object} map[string]interface{} "Berhasil mengambil data pekerjaan alumni"
// @Failure 400 {object} map[string]string "ID alumni tidak valid"
// @Failure 401 {object} map[string]string "Token tidak valid atau tidak ditemukan"
// @Failure 500 {object} map[string]string "Kesalahan server"
// @Router /api/pekerjaan/alumni/{alumni_id} [get]
func GetPekerjaanByAlumniIDService(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("alumni_id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid alumni_id"})
	}
	data, err := repository.GetPekerjaanByAlumniID(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": data})
}


// CreatePekerjaanService godoc
// @Summary Tambah pekerjaan baru
// @Description Menambahkan data pekerjaan baru untuk alumni (hanya bisa diakses user yang login)
// @Tags Pekerjaan
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.CreatePekerjaanRequest true "Data pekerjaan baru"
// @Success 201 {object} map[string]interface{} "Pekerjaan berhasil dibuat"
// @Failure 400 {object} map[string]string "Body request tidak valid"
// @Failure 401 {object} map[string]string "Token tidak valid atau tidak ditemukan"
// @Failure 500 {object} map[string]string "Gagal menyimpan data pekerjaan"
// @Router /api/pekerjaan [post]
func CreatePekerjaanService(c *fiber.Ctx) error {
	var req model.CreatePekerjaanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Body tidak valid", "detail": err.Error()})
	}

	mulai, err := time.Parse("2006-01-02", req.TanggalMulaiKerja)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format tanggal_mulai_kerja harus YYYY-MM-DD"})
	}

	var selesai *time.Time
	if req.TanggalSelesaiKerja != nil && *req.TanggalSelesaiKerja != "" {
		t, err := time.Parse("2006-01-02", *req.TanggalSelesaiKerja)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Format tanggal_selesai_kerja harus YYYY-MM-DD"})
		}
		selesai = &t
	}

	data := model.PekerjaanAlumni{
		AlumniID:            req.AlumniID,
		NamaPerusahaan:      req.NamaPerusahaan,
		PosisiJabatan:       req.PosisiJabatan,
		BidangIndustri:      req.BidangIndustri,
		LokasiKerja:         req.LokasiKerja,
		GajiRange:           req.GajiRange,
		TanggalMulaiKerja:   mulai,
		TanggalSelesaiKerja: selesai,
		StatusPekerjaan:     req.StatusPekerjaan,
		DeskripsiPekerjaan:  req.DeskripsiPekerjaan,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	newData, err := repository.CreatePekerjaan(data)
	if err != nil {
		log.Println("Service error CreatePekerjaan:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Gagal insert", "detail": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"success": true, "data": newData})
}


// UpdatePekerjaanService godoc
// @Summary Update data pekerjaan
// @Description Memperbarui data pekerjaan berdasarkan ID (hanya bisa diakses user yang login)
// @Tags Pekerjaan
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID Pekerjaan"
// @Param body body model.UpdatePekerjaanRequest true "Data pekerjaan yang diperbarui"
// @Success 200 {object} map[string]interface{} "Pekerjaan berhasil diperbarui"
// @Failure 400 {object} map[string]string "Body atau ID tidak valid"
// @Failure 401 {object} map[string]string "Token tidak valid atau tidak ditemukan"
// @Failure 500 {object} map[string]string "Gagal memperbarui data pekerjaan"
// @Router /api/pekerjaan/{id} [put]
func UpdatePekerjaanService(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid id"})
	}

	var req model.UpdatePekerjaanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Body tidak valid"})
	}

	mulai, err := time.Parse("2006-01-02", req.TanggalMulaiKerja)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format tanggal_mulai_kerja harus YYYY-MM-DD"})
	}

	var selesai *time.Time
	if req.TanggalSelesaiKerja != nil && *req.TanggalSelesaiKerja != "" {
		t, err := time.Parse("2006-01-02", *req.TanggalSelesaiKerja)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Format tanggal_selesai_kerja harus YYYY-MM-DD"})
		}
		selesai = &t
	}

	data := model.PekerjaanAlumni{
		ID:                  id,
		NamaPerusahaan:      req.NamaPerusahaan,
		PosisiJabatan:       req.PosisiJabatan,
		BidangIndustri:      req.BidangIndustri,
		LokasiKerja:         req.LokasiKerja,
		GajiRange:           req.GajiRange,
		TanggalMulaiKerja:   mulai,
		TanggalSelesaiKerja: selesai,
		StatusPekerjaan:     req.StatusPekerjaan,
		DeskripsiPekerjaan:  req.DeskripsiPekerjaan,
		UpdatedAt:           time.Now(),
	}

	updated, err := repository.UpdatePekerjaan(data)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal update"})
	}
	return c.JSON(fiber.Map{"success": true, "data": updated})
}


// DeletePekerjaanService godoc
// @Summary Hapus pekerjaan
// @Description Menghapus data pekerjaan secara permanen (hanya bisa diakses user yang login)
// @Tags Pekerjaan
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID Pekerjaan"
// @Success 200 {object} map[string]string "Pekerjaan berhasil dihapus"
// @Failure 400 {object} map[string]string "ID tidak valid"
// @Failure 401 {object} map[string]string "Token tidak valid atau tidak ditemukan"
// @Failure 500 {object} map[string]string "Gagal menghapus pekerjaan"
// @Router /api/pekerjaan/{id} [delete]
func DeletePekerjaanService(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	if err := repository.DeletePekerjaan(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan berhasil dihapus"})
}


// GetAllPekerjaanPaginationService godoc
// @Summary Ambil semua pekerjaan dengan pagination
// @Description Mengambil semua data pekerjaan dengan fitur pencarian, sorting, dan pagination (hanya bisa diakses user yang login)
// @Tags Pekerjaan
// @Security BearerAuth
// @Produce json
// @Param page query int false "Nomor halaman (default 1)"
// @Param limit query int false "Jumlah data per halaman (default 10)"
// @Param sortBy query string false "Kolom pengurutan (default created_at)"
// @Param order query string false "Urutan (asc/desc)"
// @Param search query string false "Kata kunci pencarian"
// @Success 200 {object} map[string]interface{} "Berhasil mengambil data pekerjaan"
// @Failure 401 {object} map[string]string "Token tidak valid atau tidak ditemukan"
// @Failure 500 {object} map[string]string "Kesalahan server"
// @Router /api/pekerjaan/list [get]
func GetAllPekerjaanPaginationService(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "created_at")
	order := c.Query("order", "desc")
	search := c.Query("search", "")

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	data, err := repository.GetAllPekerjaanWithPagination(search, sortBy, order, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	total, err := repository.CountPekerjaan(search)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success":    true,
		"data":       data,
		"page":       page,
		"limit":      limit,
		"total_data": total,
		"total_page": (total + limit - 1) / limit,
	})
}


// SoftDeletePekerjaanService godoc
// @Summary Soft delete pekerjaan
// @Description Menandai data pekerjaan sebagai dihapus (tidak benar-benar dihapus dari database, hanya untuk user login)
// @Tags Pekerjaan
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID Pekerjaan"
// @Success 200 {object} map[string]string "Soft delete pekerjaan berhasil"
// @Failure 400 {object} map[string]string "ID tidak valid"
// @Failure 401 {object} map[string]string "Token tidak valid atau tidak ditemukan"
// @Failure 500 {object} map[string]string "Gagal melakukan soft delete"
// @Router /api/pekerjaan/{id}/soft-delete [put]
func SoftDeletePekerjaanService(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(int)

	if role == "admin" {
		err = repository.SoftDeletePekerjaanAdmin(id)
	} else {
		err = repository.SoftDeletePekerjaanUser(id, userID)
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("Soft delete pekerjaan id %d sukses", id),
	})
}


// RestorePekerjaanService godoc
// @Summary Restore pekerjaan
// @Description Mengembalikan data pekerjaan yang sudah di-soft delete (hanya bisa diakses user yang login)
// @Tags Pekerjaan
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID Pekerjaan"
// @Success 200 {object} map[string]string "Restore pekerjaan berhasil"
// @Failure 400 {object} map[string]string "ID tidak valid"
// @Failure 401 {object} map[string]string "Token tidak valid atau tidak ditemukan"
// @Failure 500 {object} map[string]string "Gagal melakukan restore"
// @Router /api/pekerjaan/{id}/restore [put]
func RestorePekerjaanService(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(int)

	if role == "admin" {
		err = repository.RestorePekerjaanAdmin(id)
	} else {
		err = repository.RestorePekerjaanUser(id, userID)
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("Restore pekerjaan id %d sukses", id),
	})
}


// HardDeletePekerjaanService godoc
// @Summary Hard delete pekerjaan
// @Description Menghapus data pekerjaan secara permanen dari database (hanya bisa diakses user yang login)
// @Tags Pekerjaan
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID Pekerjaan"
// @Success 200 {object} map[string]string "Hard delete pekerjaan berhasil"
// @Failure 400 {object} map[string]string "ID tidak valid"
// @Failure 401 {object} map[string]string "Token tidak valid atau tidak ditemukan"
// @Failure 500 {object} map[string]string "Gagal melakukan hard delete"
// @Router /api/pekerjaan/{id}/hard-delete [delete]
func HardDeletePekerjaanService(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(int)

	if role == "admin" {
		err = repository.HardDeletePekerjaanAdmin(id)
	} else {
		err = repository.HardDeletePekerjaanUser(id, userID)
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("Hard delete pekerjaan id %d sukses", id),
	})
}



// GetTrashedPekerjaanService godoc
// @Summary Ambil data pekerjaan yang dihapus (trashed)
// @Description Menampilkan daftar pekerjaan yang sudah di-soft delete (hanya bisa diakses user yang login)
// @Tags Pekerjaan
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{} "Berhasil mengambil data pekerjaan terhapus"
// @Failure 401 {object} map[string]string "Token tidak valid atau tidak ditemukan"
// @Failure 500 {object} map[string]string "Gagal mengambil data pekerjaan terhapus"
// @Router /api/pekerjaan/trashed [get]
func GetTrashedPekerjaanService(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(int)

	var data []model.PekerjaanAlumniTrashed
	var err error

	if role == "admin" {
		data, err = repository.GetTrashedPekerjaanAdmin()
	} else {
		data, err = repository.GetTrashedPekerjaanUser(userID)
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "data": data})
}
