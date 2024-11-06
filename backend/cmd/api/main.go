package main

import (
	"fmt"

	"github.com/VasuDevrani/ticket-booking-project-v1/config"
	"github.com/VasuDevrani/ticket-booking-project-v1/db"
	"github.com/VasuDevrani/ticket-booking-project-v1/handlers"
	"github.com/VasuDevrani/ticket-booking-project-v1/middlewares"
	"github.com/VasuDevrani/ticket-booking-project-v1/repositories"
	"github.com/VasuDevrani/ticket-booking-project-v1/services"
	"github.com/gofiber/fiber/v2"
)

func main() {
	envConfig := config.NewEnvConfig()
	db := db.Init(envConfig, db.DBMigrator)

	app := fiber.New(fiber.Config{
		AppName:      "Ticket-Booking",
		ServerHeader:  "Fiber",
	})

	// Repositories
	eventRepository := repositories.NewEventRepository(db)
	ticketRepository := repositories.NewTicketRepository(db)
	authRepository := repositories.NewAuthRepository(db)

	// Services
	authService := services.NewAuthService(authRepository)

	// Routing
	server := app.Group("/api")
	handlers.NewAuthHandler(server.Group("/auth"), authService)

	privateRoutes := server.Use(middlewares.AuthProtected(db))

	//handlers
	handlers.NewEventHandler(privateRoutes.Group("/event"), eventRepository)
	handlers.NewTicketHandler(privateRoutes.Group("/ticket"), ticketRepository)

	app.Listen(fmt.Sprintf(":" + envConfig.ServerPort))
}