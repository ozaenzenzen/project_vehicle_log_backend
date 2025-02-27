package data

type SendOTPForgotPasswordResponseModel struct {
	Status  int                             `json:"status"`
	Message string                          `json:"message"`
	Data    *SendOTPForgotPasswordDataModel `json:"Data"`
}

type SendOTPForgotPasswordDataModel struct {
	OTPKey       string `json:"otp_key"`
	ResendOTPKey string `json:"resend_otp_key"`
	ForgotKey    string `json:"forgot_key"`
}
