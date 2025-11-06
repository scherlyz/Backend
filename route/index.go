package route

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	AuthRoute(api)       // login
	AlumniRoute(api)     // alumni CRUD + role
	PekerjaanRoute(api)  // pekerjaan CRUD + role
	PekerjaanMongoRoute(api)
	SetupFileRoutes(api)


}
