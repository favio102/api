package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Book represents the model for a book
// @Description A book in the library
// @Accept json
// @Produce json
type Book struct {
    ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"` // Unique identifier for the book
    Title       string             `json:"title,omitempty"`                   // Title of the book
    Author      string             `json:"author,omitempty"`                  // Author of the book
    Description string             `json:"description,omitempty"`             // Description of the book
    Year        string             `json:"year,omitempty"`                    // Publication year
    Edition     string             `json:"edition,omitempty"`                 // Edition of the book
    Language    string             `json:"language,omitempty"`                // Language of the book
    Subject     string             `json:"subject,omitempty"`                 // Subject of the book
    Publisher   string             `json:"publisher,omitempty"`               // Publisher of the book
    Format      string             `json:"format,omitempty"`                  // Format of the book (e.g., hardcover, paperback)
}
