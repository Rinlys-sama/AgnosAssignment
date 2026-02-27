package repository

import (
	"database/sql"

	"github.com/Rinlys-sama/AgnosAssignment/models"
)

type StaffRepository struct {
	db *sql.DB
}

func NewStaffRepository(db *sql.DB) *StaffRepository {
	return &StaffRepository{db: db}
}

func (r *StaffRepository) GetHospitalByCode(code string) (int, error) {
	var id int
	err := r.db.QueryRow("SELECT id FROM hospitals WHERE code = $1", code).Scan(&id)
	return id, err
}

func (r *StaffRepository) CreateStaff(username, passwordHash string, hospitalID int) (*models.Staff, error) {
	staff := &models.Staff{}
	err := r.db.QueryRow(
		`INSERT INTO staff (username, password_hash, hospital_id)
		 VALUES ($1, $2, $3)
		 RETURNING id, username, password_hash, hospital_id, created_at`,
		username, passwordHash, hospitalID,
	).Scan(&staff.ID, &staff.Username, &staff.PasswordHash, &staff.HospitalID, &staff.CreatedAt)

	if err != nil {
		return nil, err
	}
	return staff, nil
}

func (r *StaffRepository) GetStaffByUsernameAndHospital(username string, hospitalID int) (*models.Staff, error) {
	staff := &models.Staff{}
	err := r.db.QueryRow(
		`SELECT id, username, password_hash, hospital_id, created_at
		 FROM staff WHERE username = $1 AND hospital_id = $2`,
		username, hospitalID,
	).Scan(&staff.ID, &staff.Username, &staff.PasswordHash, &staff.HospitalID, &staff.CreatedAt)

	if err != nil {
		return nil, err
	}
	return staff, nil
}
