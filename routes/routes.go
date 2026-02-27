package routes

import (
	"database/sql"

	"github.com/Rinlys-sama/AgnosAssignment/config"
	"github.com/Rinlys-sama/AgnosAssignment/handlers"
	"github.com/Rinlys-sama/AgnosAssignment/middleware"
	"github.com/Rinlys-sama/AgnosAssignment/repository"
	"github.com/Rinlys-sama/AgnosAssignment/services"
	"github.com/gin-gonic/gin"
)

func SetupRouter(db *sql.DB, cfg *config.Config) *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Staff
	staffRepo := repository.NewStaffRepository(db)
	staffService := services.NewStaffService(staffRepo, cfg)
	staffHandler := handlers.NewStaffHandler(staffService)

	// Patient
	patientRepo := repository.NewPatientRepository(db)
	patientService := services.NewPatientService(patientRepo)
	patientHandler := handlers.NewPatientHandler(patientService)

	// Staff routes (no auth required)
	staffGroup := router.Group("/staff")
	{
		staffGroup.POST("/create", staffHandler.CreateStaff)
		staffGroup.POST("/login", staffHandler.Login)
	}

	// Patient routes (auth required)
	patientGroup := router.Group("/patient")
	patientGroup.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		patientGroup.GET("/search", patientHandler.SearchPatients)
	}

	return router
}
