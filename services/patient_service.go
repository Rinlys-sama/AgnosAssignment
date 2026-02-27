package services

import (
	"github.com/agnos/hospital-middleware/models"
	"github.com/agnos/hospital-middleware/repository"
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
// The hospitalID ensures staff can only see patients from their own hospital.
// This is a critical security rule from the assignment requirements!
func (s *PatientService) SearchPatients(req models.PatientSearchRequest, hospitalID int) ([]models.Patient, error) {
	return s.repo.SearchPatients(req, hospitalID)
}
