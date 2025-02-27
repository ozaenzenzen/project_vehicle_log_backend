package data

type ResendOTPResponseModel struct {
	Status  int                         `json:"status"`
	Message string                      `json:"message"`
	Data    *ResendOTPResponseDataModel `json:"Data"`
}

type ResendOTPResponseDataModel struct {
	OTPKey       string `json:"otp_key"`
	ResendOTPKey string `json:"resend_otp_key"`
}
