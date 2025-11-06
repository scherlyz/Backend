package serviceMongo

import (
	"backendgo/app/modelmongo"
	"backendgo/app/repositoryMongo"
	"backendgo/database"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const uploadBasePath = "uploads"

// UploadFile godoc
// @Summary Upload file (foto atau sertifikat)
// @Description Mengunggah file (foto atau sertifikat) ke server dan menyimpannya ke MongoDB. Hanya bisa diakses user yang login. Admin dapat mengupload file untuk user lain dengan menambahkan form field `user_id`.
// @Tags File
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File yang akan diupload (foto atau sertifikat)"
// @Param category formData string false "Kategori file (foto / sertifikat)"
// @Param user_id formData int false "Hanya admin: ID user lain yang ingin diuploadkan file"
// @Success 200 {object} map[string]interface{} "File berhasil diupload"
// @Failure 400 {object} map[string]string "Request tidak valid"
// @Failure 401 {object} map[string]string "Token tidak valid atau tidak ditemukan"
// @Failure 500 {object} map[string]string "Kesalahan server"
// @Router /api/files/upload [post]
func UploadFile(c *fiber.Ctx) error {
	if database.MongoDB == nil {
		return c.Status(500).JSON(fiber.Map{"error": "MongoDB belum terhubung"})
	}
	fileRepo := repositoryMongo.NewFileRepository(database.MongoDB)

	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(int)

	if role == "admin" {
		if u := c.FormValue("user_id"); u != "" {
			if parsed, err := strconv.Atoi(u); err == nil {
				userID = parsed
			}
		}
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "file wajib diupload"})
	}

	category := c.FormValue("category")
	if category == "" {
		category = "foto"
	}

	contentType := fileHeader.Header.Get("Content-Type")
	switch category {
	case "foto":
		if fileHeader.Size > 1*1024*1024 {
			return c.Status(400).JSON(fiber.Map{"error": "ukuran foto maksimal 1MB"})
		}
		allowed := map[string]bool{"image/jpeg": true, "image/png": true, "image/jpg": true}
		if !allowed[contentType] {
			return c.Status(400).JSON(fiber.Map{"error": "format foto hanya jpeg/png/jpg"})
		}
	case "sertifikat":
		if fileHeader.Size > 2*1024*1024 {
			return c.Status(400).JSON(fiber.Map{"error": "ukuran sertifikat maksimal 2MB"})
		}
		if contentType != "application/pdf" {
			return c.Status(400).JSON(fiber.Map{"error": "format sertifikat hanya PDF"})
		}
	}

	newName := uuid.New().String() + filepath.Ext(fileHeader.Filename)
	path := filepath.Join(uploadBasePath, category, newName)
	os.MkdirAll(filepath.Dir(path), os.ModePerm)
	c.SaveFile(fileHeader, path)

	data := modelmongo.File{
		UserID:       userID,
		FileName:     newName,
		OriginalName: fileHeader.Filename,
		FilePath:     path,
		FileSize:     fileHeader.Size,
		FileType:     contentType,
		FileCategory: category,
		UploadedAt:   time.Now(),
	}

	fileRepo.Create(&data)
	return c.JSON(fiber.Map{"success": true, "data": data})
}

// GetAllFiles godoc
// @Summary Ambil semua file
// @Description Menampilkan semua file milik user login. Jika user adalah admin, maka akan menampilkan semua file dari seluruh user.
// @Tags File
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{} "Berhasil mengambil daftar file"
// @Failure 401 {object} map[string]string "Token tidak valid atau tidak ditemukan"
// @Failure 500 {object} map[string]string "Kesalahan server"
// @Router /api/files [get]
func GetAllFiles(c *fiber.Ctx) error {
	if database.MongoDB == nil {
		return c.Status(500).JSON(fiber.Map{"error": "MongoDB belum terhubung"})
	}
	fileRepo := repositoryMongo.NewFileRepository(database.MongoDB)

	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(int)

	var files []modelmongo.File
	var err error

	if role == "admin" {
		files, err = fileRepo.FindAll()
	} else {
		files, err = fileRepo.FindByUser(userID)
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    files,
	})
}

// GetFileByID godoc
// @Summary Ambil file berdasarkan ID
// @Description Mengambil file berdasarkan ID. User hanya bisa mengakses file miliknya sendiri, kecuali admin yang bisa mengakses semua file.
// @Tags File
// @Security BearerAuth
// @Produce json
// @Param id path string true "ID File (ObjectID MongoDB)"
// @Success 200 {object} map[string]interface{} "Berhasil mengambil data file"
// @Failure 400 {object} map[string]string "ID tidak valid"
// @Failure 401 {object} map[string]string "Token tidak valid atau tidak ditemukan"
// @Failure 403 {object} map[string]string "Akses ditolak"
// @Failure 404 {object} map[string]string "File tidak ditemukan"
// @Failure 500 {object} map[string]string "Kesalahan server"
// @Router /api/files/{id} [get]
func GetFileByID(c *fiber.Ctx) error {
	if database.MongoDB == nil {
		return c.Status(500).JSON(fiber.Map{"error": "MongoDB belum terhubung"})
	}
	fileRepo := repositoryMongo.NewFileRepository(database.MongoDB)

	id := c.Params("id")
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(int)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	allFiles, _ := fileRepo.FindAll()
	var found *modelmongo.File

	for _, f := range allFiles {
		if f.ID == objID {
			found = &f
			break
		}
	}

	if found == nil {
		return c.Status(404).JSON(fiber.Map{"error": "file tidak ditemukan"})
	}

	if role != "admin" && found.UserID != userID {
		return c.Status(403).JSON(fiber.Map{"error": "akses ditolak"})
	}

	return c.JSON(fiber.Map{"success": true, "data": found})
}

// DeleteFile godoc
// @Summary Hapus file
// @Description Menghapus file berdasarkan ID. User hanya bisa menghapus file miliknya sendiri, admin dapat menghapus semua file.
// @Tags File
// @Security BearerAuth
// @Produce json
// @Param id path string true "ID File (ObjectID MongoDB)"
// @Success 200 {object} map[string]string "File berhasil dihapus"
// @Failure 400 {object} map[string]string "ID tidak valid"
// @Failure 401 {object} map[string]string "Token tidak valid atau tidak ditemukan"
// @Failure 403 {object} map[string]string "Tidak boleh hapus file milik user lain"
// @Failure 404 {object} map[string]string "File tidak ditemukan"
// @Failure 500 {object} map[string]string "Gagal menghapus file"
// @Router /api/files/{id} [delete]
func DeleteFile(c *fiber.Ctx) error {
	if database.MongoDB == nil {
		return c.Status(500).JSON(fiber.Map{"error": "MongoDB belum terhubung"})
	}
	fileRepo := repositoryMongo.NewFileRepository(database.MongoDB)

	id := c.Params("id")
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(int)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	allFiles, _ := fileRepo.FindAll()
	var target *modelmongo.File

	for _, f := range allFiles {
		if f.ID == objID {
			target = &f
			break
		}
	}

	if target == nil {
		return c.Status(404).JSON(fiber.Map{"error": "file tidak ditemukan"})
	}

	if role != "admin" && target.UserID != userID {
		return c.Status(403).JSON(fiber.Map{"error": "tidak boleh hapus file milik user lain"})
	}

	err = fileRepo.Delete(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	os.Remove(target.FilePath)
	return c.JSON(fiber.Map{"success": true, "message": "file berhasil dihapus"})
}
