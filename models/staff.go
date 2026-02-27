package models

import "time"

type Staff struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	HospitalID   int       `json:"hospital_id"`
	CreatedAt    time.Time `json:"created_at"`
}

type StaffCreateRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Hospital string `json:"hospital" binding:"required"`
}

type StaffLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Hospital string `json:"hospital" binding:"required"`
}

type StaffCreateResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Hospital string `json:"hospital"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
