package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	v1 "zeelso.com/backend/clients/user/v1"
	"zeelso.com/backend/libs/models"
)

func GetUsers(c *fiber.Ctx) error {
	res := v1.GetUsers()
	if res.Message == "failed" {
		return c.Status(404).JSON(fiber.Map{
			"error": true,
			"msg":   "No User",
		})
	}
	if res.Message == "Not Connect to User Services [1]" {
		return c.Status(505).JSON(fiber.Map{
			"responseCode": 505,
			"message":      res.Message,
			"content":      nil,
		})
	}
	if res.Message == "Not Connect to User Services [2]" {
		return c.Status(505).JSON(fiber.Map{
			"responseCode": 505,
			"message":      res.Message,
			"content":      nil,
		})
	}
	return c.Status(201).JSON(fiber.Map{
		"responseCode": 201,
		"message":      res.Message,
		"content":      res.Users,
	})
}

func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	res := v1.GetUserByID(id)
	if res.User == nil {
		return c.Status(404).JSON(fiber.Map{
			"error": true,
			"msg":   "User not found",
		})
	}
	return c.Status(201).JSON(fiber.Map{
		"responseCode": 201,
		"message":      res.Message,
		"content":      res.User,
	})
}

func CreateUser(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": true,
			"msg":   "Unable to parse user data",
		})
	}

	res := v1.CreateUser(user)

	if res.Message != "success" {
		return c.Status(500).JSON(fiber.Map{
			"error": true,
			"msg":   "Failed to create user",
		})
	}
	return c.Status(201).JSON(fiber.Map{
		"responseCode": 201,
		"message":      res.Message,
		"content":      res.User,
	})
}

// suspend user
func SuspendUser(c *fiber.Ctx) error {
	id := c.Params("id")

	res, err := v1.SuspendUser(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": true,
			"msg":   "Failed to suspend user",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"responseCode": 201,
		"message":      res.Message,
		"content":      res.Success,
	})
}

// update user
func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	request := new(models.UpdateUserRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": true,
			"msg":   "Unable to parse request",
		})
	}

	res, err := v1.UpdateUser(id, v1.UserUpdateData(request.UserData))
	if err != nil {
		log.Printf("Error during updating user: %v. Response Message: %v", err, res.GetMessage()) // Verbose error logging
		return c.Status(500).JSON(fiber.Map{
			"error": true,
			"msg":   "Failed to update user",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"responseCode": 200,
		"message":      res.Message,
		"content":      res.User,
	})
}

// DeleteUser - Kullanıcıyı siler
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	res, err := v1.DeleteUser(id)
	if err != nil {
		log.Printf("Error during deleting user: %v. Response Message: %v", err, res.GetMessage()) // Verbose error logging
		return c.Status(500).JSON(fiber.Map{
			"error": true,
			"msg":   "Failed to delete user",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"responseCode": 200,
		"message":      res.Message,
		"content":      res.Success,
	})
}

func BulkUpdate(c *fiber.Ctx) error {
	request := new(models.BulkUpdateUsersRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": true,
			"msg":   "Unable to parse request",
		})
	}

	res, err := v1.BulkUpdate(request.IDs, v1.UserUpdateData(request.UserData))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": true,
			"msg":   "Failed to bulk update users",
		})
	}

	if !res.Success { // Eğer gRPC servisi başarısız bir durum bildirirse
		return c.Status(500).JSON(fiber.Map{
			"error": true,
			"msg":   res.Message,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"responseCode": 200,
		"message":      res.Message,
		"content":      nil,
	})
}

// bulk delete
func BulkDelete(c *fiber.Ctx) error {
	request := new(models.BulkDeleteUsersRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": true,
			"msg":   "Unable to parse request",
		})
	}

	res, err := v1.BulkDelete(request.IDs)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": true,
			"msg":   "Failed to bulk delete users",
		})
	}

	if !res.Success { // Eğer gRPC servisi başarısız bir durum bildirirse
		return c.Status(500).JSON(fiber.Map{
			"error": true,
			"msg":   res.Message,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"responseCode": 200,
		"message":      res.Message,
		"content":      nil,
	})
}
