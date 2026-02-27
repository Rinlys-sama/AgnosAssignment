package handlers

import (
	"net/http"

	"github.com/agnos/hospital-middleware/models"
	"github.com/agnos/hospital-middleware/services"
	"github.com/gin-gonic/gin"
)

type PatientHandler struct {
	service *services.PatientService
}

func NewPatientHandler(service *services.PatientService) *PatientHandler {
	return &PatientHandler{service: service}
}

// SearchPatients handles GET /patient/search
func (h *PatientHandler) SearchPatients(c *gin.Context) {
	var req models.PatientSearchRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
		return
	}

	hospitalID, exists := c.Get("hospital_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	patients, err := h.service.SearchPatients(req, hospitalID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search patients"})
		return
	}

	if patients == nil {
		patients = []models.Patient{}
	}

	c.JSON(http.StatusOK, gin.H{"patients": patients})
}
