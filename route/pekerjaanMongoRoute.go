package route

import (
	"backendgo/app/serviceMongo"
	"backendgo/middleware"
	"github.com/gofiber/fiber/v2"
)

func PekerjaanMongoRoute(api fiber.Router) {
    pekerjaan := api.Group("/pekerjaan-mongo") // ← beda prefix

    pekerjaan.Get("/", middleware.AuthRequired(), serviceMongo.GetAllPekerjaanMongoService)
    pekerjaan.Get("/:id", middleware.AuthRequired(), serviceMongo.GetPekerjaanByIDMongoService)
    pekerjaan.Get("/alumni/:alumni_id", middleware.AuthRequired(), middleware.AdminOnly(), serviceMongo.GetPekerjaanByAlumniMongoService)
    pekerjaan.Post("/", middleware.AuthRequired(), middleware.AdminOnly(), serviceMongo.CreatePekerjaanMongoService)
    pekerjaan.Put("/:id", middleware.AuthRequired(), middleware.AdminOnly(), serviceMongo.UpdatePekerjaanMongoService)
    pekerjaan.Delete("/:id", middleware.AuthRequired(), middleware.AdminOnly(), serviceMongo.DeletePekerjaanMongoService)
}

