package services

import (
	"crypto/tls"
	"hitenok/pkg/config"
	"hitenok/pkg/domain"
	"hitenok/pkg/repository"
	"log"
	"math/rand"
	"net/smtp"
	"time"
)

const (
	OTPCharset = "0123456789"
)

type OTPServiceI interface {
	GenerateOTP(user *domain.User) *domain.MyError
	VerifyOTP(user domain.User, otp string) (bool, *domain.MyError)
	SendOTP(user domain.User)
	ClearOTP(user *domain.User) *domain.MyError
}

type mailOTPService struct {
	repo      repository.UserRepositoryI
	appConfig *config.AppConfig
}

func NewMailOTPService(repo repository.UserRepositoryI, appConfig *config.AppConfig) OTPServiceI {
	return &mailOTPService{
		repo:      repo,
		appConfig: appConfig,
	}
}

func (mailOTPService *mailOTPService) ClearOTP(user *domain.User) *domain.MyError {
	user.OTP = ""
	user.OTPAttempts = 0
	user.OTPSpawnedAt = time.Time{}
	err := mailOTPService.repo.SaveUser(user)
	if err != nil {
		err.Module = "mailOTPService.ClearOTP" + err.Module
		return err
	}
	return nil
}

func (mailOTPService *mailOTPService) GenerateOTP(user *domain.User) *domain.MyError {
	otp := make([]byte, 4)
	for i := range otp {
		otp[i] = OTPCharset[rand.Intn(len(OTPCharset))]
	}
	user.OTP = string(otp)
	user.OTPSpawnedAt = time.Now()
	user.OTPAttempts = 3
	err := mailOTPService.repo.SaveUser(user)
	if err != nil {
		err.Module = "mailOTPService.GenerateOTP" + err.Module
		return err
	}
	return nil
}

func (mailOTPService *mailOTPService) VerifyOTP(user domain.User, otp string) (bool, *domain.MyError) {
	if user.OTP == "" {
		return false, nil
	}

	if user.OTPSpawnedAt.Add(5 * time.Minute).Before(time.Now()) {
		return false, nil
	}
	if user.OTPAttempts <= 0 {
		return false, nil
	}
	if user.OTP != otp {
		user.OTPAttempts -= 1
		err := mailOTPService.repo.SaveUser(&user)
		if err != nil {
			err.Module = "mailOTPService.VerifyOTP" + err.Module
			return false, err
		}
		return false, nil
	}

	return true, nil
}

func (mailOTPService *mailOTPService) SendOTP(user domain.User) {
	from := mailOTPService.appConfig.Email
	password := mailOTPService.appConfig.EmailToken

	to := []string{user.Email}
	smtpHost := "smtp.mail.ru"
	smtpPort := "465"

	subject := "Subject: Благодарим за регистрацию на сайте\r\n"
	body := "\nВаш пароль: " + user.OTP + "\n"
	message := []byte(subject + "\r\n" + body)
	conn, err := tls.Dial("tcp", smtpHost+":"+smtpPort, &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	})
	if err != nil {
		log.Println(err)
		return
	}

	c, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		log.Println(err)
		return
	}

	auth := smtp.PlainAuth("", from, password, smtpHost)
	if err = c.Auth(auth); err != nil {
		log.Println(err)
		return
	}

	if err = c.Mail(from); err != nil {
		log.Println(err)
		return
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			log.Println(err)
			return
		}
	}

	w, err := c.Data()
	if err != nil {
		log.Println(err)
		return
	}

	_, err = w.Write(message)
	if err != nil {
		log.Println(err)
		return
	}

	err = w.Close()
	if err != nil {
		log.Println(err)
		return
	}
	c.Quit()
}
