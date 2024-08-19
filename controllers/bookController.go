package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"library-management/models"
)

// BookController struct
type BookController struct {
	bookCollection *mongo.Collection
}

// NewBookController creates a new BookController
func NewBookController(client *mongo.Client) *BookController {
	bookCollection := client.Database("library").Collection("books")
	return &BookController{bookCollection}
}

// validateBook validates the book fields
func validateBook(book *models.Book) error {
	if book.Title == "" {
		return fmt.Errorf("title cannot be empty")
	}
	return nil
}

// GetBooks retrieves all books
// @Summary Get all books
// @Description Retrieve a list of all books
// @Tags books
// @Accept json
// @Produce json
// @Success 200 {array} models.Book
// @Router /books [get]
func (bc *BookController) GetBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var books []models.Book
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := bc.bookCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Error fetching books:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var book models.Book
		if err = cursor.Decode(&book); err != nil {
			log.Println("Error decoding book:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}

	if err := cursor.Err(); err != nil {
		log.Println("Cursor error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(books); err != nil {
		log.Println("Error encoding books to JSON:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetBook retrieves a single book by ID
// @Summary Get a book by ID
// @Description Retrieve a single book by its ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} models.Book
// @Router /books/{id} [get]
func (bc *BookController) GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	objID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		log.Println("Invalid book ID:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var book models.Book
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = bc.bookCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&book)
	if err != nil {
		log.Println("Error fetching book:", err)
		http.NotFound(w, r)
		return
	}

	if err := json.NewEncoder(w).Encode(book); err != nil {
		log.Println("Error encoding book to JSON:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// CreateBook creates a new book with a unique ID
// @Summary Create a new book
// @Description Create a new book with a unique ID
// @Tags books
// @Accept json
// @Produce json
// @Param book body models.Book true "Book data"
// @Success 201 {object} models.Book
// @Router /books [post]
func (bc *BookController) CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book models.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate book fields
	if err := validateBook(&book); err != nil {
		log.Println("Validation error:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book.ID = primitive.NewObjectID()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = bc.bookCollection.InsertOne(ctx, book)
	if err != nil {
		log.Println("Error inserting book:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(book); err != nil {
		log.Println("Error encoding book to JSON:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// UpdateBook updates an existing book by ID
// @Summary Update an existing book by ID
// @Description Update an existing book by its ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Param book body models.Book true "Updated book data"
// @Success 200 {object} models.Book
// @Router /books/{id} [put]
func (bc *BookController) UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	objID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		log.Println("Invalid book ID:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var book models.Book
	err = json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate book fields
	if err := validateBook(&book); err != nil {
		log.Println("Validation error:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": book,
	}

	_, err = bc.bookCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		log.Println("Error updating book:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(book); err != nil {
		log.Println("Error encoding book to JSON:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// DeleteBook deletes a book by ID
// @Summary Delete a book by ID
// @Description Delete a book by its ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {string} string "Book deleted"
// @Router /books/{id} [delete]
func (bc *BookController) DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	objID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		log.Println("Invalid book ID:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = bc.bookCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		log.Println("Error deleting book:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := "Book deleted"
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error encoding response to JSON:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
