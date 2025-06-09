package main

import (
	"log"
	"os"

	controllers "cognivia-api/Delivery/controllers"
	"cognivia-api/Delivery/routers"
	mongodb "cognivia-api/Repositories"
	usecase "cognivia-api/Usecase"
	"cognivia-api/database"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize MongoDB connection
	db, err := database.NewMongoDBConnection()
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Initialize repositories
	userRepo := mongodb.NewUserRepository(db)
	notebookRepo := mongodb.NewNotebookRepository(db)
	snapnotesRepo := mongodb.NewSnapnotesRepository(db)
	prepPilotRepo := mongodb.NewPrepPilotRepository(db)
	testResultRepo := mongodb.NewTestResultRepository(db)

	// Initialize use cases
	userUseCase := usecase.NewUserUseCase(userRepo)
	notebookUseCase := usecase.NewNotebookUseCase(notebookRepo, snapnotesRepo, prepPilotRepo)
	testResultUseCase := usecase.NewTestResultUseCase(testResultRepo, notebookRepo, prepPilotRepo)

	userHandler := controllers.NewUserHandler(userUseCase)
	notebookHandler := controllers.NewNotebookHandler(notebookUseCase)
	testResultHandler := controllers.NewTestResultHandler(testResultUseCase)
	router := routers.SetupRouter(userHandler, notebookHandler, testResultHandler)

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
