package external

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net/smtp"
	"strings"

	"github.com/yourusername/gin-collection-saas/pkg/logger"
)

// EmailConfig holds email service configuration
type EmailConfig struct {
	Host       string
	Port       int
	Username   string
	Password   string
	FromEmail  string
	FromName   string
	TLS        bool
	SkipVerify bool
}

// EmailClient handles sending emails
type EmailClient struct {
	config    *EmailConfig
	templates map[string]*template.Template
}

// NewEmailClient creates a new email client
func NewEmailClient(config *EmailConfig) *EmailClient {
	client := &EmailClient{
		config:    config,
		templates: make(map[string]*template.Template),
	}

	// Load email templates
	client.loadTemplates()

	return client
}

// loadTemplates loads all email templates
func (c *EmailClient) loadTemplates() {
	// User invitation template
	c.templates["user_invitation"] = template.Must(template.New("user_invitation").Parse(userInvitationTemplate))

	// Password reset template
	c.templates["password_reset"] = template.Must(template.New("password_reset").Parse(passwordResetTemplate))

	// Welcome template
	c.templates["welcome"] = template.Must(template.New("welcome").Parse(welcomeTemplate))

	// Subscription confirmation template
	c.templates["subscription_confirmation"] = template.Must(template.New("subscription_confirmation").Parse(subscriptionConfirmationTemplate))
}

// EmailData holds common email data
type EmailData struct {
	To          string
	Subject     string
	TemplateKey string
	Data        interface{}
}

// UserInvitationData holds data for user invitation emails
type UserInvitationData struct {
	RecipientName   string
	RecipientEmail  string
	InviterName     string
	TenantName      string
	Role            string
	InviteLink      string
	ExpiresIn       string
}

// PasswordResetData holds data for password reset emails
type PasswordResetData struct {
	RecipientName  string
	ResetLink      string
	ExpiresIn      string
}

// WelcomeData holds data for welcome emails
type WelcomeData struct {
	RecipientName string
	TenantName    string
	LoginLink     string
}

// SubscriptionData holds data for subscription emails
type SubscriptionData struct {
	RecipientName string
	PlanName      string
	Amount        string
	BillingCycle  string
	NextBilling   string
}

// Send sends an email using the configured SMTP server
func (c *EmailClient) Send(email *EmailData) error {
	if c.config.Host == "" || c.config.Host == "localhost" {
		// Log but don't fail if SMTP is not configured
		logger.Info("Email sending skipped (SMTP not configured)",
			"to", email.To,
			"subject", email.Subject,
			"template", email.TemplateKey)
		return nil
	}

	// Render template
	tmpl, ok := c.templates[email.TemplateKey]
	if !ok {
		return fmt.Errorf("email template not found: %s", email.TemplateKey)
	}

	var bodyBuffer bytes.Buffer
	if err := tmpl.Execute(&bodyBuffer, email.Data); err != nil {
		return fmt.Errorf("failed to render email template: %w", err)
	}

	// Build email message
	from := fmt.Sprintf("%s <%s>", c.config.FromName, c.config.FromEmail)
	msg := c.buildMessage(from, email.To, email.Subject, bodyBuffer.String())

	// Send email
	if err := c.sendMail(email.To, msg); err != nil {
		logger.Error("Failed to send email", "error", err.Error(), "to", email.To)
		return fmt.Errorf("failed to send email: %w", err)
	}

	logger.Info("Email sent successfully", "to", email.To, "subject", email.Subject)
	return nil
}

// buildMessage builds the email message with headers
func (c *EmailClient) buildMessage(from, to, subject, body string) []byte {
	var msg strings.Builder

	msg.WriteString(fmt.Sprintf("From: %s\r\n", from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(body)

	return []byte(msg.String())
}

// sendMail sends the email via SMTP
func (c *EmailClient) sendMail(to string, msg []byte) error {
	addr := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)

	var auth smtp.Auth
	if c.config.Username != "" {
		auth = smtp.PlainAuth("", c.config.Username, c.config.Password, c.config.Host)
	}

	if c.config.TLS {
		return c.sendMailTLS(addr, auth, to, msg)
	}

	return smtp.SendMail(addr, auth, c.config.FromEmail, []string{to}, msg)
}

// sendMailTLS sends email using TLS
func (c *EmailClient) sendMailTLS(addr string, auth smtp.Auth, to string, msg []byte) error {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: c.config.SkipVerify,
		ServerName:         c.config.Host,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, c.config.Host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	if auth != nil {
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP authentication failed: %w", err)
		}
	}

	if err := client.Mail(c.config.FromEmail); err != nil {
		return fmt.Errorf("SMTP MAIL command failed: %w", err)
	}

	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("SMTP RCPT command failed: %w", err)
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("SMTP DATA command failed: %w", err)
	}

	if _, err := writer.Write(msg); err != nil {
		return fmt.Errorf("failed to write email body: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close email writer: %w", err)
	}

	return client.Quit()
}

// SendUserInvitation sends a user invitation email
func (c *EmailClient) SendUserInvitation(data *UserInvitationData) error {
	return c.Send(&EmailData{
		To:          data.RecipientEmail,
		Subject:     fmt.Sprintf("Einladung zu %s - GinVault", data.TenantName),
		TemplateKey: "user_invitation",
		Data:        data,
	})
}

// SendPasswordReset sends a password reset email
func (c *EmailClient) SendPasswordReset(to string, data *PasswordResetData) error {
	return c.Send(&EmailData{
		To:          to,
		Subject:     "Passwort zurücksetzen - GinVault",
		TemplateKey: "password_reset",
		Data:        data,
	})
}

// SendWelcome sends a welcome email
func (c *EmailClient) SendWelcome(to string, data *WelcomeData) error {
	return c.Send(&EmailData{
		To:          to,
		Subject:     fmt.Sprintf("Willkommen bei %s!", data.TenantName),
		TemplateKey: "welcome",
		Data:        data,
	})
}

// SendSubscriptionConfirmation sends a subscription confirmation email
func (c *EmailClient) SendSubscriptionConfirmation(to string, data *SubscriptionData) error {
	return c.Send(&EmailData{
		To:          to,
		Subject:     fmt.Sprintf("Dein %s-Abonnement ist aktiv", data.PlanName),
		TemplateKey: "subscription_confirmation",
		Data:        data,
	})
}

// Email templates
const userInvitationTemplate = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Team-Einladung</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { text-align: center; padding: 20px 0; border-bottom: 2px solid #10b981; }
        .logo { font-size: 24px; font-weight: bold; color: #10b981; }
        .content { padding: 30px 0; }
        .button { display: inline-block; background: #10b981; color: white; padding: 14px 28px; text-decoration: none; border-radius: 8px; font-weight: 600; margin: 20px 0; }
        .button:hover { background: #059669; }
        .footer { text-align: center; padding-top: 20px; border-top: 1px solid #e5e7eb; color: #6b7280; font-size: 14px; }
        .role-badge { display: inline-block; background: #dbeafe; color: #1d4ed8; padding: 4px 12px; border-radius: 9999px; font-size: 14px; font-weight: 500; }
    </style>
</head>
<body>
    <div class="header">
        <div class="logo">🍸 GinVault</div>
    </div>
    <div class="content">
        <h2>Hallo{{if .RecipientName}} {{.RecipientName}}{{end}},</h2>
        <p><strong>{{.InviterName}}</strong> hat dich eingeladen, dem Team von <strong>{{.TenantName}}</strong> beizutreten.</p>
        <p>Deine Rolle: <span class="role-badge">{{.Role}}</span></p>
        <p>Klicke auf den Button unten, um die Einladung anzunehmen und dein Konto einzurichten:</p>
        <p style="text-align: center;">
            <a href="{{.InviteLink}}" class="button">Einladung annehmen</a>
        </p>
        <p style="color: #6b7280; font-size: 14px;">Dieser Link ist {{.ExpiresIn}} gültig.</p>
        <p>Falls du diese Einladung nicht erwartet hast, kannst du diese E-Mail ignorieren.</p>
    </div>
    <div class="footer">
        <p>&copy; 2026 GinVault. Alle Rechte vorbehalten.</p>
    </div>
</body>
</html>`

const passwordResetTemplate = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Passwort zurücksetzen</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { text-align: center; padding: 20px 0; border-bottom: 2px solid #10b981; }
        .logo { font-size: 24px; font-weight: bold; color: #10b981; }
        .content { padding: 30px 0; }
        .button { display: inline-block; background: #10b981; color: white; padding: 14px 28px; text-decoration: none; border-radius: 8px; font-weight: 600; margin: 20px 0; }
        .footer { text-align: center; padding-top: 20px; border-top: 1px solid #e5e7eb; color: #6b7280; font-size: 14px; }
    </style>
</head>
<body>
    <div class="header">
        <div class="logo">🍸 GinVault</div>
    </div>
    <div class="content">
        <h2>Passwort zurücksetzen</h2>
        <p>Hallo{{if .RecipientName}} {{.RecipientName}}{{end}},</p>
        <p>Du hast angefordert, dein Passwort zurückzusetzen. Klicke auf den Button unten, um ein neues Passwort zu erstellen:</p>
        <p style="text-align: center;">
            <a href="{{.ResetLink}}" class="button">Neues Passwort erstellen</a>
        </p>
        <p style="color: #6b7280; font-size: 14px;">Dieser Link ist {{.ExpiresIn}} gültig.</p>
        <p>Falls du kein neues Passwort angefordert hast, kannst du diese E-Mail ignorieren. Dein Passwort bleibt unverändert.</p>
    </div>
    <div class="footer">
        <p>&copy; 2026 GinVault. Alle Rechte vorbehalten.</p>
    </div>
</body>
</html>`

const welcomeTemplate = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Willkommen</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { text-align: center; padding: 20px 0; border-bottom: 2px solid #10b981; }
        .logo { font-size: 24px; font-weight: bold; color: #10b981; }
        .content { padding: 30px 0; }
        .button { display: inline-block; background: #10b981; color: white; padding: 14px 28px; text-decoration: none; border-radius: 8px; font-weight: 600; margin: 20px 0; }
        .feature { padding: 10px 0; border-bottom: 1px solid #e5e7eb; }
        .footer { text-align: center; padding-top: 20px; border-top: 1px solid #e5e7eb; color: #6b7280; font-size: 14px; }
    </style>
</head>
<body>
    <div class="header">
        <div class="logo">🍸 GinVault</div>
    </div>
    <div class="content">
        <h2>Willkommen bei {{.TenantName}}!</h2>
        <p>Hallo {{.RecipientName}},</p>
        <p>Dein Konto wurde erfolgreich erstellt. Du kannst dich jetzt anmelden und deine Gin-Sammlung verwalten.</p>
        <p style="text-align: center;">
            <a href="{{.LoginLink}}" class="button">Jetzt anmelden</a>
        </p>
        <h3>Was du tun kannst:</h3>
        <div class="feature">✓ Gins katalogisieren mit Fotos und Details</div>
        <div class="feature">✓ Verkostungsnotizen hinzufügen</div>
        <div class="feature">✓ Statistiken über deine Sammlung anzeigen</div>
        <div class="feature">✓ Deine Sammlung durchsuchen und filtern</div>
    </div>
    <div class="footer">
        <p>&copy; 2026 GinVault. Alle Rechte vorbehalten.</p>
    </div>
</body>
</html>`

const subscriptionConfirmationTemplate = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Abonnement bestätigt</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { text-align: center; padding: 20px 0; border-bottom: 2px solid #10b981; }
        .logo { font-size: 24px; font-weight: bold; color: #10b981; }
        .content { padding: 30px 0; }
        .plan-box { background: #f0fdf4; border: 2px solid #10b981; border-radius: 12px; padding: 20px; margin: 20px 0; text-align: center; }
        .plan-name { font-size: 24px; font-weight: bold; color: #10b981; }
        .plan-price { font-size: 18px; color: #333; margin-top: 8px; }
        .detail-row { display: flex; justify-content: space-between; padding: 8px 0; border-bottom: 1px solid #e5e7eb; }
        .footer { text-align: center; padding-top: 20px; border-top: 1px solid #e5e7eb; color: #6b7280; font-size: 14px; }
    </style>
</head>
<body>
    <div class="header">
        <div class="logo">🍸 GinVault</div>
    </div>
    <div class="content">
        <h2>Dein Abonnement ist aktiv! 🎉</h2>
        <p>Hallo {{.RecipientName}},</p>
        <p>Vielen Dank für dein Upgrade! Dein neues Abonnement ist jetzt aktiv.</p>
        <div class="plan-box">
            <div class="plan-name">{{.PlanName}}</div>
            <div class="plan-price">{{.Amount}} / {{.BillingCycle}}</div>
        </div>
        <h3>Abrechnungsdetails</h3>
        <div class="detail-row">
            <span>Nächste Abrechnung:</span>
            <strong>{{.NextBilling}}</strong>
        </div>
        <p style="margin-top: 20px;">Du hast jetzt Zugang zu allen Premium-Funktionen. Viel Spaß mit deiner erweiterten Gin-Sammlung!</p>
    </div>
    <div class="footer">
        <p>&copy; 2026 GinVault. Alle Rechte vorbehalten.</p>
    </div>
</body>
</html>`
