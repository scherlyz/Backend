package serviceMongo

import (
	"backendgo/app/modelmongo"
	"backendgo/app/repositoryMongo"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GetAllPekerjaanMongoService godoc
// @Summary Ambil semua data pekerjaan (MongoDB)
// @Description Mengambil semua data pekerjaan dari MongoDB. Hanya bisa diakses oleh user yang login.
// @Tags Pekerjaan Mongo
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{} "Berhasil mengambil semua data pekerjaan"
// @Failure 401 {object} map[string]string "Token tidak valid atau tidak ditemukan"
// @Failure 500 {object} map[string]string "Gagal mengambil data"
// @Router /api/pekerjaan-mongo [get]
func GetAllPekerjaanMongoService(c *fiber.Ctx) error {
	data, err := repositoryMongo.GetAllPekerjaanMongo()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

// GetPekerjaanByIDMongoService godoc
// @Summary Ambil pekerjaan berdasarkan ID (MongoDB)
// @Description Mengambil satu data pekerjaan berdasarkan ID di MongoDB. Hanya bisa diakses user login.
// @Tags Pekerjaan Mongo
// @Security BearerAuth
// @Produce json
// @Param id path string true "ID Pekerjaan (ObjectID MongoDB)"
// @Success 200 {object} map[string]interface{} "Berhasil mengambil data pekerjaan"
// @Failure 400 {object} map[string]string "ID tidak valid"
// @Failure 401 {object} map[string]string "Token tidak valid atau tidak ditemukan"
// @Failure 404 {object} map[string]string "Data tidak ditemukan"
// @Failure 500 {object} map[string]string "Kesalahan server"
// @Router /api/pekerjaan-mongo/{id} [get]
func GetPekerjaanByIDMongoService(c *fiber.Ctx) error {
	id := c.Params("id")

	data, err := repositoryMongo.GetPekerjaanByIDMongo(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"error":   "Data tidak ditemukan",
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

// GetPekerjaanByAlumniMongoService godoc
// @Summary Ambil pekerjaan berdasarkan ID alumni (MongoDB)
// @Description Mengambil semua pekerjaan milik alumni tertentu dari MongoDB. Hanya bisa diakses user yang login.
// @Tags Pekerjaan Mongo
// @Security BearerAuth
// @Produce json
// @Param alumni_id path string true "ID Alumni"
// @Success 200 {object} map[string]interface{} "Berhasil mengambil data pekerjaan"
// @Failure 401 {object} map[string]string "Token tidak valid atau tidak ditemukan"
// @Failure 500 {object} map[string]string "Kesalahan server"
// @Router /api/pekerjaan-mongo/alumni/{alumni_id} [get]
func GetPekerjaanByAlumniMongoService(c *fiber.Ctx) error {
	alumniID := c.Params("alumni_id")

	data, err := repositoryMongo.GetPekerjaanByAlumniMongo(alumniID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

// CreatePekerjaanMongoService godoc
// @Summary Tambah data pekerjaan (MongoDB)
// @Description Menambahkan data pekerjaan baru ke MongoDB. Hanya bisa diakses user yang login.
// @Tags Pekerjaan Mongo
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body modelmongo.CreatePekerjaanRequest true "Data pekerjaan baru"
// @Success 201 {object} map[string]interface{} "Data pekerjaan berhasil ditambahkan"
// @Failure 400 {object} map[string]string "Body request tidak valid"
// @Failure 401 {object} map[string]string "Token tidak valid atau tidak ditemukan"
// @Failure 500 {object} map[string]string "Gagal menyimpan data"
// @Router /api/pekerjaan-mongo [post]
func CreatePekerjaanMongoService(c *fiber.Ctx) error {
	var req modelmongo.CreatePekerjaanRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Request body tidak valid",
		})
	}

	newData := modelmongo.PekerjaanAlumni{
		AlumniID:           req.AlumniID,
		NamaPerusahaan:     req.NamaPerusahaan,
		PosisiJabatan:      req.PosisiJabatan,
		BidangIndustri:     req.BidangIndustri,
		LokasiKerja:        req.LokasiKerja,
		GajiRange:          req.GajiRange,
		StatusPekerjaan:    req.StatusPekerjaan,
		DeskripsiPekerjaan: req.DeskripsiPekerjaan,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		IsDeleted:          false,
	}

	data, err := repositoryMongo.CreatePekerjaanMongo(newData)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "Data pekerjaan berhasil ditambahkan",
		"data":    data,
	})
}

// UpdatePekerjaanMongoService godoc
// @Summary Update data pekerjaan (MongoDB)
// @Description Memperbarui data pekerjaan berdasarkan ID di MongoDB. Hanya bisa diakses user login.
// @Tags Pekerjaan Mongo
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID Pekerjaan (ObjectID MongoDB)"
// @Param body body modelmongo.UpdatePekerjaanRequest true "Data pekerjaan yang akan diperbarui"
// @Success 200 {object} map[string]interface{} "Data pekerjaan berhasil diupdate"
// @Failure 400 {object} map[string]string "Request tidak valid"
// @Failure 401 {object} map[string]string "Token tidak valid atau tidak ditemukan"
// @Failure 500 {object} map[string]string "Gagal memperbarui data"
// @Router /api/pekerjaan-mongo/{id} [put]
func UpdatePekerjaanMongoService(c *fiber.Ctx) error {
	id := c.Params("id")
	var req modelmongo.UpdatePekerjaanRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Request body tidak valid",
		})
	}

	err := repositoryMongo.UpdatePekerjaanMongo(id, req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Data pekerjaan berhasil diupdate",
	})
}

// DeletePekerjaanMongoService godoc
// @Summary Hapus data pekerjaan (MongoDB)
// @Description Menghapus data pekerjaan secara permanen dari MongoDB. Hanya bisa diakses user login.
// @Tags Pekerjaan Mongo
// @Security BearerAuth
// @Produce json
// @Param id path string true "ID Pekerjaan (ObjectID MongoDB)"
// @Success 200 {object} map[string]string "Data pekerjaan berhasil dihapus"
// @Failure 400 {object} map[string]string "ID tidak valid"
// @Failure 401 {object} map[string]string "Token tidak valid atau tidak ditemukan"
// @Failure 500 {object} map[string]string "Gagal menghapus data"
// @Router /api/pekerjaan-mongo/{id} [delete]
func DeletePekerjaanMongoService(c *fiber.Ctx) error {
	id := c.Params("id")

	err := repositoryMongo.HardDeletePekerjaanMongo(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Data pekerjaan berhasil dihapus",
	})
}
