# Backend

## Overview

This directory contains the Go application for the Library Management System.

## Setup Instructions

1. **Install Go Dependencies:**

   ```
   go mod tidy
   ```

2. **Set Up Environment Variables:**

   Create a .env file in the backend directory and configure your environment variables as needed

   Example:

   ```bash
   MONGO_URI=ongodb+srv://user-name:u8sdfx@cluster0.gzxcv55sdf.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0

   FRONTEND_URL=http://localhost:3000/
   ```

3. **Run the Backend Server:**
   ```bash
   go run main.go
   ```
4. **Swagger Documentation:**

   Access the Swagger API documentation to below link to explore and test API endpoints.

   ```
    http://localhost:8080/swagger/index.html
   ```

5. **Run Unit Test:**
   
   Under /Backend run the following command to run tests
   ```    
   go test ./...   
   ```

## Directory Structure

- config/: Contains configuration files.
        config.go: Handles environment variable loading and configuration setup.
- controllers/: Contains API request handlers.
        bookController.go: Manages book-related API operations.
- models/: Contains data models.
        book.go: Defines the Book model used in MongoDB.
- routes/: Contains API route definitions.
        routes.go: Registers API routes and connects them with the controllers.
- docs/: Contains Swagger documentation files.
        docs.go: Auto-generated file for Swagger documentation setup.

        swagger.json and swagger.yaml: Swagger specifications for the API.

## API Documentation

It provides interactive API documentation for testing and exploring the endpoints. You can test in the following link after backend run dev server.
`     http://localhost:8080/swagger/index.html 
    `

## Screenshots

<!-- Include backend-specific screenshots or API documentation screenshots here. -->
<img width="800" alt="Screenshot 2024-07-29 at 7 56 57 AM" src="https://github.com/user-attachments/assets/9c16ad10-e459-48ff-8b2c-7a07db0418d5">
<img width="800" alt="Screenshot 2024-07-29 at 7 57 39 AM" src="https://github.com/user-attachments/assets/37609671-ab5c-4a04-b7a1-4664806a4541">

<img width="1218" alt="Screenshot 2024-07-29 at 7 57 59 AM" src="https://github.com/user-attachments/assets/f0be33fd-706a-4c4b-b1d5-5722aa68e931">

<img width="1250" alt="Screenshot 2024-07-29 at 7 58 14 AM" src="https://github.com/user-attachments/assets/df7f52af-855a-4421-9286-f7c24f4b15d4">

## Additional Notes

- Ensure MongoDB is running and configured correctly. The backend connects to MongoDB for data storage.
- Update environment variables in the .env file as needed for local development and testing.
- Refer to the frontend README.md for details on integration and usage.

For any issues or further assistance, please refer to the Swagger documentation or contact me.
