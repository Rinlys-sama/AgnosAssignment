package services

import (
	"time"

	"github.com/Rinlys-sama/AgnosAssignment/config"
	"github.com/Rinlys-sama/AgnosAssignment/models"
	"github.com/Rinlys-sama/AgnosAssignment/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type StaffService struct {
	repo *repository.StaffRepository
	cfg  *config.Config
}

func NewStaffService(repo *repository.StaffRepository, cfg *config.Config) *StaffService {
	return &StaffService{repo: repo, cfg: cfg}
}

func (s *StaffService) CreateStaff(req models.StaffCreateRequest) (*models.StaffCreateResponse, error) {
	hospitalID, err := s.repo.GetHospitalByCode(req.Hospital)
	if err != nil {
		return nil, ErrHospitalNotFound
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, ErrHashFailed
	}

	staff, err := s.repo.CreateStaff(req.Username, string(hashedPassword), hospitalID)
	if err != nil {
		return nil, ErrUsernameExists
	}

	return &models.StaffCreateResponse{
		ID:       staff.ID,
		Username: staff.Username,
		Hospital: req.Hospital,
	}, nil
}

func (s *StaffService) Login(req models.StaffLoginRequest) (*models.LoginResponse, error) {
	hospitalID, err := s.repo.GetHospitalByCode(req.Hospital)
	if err != nil {
		return nil, ErrInvalidCreds
	}

	staff, err := s.repo.GetStaffByUsernameAndHospital(req.Username, hospitalID)
	if err != nil {
		return nil, ErrInvalidCreds
	}

	if err := bcrypt.CompareHashAndPassword([]byte(staff.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCreds
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"staff_id":    staff.ID,
		"hospital_id": staff.HospitalID,
		"username":    staff.Username,
		"exp":         time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, ErrTokenFailed
	}

	return &models.LoginResponse{Token: tokenString}, nil
}
