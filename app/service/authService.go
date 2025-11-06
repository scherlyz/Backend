package service

import (
	"backendgo/app/model"
	"backendgo/app/repository"
	"backendgo/utils"

	"github.com/gofiber/fiber/v2"
)

// LoginService godoc
// @Summary Login user
// @Description Autentikasi user menggunakan username/email dan password, kemudian menghasilkan JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body model.LoginRequest true "Data login (username dan password)"
// @Success 200 {object} model.LoginResponse
// @Failure 400 {object} map[string]string "Request tidak valid"
// @Failure 401 {object} map[string]string "Username atau password salah"
// @Failure 500 {object} map[string]string "Kesalahan server"
// @Router /api/login [post]
func LoginService(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Request body tidak valid",
		})
	}

	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username dan password harus diisi",
		})
	}

	// Ambil user dari DB via repository
	user, err := repository.GetUserByUsernameOrEmail(req.Username)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Username atau password salah",
		})
	}

	// Validasi password
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Username atau password salah",
		})
	}

	// Generate JWT
	token, err := utils.GenerateToken(*user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal generate token",
		})
	}

	resp := model.LoginResponse{
		User: model.User{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		},
		Token: token,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Login berhasil",
		"data":    resp,
	})
}

// GetProfileService godoc
// @Summary Ambil profil user
// @Description Mengambil profil user yang sedang login berdasarkan token JWT
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{} "Berhasil mengambil profil"
// @Failure 401 {object} map[string]string "Token tidak valid atau tidak ditemukan"
// @Router /api/profile [get]
func GetProfileService(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	username := c.Locals("username").(string)
	role := c.Locals("role").(string)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Profile berhasil diambil",
		"data": fiber.Map{
			"user_id":  userID,
			"username": username,
			"role":     role,
		},
	})
}
