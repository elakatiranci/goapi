package endpoints

import (
	"github.com/gofiber/fiber/v2"
	"zeelso.com/api/v1/lib/controllers"
)

func RegisterUserRouter(a *fiber.App) {

	// Group to user route
	user := a.Group("/api/v1/users")

	user.Get("", controllers.GetUsers)    // get all users
	user.Get("/:id", controllers.GetUser) // get a user by id
	user.Post("", controllers.CreateUser) // create a user
	// user.Post("/suspend/:id", controllers.SuspendUser) // suspend a user by id")
	user.Post("/bulk_update", controllers.BulkUpdate) // update multiple users (only updates suspend rn)
	user.Put("/bulk_delete", controllers.BulkDelete)  // delete multiple users with ids
	user.Put("/:id", controllers.UpdateUser)          // update a user by id
	user.Delete("/:id", controllers.DeleteUser)       //delete a user by id
	user.Put("/:id/suspend", controllers.SuspendUser) // suspend a user by id
}
