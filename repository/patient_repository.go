package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/agnos/hospital-middleware/models"
)

type PatientRepository struct {
	db *sql.DB
}

func NewPatientRepository(db *sql.DB) *PatientRepository {
	return &PatientRepository{db: db}
}

// SearchPatients dynamically builds a query based on provided search fields.
// Always filters by hospital_id to enforce hospital isolation.
func (r *PatientRepository) SearchPatients(req models.PatientSearchRequest, hospitalID int) ([]models.Patient, error) {
	query := `SELECT id, first_name_th, middle_name_th, last_name_th,
	           first_name_en, middle_name_en, last_name_en,
	           date_of_birth, patient_hn, national_id, passport_id,
	           phone_number, email, gender, hospital_id, created_at
	           FROM patients WHERE hospital_id = $1`

	args := []interface{}{hospitalID}
	argCount := 1

	if req.NationalID != "" {
		argCount++
		query += fmt.Sprintf(" AND national_id = $%d", argCount)
		args = append(args, req.NationalID)
	}
	if req.PassportID != "" {
		argCount++
		query += fmt.Sprintf(" AND passport_id = $%d", argCount)
		args = append(args, req.PassportID)
	}
	if req.FirstName != "" {
		argCount++
		query += fmt.Sprintf(" AND (first_name_en ILIKE $%d OR first_name_th ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+req.FirstName+"%")
	}
	if req.MiddleName != "" {
		argCount++
		query += fmt.Sprintf(" AND (middle_name_en ILIKE $%d OR middle_name_th ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+req.MiddleName+"%")
	}
	if req.LastName != "" {
		argCount++
		query += fmt.Sprintf(" AND (last_name_en ILIKE $%d OR last_name_th ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+req.LastName+"%")
	}
	if req.DateOfBirth != "" {
		argCount++
		query += fmt.Sprintf(" AND date_of_birth = $%d", argCount)
		args = append(args, req.DateOfBirth)
	}
	if req.PhoneNumber != "" {
		argCount++
		query += fmt.Sprintf(" AND phone_number = $%d", argCount)
		args = append(args, req.PhoneNumber)
	}
	if req.Email != "" {
		argCount++
		query += fmt.Sprintf(" AND email ILIKE $%d", argCount)
		args = append(args, "%"+req.Email+"%")
	}

	query += " ORDER BY id"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanPatients(rows)
}

func scanPatients(rows *sql.Rows) ([]models.Patient, error) {
	var patients []models.Patient

	for rows.Next() {
		var p models.Patient
		var (
			middleNameTH sql.NullString
			middleNameEN sql.NullString
			nationalID   sql.NullString
			passportID   sql.NullString
			phoneNumber  sql.NullString
			email        sql.NullString
			dateOfBirth  sql.NullString
		)

		err := rows.Scan(
			&p.ID, &p.FirstNameTH, &middleNameTH, &p.LastNameTH,
			&p.FirstNameEN, &middleNameEN, &p.LastNameEN,
			&dateOfBirth, &p.PatientHN, &nationalID, &passportID,
			&phoneNumber, &email, &p.Gender, &p.HospitalID, &p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		p.MiddleNameTH = nullStringToString(middleNameTH)
		p.MiddleNameEN = nullStringToString(middleNameEN)
		p.NationalID = nullStringToString(nationalID)
		p.PassportID = nullStringToString(passportID)
		p.PhoneNumber = nullStringToString(phoneNumber)
		p.Email = nullStringToString(email)
		if dateOfBirth.Valid {
			p.DateOfBirth = strings.Split(dateOfBirth.String, "T")[0]
		}

		patients = append(patients, p)
	}

	return patients, rows.Err()
}

func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
