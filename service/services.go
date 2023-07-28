package service

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/joewilson27/fiber-x-postgre/models"
	"gorm.io/gorm"
)

type Book struct {
	Author    string `json:"author" gorm:"column: author_name"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateBook(context *fiber.Ctx) error {

	book := Book{}
	err := context.BodyParser(&book)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "request failed"})
		return err

	}

	// validator := validator.New()

	// err = validator.Struct(Book{})

	// if err != nil {

	// 	context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": err})
	// 	return err

	// }

	err = r.DB.Create(&book).Error

	if err != nil {

		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not create book"})
		return err

	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "book has been successfully added"})
	return nil

}

func (r *Repository) UpdateBook(context *fiber.Ctx) error {

	id := context.Params("id")

	if id == "" {

		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "id cannot be empty"})
		return nil

	}

	bookModel := &models.Books{}

	book := Book{}

	err := context.BodyParser(&book)

	if err != nil {

		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "request failed"})
		return err

	}

	err = r.DB.Model(bookModel).Where("id = ?", id).Updates(book).Error

	if err != nil {

		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not update book"})
		return err

	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "book has been successfully updated"})
	return nil

}

func (r *Repository) DeleteBook(context *fiber.Ctx) error {

	bookModel := &models.Books{}

	id := context.Params("id")

	if id == "" {

		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "id cannot be empty"})
		return nil

	}

	err := r.DB.Delete(bookModel, id)

	if err.Error != nil {

		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not delete book"})
		return err.Error

	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "book has been successfully deleted"})
	return nil
}

func (r *Repository) GetBooks(context *fiber.Ctx) error {

	bookModels := &[]models.Books{}

	err := r.DB.Find(bookModels).Error

	if err != nil {

		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not get books"})

		return err

	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "books gotten successfully",
		"data":    bookModels,
	})

	return nil

}

func (r *Repository) GetBookByID(context *fiber.Ctx) error {

	id := context.Params("id")

	bookModel := &models.Books{}

	if id == "" {

		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "id cannot be empty"})

		return nil

	}

	err := r.DB.Where("id = ?", id).First(bookModel).Error

	if err != nil {

		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not get book"})

		return err

	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "books id gotten successfully",
		"data":    bookModel,
	})

	return nil

}

func (r *Repository) SetupRoutes(app *fiber.App) {

	api := app.Group("/api")

	api.Post("/create_books", r.CreateBook)

	api.Delete("/delete_book/:id", r.DeleteBook)

	api.Put("/update_book/:id", r.UpdateBook)

	api.Get("/get_books/:id", r.GetBookByID)

	api.Get("/books", r.GetBooks)

}
