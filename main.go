package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Book structure
type Book struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title    string             `json:"title,omitempty" bson:"title,omitempty"`
	Author   string             `json:"author,omitempty" bson:"author,omitempty"`
	ISBN     string             `json:"isbn,omitempty" bson:"isbn,omitempty"`
	Released time.Time          `json:"released,omitempty" bson:"released,omitempty"`
}

var client *mongo.Client

// CreateBookEndpoint creates a new book
func CreateBook(response http.ResponseWriter, request *http.Request) {
	//Log that the endpoint has been hit
	log.Println("Received request to create a new book")

	response.Header().Set("content-type", "application/json")
	var book Book
	_ = json.NewDecoder(request.Body).Decode(&book)
	collection := client.Database("books").Collection("books")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.InsertOne(ctx, book)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "Error creating book"}`))
		return
	}
	json.NewEncoder(response).Encode(result)

	//Log that the book has been created successfully
	log.Println("Book created successfully")
}

// GetAllBooksEndpoint retrieves all books
func GetAllBooks(response http.ResponseWriter, request *http.Request) {
	//Log that the endpoint has been hit
	log.Println("Received request to get all books")
	response.Header().Set("content-type", "application/json")
	var books []Book
	collection := client.Database("books").Collection("books")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "Error fetching books"}`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var book Book
		cursor.Decode(&book)
		books = append(books, book)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "Error iterating through books"}`))
		return
	}
	json.NewEncoder(response).Encode(books)

	//Log that all the books have been fetched successfully
	log.Println("All Books Fetched successfully")
}

// GetBookByYearEndpoint retrieves books by a specific year of release
func GetBookByYear(response http.ResponseWriter, request *http.Request) {
	//Log that the endpoint has been hit
	log.Println("Received request to get book by year")
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	year := params["year"]

	parsedYear, err := time.Parse("2006", year) // Assuming the year is provided in YYYY format
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(`{"message": "Invalid year format. Please provide year in YYYY format"}`))
		return
	}

	startDate := parsedYear
	endDate := parsedYear.AddDate(1, 0, 0) // Adding 1 year to the provided year

	var books []Book
	collection := client.Database("books").Collection("books")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{"released": bson.M{"$gte": startDate, "$lt": endDate}})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "Error fetching books by year"}`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var book Book
		cursor.Decode(&book)
		books = append(books, book)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "Error iterating through books by year"}`))
		return
	}
	json.NewEncoder(response).Encode(books)
	//Log that all the books have been fetched successfully
	log.Println("Books By Year Fetched")
}

// Get Book name by author
func GetBookByAuthor(response http.ResponseWriter, request *http.Request) {
	//Log that the endpoint has been hit
	log.Println("Received request to get books by author")
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	author := params["author"]
	var books []Book
	collection := client.Database("books").Collection("books")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{"author": author})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "Error fetching books by author"}`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var book Book
		cursor.Decode(&book)
		books = append(books, book)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "Error iterating through books by author"}`))
		return
	}
	json.NewEncoder(response).Encode(books)
	//Log that all the books have been fetched by author name
	log.Println("Books By Author Fetched")
}

// GetBookEndpoint retrieves a single book
func GetBook(response http.ResponseWriter, request *http.Request) {
	//Log that the endpoint has been hit
	log.Println("Received request to get the desired book")
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var book Book
	collection := client.Database("books").Collection("books")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Book{ID: id}).Decode(&book)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "Error fetching book"}`))
		return
	}
	json.NewEncoder(response).Encode(book)
	//Log that the book has been fetched successfully
	log.Println("Book Fetched successfully")
}

// UpdateBookEndpoint updates a book
func UpdateBook(response http.ResponseWriter, request *http.Request) {
	//Log that the endpoint has been hit
	log.Println("Received request to update the book")
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var book Book
	_ = json.NewDecoder(request.Body).Decode(&book)
	collection := client.Database("books").Collection("books")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.ReplaceOne(ctx, Book{ID: id}, book)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "Error updating book"}`))
		return
	}
	json.NewEncoder(response).Encode(result)
	//Log that the book has been updated successfully
	log.Println("Book Updated successfully")
}

// DeleteBookEndpoint deletes a book
func DeleteBook(response http.ResponseWriter, request *http.Request) {
	//Log that the endpoint has been hit
	log.Println("Received request to delete the book")
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	collection := client.Database("books").Collection("books")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.DeleteOne(ctx, Book{ID: id})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "Error deleting book"}`))
		return
	}
	json.NewEncoder(response).Encode(result)
	//Log that the book has been updated successfully
	log.Println("Book Deleted successfully")
}

func main() {
	// Set up the MongoDB client
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(context.Background(), clientOptions)

	// Check if the connection to MongoDB is successful
	err := client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	} else {
		log.Println("Connected to MongoDB!")
	}

	// Initialize the router
	router := mux.NewRouter()

	// Define endpoints
	router.HandleFunc("/books", CreateBook).Methods("POST")
	router.HandleFunc("/books/{id}", GetBook).Methods("GET")
	router.HandleFunc("/books", GetAllBooks).Methods("GET")
	router.HandleFunc("/books/{id}", UpdateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", DeleteBook).Methods("DELETE")
	router.HandleFunc("/books/author/{author}", GetBookByAuthor).Methods("GET")
	router.HandleFunc("/books/year/{year}", GetBookByYear).Methods("GET")

	// Start the server
	log.Println("Server started on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
