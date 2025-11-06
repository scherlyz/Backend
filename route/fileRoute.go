package route

import (
	"backendgo/app/serviceMongo"
	"backendgo/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupFileRoutes(api fiber.Router) {
	files := api.Group("/files")

	files.Post("/upload", middleware.AuthRequired(), serviceMongo.UploadFile)
	files.Get("/", middleware.AuthRequired(), serviceMongo.GetAllFiles)
	files.Get("/:id", middleware.AuthRequired(), serviceMongo.GetFileByID)
	files.Delete("/:id", middleware.AuthRequired(), serviceMongo.DeleteFile)
}
