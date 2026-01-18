// +build ignore

package main

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

func main() {
	// Get target email from command line or use default
	targetEmail := "stefan.buell@outlook.de"
	if len(os.Args) > 1 {
		targetEmail = os.Args[1]
	}

	// SMTP Configuration from environment
	host := getEnv("SMTP_HOST", "smtp.hostinger.com")
	port := getEnv("SMTP_PORT", "465")
	username := getEnv("SMTP_USERNAME", "info@ginvault.cloud")
	password := os.Getenv("SMTP_PASSWORD")
	fromEmail := getEnv("SMTP_FROM_EMAIL", "info@ginvault.cloud")
	fromName := getEnv("SMTP_FROM_NAME", "GinVault")

	if password == "" {
		fmt.Println("ERROR: SMTP_PASSWORD environment variable is required")
		os.Exit(1)
	}

	fmt.Printf("Sending test email to: %s\n", targetEmail)
	fmt.Printf("SMTP Host: %s:%s\n", host, port)
	fmt.Printf("From: %s <%s>\n", fromName, fromEmail)

	// Build email message
	subject := "GinVault - Test E-Mail"
	body := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Test E-Mail</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { text-align: center; padding: 20px 0; border-bottom: 2px solid #D4A857; background: #0F3D2E; border-radius: 8px 8px 0 0; }
        .logo { font-size: 24px; font-weight: bold; color: #D4A857; }
        .content { padding: 30px; background: #1a1a2e; color: #fff; }
        .success { background: #10b981; color: white; padding: 15px; border-radius: 8px; text-align: center; margin: 20px 0; }
        .footer { text-align: center; padding: 20px; background: #0F3D2E; color: #D4A857; font-size: 14px; border-radius: 0 0 8px 8px; }
    </style>
</head>
<body style="background: #0a0a14; padding: 20px;">
    <div style="max-width: 600px; margin: 0 auto;">
        <div class="header">
            <div class="logo">🍸 GinVault</div>
        </div>
        <div class="content">
            <h2 style="color: #D4A857;">E-Mail-Test erfolgreich!</h2>
            <div class="success">
                ✓ SMTP-Konfiguration funktioniert einwandfrei
            </div>
            <p>Diese Test-E-Mail bestätigt, dass die E-Mail-Konfiguration für GinVault korrekt eingerichtet ist.</p>
            <p><strong>Konfiguration:</strong></p>
            <ul>
                <li>SMTP Host: smtp.hostinger.com</li>
                <li>Port: 465 (TLS)</li>
                <li>Absender: info@ginvault.cloud</li>
            </ul>
        </div>
        <div class="footer">
            <p>&copy; 2026 GinVault. Alle Rechte vorbehalten.</p>
        </div>
    </div>
</body>
</html>`

	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("From: %s <%s>\r\n", fromName, fromEmail))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", targetEmail))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(body)

	// Send via TLS
	addr := fmt.Sprintf("%s:%s", host, port)

	tlsConfig := &tls.Config{
		ServerName: host,
	}

	fmt.Println("Connecting to SMTP server...")
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		fmt.Printf("ERROR: Failed to connect: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		fmt.Printf("ERROR: Failed to create SMTP client: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	fmt.Println("Authenticating...")
	auth := smtp.PlainAuth("", username, password, host)
	if err := client.Auth(auth); err != nil {
		fmt.Printf("ERROR: Authentication failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Setting sender...")
	if err := client.Mail(fromEmail); err != nil {
		fmt.Printf("ERROR: MAIL command failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Setting recipient...")
	if err := client.Rcpt(targetEmail); err != nil {
		fmt.Printf("ERROR: RCPT command failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Sending message...")
	writer, err := client.Data()
	if err != nil {
		fmt.Printf("ERROR: DATA command failed: %v\n", err)
		os.Exit(1)
	}

	if _, err := writer.Write([]byte(msg.String())); err != nil {
		fmt.Printf("ERROR: Failed to write message: %v\n", err)
		os.Exit(1)
	}

	if err := writer.Close(); err != nil {
		fmt.Printf("ERROR: Failed to close writer: %v\n", err)
		os.Exit(1)
	}

	if err := client.Quit(); err != nil {
		fmt.Printf("WARNING: Quit failed: %v\n", err)
	}

	fmt.Println("")
	fmt.Println("✓ SUCCESS! Test email sent to:", targetEmail)
	fmt.Println("  Check your inbox (and spam folder)")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
