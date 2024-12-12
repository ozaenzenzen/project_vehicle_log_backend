package service

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"sync"
	"time"
)

// OTPEntry holds the OTP and its expiration time.
type OTPEntry struct {
	OTP       string
	ExpiresAt time.Time
}

// OTPStore stores OTPs with expiration times.
var otpStore = struct {
	sync.RWMutex
	data map[string]OTPEntry // Email to OTP mapping
}{data: make(map[string]OTPEntry)}

// registerUser handles user registration and sends OTP to their email.
func SendEmailRegisterUser(email string) error {
	otp := generateOTP(6)
	expiration := time.Now().Add(5 * time.Minute) // OTP expires in 5 minutes
	storeOTP(email, otp, expiration)

	message := fmt.Sprintf("Subject: Account Verification\n\nYour OTP is: %s", otp)
	err := sendEmail(email, message)
	if err != nil {
		return fmt.Errorf("failed to send OTP email: %v", err)
	}
	fmt.Println("OTP sent to:", email)
	return nil
}

// verifyOTP checks if the OTP entered by the user is valid and not expired.
func verifyOTP(email, inputOTP string) bool {
	otpStore.RLock()
	defer otpStore.RUnlock()
	entry, exists := otpStore.data[email]
	if !exists {
		return false // OTP not found
	}

	// Check if OTP is expired
	if time.Now().After(entry.ExpiresAt) {
		deleteOTP(email) // Cleanup expired OTP
		return false     // OTP expired
	}

	// Check if OTP matches
	return entry.OTP == inputOTP
}

// storeOTP saves the OTP and its expiration time for the given email.
func storeOTP(email, otp string, expiresAt time.Time) {
	otpStore.Lock()
	defer otpStore.Unlock()
	otpStore.data[email] = OTPEntry{
		OTP:       otp,
		ExpiresAt: expiresAt,
	}
}

// deleteOTP removes an OTP entry for a given email.
func deleteOTP(email string) {
	otpStore.Lock()
	defer otpStore.Unlock()
	delete(otpStore.data, email)
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
