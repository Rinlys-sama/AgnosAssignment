package models

import "time"

type Patient struct {
	ID           int       `json:"id"`
	FirstNameTH  string    `json:"first_name_th"`
	MiddleNameTH string    `json:"middle_name_th,omitempty"`
	LastNameTH   string    `json:"last_name_th"`
	FirstNameEN  string    `json:"first_name_en"`
	MiddleNameEN string    `json:"middle_name_en,omitempty"`
	LastNameEN   string    `json:"last_name_en"`
	DateOfBirth  string    `json:"date_of_birth"`
	PatientHN    string    `json:"patient_hn"`
	NationalID   string    `json:"national_id,omitempty"`
	PassportID   string    `json:"passport_id,omitempty"`
	PhoneNumber  string    `json:"phone_number,omitempty"`
	Email        string    `json:"email,omitempty"`
	Gender       string    `json:"gender"`
	HospitalID   int       `json:"hospital_id"`
	CreatedAt    time.Time `json:"created_at"`
}

type PatientSearchRequest struct {
	NationalID  string `form:"national_id"`
	PassportID  string `form:"passport_id"`
	FirstName   string `form:"first_name"`
	MiddleName  string `form:"middle_name"`
	LastName    string `form:"last_name"`
	DateOfBirth string `form:"date_of_birth"`
	PhoneNumber string `form:"phone_number"`
	Email       string `form:"email"`
}
