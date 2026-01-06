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
	FirstName string `json:"first_name" validate:"required,max=100,alpha"`
	LastName  string `json:"last_name" validate:"required,max=100"`
	Username  string `json:"username" validate:"required,max=100,alpha"`
	Email     string `json:"email" validate:"required,max=100,email"`
	Password  string `json:"password" validate:"required,min=8,max=100"`
}

type UserValidateAccRequest struct {
	ReferenceID string `json:"reference_id"`
	OTP         string `json:"otp"`
}

type UserLoginRequest struct {
	Username string `json:"username" validate:"required,max=100,alpha"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}
