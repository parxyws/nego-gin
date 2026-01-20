package dto

import "time"

type RoleResponse struct {
	ID          int32     `json:"id"`
	RoleName    string    `json:"role_name"`
	Description string    `json:"description"`
	UserCount   int64     `json:"user_count,omitempty"`
	IsUsed      bool      `json:"is_used,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type UserResponse struct {
	UserID      string   `json:"user_id"`
	Username    string   `json:"username"`
	FirstName   string   `json:"first_name"`
	LastName    string   `json:"last_name"`
	Email       string   `json:"email"`
	Avatar      string   `json:"avatar,omitempty"`
	PhoneNumber string   `json:"phone_number,omitempty"`
	Jwt         JwtToken `json:"jwt,omitempty"`
}

type UserProfileResponse struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Avatar    string `json:"avatar,omitempty"`
}

type JwtToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserRegisterResponse struct {
	ReferenceID string `json:"reference_id"`
}

type UserForgotPasswordResponse struct {
	ReferenceID string `json:"reference_id"`
}

type UserAddressResponse struct {
	AddressID     int32     `json:"address_id"`
	UserID        string    `json:"user_id"`
	AddressType   string    `json:"address_type"`
	IsDefault     bool      `json:"is_default"`
	RecipientName string    `json:"recipient_name"`
	AddressLine1  string    `json:"address_line1"`
	AddressLine2  string    `json:"address_line2,omitempty"`
	City          string    `json:"city"`
	StateProvince string    `json:"state_province,omitempty"`
	PostalCode    string    `json:"postal_code"`
	CountryCode   string    `json:"country_code"`
	Phone         string    `json:"phone,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
