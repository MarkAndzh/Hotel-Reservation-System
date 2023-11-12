package api

import (
	"errors"
	"github.com/MarkAndzh/hotel-reservation/db"
	"github.com/MarkAndzh/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	var (
		id = c.Params("id")
	)
	var params types.UpdateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if err := h.userStore.PutUser(c.Context(), id, params); err != nil {
		return err
	}
	return c.JSON(map[string]string{"updated": id})
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	var (
		id = c.Params("id")
	)
	err := h.userStore.DeleteUser(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(map[string]string{"deleted": id})
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}

	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}

	user, err = h.userStore.CreateUser(c.Context(), user)
	if err != nil {
		return err
	}

	return c.JSON(user)
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	var (
		id = c.Params("id")
	)
	user, err := h.userStore.GetUserByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"error": "not found"})
		}
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(users)
}
