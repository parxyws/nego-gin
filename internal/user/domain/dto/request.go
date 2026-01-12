package dto

type RoleRegisterRequest struct {
	RoleName    string `json:"role_name"`
	Description string `json:"description"`
}

type RoleDeleteRequest struct {
	ID           int32  `json:"id"`
	RoleName     string `json:"role_name"`
	DropUserRole bool   `json:"drop_user_role"`
}

type UserRegisterRequest struct {
	FirstName string `json:"first_name" validate:"required,max=100"`
	LastName  string `json:"last_name" validate:"required,max=100"`
	Username  string `json:"username" validate:"required,max=100"`
	Email     string `json:"email" validate:"required,max=100,email"`
	Password  string `json:"password" validate:"required,min=8,max=100"`
}

type UserValidateAccRequest struct {
	ReferenceID string `json:"reference_id" validate:"required,min=10"`
	OTP         string `json:"otp" validate:"required,max=8"`
}

type UserLoginRequest struct {
	Username string `json:"username" validate:"required,max=100"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Email           string `json:"email" validate:"required,email"`
	OTP             string `json:"otp" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

type UpdateCurrentUserRequest struct {
	FirstName   string `json:"first_name" validate:"required,max=100"`
	LastName    string `json:"last_name" validate:"required,max=100"`
	Username    string `json:"username" validate:"required,max=100"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

type GetUserProfileRequest struct {
	UserID string `json:"user_id" validate:"required,max=100"`
}

type UserAddressCreateRequest struct {
	AddressType   string `json:"address_type" validate:"required"`
	RecipientName string `json:"recipient_name" validate:"required"`
	AddressLine1  string `json:"address_line1" validate:"required"`
	AddressLine2  string `json:"address_line2"`
	City          string `json:"city" validate:"required"`
	StateProvince string `json:"state_province"`
	PostalCode    string `json:"postal_code" validate:"required"`
	CountryCode   string `json:"country_code" validate:"required"`
	Phone         string `json:"phone"`
	IsDefault     bool   `json:"is_default"`
}

type UserAddressUpdateRequest struct {
	AddressType   string `json:"address_type"`
	RecipientName string `json:"recipient_name"`
	AddressLine1  string `json:"address_line1"`
	AddressLine2  string `json:"address_line2"`
	City          string `json:"city"`
	StateProvince string `json:"state_province"`
	PostalCode    string `json:"postal_code"`
	CountryCode   string `json:"country_code"`
	Phone         string `json:"phone"`
	IsDefault     bool   `json:"is_default"`
}
