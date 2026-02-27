package handlers

import (
	"net/http"

	"github.com/agnos/hospital-middleware/models"
	"github.com/agnos/hospital-middleware/services"
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
		status := http.StatusInternalServerError
		if err.Error() == "hospital not found" {
			status = http.StatusBadRequest
		} else if err.Error() == "username already exists" {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
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
