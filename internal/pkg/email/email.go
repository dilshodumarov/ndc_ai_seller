package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"

	"sugurta/internal/pkg/config"
)

// EmailType represents the type of email to send
type EmailType string

const (
	// VerificationEmail is an email for verifying user registration
	VerificationEmail EmailType = "verification"
	// ResetPasswordEmail is an email for resetting password
	ResetPasswordEmail EmailType = "reset_password"
)

// SendEmailRequest represents the request to send an email
type SendEmailRequest struct {
	To      []string          // Recipients
	Subject string            // Email subject
	Body    map[string]string // Email body content
	Type    EmailType         // Type of email
}

// SendEmail sends an email using the provided configuration and request
func SendEmail(cfg *config.EmailConfig, req *SendEmailRequest) error {
	// Prepare authentication
	auth := smtp.PlainAuth("", cfg.From, cfg.Password, "smtp.gmail.com")

	// Prepare templates based on email type
	var templateContent string
	switch req.Type {
	case VerificationEmail:
		templateContent = verificationEmailTemplate
	case ResetPasswordEmail:
		templateContent = resetPasswordEmailTemplate
	default:
		return fmt.Errorf("unknown email type: %s", req.Type)
	}

	// Parse template
	tmpl, err := template.New("email").Parse(templateContent)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	// Create buffer to hold the rendered template
	var body bytes.Buffer

	// Execute template with the provided data
	if err := tmpl.Execute(&body, req.Body); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	// Prepare email headers
	headers := make(map[string]string)
	headers["From"] = cfg.From
	headers["To"] = req.To[0]
	headers["Subject"] = req.Subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	// Construct email message
	var message bytes.Buffer
	for k, v := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	message.WriteString("\r\n")
	message.Write(body.Bytes())

	// Send email
	err = smtp.SendMail("smtp.gmail.com:587", auth, cfg.From, req.To, message.Bytes())
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// Email templates
const (
	verificationEmailTemplate = `
<!DOCTYPE html>
<html>
<head>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 20px;
            color: #333;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            background-color: #f9f9f9;
            padding: 20px;
            border-radius: 5px;
        }
        .header {
            text-align: center;
            padding-bottom: 20px;
            border-bottom: 1px solid #ddd;
        }
        .content {
            padding: 20px 0;
        }
        .verification-code {
            font-size: 24px;
            font-weight: bold;
            text-align: center;
            padding: 10px;
            margin: 20px 0;
            background-color: #eee;
            border-radius: 5px;
        }
        .footer {
            padding-top: 20px;
            border-top: 1px solid #ddd;
            text-align: center;
            font-size: 12px;
            color: #777;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>Email Verification</h2>
        </div>
        <div class="content">
            <p>Thank you for registering. To complete your registration, please use the following verification code:</p>
            <div class="verification-code">{{.code}}</div>
            <p>This code will expire in 5 minutes.</p>
            <p>If you did not request this verification, please ignore this email.</p>
        </div>
        <div class="footer">
            <p>This is an automated message, please do not reply.</p>
        </div>
    </div>
</body>
</html>
`

	resetPasswordEmailTemplate = `
<!DOCTYPE html>
<html>
<head>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 20px;
            color: #333;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            background-color: #f9f9f9;
            padding: 20px;
            border-radius: 5px;
        }
        .header {
            text-align: center;
            padding-bottom: 20px;
            border-bottom: 1px solid #ddd;
        }
        .content {
            padding: 20px 0;
        }
        .verification-code {
            font-size: 24px;
            font-weight: bold;
            text-align: center;
            padding: 10px;
            margin: 20px 0;
            background-color: #eee;
            border-radius: 5px;
        }
        .footer {
            padding-top: 20px;
            border-top: 1px solid #ddd;
            text-align: center;
            font-size: 12px;
            color: #777;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>Password Reset</h2>
        </div>
        <div class="content">
            <p>You have requested to reset your password. Please use the following code to complete the process:</p>
            <div class="verification-code">{{.code}}</div>
            <p>This code will expire in 5 minutes.</p>
            <p>If you did not request a password reset, please ignore this email.</p>
        </div>
        <div class="footer">
            <p>This is an automated message, please do not reply.</p>
        </div>
    </div>
</body>
</html>
`
)
