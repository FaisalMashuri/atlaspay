package constant

type UserStatus string

const (
	UserStatusPending   UserStatus = "PENDING"
	UserStatusActive    UserStatus = "ACTIVE"
	UserStatusSuspended UserStatus = "SUSPENDED"
)

type OTPPurpose string

const (
	OTPPurposeRegister       OTPPurpose = "REGISTER"
	OTPPurposeForgotPassword OTPPurpose = "FORGOT_PASSWORD"
)
