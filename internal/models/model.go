package models

import "time"

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"pass,-"`
	Role     string `json:"role"`
}

// user role types
const (
	RoleEmployee  = "employee"
	RoleModerator = "moderator"
)

type PVZ struct {
	ID               string    `json:"id"`
	RegistrationDate time.Time `json:"registrtionDate"`
	City             string    `json:"city"`
}

type Reception struct {
	ID       string    `json:"id"`
	DateTime time.Time `json:"dateTime"`
	PVZID    string    `json:"pvzId"`
	Status   string    `json:"status"`
}

//reception status types
const (
	StatusInProgress = "in_progress"
	StatusClosed     = "close"
)

type Product struct {
	ID          string    `json:"id"`
	DateTime    time.Time `json:"dateTime"`
	Type        string    `json:"type"`
	ReceptionID string    `json:"receptionId"`
}

//product types
const (
	TypeElectronic = "электроника"
	TypeClothes    = "одежда"
	TypeShoes      = "обувь"
)
