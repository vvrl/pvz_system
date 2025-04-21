package models

import "time"

type User struct {
	ID       int      `json:"id"`
	Email    string   `json:"email"`
	Password string   `json:"-"`
	Role     UserRole `json:"role"`
}

// user role types
type UserRole string

const (
	RoleEmployee  UserRole = "employee"
	RoleModerator UserRole = "moderator"
)

type PVZ struct {
	ID               int       `json:"id"`
	RegistrationDate time.Time `json:"registrtionDate"`
	City             string    `json:"city"`
}

type ReceptionStatus string

type Reception struct {
	ID       int             `json:"id"`
	DateTime time.Time       `json:"dateTime"`
	PVZID    int             `json:"pvzId"`
	Status   ReceptionStatus `json:"status"`
}

//reception status types

const (
	StatusInProgress ReceptionStatus = "in_progress"
	StatusClosed     ReceptionStatus = "close"
)

type Product struct {
	ID          int         `json:"id"`
	DateTime    time.Time   `json:"dateTime"`
	Type        ProductType `json:"type"`
	ReceptionID int         `json:"receptionId"`
}

//product types
type ProductType string

const (
	TypeElectronic ProductType = "электроника"
	TypeClothes    ProductType = "одежда"
	TypeShoes      ProductType = "обувь"
)
