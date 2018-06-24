package core

import (
	"net/smtp"
	"log"
)

type EmailConfig struct{
	Email string `env:"APP_EMAIL,required"`
	Password string `env:"APP_EMAIL_PASSWORD,required"`
	Domain   string `env:"APP_DOMAIN,default=localhost:8000"`
}

func SendResetPasswordToken(toEmail string, token string, cfg EmailConfig) error{
	msg := "From: " + cfg.Email + "\n" +
		"To: " + toEmail + "\n" +
		"Subject: The App Reset Token\n\n" +
		"You requested to reset your password. \n" +
		"If you did not request a change of password, please ignore this e-mail.\n\n" +
		"To reset password follow this link: \n" + "http://"+cfg.Domain+"/reset/confirm/?t=" + token

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		smtp.PlainAuth("",  cfg.Email, cfg.Password, "smtp.gmail.com"),
		cfg.Email, []string{toEmail}, []byte(msg),
	)
	if err != nil{
		log.Printf("smtp error: %s", err)
		return err
	}

	return nil
}
