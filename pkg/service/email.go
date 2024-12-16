package service

import (
	"context"
	"fmt"

	"github.com/resend/resend-go/v2"
	"neploy.dev/config"
)

type Email interface {
	SendInvitation(ctx context.Context, to, teamName, role, inviteLink string) error
}

type email struct {
	client *resend.Client
}

func NewEmail() Email {
	client := resend.NewClient(config.Env.ResendAPIKey)
	return &email{client: client}
}

func (e *email) SendInvitation(ctx context.Context, to, teamName, role, inviteLink string) error {
	params := &resend.SendEmailRequest{
		From:    "Neploy <onboarding@resend.dev>", // You can use this for testing
		To:      []string{to},
		Subject: fmt.Sprintf("Invitation to join %s on Neploy", teamName),
		Html: fmt.Sprintf(`
			<div style="font-family: sans-serif; max-width: 600px; margin: 0 auto;">
				<h2>You've been invited to join %s</h2>
				<p>You've been invited to join %s as a %s.</p>
				<p>Click the button below to accept the invitation:</p>
				<a href="%s" style="display: inline-block; background: #0070f3; color: white; padding: 12px 24px; text-decoration: none; border-radius: 5px; margin: 16px 0;">
					Accept Invitation
				</a>
				<p>Or copy and paste this URL into your browser:</p>
				<p>%s</p>
				<p>This invitation will expire in 7 days.</p>
			</div>
		`, teamName, teamName, role, inviteLink, inviteLink),
	}

	_, err := e.client.Emails.Send(params)
	return err
}
