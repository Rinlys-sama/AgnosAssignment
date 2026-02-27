package handlers

import (
	"errors"
	"net/http"

	"github.com/Rinlys-sama/AgnosAssignment/models"
	"github.com/Rinlys-sama/AgnosAssignment/services"
	"github.com/gin-gonic/gin"
)

type StaffHandler struct {
	service *services.StaffService
}

func NewStaffHandler(service *services.StaffService) *StaffHandler {
	return &StaffHandler{service: service}
}

// CreateStaff handles POST /staff/create
func (h *StaffHandler) CreateStaff(c *gin.Context) {
	var req models.StaffCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username, password, and hospital are required"})
		return
	}

	resp, err := h.service.CreateStaff(req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrHospitalNotFound):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrUsernameExists):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// Login handles POST /staff/login
func (h *StaffHandler) Login(c *gin.Context) {
	var req models.StaffLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username, password, and hospital are required"})
		return
	}

	resp, err := h.service.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
