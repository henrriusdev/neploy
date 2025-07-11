package email

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/resend/resend-go/v2"
	"neploy.dev/config"
)

//go:embed templates/*
var emailTmpl embed.FS

type Email struct {
	client *resend.Client
}

func NewEmail() *Email {
	client := resend.NewClient(config.Env.ResendAPIKey)
	return &Email{client: client}
}

func (e *Email) SendInvitation(ctx context.Context, to, teamName, role, inviteLink string, language string) error {
	// Default to English if no language is specified
	if language == "" {
		language = "en"
	}

	// Extract recipient name from email address
	recipientName := to
	if atIndex := strings.Index(to, "@"); atIndex > 0 {
		recipientName = to[:atIndex]
	}

	// Create invitation data with translations
	data := InvitationData{
		RecipientName: recipientName,
		CompanyName:   "Neploy",
		LogoURL:       "https://lh3.googleusercontent.com/d/1McJEcUM6u69CasiERZNpf2sIh1jEg7Zz",
		TeamName:      teamName,
		Role:          role,
		InviteLink:    inviteLink,
		CurrentYear:   time.Now().Year(),
		Language:      language,
		Translations:  GetInvitationTranslations(language),
	}

	// Parse and execute the template
	tmpl, err := template.ParseFS(emailTmpl, "templates/invitation_email.gohtml")
	if err != nil {
		return fmt.Errorf("parsing invitation template: %w", err)
	}

	var htmlBody bytes.Buffer
	if err := tmpl.Execute(&htmlBody, data); err != nil {
		return fmt.Errorf("executing invitation template: %w", err)
	}

	// Create and send the email
	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("%s <%s>", config.Env.ResendFromName, config.Env.ResendFromEmail),
		To:      []string{to},
		Subject: data.Translations.Title + " " + teamName,
		Html:    htmlBody.String(),
	}

	_, err = e.client.Emails.Send(params)
	return err
}

func (e *Email) SendPasswordReset(to string, data PasswordResetData) error {
	tmpl, err := template.ParseFS(emailTmpl, "templates/password_reset_email.gohtml")
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}

	var htmlBody bytes.Buffer
	if err := tmpl.Execute(&htmlBody, data); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}

	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("%s <%s>", config.Env.ResendFromName, config.Env.ResendFromEmail),
		To:      []string{to},
		Subject: data.Translations.Title,
		Html:    htmlBody.String(),
	}

	_, err = e.client.Emails.Send(params)
	return err
}
