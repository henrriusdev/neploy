package service

import (
	"context"
	"fmt"

	"github.com/resend/resend-go/v2"
	"neploy.dev/config"
	"neploy.dev/pkg/logger"
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
		From:    fmt.Sprintf("%s <%s>", config.Env.ResendFromName, config.Env.ResendFromEmail),
		To:      []string{to},
		Subject: fmt.Sprintf("Invitation to join %s on Neploy", teamName),
		Html: fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
</head>
<body style="font-family: -apple-system, system-ui, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; margin: 0; padding: 0; background-color: #ffffff;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <img src="https://lh3.googleusercontent.com/d/1McJEcUM6u69CasiERZNpf2sIh1jEg7Zz" alt="Neploy" style="width: 48px; height: 48px; margin-bottom: 24px;">
        
        <h1 style="color: #111827; font-size: 24px; margin-bottom: 24px;">Hey %s! You've been invited!</h1>
        
        <p style="color: #374151; font-size: 16px; line-height: 24px; margin-bottom: 24px;">
            You've been invited to join <strong>%s</strong> as a <strong>%s</strong>.
        </p>
        
        <a href="%s" style="display: inline-block; background-color: #0070f3; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; font-weight: 500; margin-bottom: 24px;">
            Accept Invitation
        </a>
        
        <p style="color: #6B7280; font-size: 14px; line-height: 20px; margin-top: 32px;">
            If you didn't expect this invitation, you can safely ignore this email.
        </p>
        
        <hr style="border: none; border-top: 1px solid #E5E7EB; margin: 32px 0;">
        
        <div style="color: #6B7280; font-size: 12px; line-height: 16px;">
            <p style="margin: 0;">Best regards,</p>
            <p style="margin: 8px 0;">The Neploy Team</p>
            <img src="https://lh3.googleusercontent.com/d/1McJEcUM6u69CasiERZNpf2sIh1jEg7Zz" alt="Neploy" style="width: 32px; height: 32px; margin: 16px 0;">
            <p style="margin: 0;">Neploy - Modern Deployment Platform</p>
        </div>
    </div>
</body>
</html>`, role, teamName, role, inviteLink),
	}

	res, err := e.client.Emails.Send(params)
	logger.Info("send invitation: %v", res)
	logger.Error("send invitation: %v", err)
	return err
}
