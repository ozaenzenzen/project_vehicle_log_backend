package service

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"sync"
	"time"
)

var otpStore = struct {
	sync.RWMutex
	data map[string]string // Email to OTP mapping
}{data: make(map[string]string)}

// registerUser handles user registration and sends OTP to their email.
func SendEmailRegisterUser(email string) error {
	otp := generateOTP(6)
	storeOTP(email, otp)

	message := fmt.Sprintf("Subject: Account Verification\n\nYour OTP is: %s", otp)
	err := sendEmail(email, message)
	if err != nil {
		return fmt.Errorf("failed to send OTP email: %v", err)
	}
	fmt.Println("OTP sent to:", email)
	return nil
}

// verifyOTP checks if the OTP entered by the user is valid.
func verifyOTP(email, inputOTP string) bool {
	otpStore.RLock()
	defer otpStore.RUnlock()
	storedOTP, exists := otpStore.data[email]
	return exists && storedOTP == inputOTP
}

// storeOTP saves the OTP for the given email.
func storeOTP(email, otp string) {
	otpStore.Lock()
	defer otpStore.Unlock()
	otpStore.data[email] = otp
}

// generateOTP generates a random numeric OTP of the given length.
func generateOTP(length int) string {
	rand.Seed(time.Now().UnixNano())
	otp := ""
	for i := 0; i < length; i++ {
		otp += fmt.Sprintf("%d", rand.Intn(10))
	}
	return otp
}

// sendEmail sends an email using SMTP.
func sendEmail(to, message string) error {
	from := "fauzanamahdi@gmail.com"
	password := "boxckmwjdmufhckd"
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(message))
}
