package services

import (
	"github.com/Rinlys-sama/AgnosAssignment/models"
	"github.com/Rinlys-sama/AgnosAssignment/repository"
)

// PatientService contains business logic for patient operations.
type PatientService struct {
	repo *repository.PatientRepository
}

// NewPatientService creates a new PatientService
func NewPatientService(repo *repository.PatientRepository) *PatientService {
	return &PatientService{repo: repo}
}

// SearchPatients searches for patients matching the criteria.
func (s *PatientService) SearchPatients(req models.PatientSearchRequest, hospitalID int) ([]models.Patient, error) {
	return s.repo.SearchPatients(req, hospitalID)
}
